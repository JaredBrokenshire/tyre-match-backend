package services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/random"
	"path/filepath"
	m "tyre-match-backend/db/models"
	"tyre-match-backend/db/repositories"
	"tyre-match-backend/pkg/file_storage"
)

type TyreImpressionDTO struct {
	PixelsPerInch float32 `json:"pixels_per_inch"`
	ROITop        int     `json:"roi_top"`
	ROILeft       int     `json:"roi_left"`
	ROIRight      int     `json:"roi_right"`
	ROIBottom     int     `json:"roi_bottom"`
}

type TyreImpressionService struct {
	repo                   *repositories.TyreImpressionRepository
	fileRepo               *repositories.FileRepository
	fileStore              file_storage.Store
	validator              FileValidator
	imageProcessingService ImageProcessingServiceInterface
}

func NewTyreImpressionService(
	repo *repositories.TyreImpressionRepository,
	fileRepo *repositories.FileRepository,
	fileStore file_storage.Store,
	validator FileValidator,
	imageProcessingService ImageProcessingServiceInterface,
) *TyreImpressionService {
	return &TyreImpressionService{
		repo:                   repo,
		fileRepo:               fileRepo,
		fileStore:              fileStore,
		validator:              validator,
		imageProcessingService: imageProcessingService,
	}
}

func (s *TyreImpressionService) List(page, pageSize int) ([]*m.TyreImpression, int, int, int) {
	var scopes []func(db *gorm.DB) *gorm.DB

	return s.repo.List(page, pageSize, scopes)
}

func (s *TyreImpressionService) Get(id uint) (*m.TyreImpression, error) {
	tyreImpression := s.repo.GetByID(id)
	if tyreImpression == nil {
		return nil, NotFoundError
	}

	return tyreImpression, nil
}

func (s *TyreImpressionService) Create(dto TyreImpressionDTO) (*m.TyreImpression, error) {
	tyreImpression := &m.TyreImpression{
		PixelsPerInch: dto.PixelsPerInch,
		ROITop:        dto.ROITop,
		ROILeft:       dto.ROILeft,
		ROIRight:      dto.ROIRight,
		ROIBottom:     dto.ROIBottom,
		Status:        m.ProcessingStatusUploaded,
	}
	if err := s.repo.Create(tyreImpression); err != nil {
		return nil, err
	}

	return tyreImpression, nil
}

func (s *TyreImpressionService) Upload(id uint, file UploadedFile) (*m.TyreImpression, error) {
	tyreImpression := s.repo.GetByID(id)
	if tyreImpression == nil {
		return nil, NotFoundError
	}
	if len(tyreImpression.Images) != 0 {
		return nil, AlreadyExistsError
	}

	if err := s.validator.Validate(file); err != nil {
		return nil, fmt.Errorf("%w: %w", InvalidUploadError, err)
	}

	extension := filepath.Ext(file.Name)
	filename := random.String(32) + extension
	path := fmt.Sprintf("tyre-impressions/%v/%v", id, m.FileTypeOriginal)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open upload: %w", err)
	}
	defer src.Close()

	if err := s.fileStore.Save(path, src, filename); err != nil {
		return nil, fmt.Errorf("%w: %w", FileStoreError, err)
	}

	fileRecord := &m.File{
		Model: m.FileModelTyreImpression, ModelId: tyreImpression.ID,
		Name: filename, Location: path, FileType: m.FileTypeOriginal,
	}
	if err := s.fileRepo.Create(fileRecord); err != nil {
		return nil, fmt.Errorf("%w: %w", DatabaseError, err)
	}

	tyreImpression.Images[m.FileTypeOriginal] = fileRecord
	tyreImpression.Status = m.ProcessingStatusProcessing
	if err := s.repo.Update(tyreImpression); err != nil {
		return nil, fmt.Errorf("%w: %w", DatabaseError, err)
	}

	if err := s.imageProcessingService.ProcessTyreImpression(tyreImpression); err != nil {
		tyreImpression.Status = m.ProcessingStatusFailed
		if updateErr := s.repo.Update(tyreImpression); updateErr != nil {
			return nil, fmt.Errorf("processing failed: %w (also unable to mark as failed: %w)", err, updateErr)
		}
		return nil, fmt.Errorf("%w: %w", ProcessingError, err)
	}

	// Reload so the response includes the generated files
	tyreImpression = s.repo.GetByID(id)

	return tyreImpression, nil
}

func (s *TyreImpressionService) Delete(id uint) error {
	tyreImpression := s.repo.GetByID(id)
	if tyreImpression == nil {
		return NotFoundError
	}

	return s.repo.Delete(tyreImpression)
}

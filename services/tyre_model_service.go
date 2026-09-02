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

type TyreModelDTO struct {
	Manufacturer      string
	ModelName         string
	Category          string
	VehicleType       string
	WidthMm           int
	AspectRatio       int
	RimDiameterInches int
	GrooveCount       int
	PatternType       string
	PixelsPerInch     float32
	ROITop            int
	ROILeft           int
	ROIRight          int
	ROIBottom         int
}

type TyreModelService struct {
	repo                   *repositories.TyreModelRepository
	fileRepo               *repositories.FileRepository
	fileStore              file_storage.Store
	validator              FileValidator
	imageProcessingService ImageProcessingServiceInterface
}

func NewTyreModelService(
	repo *repositories.TyreModelRepository,
	fileRepo *repositories.FileRepository,
	fileStore file_storage.Store,
	validator FileValidator,
	imageProcessingService ImageProcessingServiceInterface,
) *TyreModelService {
	return &TyreModelService{
		repo:                   repo,
		fileRepo:               fileRepo,
		fileStore:              fileStore,
		validator:              validator,
		imageProcessingService: imageProcessingService,
	}
}

func (s *TyreModelService) List(page, pageSize int, search string) ([]*m.TyreModel, int, int, int) {
	var scopes []func(db *gorm.DB) *gorm.DB

	if search != "" {
		searchTerm := fmt.Sprintf("%%%s%%", search)
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("manufacturer LIKE ? OR model_name LIKE ?", searchTerm, searchTerm)
		})
	}

	return s.repo.List(page, pageSize, scopes)
}

func (s *TyreModelService) Get(id uint) (*m.TyreModel, error) {
	tyreModel := s.repo.GetByID(id)
	if tyreModel == nil {
		return nil, NotFoundError
	}

	return tyreModel, nil
}

func (s *TyreModelService) Create(dto TyreModelDTO) (*m.TyreModel, error) {
	tyreModel := &m.TyreModel{
		Manufacturer:      dto.Manufacturer,
		ModelName:         dto.ModelName,
		WidthMm:           dto.WidthMm,
		AspectRatio:       dto.AspectRatio,
		RimDiameterInches: dto.RimDiameterInches,
		GrooveCount:       dto.GrooveCount,
		ROITop:            dto.ROITop,
		ROILeft:           dto.ROILeft,
		ROIRight:          dto.ROIRight,
		ROIBottom:         dto.ROIBottom,
		PixelsPerInch:     dto.PixelsPerInch,
		Status:            m.ProcessingStatusUploaded,
	}

	if err := s.repo.Create(tyreModel); err != nil {
		tyreModel.Status = m.ProcessingStatusFailed
		_ = s.repo.Update(tyreModel)

		return nil, err
	}

	return tyreModel, nil
}

func (s *TyreModelService) Update(id uint, dto TyreModelDTO) (*m.TyreModel, error) {
	tyreModel := s.repo.GetByID(id)
	if tyreModel == nil {
		return nil, NotFoundError
	}

	tyreModel.Manufacturer = dto.Manufacturer
	tyreModel.ModelName = dto.ModelName
	tyreModel.WidthMm = dto.WidthMm
	tyreModel.AspectRatio = dto.AspectRatio
	tyreModel.RimDiameterInches = dto.RimDiameterInches
	tyreModel.GrooveCount = dto.GrooveCount

	if err := s.repo.Update(tyreModel); err != nil {
		return nil, err
	}

	return tyreModel, nil
}

func (s *TyreModelService) Upload(id uint, file UploadedFile) (*m.TyreModel, error) {
	tyreModel := s.repo.GetByID(id)
	if tyreModel == nil {
		return nil, NotFoundError
	}
	if len(tyreModel.Images) != 0 {
		return nil, AlreadyExistsError
	}

	if err := s.validator.Validate(file); err != nil {
		return nil, InvalidUploadError
	}

	extension := filepath.Ext(file.Name)
	filename := random.String(32) + extension
	path := fmt.Sprintf("tyre-models/%v/%v", id, m.FileTypeOriginal)

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	if err := s.fileStore.Save(path, src, filename); err != nil {
		return nil, FileStoreError
	}

	record := &m.File{
		Model:    m.FileModelTyreModel,
		ModelId:  tyreModel.ID,
		Name:     filename,
		Location: path,
		FileType: m.FileTypeOriginal,
	}
	if err := s.fileRepo.Create(record); err != nil {
		return nil, DatabaseError
	}

	tyreModel.Images[m.FileTypeOriginal] = record
	tyreModel.Status = m.ProcessingStatusProcessing
	if err := s.repo.Update(tyreModel); err != nil {
		return nil, fmt.Errorf("%w: %w", DatabaseError, err)
	}

	if err := s.imageProcessingService.Process(tyreModel.ID, m.FileModelTyreModel); err != nil {
		tyreModel.Status = m.ProcessingStatusFailed
		if updateErr := s.repo.Update(tyreModel); updateErr != nil {
			return nil, fmt.Errorf("processing failed: %w (also unable to mark as failed: %w)", err, updateErr)
		}
		return nil, fmt.Errorf("%w: %w", ProcessingError, err)
	}

	// Reload so the response includes the generated files
	tyreModel = s.repo.GetByID(id)

	return tyreModel, nil
}

func (s *TyreModelService) Delete(id uint) error {
	tyreModel := s.repo.GetByID(id)
	if tyreModel == nil {
		return NotFoundError
	}

	return s.repo.Delete(tyreModel)
}

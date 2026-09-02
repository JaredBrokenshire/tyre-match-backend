package services

import (
	"fmt"
	"github.com/labstack/gommon/random"
	cv "gocv.io/x/gocv"
	_ "golang.org/x/image/webp"
	_ "image/jpeg"
	"path/filepath"
	"tyre-match-backend/cv/processors"
	m "tyre-match-backend/db/models"
	"tyre-match-backend/db/repositories"
	"tyre-match-backend/pkg/file_storage"
)

type ImageProcessingServiceInterface interface {
	Process(id uint, model string) error
	SaveStage(id uint, model, fileType string, resultImage *cv.Mat) error
}

type ImageProcessor interface {
	Process(id uint, model string) error
}

type ImageProcessingService struct {
	TyreImpressionRepo *repositories.TyreImpressionRepository
	TyreModelRepo      *repositories.TyreModelRepository
	*FileService
	*repositories.FileRepository
	FileStore  file_storage.Store
	Processors []processors.Processor
}

func NewImageProcessingService(impressionRepo *repositories.TyreImpressionRepository, modelRepo *repositories.TyreModelRepository, fileService *FileService, fileRepo *repositories.FileRepository, fileStore file_storage.Store) *ImageProcessingService {
	return &ImageProcessingService{
		TyreImpressionRepo: impressionRepo,
		TyreModelRepo:      modelRepo,
		FileService:        fileService,
		FileRepository:     fileRepo,
		FileStore:          fileStore,
		Processors:         []processors.Processor{},
	}
}

func (s *ImageProcessingService) Process(id uint, model string) error {
	original, ROITop, ROIRight, ROIBottom, ROILeft, err := s.extractROIAndOriginalImage(id, model)
	if err != nil {
		return err
	}

	imagePath := filepath.Join(s.FileStore.GetStorageLocation(), original.Location, original.Name)
	grayscaleImage := cv.IMRead(imagePath, cv.IMReadGrayScale)
	if grayscaleImage.Empty() {
		return fmt.Errorf("image is empty")
	}
	defer grayscaleImage.Close()

	// Build the full processor pipeline for this impression. NormalisationProcessor
	// is constructed here because it requires the impression's per-impression ROI
	// values. MorphologyProcessor is also constructed here because it requires the
	// impression's calibrated PixelsPerInch value.
	pipeline := []processors.Processor{
		processors.NewNormalisationProcessor(ROITop, ROILeft, ROIRight, ROIBottom),
		processors.NewEnhancementProcessor(),
		processors.NewBinaryProcessor(),
	}

	currentImage := &grayscaleImage

	for _, processor := range pipeline {
		result, err := processor.Process(currentImage)
		if err != nil {
			return fmt.Errorf("%s stage: %v", processor.GetName(), err)
		}
		defer result.Close()

		if err := s.SaveStage(id, model, processor.GetFileType(), result); err != nil {
			return err
		}

		// The result becomes the input for the next processor.
		if currentImage != &grayscaleImage {
			currentImage.Close()
		}

		currentImage = result
	}

	switch model {
	case m.FileModelTyreImpression:
		tyreImpression := s.TyreImpressionRepo.GetByID(id)
		if tyreImpression == nil {
			return fmt.Errorf("tyre impression with id %v not found", id)
		}

		tyreImpression.Status = m.ProcessingStatusProcessed
		if err := s.TyreImpressionRepo.Update(tyreImpression); err != nil {
			return ProcessingError
		}
	case m.FileModelTyreModel:
		tyreModel := s.TyreModelRepo.GetByID(id)
		if tyreModel == nil {
			return fmt.Errorf("tyre model with id %v not found", id)
		}

		tyreModel.Status = m.ProcessingStatusProcessed
		if err := s.TyreModelRepo.Update(tyreModel); err != nil {
			return ProcessingError
		}
	}

	return nil
}

func (s *ImageProcessingService) SaveStage(id uint, model, fileType string, resultImage *cv.Mat) error {
	if resultImage == nil || resultImage.Empty() {
		return fmt.Errorf("%w: can not save empty result image", ProcessingError)
	}

	var targetDir string
	switch model {
	case m.FileModelTyreImpression:
		targetDir = "tyre-impressions"
	case m.FileModelTyreModel:
		targetDir = "tyre-models"
	default:
		return fmt.Errorf("unsupported model %q", model)
	}

	encoded, err := cv.IMEncode(".png", *resultImage)
	if err != nil {
		return fmt.Errorf("error encoding result image: %w", err)
	}
	defer encoded.Close()

	request := SaveFileRequest{
		Data:            encoded.GetBytes(),
		Name:            fmt.Sprintf("%v.png", random.String(32)),
		TargetDirectory: fmt.Sprintf("%v/%v/%v", targetDir, id, fileType),
		Model:           model,
		ModelId:         id,
		FileType:        fileType,
		Extension:       "png",
	}

	// Create file store instance
	err = s.FileService.SaveFile(request)
	if err != nil {
		return fmt.Errorf("error saving image processing stage: %w", err)
	}

	// Create DB record
	fileRecord := &m.File{
		Model:    model,
		ModelId:  id,
		FileType: fileType,
		Name:     request.Name,
		Location: request.TargetDirectory,
	}
	err = s.FileRepository.Create(fileRecord)
	if err != nil {
		return fmt.Errorf("error saving image processing stage to db: %w", err)
	}

	return nil
}

func (s *ImageProcessingService) extractROIAndOriginalImage(id uint, model string) (*m.File, int, int, int, int, error) {
	var tyreImpression *m.TyreImpression
	var tyreModel *m.TyreModel
	var original *m.File
	var ROITop, ROIRight, ROIBottom, ROILeft int

	switch model {
	case m.FileModelTyreImpression:
		tyreImpression = s.TyreImpressionRepo.GetByID(id)
		if tyreImpression == nil {
			return nil, 0, 0, 0, 0, fmt.Errorf("no tyre impression found with id %d", id)
		}

		original = tyreImpression.Images[m.FileTypeOriginal]

		ROITop, ROIRight, ROIBottom, ROILeft = tyreImpression.ROITop, tyreImpression.ROIRight, tyreImpression.ROIBottom, tyreImpression.ROILeft
	case m.FileModelTyreModel:
		tyreModel = s.TyreModelRepo.GetByID(id)
		if tyreModel == nil {
			return nil, 0, 0, 0, 0, fmt.Errorf("no tyre impression found with id %d", id)
		}

		original = tyreModel.Images[m.FileTypeOriginal]
		ROITop, ROIRight, ROIBottom, ROILeft = tyreModel.ROITop, tyreModel.ROIRight, tyreModel.ROIBottom, tyreModel.ROILeft
	default:
		return nil, 0, 0, 0, 0, fmt.Errorf("unsupported model %q", model)
	}

	if original == nil {
		return nil, 0, 0, 0, 0, fmt.Errorf("original image is missing")
	}

	return original, ROITop, ROIRight, ROIBottom, ROILeft, nil
}

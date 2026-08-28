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
	*FileService
	*repositories.FileRepository
	FileStore  file_storage.Store
	Processors []processors.Processor
}

func NewImageProcessingService(impressionRepo *repositories.TyreImpressionRepository, fileService *FileService, fileRepo *repositories.FileRepository, fileStore file_storage.Store) *ImageProcessingService {
	return &ImageProcessingService{
		TyreImpressionRepo: impressionRepo,
		FileService:        fileService,
		FileRepository:     fileRepo,
		FileStore:          fileStore,
		Processors: []processors.Processor{
			processors.NewEnhancementProcessor(),
		},
	}
}

func (s *ImageProcessingService) Process(id uint, model string) error {
	if model != m.FileModelTyreImpression {
		return fmt.Errorf("unsupported model %q", model)
	}

	impression := s.TyreImpressionRepo.GetByID(id)
	if impression == nil {
		return NotFoundError
	}

	original := impression.Images[m.FileTypeOriginal]
	if original == nil {
		return fmt.Errorf("original image is missing")
	}

	imagePath := filepath.Join(s.FileStore.GetStorageLocation(), original.Location, original.Name)
	grayscaleImage := cv.IMRead(imagePath, cv.IMReadGrayScale)
	if grayscaleImage.Empty() {
		return fmt.Errorf("image is empty")
	}
	defer grayscaleImage.Close()

	currentImage := &grayscaleImage

	for _, processor := range s.Processors {
		result, err := processor.Process(currentImage)
		if err != nil {
			return fmt.Errorf("%s stage: %v", processor.GetName(), err)
		}
		defer result.Close()

		if err := s.SaveStage(impression.ID, m.FileModelTyreImpression, processor.GetFileType(), result); err != nil {
			return err
		}

		// The result becomes the input for the next processor.
		if currentImage != &grayscaleImage {
			currentImage.Close()
		}

		currentImage = result
	}

	impression.Status = m.TyreImpressionStatusProcessed
	if err := s.TyreImpressionRepo.Update(impression); err != nil {
		return ProcessingError
	}

	return nil
}

func (s *ImageProcessingService) SaveStage(id uint, model, fileType string, resultImage *cv.Mat) error {
	if resultImage == nil || resultImage.Empty() {
		return fmt.Errorf("%w: can not save empty result image", ProcessingError)
	}

	encoded, err := cv.IMEncode(".png", *resultImage)
	if err != nil {
		return fmt.Errorf("error encoding result image: %w", err)
	}
	defer encoded.Close()

	request := SaveFileRequest{
		Data:            encoded.GetBytes(),
		Name:            fmt.Sprintf("%v.png", random.String(32)),
		TargetDirectory: fmt.Sprintf("tyre-impressions/%v/%v", id, fileType),
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

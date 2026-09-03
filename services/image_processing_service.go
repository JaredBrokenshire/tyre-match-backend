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
	ProcessTyreImpression(tyreImpression *m.TyreImpression) error
	ProcessTyreModel(tyreModel *m.TyreModel) error
	SaveStage(id uint, model, fileType string, image *cv.Mat) error
}

type ImageProcessingService struct {
	TyreImpressionRepo *repositories.TyreImpressionRepository
	TyreModelRepo      *repositories.TyreModelRepository
	*FileService
	*repositories.FileRepository
	FileStore file_storage.Store
}

func NewImageProcessingService(impressionRepo *repositories.TyreImpressionRepository, modelRepo *repositories.TyreModelRepository, fileService *FileService, fileRepo *repositories.FileRepository, fileStore file_storage.Store) *ImageProcessingService {
	return &ImageProcessingService{
		TyreImpressionRepo: impressionRepo,
		TyreModelRepo:      modelRepo,
		FileService:        fileService,
		FileRepository:     fileRepo,
		FileStore:          fileStore,
	}
}

func (s *ImageProcessingService) ProcessTyreImpression(tyreImpression *m.TyreImpression) error {
	if len(tyreImpression.Images) == 0 {
		return fmt.Errorf("tyre impression original image missing")
	}

	original := tyreImpression.Images[m.FileTypeOriginal]

	grayscaleImage, err := s.readGrayscale(original)
	if err != nil {
		return err
	}
	defer grayscaleImage.Close()

	pipeline := []processors.Processor{
		processors.NewNormalisationProcessor(tyreImpression.ROITop, tyreImpression.ROILeft, tyreImpression.ROIRight, tyreImpression.ROIBottom),
		processors.NewEnhancementProcessor(),
		processors.NewBinaryProcessor(),
	}

	err = s.runPipeline(pipeline, tyreImpression.ID, m.FileModelTyreImpression, grayscaleImage)
	if err != nil {
		return err
	}

	tyreImpression.Status = m.ProcessingStatusProcessed
	if err := s.TyreImpressionRepo.Update(tyreImpression); err != nil {
		return ProcessingError
	}

	return nil
}

func (s *ImageProcessingService) ProcessTyreModel(tyreModel *m.TyreModel) error {
	if len(tyreModel.Images) == 0 {
		return fmt.Errorf("tyre impression original image missing")
	}

	original := tyreModel.Images[m.FileTypeOriginal]

	grayscaleImage, err := s.readGrayscale(original)
	if err != nil {
		return err
	}
	defer grayscaleImage.Close()

	pipeline := []processors.Processor{
		processors.NewNormalisationProcessor(tyreModel.ROITop, tyreModel.ROILeft, tyreModel.ROIRight, tyreModel.ROIBottom),
		processors.NewEnhancementProcessor(),
		processors.NewBinaryProcessor(),
	}

	err = s.runPipeline(pipeline, tyreModel.ID, m.FileModelTyreModel, grayscaleImage)
	if err != nil {
		return err
	}

	tyreModel.Status = m.ProcessingStatusProcessed
	if err := s.TyreModelRepo.Update(tyreModel); err != nil {
		return ProcessingError
	}

	return nil
}

func (s *ImageProcessingService) SaveStage(id uint, model, fileType string, image *cv.Mat) error {
	if image == nil || image.Empty() {
		return fmt.Errorf("%w: can not save empty image", ProcessingError)
	}

	var targetDir string
	switch model {
	case m.FileModelTyreImpression:
		targetDir = "tyre-impressions"
	case m.FileModelTyreModel:
		targetDir = "tyre-models"
	default:
		return fmt.Errorf("%w: unsupported model %v", ProcessingError, model)
	}

	encoded, err := cv.IMEncode(".png", *image)
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

func (s *ImageProcessingService) readGrayscale(original *m.File) (*cv.Mat, error) {
	imagePath := filepath.Join(s.FileStore.GetStorageLocation(), original.Location, original.Name)
	grayscaleImage := cv.IMRead(imagePath, cv.IMReadGrayScale)
	if grayscaleImage.Empty() {
		return nil, fmt.Errorf("image is empty")
	}

	return &grayscaleImage, nil
}

func (s *ImageProcessingService) runPipeline(pipeline []processors.Processor, id uint, model string, image *cv.Mat) error {
	currentImage := image
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
		if currentImage != image {
			currentImage.Close()
		}

		currentImage = result
	}

	return nil
}

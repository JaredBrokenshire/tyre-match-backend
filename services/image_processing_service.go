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
	if tyreImpression == nil {
		return errors.New("tyre impression is nil")
	}
	if len(tyreImpression.Images) == 0 {
		return fmt.Errorf("tyre impression original image missing")
	}

	original := tyreImpression.Images[m.FileTypeOriginal]
	if original == nil {
		return errors.New("tyre impression original image missing")
	}
	grayscaleImage, err := s.readGrayscale(original)
	if err != nil {
		return err
	}
	defer grayscaleImage.Close()

	pipeline := []processors.Processor{
		processors.NewNormalisationProcessor(tyreImpression.ROITop, tyreImpression.ROILeft, tyreImpression.ROIRight, tyreImpression.ROIBottom),
		processors.NewEnhancementProcessor(),
	}

	enhanced, err := s.runPipeline(pipeline, tyreImpression.ID, m.FileModelTyreImpression, grayscaleImage)
	if err != nil {
		return err
	}
	defer enhanced.Close()

	treadMask, err := s.runBinaryStages(
		enhanced,
		tyreImpression.ID,
		m.FileModelTyreImpression,
		tyreImpression.PixelsPerInch,
	)
	if err != nil {
		return err
	}
	defer treadMask.Close()

	//if err := s.extractImpressionFeaturesAndMatches(tyreImpression, treadMask); err != nil {
	//	return err
	//}

	if tyreImpression.Status != m.ProcessingStatusMatched {
		tyreImpression.Status = m.ProcessingStatusProcessed
	}
	if err := s.TyreImpressionRepo.Update(tyreImpression); err != nil {
		return ProcessingError
	}
	return nil
}

func (s *ImageProcessingService) ProcessTyreModel(tyreModel *m.TyreModel) error {
	if tyreModel == nil {
		return errors.New("tyre model is nil")
	}
	if len(tyreModel.Images) == 0 {
		return fmt.Errorf("tyre model original image missing")
	}

	original := tyreModel.Images[m.FileTypeOriginal]
	if original == nil {
		return errors.New("tyre model original image missing")
	}
	grayscaleImage, err := s.readGrayscale(original)
	if err != nil {
		return err
	}
	defer grayscaleImage.Close()

	pipeline := []processors.Processor{
		processors.NewNormalisationProcessor(tyreModel.ROITop, tyreModel.ROILeft, tyreModel.ROIRight, tyreModel.ROIBottom),
		processors.NewEnhancementProcessor(),
	}

	enhanced, err := s.runPipeline(pipeline, tyreModel.ID, m.FileModelTyreModel, grayscaleImage)
	if err != nil {
		return err
	}
	defer enhanced.Close()

	treadMask, err := s.runBinaryStages(
		enhanced,
		tyreModel.ID,
		m.FileModelTyreModel,
		tyreModel.PixelsPerInch,
	)
	if err != nil {
		return err
	}
	defer treadMask.Close()

	if s.FeatureRepo == nil {
		return fmt.Errorf("tyre model feature repository is not configured")
	}
	if s.FeatureExtractor == nil {
		return fmt.Errorf("tyre model feature extractor is not configured")
	}

	//fingerprint, err := s.FeatureExtractor.ExtractWithContext(treadMask, fe.ExtractionContext{
	//	ROITop:        0,
	//	ROILeft:       0,
	//	ROIRight:      treadMask.Cols(),
	//	ROIBottom:     treadMask.Rows(),
	//	PixelsPerInch: tyreModel.PixelsPerInch,
	//})
	//if err != nil {
	//	return fmt.Errorf("extract tyre model features: %w", err)
	//}
	//encoded, err := fe.Encode(fingerprint)
	//if err != nil {
	//	return fmt.Errorf("encode tyre model features: %w", err)
	//}
	//if err := s.FeatureRepo.Upsert(&m.TyreModelFeature{
	//	TyreModelID: tyreModel.ID,
	//	FeatureJSON: encoded,
	//}); err != nil {
	//	return fmt.Errorf("store tyre model features: %w", err)
	//}

	tyreModel.Status = m.ProcessingStatusProcessed
	if err := s.TyreModelRepo.Update(tyreModel); err != nil {
		return ProcessingError
	}
	return nil
}

type TyreMatch struct {
	Rank                     int     `json:"rank"`
	TyreModelID              uint    `json:"tyre_model_id"`
	Manufacturer             string  `json:"manufacturer"`
	ModelName                string  `json:"model_name"`
	Similarity               float32 `json:"similarity"`
	Distance                 float32 `json:"distance"`
	PatternSimilarity        float32 `json:"pattern_similarity"`
	LongitudinalSimilarity   float32 `json:"longitudinal_similarity"`
	StructureCountSimilarity float32 `json:"structure_count_similarity"`
	ScaleCompatibility       float32 `json:"scale_compatibility"`
}

func (s *ImageProcessingService) extractImpressionFeaturesAndMatches(impression *m.TyreImpression, treadMask *cv.Mat) error {
	if s.FeatureRepo == nil {
		return fmt.Errorf("tyre model feature repository is not configured")
	}
	if s.FeatureExtractor == nil {
		return fmt.Errorf("tyre model feature extractor is not configured")
	}
	fingerprint, err := s.FeatureExtractor.ExtractWithContext(treadMask, fe.ExtractionContext{
		ROITop:        0,
		ROILeft:       0,
		ROIRight:      treadMask.Cols(),
		ROIBottom:     treadMask.Rows(),
		PixelsPerInch: impression.PixelsPerInch,
	})
	if err != nil {
		return fmt.Errorf("extract tyre impression features: %w", err)
	}
	encoded, err := fe.Encode(fingerprint)
	if err != nil {
		return fmt.Errorf("encode tyre impression features: %w", err)
	}
	impression.FeatureJSON = encoded

	modelFeatures := s.FeatureRepo.List()
	models := s.TyreModelRepo.ListAll()
	modelByID := make(map[uint]*m.TyreModel, len(models))
	for _, model := range models {
		modelByID[model.ID] = model
	}
	matches := make([]TyreMatch, 0, len(modelFeatures))
	for _, stored := range modelFeatures {
		model := modelByID[stored.TyreModelID]
		if model == nil {
			continue
		}
		modelFingerprint, err := fe.Decode(stored.FeatureJSON)
		if err != nil {
			return fmt.Errorf("decode stored features for tyre model %d: %w", stored.TyreModelID, err)
		}
		score := s.FeatureExtractor.Compare(fingerprint, modelFingerprint)
		matches = append(matches, TyreMatch{
			TyreModelID:              stored.TyreModelID,
			Manufacturer:             model.Manufacturer,
			ModelName:                model.ModelName,
			Similarity:               score.Final,
			Distance:                 1 - score.Final,
			PatternSimilarity:        score.Pattern,
			LongitudinalSimilarity:   score.Longitudinal,
			StructureCountSimilarity: score.StructureCount,
			ScaleCompatibility:       score.Scale,
		})
	}

	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].Similarity == matches[j].Similarity {
			return matches[i].TyreModelID < matches[j].TyreModelID
		}
		return matches[i].Similarity > matches[j].Similarity
	})
	for i := range matches {
		matches[i].Rank = i + 1
	}

	matchesJSON, err := json.Marshal(matches)
	if err != nil {
		return fmt.Errorf("encode tyre impression matches: %w", err)
	}
	impression.MatchesJSON = string(matchesJSON)
	if len(matches) > 0 {
		impression.Status = m.ProcessingStatusMatched
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
	if original == nil {
		return nil, errors.New("original image file is nil")
	}
	imagePath := filepath.Join(s.FileStore.GetStorageLocation(), original.Location, original.Name)
	grayscaleImage := cv.IMRead(imagePath, cv.IMReadGrayScale)
	if grayscaleImage.Empty() {
		return nil, fmt.Errorf("image is empty")
	}

	return &grayscaleImage, nil
}

func (s *ImageProcessingService) runBinaryStages(source *cv.Mat, id uint, model string, pixelsPerInch float32) (*cv.Mat, error) {
	if source == nil || source.Empty() {
		return nil, fmt.Errorf("binary stage source image is empty")
	}

	processor := processors.NewBinaryProcessor(pixelsPerInch)
	binary := cv.NewMat()
	if err := processor.Segment(source, &binary); err != nil {
		binary.Close()
		return nil, fmt.Errorf("%s stage: %v", processor.GetName(), err)
	}

	if err := s.SaveStage(id, model, processor.GetFileType(), &binary); err != nil {
		binary.Close()
		return nil, err
	}
	binary.Close()

	treadMask := cv.NewMat()
	if err := processor.SegmentTreadMask(source, &treadMask); err != nil {
		treadMask.Close()
		return nil, fmt.Errorf("%s tread mask stage: %v", processor.GetName(), err)
	}

	if err := s.SaveStage(id, model, m.FileTypeTreadMask, &treadMask); err != nil {
		treadMask.Close()
		return nil, err
	}

	return &treadMask, nil
}

func (s *ImageProcessingService) runPipeline(pipeline []processors.Processor, id uint, model string, image *cv.Mat) (*cv.Mat, error) {
	currentImage := image
	for _, processor := range pipeline {
		result, err := processor.Process(currentImage)
		if err != nil {
			return nil, fmt.Errorf("%s stage: %v", processor.GetName(), err)
		}

		if err := s.SaveStage(id, model, processor.GetFileType(), result); err != nil {
			result.Close()
			return nil, err
		}

		if currentImage != image {
			currentImage.Close()
		}
		currentImage = result
	}

	return currentImage, nil
}

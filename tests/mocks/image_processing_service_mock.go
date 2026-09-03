package mocks

import (
	"fmt"
	cv "gocv.io/x/gocv"
	m "tyre-match-backend/db/models"
)

type ImageProcessingServiceMock struct {
	ProcessTyreImpressionCalls []*m.TyreImpression
	ProcessTyreImpressionError error

	ProcessTyreModelCalls []*m.TyreModel
	ProcessTyreModelError error

	SaveStageCalls []struct {
		ID       uint
		Model    string
		FileType string
		Image    *cv.Mat
	}
	SaveStageError error
}

func NewImageProcessingServiceMock() *ImageProcessingServiceMock {
	return &ImageProcessingServiceMock{}
}

func (m *ImageProcessingServiceMock) ProcessTyreImpression(tyreImpression *m.TyreImpression) error {
	m.ProcessTyreImpressionCalls = append(m.ProcessTyreImpressionCalls, tyreImpression)
	return m.ProcessTyreImpressionError
}

func (m *ImageProcessingServiceMock) ProcessTyreModel(tyreModel *m.TyreModel) error {

	fmt.Println(">>> TEST")
	m.ProcessTyreModelCalls = append(m.ProcessTyreModelCalls, tyreModel)
	return m.ProcessTyreModelError
}

func (m *ImageProcessingServiceMock) SaveStage(id uint, model, fileType string, image *cv.Mat) error {
	m.SaveStageCalls = append(m.SaveStageCalls, struct {
		ID       uint
		Model    string
		FileType string
		Image    *cv.Mat
	}{ID: id, Model: model, FileType: fileType, Image: image})
	return m.SaveStageError
}

func (m *ImageProcessingServiceMock) Reset() {
	m.ProcessTyreImpressionCalls = nil
	m.ProcessTyreImpressionError = nil

	m.ProcessTyreModelCalls = nil
	m.ProcessTyreModelError = nil

	m.SaveStageCalls = nil
	m.SaveStageError = nil
}

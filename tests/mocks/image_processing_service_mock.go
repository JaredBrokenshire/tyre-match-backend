package mocks

import cv "gocv.io/x/gocv"

type ImageProcessingServiceMock struct {
	ProcessCalls []struct {
		ID    uint
		Model string
	}
	ProcessError error

	SaveStageCalls []struct {
		ID          uint
		Model       string
		FileType    string
		ResultImage *cv.Mat
	}
	SaveStageError error
}

func NewImageProcessingServiceMock() *ImageProcessingServiceMock {
	return &ImageProcessingServiceMock{}
}

func (m *ImageProcessingServiceMock) Process(id uint, model string) error {
	m.ProcessCalls = append(m.ProcessCalls, struct {
		ID    uint
		Model string
	}{ID: id, Model: model})
	return m.ProcessError
}

func (m *ImageProcessingServiceMock) SaveStage(id uint, model, fileType string, resultImage *cv.Mat) error {
	m.SaveStageCalls = append(m.SaveStageCalls, struct {
		ID          uint
		Model       string
		FileType    string
		ResultImage *cv.Mat
	}{ID: id, Model: model, FileType: fileType, ResultImage: resultImage})
	return m.SaveStageError
}

func (m *ImageProcessingServiceMock) Reset() {
	m.ProcessCalls = nil
	m.ProcessError = nil
}

package processors

import (
	"errors"
	"fmt"
	cv "gocv.io/x/gocv"
)

type Processor interface {
	GetName() string
	GetFileType() string
	Process(*cv.Mat) (*cv.Mat, error)
}

type BaseProcessor struct {
	Name            string
	FileType        string
	ProcessingSteps []ProcessingStep
}

func NewBaseProcessor(name string, fileType string) *BaseProcessor {
	return &BaseProcessor{
		Name:     name,
		FileType: fileType,
	}
}

type ProcessingStep func(*cv.Mat, *cv.Mat) error

func (p *BaseProcessor) Process(source *cv.Mat) (*cv.Mat, error) {
	if source == nil {
		return nil, errors.New("process received a nil image")
	}

	if len(p.ProcessingSteps) == 0 {
		return source, nil
	}

	currentImage := source

	for i, step := range p.ProcessingSteps {
		result := cv.NewMat()

		if err := step(currentImage, &result); err != nil {
			result.Close()

			if currentImage != source {
				currentImage.Close()
			}

			return nil, fmt.Errorf(`processor "%s" step %d: %w`, p.Name, i+1, err)
		}

		if result.Empty() {
			if currentImage != source {
				currentImage.Close()
			}

			return nil, fmt.Errorf(`processor "%s" step %d returned a nil image`, p.Name, i+1)
		}

		if currentImage != source {
			currentImage.Close()
		}

		currentImage = &result
	}

	return currentImage, nil
}

func (p *BaseProcessor) GetName() string {
	return p.Name
}

func (p *BaseProcessor) GetFileType() string {
	return p.FileType
}

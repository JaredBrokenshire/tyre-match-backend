package processors

import (
	"fmt"
	cv "gocv.io/x/gocv"
	"image"
	m "tyre-match-backend/db/models"
)

type EnhancementProcessor struct {
	BaseProcessor

	DenoiseH                  float32
	DenoiseTemplateWindowSize int
	DenoiseSearchWindowSize   int

	BlurSigma      float64
	BlurKSize      image.Point
	BlurBorderType cv.BorderType

	SharpenStrength float64

	CLAHEClipLimit    float64
	CLAHETileGridSize image.Point
}

func NewEnhancementProcessor() *EnhancementProcessor {
	enhancementProcessor := &EnhancementProcessor{
		BaseProcessor: BaseProcessor{
			Name:     "enhancement",
			FileType: m.FileTypeEnhanced,
		},

		DenoiseH:                  10,
		DenoiseTemplateWindowSize: 7,
		DenoiseSearchWindowSize:   21,

		BlurSigma:      1.0,
		BlurKSize:      image.Pt(0, 0),
		BlurBorderType: cv.BorderDefault,

		SharpenStrength: 1.2,

		CLAHEClipLimit:    3.0,
		CLAHETileGridSize: image.Pt(4, 4),
	}

	enhancementProcessor.ProcessingSteps = []ProcessingStep{
		enhancementProcessor.Denoise,
		enhancementProcessor.ApplyCLAHE,
		//enhancementProcessor.Sharpen,
	}

	return enhancementProcessor
}

func (p *EnhancementProcessor) ApplyCLAHE(source, destination *cv.Mat) error {
	err := p.ValidateSourceImage(source)
	if err != nil {
		return fmt.Errorf("apply clahe %v", err)
	}

	clahe := cv.NewCLAHEWithParams(p.CLAHEClipLimit, p.CLAHETileGridSize)
	err = clahe.Apply(*source, destination)
	if err != nil {
		return fmt.Errorf("apply clahe - apply failed: %v", err)
	}

	return nil
}

func (p *EnhancementProcessor) Denoise(source, destination *cv.Mat) error {
	err := p.ValidateSourceImage(source)
	if err != nil {
		return fmt.Errorf("denoise %v", err)
	}

	err = cv.FastNlMeansDenoisingWithParams(*source, destination, p.DenoiseH, p.DenoiseTemplateWindowSize, p.DenoiseSearchWindowSize)
	if err != nil {
		return fmt.Errorf("denoise - fast nl means denoising failed: %v", err)
	}

	return nil
}

func (p *EnhancementProcessor) Sharpen(source, destination *cv.Mat) error {
	err := p.ValidateSourceImage(source)
	if err != nil {
		return fmt.Errorf("sharpen %v", err)
	}

	blurred := cv.NewMat()
	defer blurred.Close()

	err = cv.GaussianBlur(*source, &blurred, p.BlurKSize, p.BlurSigma, p.BlurSigma, p.BlurBorderType)
	if err != nil {
		return fmt.Errorf("sharpen - gaussian blur failed: %v", err)
	}
	if blurred.Empty() {
		return fmt.Errorf("sharpen - gaussian blur produced an empty image")
	}

	err = cv.AddWeighted(*source, 1.0+p.SharpenStrength, blurred, -p.SharpenStrength, 0, destination)
	if err != nil {
		return fmt.Errorf("sharpen - add weighted failed: %v", err)
	}
	if destination.Empty() {
		return fmt.Errorf("sharpen - add weighted produced an empty image")
	}

	return nil
}

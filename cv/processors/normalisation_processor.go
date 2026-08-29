package processors

import (
	"errors"
	"fmt"
	"image"
	"math"

	m "tyre-match-backend/db/models"

	cv "gocv.io/x/gocv"
)

// NormalisationProcessor removes low-frequency illumination variation from a
// grayscale image without changing its dimensions or storage depth.
//
// The illumination field is estimated with a large Gaussian blur and the
// source is corrected multiplicatively:
//
//	corrected = source / illumination * reference
//
// Processing is performed in CV32F so that division does not introduce the
// integer rounding that would occur if the source were divided at its native
// integer depth. The result is converted back to the original image type.
//
// This processor deliberately does not resize, threshold, denoise, sharpen or
// otherwise alter high-frequency tread detail. Those operations belong to
// later stages of the pipeline.
type NormalisationProcessor struct {
	BaseProcessor

	// Fraction of the smaller image dimension used as the illumination
	// estimation scale.
	IlluminationSigmaFraction float64

	// Minimum Gaussian sigma, in pixels.
	MinimumIlluminationSigma float64

	// Prevents division by zero in very dark regions.
	IlluminationEpsilon float32
}

func NewNormalisationProcessor() *NormalisationProcessor {
	processor := &NormalisationProcessor{
		BaseProcessor: BaseProcessor{
			Name:     "normalisation",
			FileType: m.FileTypeNormalised,
		},

		IlluminationSigmaFraction: 0.02,
		MinimumIlluminationSigma:  25.0,
		IlluminationEpsilon:       1e-6,
	}

	processor.ProcessingSteps = []ProcessingStep{
		processor.CorrectIllumination,
	}

	return processor
}

func (p *NormalisationProcessor) CorrectIllumination(source, destination *cv.Mat) error {
	err := p.ValidateSourceImage(source)
	if err != nil {
		return fmt.Errorf("correct illumination %v", err)
	}

	sigma := p.illuminationSigma(source)
	if sigma <= 0 || math.IsNaN(sigma) || math.IsInf(sigma, 0) {
		return errors.New("correct illumination calculated an invalid illumination sigma")
	}

	sourceFloat := cv.NewMat()
	defer sourceFloat.Close()

	// Increase type precision for division
	err = source.ConvertTo(&sourceFloat, cv.MatTypeCV32FC1)
	if err != nil {
		return fmt.Errorf("convert source to float failed: %w", err)
	}

	illumination := cv.NewMat()
	defer illumination.Close()

	// Estimate illumination
	err = cv.GaussianBlur(
		sourceFloat,
		&illumination,
		image.Pt(0, 0),
		sigma,
		sigma,
		cv.BorderReflect101,
	)
	if err != nil {
		return fmt.Errorf("estimate illumination failed: %w", err)
	}
	if illumination.Empty() {
		return errors.New("estimate illumination produced an empty image")
	}

	corrected := cv.NewMat()
	defer corrected.Close()

	err = cv.Divide(sourceFloat, illumination, &corrected)
	if err != nil {
		return fmt.Errorf("divide by illumination failed: %w", err)
	}

	referenceLevel, err := p.meanFloat32(&illumination)
	if err != nil {
		return fmt.Errorf("mean illumination failed: %w", err)
	}

	corrected.MultiplyFloat(float32(referenceLevel))
	if corrected.Empty() {
		return errors.New("correct illumination produced an empty image")
	}

	err = corrected.ConvertTo(
		destination,
		source.Type(),
	)
	if err != nil {
		return fmt.Errorf("convert corrected image to source depth failed: %w", err)
	}
	if destination.Empty() {
		return errors.New("correct illumination produced an empty destination image")
	}

	return nil
}

func (p *NormalisationProcessor) illuminationSigma(source *cv.Mat) float64 {
	minDimension := source.Rows()

	if source.Cols() < minDimension {
		minDimension = source.Cols()
	}

	sigma := float64(minDimension) * p.IlluminationSigmaFraction

	if sigma < p.MinimumIlluminationSigma {
		sigma = p.MinimumIlluminationSigma
	}

	return sigma
}

func (p *NormalisationProcessor) meanFloat32(mat *cv.Mat) (float64, error) {
	values, err := mat.DataPtrFloat32()
	if err != nil {
		return 0, err
	}

	if len(values) == 0 {
		return 0, errors.New("cannot calculate mean of empty matrix")
	}

	var sum float64
	for _, value := range values {
		sum += float64(value)
	}

	return sum / float64(len(values)), nil
}

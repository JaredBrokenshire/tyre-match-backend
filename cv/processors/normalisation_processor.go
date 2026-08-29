package processors

import (
	"errors"
	"fmt"
	"image"
	"image/color"
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

	// ROI coordinates are fractions of image dimensions. They are intentionally
	// configurable because the photographic rig/dataset framing may change.
	//
	// These defaults are based on the supplied standardised dirt-impression
	// photograph: the impression occupies the central/lower part of the frame,
	// with the ruler below it and asphalt on either side.
	ROILeftFraction   float64
	ROIRightFraction  float64
	ROITopFraction    float64
	ROIBottomFraction float64

	// Everything outside the ROI is blurred with this sigma. A relatively large
	// sigma suppresses the high-frequency asphalt texture without introducing
	// high-frequency structure into the downstream enhancement/segmentation
	// stages.
	OutsideBlurSigma float64

	// Fraction of the smaller image dimension used as the illumination
	// estimation scale.
	IlluminationSigmaFraction float64

	// Minimum Gaussian sigma, in pixels.
	MinimumIlluminationSigma float64
}

func NewNormalisationProcessor() *NormalisationProcessor {
	processor := &NormalisationProcessor{
		BaseProcessor: BaseProcessor{
			Name:     "normalisation",
			FileType: m.FileTypeNormalised,
		},

		// TODO: Have these values stored for each impression at upload
		ROILeftFraction:   0.075,
		ROIRightFraction:  0.92,
		ROITopFraction:    0.27,
		ROIBottomFraction: 0.86,

		OutsideBlurSigma: 101.0,

		IlluminationSigmaFraction: 0.02,
		MinimumIlluminationSigma:  25.0,
	}

	processor.ProcessingSteps = []ProcessingStep{
		processor.IsolateRegionOfInterest,
		processor.CorrectIllumination,
	}

	return processor
}

func (p *NormalisationProcessor) IsolateRegionOfInterest(source, destination *cv.Mat) error {
	if err := p.ValidateSourceImage(source); err != nil {
		return fmt.Errorf("isolate roi %v", err)
	}

	left := int(math.Round(float64(source.Cols()) * p.ROILeftFraction))
	right := int(math.Round(float64(source.Cols()) * p.ROIRightFraction))
	top := int(math.Round(float64(source.Rows()) * p.ROITopFraction))
	bottom := int(math.Round(float64(source.Rows()) * p.ROIBottomFraction))

	roi := image.Rect(left, top, right, bottom)

	// Start with a blurred copy. The original source is then copied over the
	// ROI, leaving the dimensions and pixels inside the ROI unchanged.
	background := cv.NewMat()
	defer background.Close()

	if err := cv.GaussianBlur(
		*source,
		&background,
		image.Pt(0, 0),
		p.OutsideBlurSigma,
		p.OutsideBlurSigma,
		cv.BorderReflect101,
	); err != nil {
		return fmt.Errorf("blur background failed: %w", err)
	}

	if background.Empty() {
		return errors.New("blur background produced an empty image")
	}

	mask := cv.NewMatWithSize(source.Rows(), source.Cols(), cv.MatTypeCV8UC1)
	defer mask.Close()

	// Fill the ROI with 255. The remainder stays zero.
	cv.Rectangle(&mask, roi, color.RGBA{R: 255, G: 255, B: 255, A: 255}, -1)

	background.CopyTo(destination)
	if err := source.CopyToWithMask(destination, mask); err != nil {
		return fmt.Errorf("copy roi over blurred background failed: %w", err)
	}

	if destination.Empty() {
		return errors.New("isolate roi produced an empty image")
	}

	return nil
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

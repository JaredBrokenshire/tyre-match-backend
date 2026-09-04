package processors

import (
	"errors"
	"fmt"
	"image"
	"math"

	m "tyre-match-backend/db/models"

	cv "gocv.io/x/gocv"
)

// NormalisationProcessor crops an input image to the supplied region of
// interest and removes low-frequency illumination variation from that region.
//
// The ROI crop is deliberately performed before illumination correction. The
// downstream processors therefore receive only the evidence selected by the
// user, rather than an image containing the original photograph's surrounding
// background, rulers or other objects.
//
// Processing is performed in CV32F so that illumination correction does not
// introduce the integer rounding that would occur if the source were divided
// at its native integer depth. The result is converted back to the original
// image type.
type NormalisationProcessor struct {
	BaseProcessor

	// ROI coordinates are absolute pixel values in the source image. The right
	// and bottom coordinates are exclusive, matching image.Rectangle semantics.
	ROITop    int
	ROILeft   int
	ROIRight  int
	ROIBottom int

	IlluminationSigmaFraction float64
	MinimumIlluminationSigma  float64
}

func NewNormalisationProcessor(roiTop, roiLeft, roiRight, roiBottom int) *NormalisationProcessor {
	processor := &NormalisationProcessor{
		BaseProcessor: BaseProcessor{
			Name:     "normalisation",
			FileType: m.FileTypeNormalised,
		},

		ROITop:    roiTop,
		ROILeft:   roiLeft,
		ROIRight:  roiRight,
		ROIBottom: roiBottom,

		IlluminationSigmaFraction: 0.02,
		MinimumIlluminationSigma:  25.0,
	}

	processor.ProcessingSteps = []ProcessingStep{
		processor.CropToRegionOfInterest,
		processor.CorrectIllumination,
	}

	return processor
}

// CropToRegionOfInterest validates the configured ROI and returns a copy of
// exactly that region. The source image is never mutated.
func (p *NormalisationProcessor) CropToRegionOfInterest(source, destination *cv.Mat) error {
	if err := p.ValidateSourceImage(source); err != nil {
		return fmt.Errorf("crop roi %v", err)
	}

	roi, err := p.regionOfInterest(source)
	if err != nil {
		return err
	}

	cropped := source.Region(roi)
	defer cropped.Close()

	if cropped.Empty() {
		return errors.New("crop roi produced an empty region")
	}

	cropped.CopyTo(destination)
	if destination.Empty() {
		return errors.New("crop roi produced an empty destination image")
	}

	return nil
}

// IsolateRegionOfInterest is retained as the named processing-step entry point
// used by earlier callers. It now performs a true crop rather than blurring
// and retaining the original image dimensions.
func (p *NormalisationProcessor) IsolateRegionOfInterest(source, destination *cv.Mat) error {
	return p.CropToRegionOfInterest(source, destination)
}

func (p *NormalisationProcessor) regionOfInterest(source *cv.Mat) (image.Rectangle, error) {
	if p.ROILeft < 0 || p.ROITop < 0 {
		return image.Rectangle{}, errors.New("crop roi coordinates cannot be negative")
	}
	if p.ROIRight <= p.ROILeft || p.ROIBottom <= p.ROITop {
		return image.Rectangle{}, errors.New("crop roi must have positive width and height")
	}
	if p.ROIRight > source.Cols() || p.ROIBottom > source.Rows() {
		return image.Rectangle{}, fmt.Errorf(
			"crop roi %d,%d,%d,%d exceeds source dimensions %dx%d",
			p.ROILeft,
			p.ROITop,
			p.ROIRight,
			p.ROIBottom,
			source.Cols(),
			source.Rows(),
		)
	}

	return image.Rect(p.ROILeft, p.ROITop, p.ROIRight, p.ROIBottom), nil
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

	if err = source.ConvertTo(&sourceFloat, cv.MatTypeCV32FC1); err != nil {
		return fmt.Errorf("convert source to float failed: %w", err)
	}

	illumination := cv.NewMat()
	defer illumination.Close()

	if err = cv.GaussianBlur(
		sourceFloat,
		&illumination,
		image.Pt(0, 0),
		sigma,
		sigma,
		cv.BorderReflect101,
	); err != nil {
		return fmt.Errorf("estimate illumination failed: %w", err)
	}
	if illumination.Empty() {
		return errors.New("estimate illumination produced an empty image")
	}

	corrected := cv.NewMat()
	defer corrected.Close()

	if err = cv.Divide(sourceFloat, illumination, &corrected); err != nil {
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

	if err = corrected.ConvertTo(destination, source.Type()); err != nil {
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

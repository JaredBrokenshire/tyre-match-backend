package processors

import (
	"errors"
	"fmt"
	"image"
	"math"
	m "tyre-match-backend/db/models"

	cv "gocv.io/x/gocv"
)

// BinaryProcessor creates two complementary binary representations from the
// enhanced ROI image:
//
//   - Segment produces the existing boundary/groove-oriented binary. It uses
//     inverted adaptive thresholding and is retained for backwards-compatible
//     diagnostics and fine boundary detail.
//   - SegmentTreadMask produces a tread/contact mask. It uses Otsu global
//     thresholding so the brighter tread/contact regions become white
//     foreground and the darker grooves become black background. A small,
//     PPI-aware morphological closing bridges segmentation breaks without
//     deliberately reconstructing large missing tread regions.
//
// Both representations are 8-bit, single-channel images with 0/255 values.
type BinaryProcessor struct {
	BaseProcessor

	AdaptiveBlockSize int
	AdaptiveC         float32

	PixelsPerInch           float32
	TreadClosingMillimetres float64
}

// NewBinaryProcessor accepts an optional PPI for scale-aware tread-mask
// cleanup. The optional form keeps existing processor tests and callers valid;
// production image processing supplies the user-calibrated PPI.
func NewBinaryProcessor(pixelsPerInch ...float32) *BinaryProcessor {
	ppi := float32(0)
	if len(pixelsPerInch) > 0 {
		ppi = pixelsPerInch[0]
	}

	processor := &BinaryProcessor{
		BaseProcessor: BaseProcessor{
			Name:     "binary",
			FileType: m.FileTypeBinary,
		},

		AdaptiveBlockSize:       101,
		AdaptiveC:               8.0,
		PixelsPerInch:           ppi,
		TreadClosingMillimetres: 0.30,
	}

	processor.ProcessingSteps = []ProcessingStep{
		processor.Segment,
	}

	return processor
}

// Segment creates the existing inverted binary representation. This output
// remains the processor's normal Process result and is saved as the existing
// "binary" processing stage.
func (p *BinaryProcessor) Segment(source, destination *cv.Mat) error {
	if err := p.ValidateSourceImage(source); err != nil {
		return fmt.Errorf("segment %v", err)
	}

	if err := p.validateAdaptiveParameters(); err != nil {
		return fmt.Errorf("segment %v", err)
	}

	if err := cv.AdaptiveThreshold(
		*source,
		destination,
		255,
		cv.AdaptiveThresholdGaussian,
		cv.ThresholdBinaryInv,
		p.AdaptiveBlockSize,
		p.AdaptiveC,
	); err != nil {
		return fmt.Errorf("segment adaptive threshold failed: %w", err)
	}

	if destination.Empty() {
		return errors.New("segment produced an empty image")
	}

	return ensureBinary(destination, "segment")
}

// SegmentTreadMask creates the new tread/contact representation. Unlike the
// existing adaptive binary, this deliberately uses Otsu thresholding rather
// than merely reversing the adaptive threshold polarity. Adaptive thresholding
// can classify broad, nearly uniform grooves as foreground because its local
// threshold follows the groove intensity. Otsu instead selects a single
// intensity threshold from the image histogram, producing filled bright
// tread/contact regions while retaining darker grooves as background.
//
// The existing boundary mask remains independent and available for diagnostics.
func (p *BinaryProcessor) SegmentTreadMask(source, destination *cv.Mat) error {
	if err := p.ValidateSourceImage(source); err != nil {
		return fmt.Errorf("segment tread mask %v", err)
	}

	thresholded := cv.NewMat()
	defer thresholded.Close()

	// Otsu chooses the threshold from the grayscale histogram. The explicit
	// zero threshold is required by OpenCV when the Otsu flag is supplied.
	// A returned threshold of zero is valid, so it is not treated as an error.
	_ = cv.Threshold(
		*source,
		&thresholded,
		0,
		255,
		cv.ThresholdBinary|cv.ThresholdOtsu,
	)

	if thresholded.Empty() {
		return errors.New("segment tread mask Otsu thresholding produced an empty image")
	}

	kernelSize := p.treadClosingKernelSize()
	if kernelSize > 1 {
		kernel := cv.GetStructuringElement(cv.MorphRect, image.Pt(kernelSize, kernelSize))
		defer kernel.Close()

		closed := cv.NewMat()
		defer closed.Close()

		if err := cv.MorphologyEx(thresholded, &closed, cv.MorphClose, kernel); err != nil {
			return fmt.Errorf("segment tread mask morphological closing failed: %w", err)
		}
		if closed.Empty() {
			return errors.New("segment tread mask morphological closing produced an empty image")
		}

		closed.CopyTo(destination)
	} else {
		thresholded.CopyTo(destination)
	}

	if destination.Empty() {
		return errors.New("segment tread mask produced an empty image")
	}

	if err := ensureBinary(destination, "segment tread mask"); err != nil {
		return err
	}

	return nil
}

func (p *BinaryProcessor) validateAdaptiveParameters() error {
	if p.AdaptiveBlockSize < 3 || p.AdaptiveBlockSize%2 == 0 {
		return errors.New("adaptive block size must be an odd integer greater than or equal to 3")
	}
	if math.IsNaN(float64(p.AdaptiveC)) || math.IsInf(float64(p.AdaptiveC), 0) {
		return errors.New("adaptive C must be finite")
	}
	return nil
}

func (p *BinaryProcessor) treadClosingKernelSize() int {
	if p.PixelsPerInch <= 0 || p.TreadClosingMillimetres <= 0 {
		return 1
	}

	pixels := p.TreadClosingMillimetres * float64(p.PixelsPerInch) / 25.4
	size := int(math.Round(pixels))
	if size < 1 {
		size = 1
	}
	if size%2 == 0 {
		size++
	}

	return size
}

func ensureBinary(mat *cv.Mat, operation string) error {
	values := mat.ToBytes()
	for _, value := range values {
		if value != 0 && value != 255 {
			return fmt.Errorf("%s produced a non-binary pixel value %d", operation, value)
		}
	}
	return nil
}

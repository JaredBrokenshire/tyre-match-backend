package processors

import (
	"errors"
	"fmt"
	m "tyre-match-backend/db/models"

	cv "gocv.io/x/gocv"
)

// SegmentationProcessor creates a binary representation of the tread
// pattern using adaptive thresholding. The input is expected to be a grayscale
// image that has already passed through the project's enhancement stage.
//
// Adaptive thresholding is used instead of a single global threshold because
// the dataset contains different illumination/contrast conditions and because
// the ground impressions contain local substrate variation.
//
// ThresholdBinaryInv is used because in the ground-impression photographs the
// tread contact areas are darker than the surrounding substrate — the tyre
// compresses the dirt, displacing it into the groove channels which appear as
// lighter raised ridges. Inverted thresholding therefore classifies the darker
// tread contact areas as foreground (white) and the lighter grooves/substrate
// as background (black), which is the correct convention for downstream
// morphological cleaning and feature extraction.
type BinaryProcessor struct {
	BaseProcessor

	AdaptiveBlockSize int
	AdaptiveC         float32
}

func NewBinaryProcessor() *BinaryProcessor {
	processor := &BinaryProcessor{
		BaseProcessor: BaseProcessor{
			Name:     "binary",
			FileType: m.FileTypeBinary,
		},

		// AdaptiveBlockSize defines the size of the local neighbourhood used when
		// calculating the threshold for each pixel. A larger value considers a wider
		// area of the image when determining whether a pixel belongs to the foreground
		// (tread) or background (groove/substrate), while a smaller value makes the
		// segmentation more responsive to local changes in illumination and texture.
		//
		// The value must be an odd integer greater than or equal to 3, as required by
		// OpenCV's adaptive thresholding algorithm.
		//
		// At the dataset's calibrated resolution of ~294 PPI, 101 pixels represents
		// approximately 8.7mm — wide enough to span a full tread block and capture
		// genuine rib-vs-groove contrast rather than sub-millimetre substrate texture.
		AdaptiveBlockSize: 101,
		// AdaptiveC is a constant subtracted from the locally calculated threshold
		// for each pixel. Increasing this value makes it harder for pixels to be
		// classified as foreground (tread contact), which suppresses mid-tone
		// substrate texture that would otherwise be misclassified.
		//
		// The value may be positive, zero, or negative. For ground-impression
		// photographs of the project dataset, a value of 8 preserves genuine tread
		// block foreground while suppressing the most obvious substrate texture.
		// Values much higher than this cause the tread blocks themselves to be
		// under-segmented. Residual substrate noise is handled by the subsequent
		// morphological cleaning stage rather than by making this threshold more
		// aggressive.
		//
		// This value should be re-evaluated experimentally if the substrate
		// material, illumination conditions, or camera setup changes.
		AdaptiveC: 8.0,
	}

	processor.ProcessingSteps = []ProcessingStep{
		processor.Segment,
	}

	return processor
}

func (p *BinaryProcessor) Segment(source, destination *cv.Mat) error {
	err := p.ValidateSourceImage(source)
	if err != nil {
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

	return nil
}

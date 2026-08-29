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
		AdaptiveBlockSize: 101,
		// AdaptiveC is a constant subtracted from the locally calculated threshold
		// for each pixel. Increasing this value makes it easier for pixels to be
		// classified as foreground, while decreasing it makes the foreground
		// classification more restrictive.
		//
		// The value may be positive, zero, or negative. Its useful range is dependent
		// on the characteristics of the image and should therefore be determined
		// experimentally using the project dataset.
		AdaptiveC: 5.0,
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
		cv.ThresholdBinary,
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

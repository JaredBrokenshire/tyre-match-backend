package cv_tests_test

import (
	"github.com/stretchr/testify/assert"
	cv "gocv.io/x/gocv"
	"image"
	"math/rand"
	"testing"
	"tyre-match-backend/cv/processors"
	"tyre-match-backend/tests/helpers"
)

func TestBinaryProcessorSegment(t *testing.T) {
	processor := processors.NewBinaryProcessor()

	realisticImage := cv.IMRead(
		"../../assets/example.jpg",
		cv.IMReadGrayScale,
	)
	assert.False(t, realisticImage.Empty(), "unable to load realistic image")
	defer realisticImage.Close()

	cases := append(
		helpers.ProcessingStepTests("segment"),
		[]helpers.ProcessingStepTest{
			{
				Name:   "Segments noisy image",
				Source: noisySegmentationImage(),
			},
			{
				Name:   "Segments gradient image",
				Source: gradientSegmentationImage(),
			},
			{
				Name:   "Segments isolated bright pixel",
				Source: isolatedBrightPixelImage(),
			},
			{
				Name:   "Segments realistic tyre impression casting",
				Source: &realisticImage,
			},
		}...,
	)

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			if test.Source != nil && test.Source != &realisticImage {
				defer test.Source.Close()
			}

			if test.ExpectedResult != nil {
				defer test.ExpectedResult.Close()
			}

			result := cv.NewMat()
			defer result.Close()

			err := processor.Segment(test.Source, &result)
			if test.ExpectedError != "" {
				assert.EqualError(t, err, test.ExpectedError, "error does not match")
				assert.True(t, result.Empty(), "result should be empty")
				return
			}

			assertSegmentationResult(t, test.Source, &result)

			if test.AssertResult != nil {
				test.AssertResult(t, &result)
			}
		})
	}
}

func assertSegmentationResult(
	t *testing.T,
	source, result *cv.Mat,
) {
	t.Helper()

	assert.False(
		t,
		result.Empty(),
		"segmentation result should not be empty",
	)

	assert.Equal(
		t,
		source.Rows(),
		result.Rows(),
		"segmentation result should preserve input height",
	)

	assert.Equal(
		t,
		source.Cols(),
		result.Cols(),
		"segmentation result should preserve input width",
	)

	assert.Equal(
		t,
		cv.MatTypeCV8UC1,
		result.Type(),
		"segmentation result should be an 8-bit single-channel image",
	)

	values := result.ToBytes()

	assert.NotEmpty(
		t,
		values,
		"segmentation result should contain pixel data",
	)

	hasForeground := false
	hasBackground := false

	for _, value := range values {
		switch value {
		case 0:
			hasBackground = true
		case 255:
			hasForeground = true
		default:
			assert.Failf(
				t,
				"segmentation result is not binary",
				"unexpected pixel value %d",
				value,
			)
		}
	}

	assert.True(
		t,
		hasBackground,
		"segmentation should contain background pixels",
	)

	assert.True(
		t,
		hasForeground,
		"segmentation should contain foreground pixels",
	)
}

func noisySegmentationImage() *cv.Mat {
	const (
		width  = 301
		height = 301
	)

	mat := cv.NewMatWithSize(
		height,
		width,
		cv.MatTypeCV8UC1,
	)

	// Bright background (substrate).
	mat.SetTo(cv.NewScalar(200, 200, 200, 0))

	// Dark central region (tread contact) — value must be far enough below the
	// local neighbourhood mean to survive AdaptiveC=20 subtraction.
	// Background=200, tread=40: difference=160, well above AdaptiveC=20.
	region := mat.Region(
		image.Rect(75, 75, 226, 226),
	)
	region.SetTo(cv.NewScalar(40, 40, 40, 0))
	region.Close()

	// Add deterministic noise within the tread region to provide local contrast
	// that the adaptive threshold can respond to.
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 3000; i++ {
		x := 75 + rng.Intn(151)
		y := 75 + rng.Intn(151)
		value := uint8(20 + rng.Intn(41)) // noise in range 20–60
		mat.SetUCharAt(y, x, value)
	}

	return &mat
}

func gradientSegmentationImage() *cv.Mat {
	const (
		width  = 301
		height = 301
	)

	mat := cv.NewMatWithSize(
		height,
		width,
		cv.MatTypeCV8UC1,
	)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gradient := uint8(
				20 + (200 * x / (width - 1)),
			)

			mat.SetUCharAt(y, x, gradient)
		}
	}

	// Dark local structure representing a tread contact region — under
	// ThresholdBinaryInv this dark patch will be classified as foreground.
	region := mat.Region(
		image.Rect(100, 75, 201, 226),
	)
	region.SetTo(cv.NewScalar(20, 20, 20, 0))
	region.Close()

	return &mat
}

func isolatedBrightPixelImage() *cv.Mat {
	const (
		width  = 301
		height = 301
	)

	mat := cv.NewMatWithSize(
		height,
		width,
		cv.MatTypeCV8UC1,
	)

	// Bright background (substrate).
	mat.SetTo(cv.NewScalar(210, 210, 210, 0))

	// Dark central region (tread contact).
	region := mat.Region(
		image.Rect(75, 75, 226, 226),
	)
	region.SetTo(cv.NewScalar(40, 40, 40, 0))
	region.Close()

	// Deliberately introduce a single bright pixel outside the main
	// foreground structure. Segmentation is not responsible for removing
	// isolated noise; that is handled by the subsequent morphology stage.
	mat.SetUCharAt(25, 25, 255)

	return &mat
}

func assertRegionForegroundRatio(
	t *testing.T,
	result *cv.Mat,
	region image.Rectangle,
	minRatio, maxRatio float64,
) {
	t.Helper()

	roi := result.Region(region)
	defer roi.Close()

	total := roi.Rows() * roi.Cols()

	assert.Greater(
		t,
		total,
		0,
		"region should contain pixels",
	)

	foreground := cv.CountNonZero(roi)
	ratio := float64(foreground) / float64(total)

	assert.GreaterOrEqual(
		t,
		ratio,
		minRatio,
		"foreground ratio is below expected minimum",
	)

	assert.LessOrEqual(
		t,
		ratio,
		maxRatio,
		"foreground ratio is above expected maximum",
	)
}

package cv_tests_test

import (
	"github.com/stretchr/testify/assert"
	cv "gocv.io/x/gocv"
	"image"
	"math"
	"testing"
	"tyre-match-backend/cv/processors"
	"tyre-match-backend/tests/helpers"
)

func TestEnhancementProcessorApplyCLAHE(t *testing.T) {
	processor := processors.NewEnhancementProcessor()

	realisticImage := cv.IMRead(
		"../../assets/example.jpg",
		cv.IMReadGrayScale,
	)
	assert.False(t, realisticImage.Empty(), "unable to load realistic image")

	cases := append(
		helpers.ProcessingStepTests("apply clahe"),
		[]helpers.ProcessingStepTest{
			{
				Name:   "Preserves uniform image",
				Source: helpers.SolidGray(100, 100),
				// NOTE: CLAHE can change brightness values of a uniform image while maintaining uniformity
				// Therefore, ExpectedEqual can not be set to true for this test.
			},
			{
				Name:   "Enhances low contrast gradient",
				Source: helpers.GradientGray(100, 100),
			},
			{
				Name:   "Enhances low contrast regions",
				Source: helpers.LowContrast(100, 100),
			},
			{
				Name:   "Enhances realistic image",
				Source: &realisticImage,
			},
		}...,
	)

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			if test.Source != nil {
				defer test.Source.Close()
			}

			result := cv.NewMat()
			defer result.Close()

			err := processor.ApplyCLAHE(test.Source, &result)

			if test.ExpectedError != "" {
				assert.EqualError(t, err, test.ExpectedError, "error does not match")
				assert.True(t, result.Empty(), "result should be empty")
				return
			}

			assert.NoError(t, err, "should not have an error")
			assert.False(t, result.Empty(), "result should not be empty")

			assertProcessingStep(
				t,
				test.Source,
				&result,
				processor.ApplyCLAHE,
				test.ExpectedEqual,
			)

			switch test.Name {
			case "Preserves uniform image":
				assertUniformImagePreserved(t, test.Source, &result)
			default:
				assertCLAHEApplied(t, test.Source, &result)
			}
		})
	}
}

func TestEnhancementProcessorDenoise(t *testing.T) {
	processor := processors.NewEnhancementProcessor()

	realisticImage := cv.IMRead(
		"../../assets/example.jpg",
		cv.IMReadGrayScale,
	)
	assert.False(t, realisticImage.Empty(), "unable to load realistic image")

	cases := append(
		helpers.ProcessingStepTests("denoise"),
		[]helpers.ProcessingStepTest{
			{
				Name:          "Preserves uniform image",
				Source:        helpers.SolidGray(100, 100),
				ExpectedEqual: true,
			},
			{
				Name:   "Denoises gradient",
				Source: helpers.GradientGray(100, 100),
			},
			{
				Name:   "Denoises isolated bright pixel",
				Source: helpers.ImpulseGray(100, 100),
			},
			{
				Name:   "Denoises noisy image",
				Source: helpers.Noise(100),
			},
			{
				Name:   "Denoises realistic image",
				Source: &realisticImage,
			},
		}...,
	)

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			if test.Source != nil {
				defer test.Source.Close()
			}

			result := cv.NewMat()
			err := processor.Denoise(test.Source, &result)

			if test.ExpectedError != "" {
				assert.EqualError(t, err, test.ExpectedError, "error does not match")
				assert.True(t, result.Empty(), "result should be empty")
				return
			}

			assert.NoError(t, err, "should not have an error")
			assert.False(t, result.Empty(), "result should not be empty")

			assertProcessingStep(
				t,
				test.Source,
				&result,
				processor.Denoise,
				test.ExpectedEqual,
			)

			switch test.Name {
			case "Preserves uniform image":
				assertUniformImagePreserved(t, test.Source, &result)
			default:
				// Denoise assertions - calculate high-frequency variation (HVF)
				sourceHFV := calculateHighFrequencyVariation(test.Source)
				resultHVF := calculateHighFrequencyVariation(&result)

				assert.Less(t, resultHVF, sourceHFV, "high-frequency variation should be less than source")
			}
		})
	}
}

func TestEnhancementProcessorSharpen(t *testing.T) {
	processor := processors.NewEnhancementProcessor()

	realisticImage := cv.IMRead(
		"../../assets/example.jpg",
		cv.IMReadGrayScale,
	)
	assert.False(t, realisticImage.Empty(), "unable to load realistic image")

	cases := append(
		helpers.ProcessingStepTests("sharpen"),
		[]helpers.ProcessingStepTest{
			{
				Name:          "Preserves uniform image",
				Source:        helpers.SolidGray(100, 100),
				ExpectedEqual: true,
			},
			{
				Name:   "Sharpens gradient",
				Source: helpers.GradientGray(100, 100),
			},
			{
				Name:   "Sharpens isolated bright pixel",
				Source: helpers.ImpulseGray(100, 100),
			},
			{
				Name:   "Processes noisy image",
				Source: helpers.Noise(100),
			},
			{
				Name:   "Sharpens realistic image",
				Source: &realisticImage,
			},
		}...,
	)

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			if test.Source != nil {
				defer test.Source.Close()
			}

			result := cv.NewMat()
			defer result.Close()

			err := processor.Sharpen(test.Source, &result)

			if test.ExpectedError != "" {
				assert.EqualError(t, err, test.ExpectedError, "error does not match")
				assert.True(t, result.Empty(), "result should be empty")
				return
			}

			assert.NoError(t, err, "should not have an error")
			assert.False(t, result.Empty(), "result should not be empty")

			assertProcessingStep(
				t,
				test.Source,
				&result,
				processor.Sharpen,
				test.ExpectedEqual,
			)

			switch test.Name {
			case "Preserves uniform image":
				assertUniformImagePreserved(t, test.Source, &result)

			case "Sharpens gradient",
				"Sharpens isolated bright pixel",
				"Sharpens realistic image":
				assertSharpenApplied(t, test.Source, &result)

			case "Processes noisy image":
				assertSharpenProcessedNoise(t, test.Source, &result)
			}
		})
	}
}

func assertProcessingStep(
	t *testing.T,
	source *cv.Mat,
	result *cv.Mat,
	step processors.ProcessingStep,
	expectedEqual bool,
) {
	t.Helper()

	original := source.Clone()
	defer original.Close()

	repeatResult := cv.NewMat()
	defer repeatResult.Close()

	err := step(source, &repeatResult)
	assert.NoError(t, err, "should not have an error")

	if source.Empty() {
		assert.True(t, result.Empty(), "result should be empty")
		assert.True(t, repeatResult.Empty(), "repeat-result should be empty")
		return
	}

	helpers.AssertMatSameShapeAndType(t, source, result)

	// Operation should be deterministic.
	assert.Equal(t, result.ToBytes(), repeatResult.ToBytes(), "result and repeat-result should match")

	// Processing should not mutate the source.
	assert.Equal(t, original.ToBytes(), source.ToBytes(), "source should not be mutated")

	if expectedEqual {
		assert.Equal(t, source.ToBytes(), result.ToBytes(), "source and result should match")
	} else {
		assert.NotEqual(t, source.ToBytes(), result.ToBytes(), "source and result should not match")
	}
}

func assertUniformImagePreserved(t *testing.T, source, result *cv.Mat) {
	t.Helper()

	assert.False(t, source.Empty(), "source should not be empty")
	assert.False(t, result.Empty(), "result should not be empty")

	assert.True(
		t,
		isUniform(result),
		"processing a uniform image should produce a uniform image",
	)
}

func isUniform(src *cv.Mat) bool {
	if src == nil || src.Empty() {
		return false
	}

	values := src.ToBytes()

	if len(values) == 0 {
		return false
	}

	first := values[0]

	for _, value := range values[1:] {
		if value != first {
			return false
		}
	}

	return true
}

func assertCLAHEApplied(t *testing.T, source, result *cv.Mat) {
	t.Helper()

	assert.False(t, source.Empty(), "source should not be empty")
	assert.False(t, result.Empty(), "result should not be empty")

	sourceContrast := calculateLocalContrast(source, image.Pt(8, 8))
	resultContrast := calculateLocalContrast(result, image.Pt(8, 8))

	assert.Greater(
		t,
		resultContrast,
		sourceContrast,
		"CLAHE should increase local contrast",
	)

	changedPixels := countChangedPixels(source, result)

	assert.Greater(
		t,
		changedPixels,
		0,
		"CLAHE should change at least one pixel",
	)
}

func calculateLocalContrast(src *cv.Mat, grid image.Point) float64 {
	if src == nil || src.Empty() {
		return 0
	}

	if grid.X <= 0 || grid.Y <= 0 {
		return 0
	}

	width := src.Cols()
	height := src.Rows()

	cellWidth := width / grid.X
	cellHeight := height / grid.Y

	if cellWidth == 0 || cellHeight == 0 {
		return 0
	}

	var total float64
	var regions int

	for gridY := 0; gridY < grid.Y; gridY++ {
		for gridX := 0; gridX < grid.X; gridX++ {
			startX := gridX * cellWidth
			startY := gridY * cellHeight

			endX := startX + cellWidth
			endY := startY + cellHeight

			// Include the remaining pixels when dimensions are not evenly
			// divisible by the number of regions.
			if gridX == grid.X-1 {
				endX = width
			}
			if gridY == grid.Y-1 {
				endY = height
			}

			var sum float64
			var count float64

			for y := startY; y < endY; y++ {
				for x := startX; x < endX; x++ {
					value := float64(src.GetUCharAt(y, x))
					sum += value
					count++
				}
			}

			if count == 0 {
				continue
			}

			mean := sum / count

			var variance float64
			for y := startY; y < endY; y++ {
				for x := startX; x < endX; x++ {
					value := float64(src.GetUCharAt(y, x))
					difference := value - mean
					variance += difference * difference
				}
			}

			total += math.Sqrt(variance / count)
			regions++
		}
	}

	if regions == 0 {
		return 0
	}

	return total / float64(regions)
}

func countChangedPixels(source, result *cv.Mat) int {
	if source == nil || result == nil ||
		source.Empty() || result.Empty() ||
		source.Rows() != result.Rows() ||
		source.Cols() != result.Cols() ||
		source.Type() != result.Type() {
		return 0
	}

	sourceBytes := source.ToBytes()
	resultBytes := result.ToBytes()

	if len(sourceBytes) != len(resultBytes) {
		return 0
	}

	changed := 0

	for i := range sourceBytes {
		if sourceBytes[i] != resultBytes[i] {
			changed++
		}
	}

	return changed
}

// calculateHighFrequencyVariation creates a slightly blurred copy of a source image and subtracts if from the source
// This leaves only large changes in nearby pixel brightness
// By comparing this value before and after denoising we should see a reduction in variation.
func calculateHighFrequencyVariation(src *cv.Mat) float64 {
	blurred := cv.NewMat()
	defer blurred.Close()

	if err := cv.GaussianBlur(
		*src,
		&blurred,
		image.Pt(3, 3),
		0,
		0,
		cv.BorderDefault,
	); err != nil {
		return 0
	}

	difference := cv.NewMat()
	defer difference.Close()

	if err := cv.AbsDiff(*src, blurred, &difference); err != nil {
		return 0
	}

	values := difference.ToBytes()
	if len(values) == 0 {
		return 0
	}

	var total float64
	for _, value := range values {
		total += float64(value)
	}

	return total / float64(len(values))
}

func assertSharpenApplied(t *testing.T, source, result *cv.Mat) {
	t.Helper()

	assert.False(t, source.Empty(), "source should not be empty")
	assert.False(t, result.Empty(), "result should not be empty")

	sourceHFV := calculateHighFrequencyVariation(source)
	resultHFV := calculateHighFrequencyVariation(result)

	assert.Greater(
		t,
		resultHFV,
		sourceHFV,
		"sharpening should increase high-frequency variation",
	)

	changedPixels := countChangedPixels(source, result)

	assert.Greater(
		t,
		changedPixels,
		0,
		"sharpening should change at least one pixel",
	)
}

func assertSharpenProcessedNoise(t *testing.T, source, result *cv.Mat) {
	t.Helper()

	assert.False(t, source.Empty(), "source should not be empty")
	assert.False(t, result.Empty(), "result should not be empty")

	changedPixels := countChangedPixels(source, result)

	assert.Greater(
		t,
		changedPixels,
		0,
		"sharpening should change at least one pixel in a noisy image",
	)
}

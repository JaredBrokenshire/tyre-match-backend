package cv_tests_test

import (
	"github.com/stretchr/testify/assert"
	cv "gocv.io/x/gocv"
	"math"
	"testing"
	"tyre-match-backend/cv/processors"
	"tyre-match-backend/tests/helpers"
)

func TestNormalisationProcessorIsolateROI(t *testing.T) {
	processor := processors.NewNormalisationProcessor()

	source := cv.NewMatWithSize(301, 401, cv.MatTypeCV8UC1)
	defer source.Close()

	// High-frequency background noise.
	for y := 0; y < source.Rows(); y++ {
		for x := 0; x < source.Cols(); x++ {
			value := uint8((x*17 + y*31) % 256)
			source.SetUCharAt(y, x, value)
		}
	}

	// A simple tread-like structure inside the ROI.
	for y := 100; y < 250; y++ {
		for x := 80; x < 320; x++ {
			if ((x / 10) % 2) == 0 {
				source.SetUCharAt(y, x, 190)
			} else {
				source.SetUCharAt(y, x, 60)
			}
		}
	}

	result := cv.NewMat()
	defer result.Close()

	err := processor.IsolateRegionOfInterest(&source, &result)
	assert.NoError(t, err)
	assert.False(t, result.Empty())

	// Pixels comfortably inside the configured ROI must be preserved.
	assert.Equal(t, source.GetUCharAt(150, 200), result.GetUCharAt(150, 200))

	// A pixel outside the ROI should no longer retain the original
	// high-frequency value exactly.
	assert.NotEqual(t, source.GetUCharAt(40, 40), result.GetUCharAt(40, 40))
}

func TestNormalisationProcessorCorrectIllumination(t *testing.T) {
	processor := processors.NewNormalisationProcessor()

	realisticImage := cv.IMRead(
		"../../assets/example.jpg",
		cv.IMReadGrayScale,
	)
	assert.False(t, realisticImage.Empty(), "unable to load realistic tyre impression image")
	defer realisticImage.Close()

	cases := append(
		helpers.ProcessingStepTests("correct illumination"),
		[]helpers.ProcessingStepTest{
			{
				Name:          "Preserves uniform image",
				Source:        helpers.SolidGray(301, 301),
				ExpectedEqual: true,
			},
			{
				Name:   "Corrects smooth multiplicative illumination",
				Source: multiplicativelyIlluminatedImage8(401, 301),
				AssertResult: func(t *testing.T, result *cv.Mat) {
					assertIlluminationVariationReduced(t, result)
				},
			},
			{
				Name:   "Corrects smooth illumination while preserving tread-like detail",
				Source: illuminatedTreadPattern8(401, 301),
				AssertResult: func(t *testing.T, result *cv.Mat) {
					assertTreadContrastPreserved(t, result)
				},
			},
			{
				Name:   "Processes realistic tyre impression image",
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

			err := processor.CorrectIllumination(test.Source, &result)

			if test.ExpectedError != "" {
				assert.EqualError(t, err, test.ExpectedError, "error does not match")
				assert.True(t, result.Empty(), "result should be empty")
				return
			}

			assert.NoError(t, err, "should not have an error")
			assert.False(t, result.Empty(), "result should not be empty")

			assertNormalisationProcessingStep(t, test.Source, &result, processor.CorrectIllumination)

			if test.AssertResult != nil {
				test.AssertResult(t, &result)
			}
		})
	}
}

func assertNormalisationProcessingStep(
	t *testing.T,
	source, result *cv.Mat,
	step processors.ProcessingStep,
) {
	t.Helper()

	assert.False(t, source.Empty(), "source should not be empty")
	assert.False(t, result.Empty(), "result should not be empty")

	original := source.Clone()
	defer original.Close()

	repeatResult := cv.NewMat()
	defer repeatResult.Close()

	err := step(source, &repeatResult)
	assert.NoError(t, err, "repeat processing should not have an error")

	assert.Equal(
		t,
		result.ToBytes(),
		repeatResult.ToBytes(),
		"result and repeat-result should match",
	)

	assert.Equal(
		t,
		original.ToBytes(),
		source.ToBytes(),
		"source should not be mutated",
	)
}

func assertIlluminationVariationReduced(t *testing.T, result *cv.Mat) {
	t.Helper()

	before := float64(result.GetUCharAt(
		result.Rows()/2,
		40,
	))

	after := float64(result.GetUCharAt(
		result.Rows()/2,
		result.Cols()-41,
	))

	assert.Less(
		t,
		math.Abs(float64(before-after)),
		5.0,
		"illumination correction should reduce the broad intensity gradient",
	)
}

func assertTreadContrastPreserved(t *testing.T, result *cv.Mat) {
	t.Helper()

	left := result.GetUCharAt(result.Rows()/2, 75)
	right := result.GetUCharAt(result.Rows()/2, 85)

	assert.Greater(
		t,
		math.Abs(float64(left-right)),
		20.0,
		"illumination correction should preserve local tread-like contrast",
	)
}

func assert16BitInputNormalised(t *testing.T, result *cv.Mat) {
	t.Helper()

	left := result.GetUCharAt(result.Rows()/2, 40)
	right := result.GetUCharAt(result.Rows()/2, result.Cols()-41)

	assert.Less(
		t,
		math.Abs(float64(left-right)),
		1000.0,
		"16-bit illumination correction should reduce the broad intensity gradient",
	)
}

func multiplicativelyIlluminatedImage8(width, height int) *cv.Mat {
	mat := cv.NewMatWithSize(height, width, cv.MatTypeCV8UC1)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			base := 120.0
			illumination := 0.55 + 0.9*float64(x)/float64(width-1)
			value := base * illumination
			if value > 255 {
				value = 255
			}
			mat.SetUCharAt(y, x, uint8(value))
		}
	}

	return &mat
}

func illuminatedTreadPattern8(width, height int) *cv.Mat {
	mat := cv.NewMatWithSize(height, width, cv.MatTypeCV8UC1)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			base := 70.0
			if (x/10)%2 == 0 {
				base = 170
			}

			illumination := 0.60 + 0.80*float64(x)/float64(width-1)
			value := base * illumination
			if value > 255 {
				value = 255
			}

			mat.SetUCharAt(y, x, uint8(value))
		}
	}

	return &mat
}

func multiplicativelyIlluminatedImage16(width, height int) *cv.Mat {
	mat := cv.NewMatWithSize(height, width, cv.MatTypeCV16UC1)
	values, err := mat.DataPtrUint16()
	if err != nil {
		mat.Close()
		panic(err)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			base := 30000.0
			illumination := 0.55 + 0.9*float64(x)/float64(width-1)
			value := base * illumination
			if value > 65535 {
				value = 65535
			}

			values[y*width+x] = uint16(value)
		}
	}

	return &mat
}

package cv_tests_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	cv "gocv.io/x/gocv"
	"testing"
	"tyre-match-backend/cv/processors"
	"tyre-match-backend/tests/helpers"
)

func TestBaseProcessorGetName(t *testing.T) {
	cases := []struct {
		Name          string
		ProcessorName string
	}{
		{Name: "Returns enhancement name", ProcessorName: "enhancement"},
		{Name: "Returns normalisation name", ProcessorName: "normalisation"},
		{Name: "Returns empty name", ProcessorName: ""},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			processor := processors.BaseProcessor{
				Name:     test.ProcessorName,
				FileType: "",
			}

			assert.Equal(t, test.ProcessorName, processor.GetName())
		})
	}
}

func TestBaseProcessorGetFileType(t *testing.T) {
	cases := []struct {
		Name     string
		FileType string
	}{
		{Name: "Returns enhanced file type", FileType: "enhanced"},
		{Name: "Returns normalised file type", FileType: "normalised"},
		{Name: "Returns empty file type", FileType: ""},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {

			processor := processors.BaseProcessor{
				Name:     "",
				FileType: test.FileType,
			}

			assert.Equal(t, test.FileType, processor.GetFileType())
		})
	}
}

func TestBaseProcessorProcess(t *testing.T) {
	processingError := errors.New("step failed")

	cases := []struct {
		Name          string
		Source        *cv.Mat
		Steps         []processors.ProcessingStep
		ExpectedPixel uint8
		ExpectedSame  bool
		ExpectedError string
	}{
		{
			Name:          "Returns source when no steps are configured",
			Steps:         []processors.ProcessingStep{},
			Source:        helpers.SolidGray(20, 20),
			ExpectedPixel: 120,
			ExpectedSame:  true,
		},
		{
			Name:   "Runs all steps in order and returns final image",
			Source: helpers.SolidGray(20, 20),
			Steps: []processors.ProcessingStep{
				func(src, dst *cv.Mat) error {
					src.SetUCharAt(0, 0, 0)

					*dst = cv.NewMat()
					src.CopyTo(dst)

					return nil
				},
				func(src, dst *cv.Mat) error {
					if src.GetUCharAt(0, 0) != 0 {
						return errors.New("second step received the wrong image")
					}

					src.SetUCharAt(0, 0, 255)

					*dst = cv.NewMat()
					src.CopyTo(dst)

					return nil
				},
			},
			ExpectedPixel: 255,
		},
		{
			Name: "Rejects nil source image",
			Steps: []processors.ProcessingStep{
				func(src, dst *cv.Mat) error { return nil },
			},
			ExpectedError: `process received a nil image`,
		},
		{
			Name:   "Wraps processing step error",
			Source: helpers.SolidGray(20, 20),
			Steps: []processors.ProcessingStep{
				func(src, dst *cv.Mat) error { return processingError },
			},
			ExpectedError: `processor "test" step 1: step failed`,
		},
		{
			Name:   "Rejects nil image returned by step",
			Source: helpers.SolidGray(20, 20),
			Steps: []processors.ProcessingStep{
				func(src, dst *cv.Mat) error {
					*dst = cv.NewMat()
					return nil
				},
			},
			ExpectedError: `processor "test" step 1 returned a nil image`,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			processor := processors.BaseProcessor{
				Name:            "test",
				FileType:        "test-output",
				ProcessingSteps: test.Steps,
			}

			processed, err := processor.Process(test.Source)

			if test.ExpectedError != "" {
				assert.EqualError(t, err, test.ExpectedError)
				return
			}

			assert.NoError(t, err)

			if test.ExpectedSame {
				assert.Same(t, test.Source, processed)
			}

			assert.Equal(t, test.ExpectedPixel, processed.GetUCharAt(0, 0))
		})
	}
}

func TestBaseProcessorValidateSourceImage(t *testing.T) {
	baseProcessor := processors.NewBaseProcessor("test-processor", "test-file-type")

	emptyImage := cv.NewMat()
	colourImage := cv.NewMatWithSize(100, 100, cv.MatTypeCV8UC3)
	validImage := cv.NewMatWithSize(100, 100, cv.MatTypeCV8UC1)

	cases := []struct {
		Name          string
		Source        *cv.Mat
		ExpectedError error
	}{
		{
			Name:          "rejects nil image",
			Source:        nil,
			ExpectedError: errors.New("received a nil image"),
		},
		{
			Name:          "rejects empty image",
			Source:        &emptyImage,
			ExpectedError: errors.New("received an empty image"),
		},
		{
			Name:          "rejects colour image",
			Source:        &colourImage,
			ExpectedError: errors.New("received a colour image"),
		},
		{
			Name:          "accepts valid image",
			Source:        &validImage,
			ExpectedError: nil,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			err := baseProcessor.ValidateSourceImage(test.Source)
			assert.Equal(t, test.ExpectedError, err)
		})
	}
}

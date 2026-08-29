package helpers

import (
	"fmt"
	cv "gocv.io/x/gocv"
	"testing"
)

type ProcessingStepTest struct {
	Name   string
	Source *cv.Mat

	ExpectedEqual  bool
	ExpectedResult *cv.Mat
	ExpectedError  string

	AssertResult func(*testing.T, *cv.Mat)
}

func ProcessingStepTests(methodName string) []ProcessingStepTest {
	return []ProcessingStepTest{
		{
			Name:          "Reject nil image",
			Source:        nil,
			ExpectedError: fmt.Sprintf("%s received a nil image", methodName),
		},
		{
			Name:          "Rejects empty image",
			Source:        Empty(),
			ExpectedError: fmt.Sprintf("%s received an empty image", methodName),
		},
		{
			Name:          "Rejects colour image",
			Source:        Colour(),
			ExpectedError: fmt.Sprintf("%s received a colour image", methodName),
		},
		{
			Name:          "Handles zero width image",
			Source:        SolidGray(0, 10),
			ExpectedError: fmt.Sprintf("%s received an empty image", methodName),
		},
		{
			Name:          "Handles zero height image",
			Source:        SolidGray(10, 0),
			ExpectedError: fmt.Sprintf("%s received an empty image", methodName),
		},
	}
}

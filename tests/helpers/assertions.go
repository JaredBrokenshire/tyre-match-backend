package helpers

import (
	"github.com/stretchr/testify/assert"
	cv "gocv.io/x/gocv"
	"testing"
)

func AssertMatSameShapeAndType(t *testing.T, a, b *cv.Mat) {
	t.Helper()

	assert.Equal(t, a.Type(), b.Type(), "type should be the same")
	assert.Equal(t, a.Rows(), b.Rows(), "rows should be the same")
	assert.Equal(t, a.Cols(), b.Cols(), "cols should be the same")
}

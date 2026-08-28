package cv_tests

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Run the test
	code := m.Run()

	os.Exit(code)
}

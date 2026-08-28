package service_tests

import (
	"os"
	"testing"
	"tyre-match-backend/tests/helpers"
)

var (
	ts *helpers.TestServer
)

func TestMain(m *testing.M) {
	ts = helpers.NewTestServer("../../.env", false)
	// Close the database connection
	defer ts.S.Db.Close()

	// Run the test
	code := m.Run()

	os.Exit(code)
}

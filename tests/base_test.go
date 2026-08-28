package tests

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"tyre-match-backend/pkg/closer"
	"tyre-match-backend/tests/helpers"
)

var (
	ts *helpers.TestServer
)

func TestMain(m *testing.M) {

	utc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = utc

	ts = helpers.NewTestServer("../.env", true)

	// Close the database connection
	defer ts.S.Db.Close()

	// Run the test
	code := m.Run()

	os.Exit(code)
}

func RunTestCase(t *testing.T, test helpers.TestCase) {
	res := ts.ExecuteTestCase(&test)

	ValidateResults(t, test, res)
}

func ValidateResults(t *testing.T, test helpers.TestCase, res *httptest.ResponseRecorder) {

	if test.Expected.ExpectedCallBack != nil {
		test.Expected.ExpectedCallBack(res)
	}

	if res.Code != 0 {
		assert.Equal(t, test.Expected.StatusCode, res.Code)
	}

	if test.Expected.BodyPart != "" {
		isIn(t, res.Body.String(), test.Expected.BodyPart)
	}

	if len(test.Expected.BodyParts) > 0 {
		for _, expectedText := range test.Expected.BodyParts {
			isIn(t, res.Body.String(), expectedText)
		}
	}

	if test.Expected.BodyPartMissing != "" {
		isNotIn(t, res.Body.String(), test.Expected.BodyPartMissing)
	}

	if len(test.Expected.BodyPartsMissing) > 0 {
		for _, expectedText := range test.Expected.BodyPartsMissing {
			isNotIn(t, res.Body.String(), expectedText)
		}
	}

	if test.Expected.DatabaseCheck != nil {
		resultCount := CheckDatabase(test.Expected.DatabaseCheck)
		assert.Equal(t, test.Expected.DatabaseCheck.CountExpected, resultCount)

		if test.Expected.DatabaseCheck.Name != "" {
			assert.Equal(t, test.Expected.DatabaseCheck.CountExpected, resultCount, test.Expected.DatabaseCheck.Name)
		} else {
			assert.Equal(t, test.Expected.DatabaseCheck.CountExpected, resultCount)
		}
	}

	if len(test.Expected.DatabaseChecks) > 0 {
		for _, dbQuery := range test.Expected.DatabaseChecks {
			resultCount := CheckDatabase(dbQuery)
			if dbQuery.Name != "" {
				assert.Equal(t, dbQuery.CountExpected, resultCount, dbQuery.Name)
			} else {
				assert.Equal(t, dbQuery.CountExpected, resultCount)
			}
		}
	}

}

// / CheckDatabase runs a query and returns a count of rows found
func CheckDatabase(dbQuery *helpers.DatabaseCheck) int {
	var resultCount int
	var queryScopes []func(db *gorm.DB) *gorm.DB

	modelScope := func(db *gorm.DB) *gorm.DB { return db.Where(dbQuery.Model) }
	queryScopes = append(queryScopes, modelScope)

	if dbQuery.Scope != nil {
		queryScopes = append(queryScopes, dbQuery.Scope)
	}

	db := ts.S.Db
	// Enable debugging on this query if set in the dbQuery
	if dbQuery.DebugQuery {
		db = ts.S.Db.Debug()
	}

	db.Model(dbQuery.Model).Scopes(queryScopes...).Count(&resultCount)
	return resultCount
}

func isIn(t *testing.T, s, contains string, msgAndArgs ...interface{}) bool {
	t.Helper()

	ok := strings.Contains(s, contains)
	if !ok {
		return assert.Fail(t, fmt.Sprintf("%#v is not in %#v", contains, s), msgAndArgs...)
	}

	return true
}

func isNotIn(t *testing.T, s, contains string, msgAndArgs ...interface{}) bool {
	t.Helper()

	ok := strings.Contains(s, contains)
	if ok {
		return assert.Fail(t, fmt.Sprintf("%#v has been found in %#v", contains, s), msgAndArgs...)
	}

	return true
}

func createMultipartFile(t *testing.T, fieldName, path string) (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)

	fileMw := multipart.NewWriter(body)

	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}

	w, err := fileMw.CreateFormFile(fieldName, path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(w, file); err != nil {
		t.Fatal(err)
	}
	// close the writer before making the request
	assert.NoError(t, fileMw.Close())

	return body, fileMw
}

func createMultipartFiles(t *testing.T, fieldName string, paths []string) (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)

	fileMw := multipart.NewWriter(body)
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		part, err := fileMw.CreateFormFile(fieldName, path)
		if err != nil {
			closer.Close(file)
			t.Fatal(err)
		}
		if _, err := io.Copy(part, file); err != nil {
			closer.Close(file)
			t.Fatal(err)
		}
		closer.Close(file)
	}

	// close the writer before making the request
	assert.NoError(t, fileMw.Close())

	return body, fileMw
}

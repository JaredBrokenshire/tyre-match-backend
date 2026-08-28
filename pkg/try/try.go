package try

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ajg/form"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

type TestCase struct {
	TestName           string
	Request            Request
	RequestBody        interface{}
	RequestReader      io.Reader
	RequestContentType string
	RequestCookies     []*http.Cookie
	RequestHeaders     map[string]string
	Expected           ExpectedResponse
	AccessToken        string
	Setup              func(testCase *TestCase)
	Teardown           func(testCase *TestCase, res *HijackableResponseRecorder)
	DisplayResponse    bool
}

type Request struct {
	Method string
	Url    string
}

type ExpectedResponse struct {
	StatusCode       int
	BodyPart         string
	BodyParts        []string
	BodyPartMissing  string
	BodyPartsMissing []string
	DatabaseCheck    *DatabaseCheck
	DatabaseChecks   []*DatabaseCheck
	Headers          map[string]string
	ExpectedCallBack func(res *HijackableResponseRecorder)
}

type DatabaseCheck struct {
	Name          string
	Model         interface{}
	Scope         func(*gorm.DB) *gorm.DB
	CountExpected int
	DebugQuery    bool
}

func GenerateRequest(testCase *TestCase) (*http.Request, error) {

	var err error
	var testData []byte

	if testCase.RequestContentType == "" {
		testCase.RequestContentType = echo.MIMEApplicationJSON
	}

	switch testCase.RequestContentType {
	case echo.MIMEApplicationForm:
		if testCase.RequestBody != nil {
			values, err := form.EncodeToValues(testCase.RequestBody)
			if err != nil {
				return nil, err
			}

			testData = []byte(values.Encode())
		}
	case echo.MIMEApplicationJSON:
		if testCase.RequestBody != nil {
			testData, err = json.Marshal(testCase.RequestBody)
			if err != nil {
				return nil, err
			}
		}
	}

	var req *http.Request
	if testCase.RequestReader != nil {
		req, err = http.NewRequest(testCase.Request.Method, testCase.Request.Url, testCase.RequestReader)
	} else {
		req, err = http.NewRequest(testCase.Request.Method, testCase.Request.Url, bytes.NewBuffer(testData))
	}
	if err != nil {
		return nil, err
	}

	// Should always request JSON
	if testCase.RequestContentType != "" {
		req.Header.Set(echo.HeaderContentType, testCase.RequestContentType)
	} else {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	req.Header.Set(echo.HeaderXRealIP, "127.0.0.0")

	// Add cookies in if present
	if len(testCase.RequestCookies) > 0 {
		for _, cookie := range testCase.RequestCookies {
			req.AddCookie(cookie)
		}
	}

	// Add headers if present
	if len(testCase.RequestHeaders) > 0 {
		for headerKey, headerValue := range testCase.RequestHeaders {
			// Set requires to overrise content type. May need to be add if you need multiple headers
			// with the same key.
			req.Header.Set(headerKey, headerValue)
		}
	}

	// If user agent isn't set, then set explicity tpo a browser
	if req.Header.Get("User-Agent") == "" {
		SetDefaultUserAgent(req)
	}

	return req, nil
}
func SetDefaultUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
}
func ExecuteRequest(e *echo.Echo, req *http.Request) *HijackableResponseRecorder {
	// Create a new recorder then process request with server.
	rr := NewHijackableRecorder(nil)
	e.ServeHTTP(rr, req)
	return rr
}

func ValidateResults(t *testing.T, test *TestCase, res *HijackableResponseRecorder, db *gorm.DB) {

	if test.DisplayResponse {
		fmt.Println("Request Output: ")
		fmt.Println(res.Body.String())
	}

	if test.Expected.ExpectedCallBack != nil {
		test.Expected.ExpectedCallBack(res)
	}

	if res.Code != 0 {
		assert.Equal(t, test.Expected.StatusCode, res.Code)
	}

	if test.Expected.BodyPart != "" {
		assert.Contains(t, res.Body.String(), test.Expected.BodyPart)
	}
	if len(test.Expected.BodyParts) > 0 {
		for _, expectedText := range test.Expected.BodyParts {
			assert.Contains(t, res.Body.String(), expectedText)
		}
	}
	if test.Expected.BodyPartMissing != "" {
		assert.NotContains(t, res.Body.String(), test.Expected.BodyPartMissing)
	}

	if len(test.Expected.BodyPartsMissing) > 0 {
		for _, expectedText := range test.Expected.BodyPartsMissing {
			assert.NotContains(t, res.Body.String(), expectedText)
		}
	}

	if test.Expected.DatabaseCheck != nil {
		resultCount := CheckDatabase(db, test.Expected.DatabaseCheck)
		if test.Expected.DatabaseCheck.Name != "" {
			assert.Equal(t, test.Expected.DatabaseCheck.CountExpected, resultCount, test.Expected.DatabaseCheck.Name)
		} else {
			assert.Equal(t, test.Expected.DatabaseCheck.CountExpected, resultCount)
		}
	}

	if len(test.Expected.DatabaseChecks) > 0 {
		for _, dbQuery := range test.Expected.DatabaseChecks {
			resultCount := CheckDatabase(db, dbQuery)
			if dbQuery.Name != "" {
				assert.Equal(t, dbQuery.CountExpected, resultCount, dbQuery.Name)
			} else {
				assert.Equal(t, dbQuery.CountExpected, resultCount)
			}
		}
	}

	if test.Expected.Headers != nil {
		for headerKey, headerValue := range test.Expected.Headers {
			assert.Equal(t, headerValue, res.Header().Get(headerKey))
		}
	}

}

// CheckDatabase runs a query and returns a count of rows found
func CheckDatabase(db *gorm.DB, dbQuery *DatabaseCheck) int {
	var resultCount int
	var queryScopes []func(db *gorm.DB) *gorm.DB

	modelScope := func(db *gorm.DB) *gorm.DB { return db.Where(dbQuery.Model) }
	queryScopes = append(queryScopes, modelScope)

	if dbQuery.Scope != nil {
		queryScopes = append(queryScopes, dbQuery.Scope)
	}

	// Enable debugging on this query if set in the dbQuery
	if dbQuery.DebugQuery {
		db = db.Debug()
	}

	db.Model(dbQuery.Model).Scopes(queryScopes...).Count(&resultCount)
	return resultCount
}

func ExecuteTest(t *testing.T, e *echo.Echo, testCase *TestCase, db *gorm.DB) *HijackableResponseRecorder {

	// Run any setup required before we execute the request
	if testCase.Setup != nil {
		testCase.Setup(testCase)
	}
	req, err := GenerateRequest(testCase)
	if err != nil {
		t.Errorf("unable to Generate Request: %v", err)
		return nil
	}
	res := ExecuteRequest(e, req)
	ValidateResults(t, testCase, res, db)

	if testCase.Teardown != nil {
		testCase.Teardown(testCase, res)
	}

	return res
}

var BadContentTestCase = TestCase{
	TestName:           "Bad context type triggers error",
	RequestContentType: "application/protobuf",
	Expected: ExpectedResponse{
		StatusCode: 415,
		BodyPart:   "Unsupported Media Type",
	},
}

func AddBadContentTestCase(request Request) TestCase {
	badContent := BadContentTestCase
	badContent.Request = request
	return badContent
}

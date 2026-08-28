package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"tyre-match-backend/api"
	"tyre-match-backend/api/routes"
	"tyre-match-backend/db/migrations/process"
	"tyre-match-backend/db/repositories"
	"tyre-match-backend/pkg/dependencies"
	"tyre-match-backend/pkg/file_storage"
	"tyre-match-backend/pkg/validation"
	"tyre-match-backend/services"
	"tyre-match-backend/tests/mocks"
)

type TestServer struct {
	S *api.Server
}

func NewTestServer(envFileLoc string, migrate bool) *TestServer {
	err := godotenv.Load(envFileLoc)

	if err != nil {
		log.Println("Error loading .env file")
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_EXPOSE_DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Println(">>> ERROR, tried connecting with: ", dataSourceName)
		log.Fatalf("Unable to connect to test database: %v\n", err.Error())
	}

	// disable default logging of errors, we were going to be creating some on purpose!
	db.LogMode(false)

	ts := &TestServer{
		S: &api.Server{
			Echo:  echo.New(),
			Db:    db,
			Repos: repositories.NewRepos(db),
		},
	}

	// Set mock dependencies
	ts.S.Dependencies = dependencies.NewDependencyService(db)
	ts.S.Dependencies.SetFileStore(mocks.NewFileStoreMock())

	ts.S.Echo.Validator = validation.NewCustomValidator(validator.New())

	ts.S.Services = services.NewServices(
		ts.S.Repos,
		ts.S.Dependencies,
		// CV Processor
	)

	if migrate {
		ts.migrateDatabase()
	}

	routes.ConfigureRoutes(ts.S)

	return ts
}

func (ts *TestServer) migrateDatabase() {
	// This project requires DB_PORT and DB_HOST loaded from .env automatically. Let's make them match
	_ = os.Setenv("DB_HOST", os.Getenv("TEST_DB_HOST"))
	_ = os.Setenv("DB_PORT", os.Getenv("TEST_EXPOSE_DB_PORT"))

	process.Run()
}

// ExecuteTestCase Makes a request, and returns the response from ExecuteRequest.
func (ts *TestServer) ExecuteTestCase(testCase *TestCase) *httptest.ResponseRecorder {
	// Perform any setup needed for the test case
	if testCase.Setup != nil {
		testCase.Setup(testCase)
	}
	req := ts.GenerateRequest(testCase)
	res := ts.ExecuteRequest(req)
	// Perform any teardown needed for test case
	if testCase.Teardown != nil {
		testCase.Teardown()
	}
	return res
}

// ExecuteRequest Executes a request against the API. THis runs it locally against the handler
func (ts *TestServer) ExecuteRequest(req *http.Request) *httptest.ResponseRecorder {

	// Create a new recorder then process request with server.
	rr := httptest.NewRecorder()
	ts.S.Echo.ServeHTTP(rr, req)
	return rr
}

// ClearTable Clear a table and reset the autoincrement
func (ts *TestServer) ClearTable(tableName string) {
	err := ts.S.Db.Exec(fmt.Sprintf("DELETE FROM %v", tableName)).Error
	if err != nil {
		log.Fatalf("You can't clear that table. Err: %v", err)
	}
	err = ts.S.Db.Exec(fmt.Sprintf("ALTER TABLE %v AUTO_INCREMENT = 1", tableName)).Error
	if err != nil {
		log.Fatalf("Error setting autoincrement. Err: %v", err)
	}
}

// GetDb Return the database from the server
func (ts *TestServer) GetDb() *gorm.DB {
	return ts.S.Db
}

func (ts *TestServer) CreateOrDie(o interface{}) {
	err := ts.S.Db.Create(o).Error
	if err != nil {
		log.Panicf("Error creating object as part of a test: %v", err)
	}
}

func (ts *TestServer) SetDefaultTestHeaders(req *http.Request) {
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderXRealIP, "127.0.0.0")
}
func (ts *TestServer) SetDefaultUserAgent(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
}
func (ts *TestServer) GenerateRequest(testCase *TestCase) *http.Request {
	reqJson, err := json.Marshal(testCase.RequestBody)
	if err != nil {
		log.Printf("There was an error marshalling the json: %v", err)
	}

	var req *http.Request
	if testCase.RequestReader != nil {
		req, err = http.NewRequest(testCase.Request.Method, testCase.Request.Url, testCase.RequestReader)
	} else {
		req, err = http.NewRequest(testCase.Request.Method, testCase.Request.Url, bytes.NewBuffer(reqJson))
	}

	// Set IP address and default context type to JSON
	ts.SetDefaultTestHeaders(req)

	// Change Content Type if one is set
	if testCase.RequestContentType != "" {
		req.Header.Set(echo.HeaderContentType, testCase.RequestContentType)
	}

	// Add in a default user agent.
	ts.SetDefaultUserAgent(req)

	// Add cookies in if present
	if len(testCase.RequestCookies) > 0 {
		for _, cookie := range testCase.RequestCookies {
			req.AddCookie(cookie)
		}
	}

	if len(testCase.RequestHeaders) > 0 {
		for headerKey, headerValue := range testCase.RequestHeaders {
			// Set required to override content type. May need to be updated if you need multiple headers
			// with the same key.
			req.Header.Set(headerKey, headerValue)
		}
	}

	return req
}

// Dependencies

func (ts *TestServer) SetFileStore(store file_storage.Store) {
	ts.S.Dependencies.SetFileStore(store)

	ts.S.Services = services.NewServices(
		ts.S.Repos,
		ts.S.Dependencies,
	)
}

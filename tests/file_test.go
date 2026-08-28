package tests

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"tyre-match-backend/tests/helpers"
	"tyre-match-backend/tests/mocks"
)

func TestFile_Get(t *testing.T) {
	fileStoreMock := mocks.NewFileStoreMock()
	ts.S.Dependencies.SetFileStore(fileStoreMock)

	request := helpers.Request{
		Method: http.MethodGet,
		Url:    "/files/test-filepath",
	}

	cases := []helpers.TestCase{
		{
			Name: "Can't get file if there is an error from the file store",
			Setup: func(test *helpers.TestCase) {
				fileStoreMock.Reset()
				fileStoreMock.ReadFileError = errors.New("test error")
			},
			Request: request,
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusInternalServerError,
				BodyPart:   "Unable to read file",
			},
		},
		{
			Name: "Can get file from file store",
			Setup: func(test *helpers.TestCase) {
				fileStoreMock.Reset()
			},
			Request: request,
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.ReadFileCalls))
					assert.Contains(t, fileStoreMock.ReadFileCalls[0], "test-filepath")
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			RunTestCase(t, test)
		})
	}
}

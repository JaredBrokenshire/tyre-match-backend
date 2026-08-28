package try

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

type user struct {
	Name string
}

func TestExecuteRequest(t *testing.T) {

	tc := &TestCase{
		TestName: "Can create User",
		Request: Request{
			Method: http.MethodPost,
			Url:    "/users",
		},
	}

	req, err := GenerateRequest(tc)
	if err != nil {
		t.Fatalf("Error generating request: %v", err)
	}

	e := echo.New()
	res := ExecuteRequest(e, req)

	assert.Equal(t, `{"message":"Not Found"}`+"\n", res.Body.String())
}

func TestExecuteRequestAdditional(t *testing.T) {
	r := strings.NewReader("my request")
	tc := &TestCase{
		TestName:           "Can create User",
		RequestReader:      r,
		RequestContentType: "something",
		Request: Request{
			Method: http.MethodPost,
			Url:    "/users",
		},
	}

	req, err := GenerateRequest(tc)
	if err != nil {
		t.Fatalf("Error generating request: %v", err)
	}

	e := echo.New()
	res := ExecuteRequest(e, req)

	res.Closed()
	res.Hijack()

	assert.Equal(t, `{"message":"Not Found"}`+"\n", res.Body.String())
}

func TestValidateResults(t *testing.T) {
	tc := &TestCase{
		TestName: "Can create User",
		Request: Request{
			Method: http.MethodPost,
			Url:    "/users",
		},
		Expected: ExpectedResponse{
			StatusCode:       404,
			BodyPart:         "Not Found",
			BodyParts:        []string{"Not Found"},
			BodyPartMissing:  "This is Not Returned",
			BodyPartsMissing: []string{"This is Not Returned"},
		},
	}

	req, err := GenerateRequest(tc)
	if err != nil {
		t.Fatalf("Error generating request: %v", err)
	}

	e := echo.New()
	res := ExecuteRequest(e, req)

	ValidateResults(t, tc, res, nil)
}

func TestExecuteTest(t *testing.T) {
	adminRefreshCookie := &http.Cookie{
		Name: "test cookie",
	}

	u := &user{Name: "Matt Nelson"}

	tc := &TestCase{
		TestName: "Can create User",
		Request: Request{
			Method: http.MethodPost,
			Url:    "/users",
		},
		Setup:           func(testCase *TestCase) {},
		RequestBody:     u,
		RequestCookies:  []*http.Cookie{adminRefreshCookie},
		RequestHeaders:  map[string]string{"test-header": "header"},
		DisplayResponse: true,
		Expected: ExpectedResponse{
			StatusCode:       404,
			BodyPart:         "Not Found",
			BodyParts:        []string{"Not Found"},
			BodyPartMissing:  "This is Not Returned",
			BodyPartsMissing: []string{"This is Not Returned"},
		},
	}
	e := echo.New()
	ExecuteTest(t, e, tc, nil)

}

func TestExecuteTest_error(t *testing.T) {
	tc := &TestCase{
		TestName: "Can create User",
		Request: Request{
			Method: http.MethodPost,
			Url:    "/users",
		},
		RequestBody: func() {}, // You can't marshal a function
	}
	e := echo.New()

	// The test we need to pass in is a new one.
	altTest := &testing.T{}
	ExecuteTest(altTest, e, tc, nil)
	assert.True(t, altTest.Failed())
}

// TODO: make this in to a request reader as we need to read in JSON to the request and request it as a proto buffer.
//func TestBadRequest(t *testing.T) {
//
//	tc := AddBadContentTestCase(Request{
//		Method: http.MethodGet,
//		Url:    "/users",
//	})
//
//	log.Println("After Made: ", tc.RequestContentType)
//
//	req, err := GenerateRequest(&tc)
//	if err != nil {
//		t.Fatalf("Error generating request: %v", err)
//	}
//
//	// SEt up echo to bind JSON only
//	type Request struct {
//		Info string `json:"info"`
//	}
//	e := echo.New()
//	e.GET("/users", func(c echo.Context) error {
//		loginRequest := new(Request)
//		if err := c.Bind(loginRequest); err != nil {
//			return err
//		}
//		return c.String(http.StatusOK, "Hello, World!")
//	})
//	res := ExecuteRequest(e, req)
//
//	ValidateResults(t, &tc, res, nil)
//}

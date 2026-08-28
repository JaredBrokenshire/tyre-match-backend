package helpers

import (
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"net/http/httptest"
)

const UserId = 1

type TestCase struct {
	Name               string
	Request            Request
	RequestBody        interface{}
	RequestReader      io.Reader
	RequestContentType string
	RequestCookies     []*http.Cookie
	RequestHeaders     map[string]string
	Expected           ExpectedResponse
	Setup              func(test *TestCase)
	Teardown           func()
}

type PathParam struct {
	Name  string
	Value string
}

type Request struct {
	Method    string
	Url       string
	PathParam *PathParam
}

type DatabaseCheck struct {
	Name          string
	Model         interface{}
	Scope         func(*gorm.DB) *gorm.DB
	CountExpected int
	DebugQuery    bool
}

type ExpectedResponse struct {
	StatusCode       int
	BodyPart         string
	BodyParts        []string
	BodyPartMissing  string
	BodyPartsMissing []string
	DatabaseCheck    *DatabaseCheck
	DatabaseChecks   []*DatabaseCheck
	ExpectedCallBack func(res *httptest.ResponseRecorder)
}

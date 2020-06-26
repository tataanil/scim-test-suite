package util

import (
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"strings"
)

type Suite struct {
	suite.Suite
	url        string
	middleware func(req *http.Request) *http.Request
}

func (suite *Suite) IsNumber(i interface{}) {
	suite.IsType(0.0, i)
}

func (suite *Suite) StatusOK(status int) {
	suite.Equal(status, http.StatusOK)
}

func (suite *Suite) StatusCreated(status int) {
	suite.Equal(http.StatusCreated, status)
}

func (suite *Suite) StatusForbidden(status int) {
	suite.Equal(http.StatusForbidden, status)
}

func (suite *Suite) StatusNotFound(status int) {
	suite.Equal(http.StatusNotFound, status)
}

func (suite *Suite) StatusConflict(status int) {
	suite.Equal(http.StatusConflict, status)
}

func (suite *Suite) Get(path string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, suite.url + path, nil)
	return suite.Do(req)
}

func (suite *Suite) Post(path string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, suite.url + path, body)
	req.Header.Set("Content-Type", "application/scim+json")
	return suite.Do(req)
}

func (suite *Suite) Do(req *http.Request) (*http.Response, error) {
	if suite.middleware != nil {
		req = suite.middleware(req)
	}
	return http.DefaultClient.Do(req)
}

func (suite *Suite) Middleware(callback func(req *http.Request) *http.Request) {
	suite.middleware = callback
}

func (suite *Suite) BaseURL(baseURL string) {
	suite.url = strings.TrimSuffix(baseURL, "/")
}

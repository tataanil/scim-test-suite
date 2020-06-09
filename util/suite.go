package util

import (
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"strings"
)

type Suite struct {
	suite.Suite
	url string
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
	return http.Get(suite.url + path)
}

func (suite *Suite) Post(path string, body io.Reader) (*http.Response, error) {
	return http.Post(suite.url+path, "application/scim+json", body)
}

func (suite *Suite) BaseURL(baseURL string) {
	suite.url = strings.TrimSuffix(baseURL, "/")
}

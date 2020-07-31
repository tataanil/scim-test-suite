package util

import (
	"io"
	"net/http"
)

func (suite *Suite) Get(path string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, suite.url+path, nil)
	suite.Require().NoError(err)
	return suite.Do(req)
}

func (suite *Suite) GetOk(path string) *http.Response {
	resp := suite.Get(path)
	suite.Require().Equal(resp.StatusCode, http.StatusOK)
	return resp
}

func (suite *Suite) Post(path string, body io.Reader) *http.Response {
	req, err := http.NewRequest(http.MethodPost, suite.url+path, body)
	suite.Require().NoError(err)
	req.Header.Set("Content-Type", "application/scim+json")
	return suite.Do(req)
}

func (suite *Suite) Do(req *http.Request) *http.Response {
	if suite.middleware != nil {
		req = suite.middleware(req)
	}
	resp, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)
	return resp
}

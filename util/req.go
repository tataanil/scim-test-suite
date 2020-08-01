package util

import (
	"io"
	"net/http"
)

func (suite *Suite) Delete(path string) *http.Response {
	return suite.receive(path, http.MethodDelete)
}

func (suite *Suite) Get(path string) *http.Response {
	return suite.receive(path, http.MethodGet)
}

func (suite *Suite) receive(path string, method string) *http.Response {
	req, err := http.NewRequest(method, suite.url+path, nil)
	suite.Require().NoError(err)
	return suite.Do(req)
}

func (suite *Suite) GetOk(path string) *http.Response {
	resp := suite.Get(path)
	suite.Require().Equal(resp.StatusCode, http.StatusOK)
	return resp
}

func (suite *Suite) Patch(path string, body io.Reader) *http.Response {
	return suite.send(path, body, http.MethodPatch)
}

func (suite *Suite) Post(path string, body io.Reader) *http.Response {
	return suite.send(path, body, http.MethodPost)
}

func (suite *Suite) Put(path string, body io.Reader) *http.Response {
	return suite.send(path, body, http.MethodPut)
}

func (suite *Suite) send(path string, body io.Reader, method string) *http.Response {
	req, err := http.NewRequest(method, suite.url+path, body)
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

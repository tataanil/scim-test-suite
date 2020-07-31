package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (suite *Suite) ReadAllToMap(resp *http.Response) map[string]interface{} {
	var mapData map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err)
	d := json.NewDecoder(bytes.NewBuffer(raw))
	d.UseNumber()
	suite.Require().NoError(d.Decode(&mapData))
	return mapData
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

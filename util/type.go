package util

import (
	"encoding/json"
	pS "github.com/di-wu/scim-test-suite/schema"
	"github.com/elimity-com/scim/schema"
	"strconv"
)

func (suite *Suite) GetBool(name string, mapData map[string]interface{}) bool {
	suite.Require().NotNil(mapData[name])
	boolean, ok := mapData[name].(bool)
	suite.Require().True(ok)
	return boolean
}

func (suite *Suite) GetMap(name string, mapData map[string]interface{}) map[string]interface{} {
	suite.Require().NotNil(mapData[name])
	mapData, ok := mapData[name].(map[string]interface{})
	suite.Require().True(ok)
	return mapData
}

func (suite *Suite) GetSlice(name string, mapData map[string]interface{}) []interface{} {
	suite.Require().NotNil(mapData[name])
	slice, ok := mapData[name].([]interface{})
	suite.Require().True(ok)
	return slice
}

func (suite *Suite) GetSliceOfString(name string, mapData map[string]interface{}) []string {
	strings := make([]string, 0)
	for _, v := range suite.GetSlice(name, mapData) {
		strings = append(strings, suite.IsString(v))
	}
	return strings
}

func (suite *Suite) GetFloat64(name string, mapData map[string]interface{}) float64 {
	suite.Require().NotNil(mapData[name])
	n, ok := mapData[name].(json.Number)
	suite.Require().True(ok)
	f64, err := strconv.ParseFloat(string(n), 64)
	suite.Require().NoError(err)
	return f64
}

func (suite *Suite) GetInt(name string, mapData map[string]interface{}) int {
	return int(suite.GetInt64(name, mapData))
}

func (suite *Suite) GetInt64(name string, mapData map[string]interface{}) int64 {
	suite.Require().NotNil(mapData[name])
	n, ok := mapData[name].(json.Number)
	suite.Require().True(ok)
	i64, err := strconv.ParseInt(string(n), 10, 64)
	suite.Require().NoError(err)
	return i64
}

func (suite *Suite) GetString(name string, mapData map[string]interface{}) string {
	suite.Require().NotNil(mapData[name])
	return suite.IsString(mapData[name])
}

func (suite *Suite) IsMap(i interface{}) map[string]interface{} {
	resource, ok := i.(map[string]interface{})
	suite.Require().True(ok)
	return resource
}

func (suite *Suite) IsSchema(resource map[string]interface{}) schema.Schema {
	s, err := pS.ParseJSONSchema(resource)
	suite.Require().NoError(err)
	return s
}

func (suite *Suite) IsString(i interface{}) string {
	str, ok := i.(string)
	suite.Require().True(ok)
	return str
}

func (suite *Suite) IsSlice(i interface{}) []interface{} {
	slice, ok := i.([]interface{})
	suite.Require().True(ok)
	return slice
}

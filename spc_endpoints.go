package suite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// RFC: https://tools.ietf.org/html/rfc7644#section-4

func (suite *SCIMTestSuite) TestServiceProviderConfigurationEndpoints() {
	suite.Run("ServiceProviderConfig", func() {
		suite.testServiceProviderConfigEndpoint()
	})

	suite.Run("Schemas", func() {
		suite.testSchemasEndpoint()
	})

	suite.Run("ResourceTypes", func() {
		suite.testResourceTypesEndpoint()
	})

	suite.Run("ForbiddenFilter", func() {
		suite.testResourceTypesEndpoint()
	})
}

func (suite *SCIMTestSuite) testServiceProviderConfigEndpoint() {
	// An HTTP GET to "/ServiceProviderConfig" will return a JSON structure that describes the SCIM specification features
	// available on a service provider.
	resp, err := suite.Get("/ServiceProviderConfig")
	suite.Require().NoError(err)
	suite.StatusOK(resp.StatusCode)

	// This endpoint SHALL return responses with a JSON object using a "schemas" attribute of
	// "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig".
	var mapData map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err)
	suite.Require().NoError(json.Unmarshal(raw, &mapData))

	suite.NotNil(mapData["schemas"])
	schemas, ok := mapData["schemas"].([]interface{})
	suite.Require().True(ok)

	suite.Require().Len(schemas, 1)
	schema, ok := schemas[0].(string)
	suite.Require().True(ok)
	suite.Equal(schema, "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig")
}

func (suite *SCIMTestSuite) testSchemasEndpoint() {
	// An HTTP GET to "/Schemas" is used to retrieve information about resource schemas supported by a SCIM service
	// provider.
	resp, err := suite.Get("/Schemas")
	suite.Require().NoError(err)

	// An HTTP GET to the endpoint "/Schemas" SHALL return all supported schemas in ListResponse format.
	var mapData map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err)
	suite.Require().NoError(json.Unmarshal(raw, &mapData))

	suite.NotNil(mapData["totalResults"])
	totalResults, ok := mapData["totalResults"].(float64)
	suite.Require().True(ok)

	suite.NotNil(mapData["Resources"])
	resources, ok := mapData["Resources"].([]interface{})
	suite.Require().True(ok)
	suite.Equal(int(totalResults), len(resources))

	// Individual schema definitions can be returned by appending the schema URI to the "/Schemas" endpoint.
	for _, resourceJSON := range resources {
		resource, ok := resourceJSON.(map[string]interface{})
		suite.Require().True(ok)
		suite.NotNil(resource["id"])
		id, ok := resource["id"].(string)
		suite.Require().True(ok)
		resp, err := suite.Get(fmt.Sprintf("/Schemas/%s", id))
		suite.Require().NoError(err)
		suite.StatusOK(resp.StatusCode)
	}
}

func (suite *SCIMTestSuite) testResourceTypesEndpoint() {
	// An HTTP GET to "/ResourceTypes" is used to discover the types of resources available on a SCIM service provider.
	resp, err := suite.Get("/ResourceTypes")
	suite.Require().NoError(err)
	suite.StatusOK(resp.StatusCode)

	var mapData map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err)
	suite.Require().NoError(json.Unmarshal(raw, &mapData))

	suite.NotNil(mapData["totalResults"])
	totalResults, ok := mapData["totalResults"].(float64)
	suite.Require().True(ok)

	suite.NotNil(mapData["Resources"])
	resources, ok := mapData["Resources"].([]interface{})
	suite.Require().True(ok)
	suite.Equal(int(totalResults), len(resources))

	for _, resourceJSON := range resources {
		resource, ok := resourceJSON.(map[string]interface{})
		suite.Require().True(ok)
		suite.NotNil(resource["id"])
		id, ok := resource["id"].(string)
		suite.Require().True(ok)
		resp, err := suite.Get(fmt.Sprintf("/ResourceTypes/%s", id))
		suite.Require().NoError(err)
		suite.StatusOK(resp.StatusCode)
	}
}

func (suite *SCIMTestSuite) testForbiddenFilter() {
	// If a "filter" is provided, the service provider SHOULD respond with HTTP status code 403 (Forbidden) to ensure
	// that clients cannot incorrectly assume that any matching conditions specified in a filter are true.
	respS, err := suite.Get("/Schemas?startIndex=1&count=10")
	suite.Require().NoError(err)
	suite.StatusForbidden(respS.StatusCode)

	respRT, err := suite.Get("/ResourceTypes?startIndex=1&count=10")
	suite.Require().NoError(err)
	suite.StatusForbidden(respRT.StatusCode)
}

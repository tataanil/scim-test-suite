package suite

import (
	"fmt"
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
	// available on a service provider. It SHALL return responses with a JSON object using a "schemas" attribute of
	// "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig".
	var (
		resp    = suite.GetOk("/ServiceProviderConfig")
		mapData = suite.ReadAllToMap(resp)
		schemas = suite.GetSliceOfString("schemas", mapData)
	)
	suite.Require().Len(schemas, 1)
	suite.Equal(schemas[0], "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig")
}

func (suite *SCIMTestSuite) testSchemasEndpoint() {
	// An HTTP GET to "/Schemas" is used to retrieve information about resource schemas supported by a SCIM service
	// provider. It SHALL return all supported schemas in ListResponse format.
	var (
		resp         = suite.GetOk("/Schemas")
		mapData      = suite.ReadAllToMap(resp)
		totalResults = suite.GetInt("totalResults", mapData)
		resources    = suite.GetSlice("Resources", mapData)
	)
	suite.Len(resources, totalResults)

	// Individual schema definitions can be returned by appending the schema URI to the "/Schemas" endpoint.
	for _, resourceJSON := range resources {
		var (
			resource = suite.IsMap(resourceJSON)
			id       = suite.GetString("id", resource)
			_        = suite.GetOk(fmt.Sprintf("/Schemas/%s", id))
		)
	}
}

func (suite *SCIMTestSuite) testResourceTypesEndpoint() {
	// An HTTP GET to "/ResourceTypes" is used to discover the types of resources available on a SCIM service provider.
	var (
		resp         = suite.GetOk("/ResourceTypes")
		mapData      = suite.ReadAllToMap(resp)
		totalResults = suite.GetInt("totalResults", mapData)
		resources    = suite.GetSlice("Resources", mapData)
	)
	suite.Len(resources, totalResults)

	for _, resourceJSON := range resources {
		var (
			resource = suite.IsMap(resourceJSON)
			id       = suite.GetString("id", resource)
			_        = suite.GetOk(fmt.Sprintf("/ResourceTypes/%s", id))
		)
	}
}

func (suite *SCIMTestSuite) testForbiddenFilter() {
	// If a "filter" is provided, the service provider SHOULD respond with HTTP status code 403 (Forbidden) to ensure
	// that clients cannot incorrectly assume that any matching conditions specified in a filter are true.
	respS := suite.Get("/Schemas?startIndex=1&count=10")
	suite.StatusForbidden(respS.StatusCode)

	respRT := suite.Get("/ResourceTypes?startIndex=1&count=10")
	suite.StatusForbidden(respRT.StatusCode)
}

package suite

import (
	"encoding/json"
	"io/ioutil"

	s "github.com/di-wu/scim-test-suite/schema"
	"github.com/elimity-com/abnf"
	"github.com/elimity-com/abnf/core"
	"github.com/elimity-com/abnf/operators"
)

// RFC: https://tools.ietf.org/html/rfc7643#section-2

var (
	attributeNameValidator = attrName()
)

func attrName() func(s string) bool {
	g := abnf.ParserGenerator{
		RawABNF: []byte(
			"ATTRNAME   = ALPHA *(nameChar)\n" +
				"nameChar   = \"$\" / \"-\" / \"_\" / DIGIT / ALPHA\n",
		),
		ExternalABNF: map[string]operators.Operator{
			"ALPHA": core.ALPHA(),
			"DIGIT": core.DIGIT(),
		},
	}
	attrNameOperator := g.GenerateABNFAsOperators()["ATTRNAME"]
	return func(attrName string) bool {
		return attrNameOperator([]byte(attrName)).Best() != nil
	}
}

func (suite *SCIMTestSuite) TestSchemas() {
	resp, err := suite.Get("/Schemas")
	suite.Require().NoError(err)

	var mapData map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err)
	suite.Require().NoError(json.Unmarshal(raw, &mapData))

	suite.NotNil(mapData["Resources"])
	schemas, ok := mapData["Resources"].([]interface{})
	suite.Require().True(ok)
	for _, rawSchema := range schemas {
		schema, err := s.MetaSchema.Validate(rawSchema)
		suite.Require().Nil(err)
		suite.NotEmpty(schema)
	}
}

func (suite *SCIMTestSuite) TestAttributes() {
	resp, err := suite.Get("/Schemas")
	suite.Require().NoError(err)

	var mapData map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	suite.Require().NoError(err)
	suite.Require().NoError(json.Unmarshal(raw, &mapData))

	suite.NotNil(mapData["Resources"])
	schemas, ok := mapData["Resources"].([]interface{})
	suite.Require().True(ok)
	for _, rawSchema := range schemas {
		schema, ok := rawSchema.(map[string]interface{})
		suite.Require().True(ok)
		suite.NotNil(schema["attributes"])
		rawAttributes, ok := schema["attributes"].([]interface{})
		suite.Require().True(ok)
		for _, rawAttribute := range rawAttributes {
			attribute, ok := rawAttribute.(map[string]interface{})
			suite.Require().True(ok)
			suite.NotNil(attribute["name"])
			name, ok := attribute["name"].(string)
			suite.Require().True(ok)

			// Attribute names MUST conform to the following ABNF rules:
			//		ATTRNAME   = ALPHA *(nameChar)
			//		nameChar   = "$" / "-" / "_" / DIGIT / ALPHA
			suite.True(attributeNameValidator(name))
		}
	}
}

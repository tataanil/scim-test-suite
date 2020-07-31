package suite

import s "github.com/di-wu/scim-test-suite/schema"

func (suite *SCIMTestSuite) TestSchemas() {
	var (
		resp    = suite.GetOk("/Schemas")
		mapData = suite.ReadAllToMap(resp)
		schemas = suite.GetSlice("Resources", mapData)
	)

	for _, rawSchema := range schemas {
		schema, err := s.MetaSchema.Validate(rawSchema)
		suite.Require().Nil(err)
		suite.NotEmpty(schema)
	}
}

func (suite *SCIMTestSuite) TestAttributes() {
	// Attribute names MUST conform to the following ABNF rules:
	//	ATTRNAME   = ALPHA *(nameChar)
	//	nameChar   = "$" / "-" / "_" / DIGIT / ALPHA
	var (
		resp    = suite.GetOk("/Schemas")
		mapData = suite.ReadAllToMap(resp)
		schemas = suite.GetSlice("Resources", mapData)
	)

	suite.ForEachMap(schemas, func(m map[string]interface{}) {
		for k := range m {
			suite.IsValidAttributeName(k)
		}
	})
}

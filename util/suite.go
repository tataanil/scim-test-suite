package util

import (
	"github.com/elimity-com/abnf"
	"github.com/elimity-com/abnf/core"
	"github.com/elimity-com/abnf/operators"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strings"
)

type Suite struct {
	suite.Suite
	url        string
	middleware func(req *http.Request) *http.Request

	attrNameValidator operators.Operator
}

func (suite *Suite) SetupSuite() {
	rawABNF := "" +
		"ATTRNAME   = ALPHA *(nameChar)\n" +
		"nameChar   = \"$\" / \"-\" / \"_\" / DIGIT / ALPHA\n"
	g := abnf.ParserGenerator{
		RawABNF: []byte(rawABNF),
		ExternalABNF: map[string]operators.Operator{
			"ALPHA": core.ALPHA(),
			"DIGIT": core.DIGIT(),
		},
	}
	suite.attrNameValidator = g.GenerateABNFAsOperators()["ATTRNAME"]
}

func (suite *Suite) Middleware(callback func(req *http.Request) *http.Request) {
	suite.middleware = callback
}

func (suite *Suite) BaseURL(baseURL string) {
	suite.url = strings.TrimSuffix(baseURL, "/")
}

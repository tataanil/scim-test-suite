package suite

import (
	"net/http"
	"strings"

	"github.com/stretchr/testify/suite"
)

type SCIMTestSuite struct {
	suite.Suite
	url string
}

func (suite *SCIMTestSuite) StatusOK (status int) {
	suite.Equal(status, http.StatusOK)
}

func (suite *SCIMTestSuite) StatusForbidden (status int) {
	suite.Equal(status, http.StatusForbidden)
}

func (suite *SCIMTestSuite) Get(path string) (*http.Response, error) {
	return http.Get(suite.url + path)
}

func (suite *SCIMTestSuite) BaseURL(baseURL string) {
	suite.url = strings.TrimSuffix(baseURL, "/")
}

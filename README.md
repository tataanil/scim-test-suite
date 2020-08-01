# SCIM v2 Test Suite

```go
package suite

import (
	"testing"

	. "github.com/di-wu/scim-test-suite"
	"github.com/stretchr/testify/suite"
)

func TestSCIM(t *testing.T) {
	s := new(SCIMTestSuite)
	s.BaseURL("https://path.to.scim/v2")
	suite.Run(t, s)
}
```

### RFC7644 Protocol
#### Table of Contents
The following list includes all the parts of the RFC that are covered by the test suite.
- [x] 4\. Service Provider Configuration Endpoints

### [Identity Providers](./idp/)
#### [Okta](./idp/okta/)
#### [AzureAD](./idp/azure_ad/)

# Okta SCIM v2 Spec Test
Basic tests to see if your SCIM server will work with Okta.

[Official Okta Docs](https://developer.okta.com/docs/guides/build-provisioning-integration/test-scim-api/)

```go
package suite

import (
	"testing"

	"github.com/di-wu/scim-testsuite/idp/okta"
	"github.com/stretchr/testify/suite"
)

func TestOkta(t *testing.T)  {
	s := new(okta.TestSuite)
	s.BaseURL("https://path.to.scim/v2")
	
	// Optional: customize callback functions
	s.SetInvalidID(func() string {
		return "invalidID"
	})

	suite.Run(t, s)
}
``` 
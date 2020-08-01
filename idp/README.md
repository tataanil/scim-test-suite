# IdP SCIM v2 Spec Tests
Basic tests to see if your SCIM server will work with a certain identity provider.

[Official Okta Docs](https://developer.okta.com/docs/guides/build-provisioning-integration/test-scim-api/) \
[Official AzureAD Docs](https://docs.microsoft.com/en-us/azure/active-directory/app-provisioning/use-scim-to-provision-users-and-groups)

```go
package suite

import (
	"testing"

	"github.com/di-wu/scim-test-suite/idp/okta"
	"github.com/stretchr/testify/suite"
)

func TestIdP(t *testing.T)  {
	s := new(okta.TestSuite) // or w/ another test suite
	s.BaseURL("https://path.to.scim/v2")
	
	// Optional: customize callback functions (unique for every test suite)
	s.SetInvalidID(func() string {
		return "invalidID"
	})

	suite.Run(t, s)
}
``` 
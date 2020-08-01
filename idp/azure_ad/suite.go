package azure

import (
	"fmt"
	"github.com/di-wu/scim-test-suite/util"
)

type TestSuite struct {
	util.Suite
}

func (s *TestSuite) createUserBody(userName, displayName string) map[string]interface{} {
	return map[string]interface{}{
		"userName":    userName,
		"active":      true,
		"displayName": displayName,
		"schemas":     []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		"name": map[string]interface{}{
			"formatted":  "Quint Daenen",
			"familyName": "Quint",
			"givenName":  "Daenen",
		},
		"emails": []map[string]interface{}{
			{
				"primary": true,
				"type":    "work",
				"value":   "quint@elimity.com",
			},
			{
				"primary": false,
				"type":    "home",
				"value":   "me@di-wu.be",
			},
			{
				"primary": false,
				"type":    "other",
				"value":   fmt.Sprintf("%s@elimity.com", displayName),
			},
		},
	}
}

func (s *TestSuite) createEnterpriseUserBody(userName, displayName string) map[string]interface{} {
	return map[string]interface{}{
		"userName":    userName,
		"active":      true,
		"displayName": displayName,
		"schemas": []string{
			"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
			"urn:ietf:params:scim:schemas:core:2.0:User",
		},
		"name": map[string]interface{}{
			"formatted":  "Quint Daenen",
			"familyName": "Quint",
			"givenName":  "Daenen",
		},
		"emails": []map[string]interface{}{
			{
				"primary": true,
				"type":    "work",
				"value":   "quint@elimity.com",
			},
			{
				"primary": false,
				"type":    "home",
				"value":   "me@di-wu.be",
			},
		},
		"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
			"Department": "Engineering",
			"Manager": map[string]interface{}{
				"value": "M.",
			},
		},
	}
}

func (s *TestSuite) createGroup(displayName string, userIDs ...string) map[string]interface{} {
	var members []map[string]interface{}
	for _, id := range userIDs {
		members = append(members, map[string]interface{}{
			"value": id,
		})
	}

	return map[string]interface{}{
		"schemas":     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
		"displayName": displayName,
		"members":     members,
	}
}

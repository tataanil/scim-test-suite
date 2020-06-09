package okta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

// SOURCE: https://developer.okta.com/standards/SCIM/SCIMFiles/Okta-SCIM-20-SPEC-Test.json

// Required Test: Test Users endpoint.
func (s *TestSuite) TestGetFirstUser() {
	resp, err := s.Get("/Users?count=1&startIndex=1")
	s.Require().NoError(err)

	// Assertion 0
	s.Run("StatusCode", func() {
		s.StatusOK(resp.StatusCode)
	})

	var mapData map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().NoError(json.Unmarshal(raw, &mapData))

	// Assertion 1
	s.Run("ResourcesNotEmpty", func() {
		s.NotEmpty(mapData["Resources"])
	})

	// Assertion 2
	s.Run("ContainsSchema", func() {
		s.Contains(mapData["schemas"], "urn:ietf:params:scim:api:messages:2.0:ListResponse")
	})

	// Assertion 3
	s.Run("ItemsPerPageIsNumber", func() {
		s.IsNumber(mapData["itemsPerPage"])
	})

	// Assertion 4
	s.Run("StartIndexIsNumber", func() {
		s.IsNumber(mapData["startIndex"])
	})

	// Assertion 5
	s.Run("TotalResultsIsNumber", func() {
		s.IsNumber(mapData["totalResults"])
	})

	resources, ok := mapData["Resources"].([]interface{})
	s.Require().True(ok)
	s.Require().NotEmpty(resources)
	entity, ok := resources[0].(map[string]interface{})
	s.Require().True(ok)

	// Assertion 6
	s.Run("IDNotEmpty", func() {
		s.NotEmpty(entity["id"])
	})

	s.Require().NotEmpty(entity["name"])
	name, ok := entity["name"].(map[string]interface{})
	s.Require().True(ok)

	// Assertion 7
	s.Run("FamilyNameNotEmpty", func() {
		s.NotEmpty(name["familyName"])
	})

	// Assertion 8
	s.Run("GivenNameNotEmpty", func() {
		s.NotEmpty(name["givenName"])
	})

	// Assertion 9
	s.Run("UserNameNotEmpty", func() {
		s.NotEmpty(entity["userName"])
	})

	// Assertion 10
	s.Run("ActiveNotEmpty", func() {
		s.NotEmpty(entity["active"])
	})

	s.Require().NotEmpty(entity["emails"])
	emails, ok := entity["emails"].([]interface{})
	s.Require().True(ok)
	s.Require().NotEmpty(emails)
	email, ok := emails[0].(map[string]interface{})
	s.Require().True(ok)

	// Assertion 11
	s.Run("FirstEmailValueNotEmpty", func() {
		s.NotEmpty(email["value"])
	})
}

// Required Test: Get Users/{{id}}.
func (s *TestSuite) TestGetExistingUser() {
	_resp, err := s.Get("/Users?count=1&startIndex=1")
	s.Require().NoError(err)
	var _map map[string]interface{}
	_raw, err := ioutil.ReadAll(_resp.Body)
	s.Require().NoError(err)
	s.Require().NoError(json.Unmarshal(_raw, &_map))
	_resources, ok := _map["Resources"].([]interface{})
	s.Require().True(ok)
	s.Require().NotEmpty(_resources)
	_entity, ok := _resources[0].(map[string]interface{})
	s.Require().True(ok)
	s.NotEmpty(_entity["id"])
	id, ok := _entity["id"]
	s.Require().True(ok)

	resp, err := s.Get(fmt.Sprintf("/Users/%s", id))
	s.Require().NoError(err)

	// Assertion 0
	s.Run("StatusCode", func() {
		s.StatusOK(resp.StatusCode)
	})

	var entity map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().NoError(json.Unmarshal(raw, &entity))

	// Assertion 1
	s.Run("IDNotEmpty", func() {
		s.NotEmpty(entity["id"])
	})

	s.Require().NotEmpty(entity["name"])
	name, ok := entity["name"].(map[string]interface{})
	s.Require().True(ok)

	// Assertion 2
	s.Run("FamilyNameNotEmpty", func() {
		s.NotEmpty(name["familyName"])
	})

	// Assertion 3
	s.Run("GivenNameNotEmpty", func() {
		s.NotEmpty(name["givenName"])
	})

	// Assertion 4
	s.Run("UserNameNotEmpty", func() {
		s.NotEmpty(entity["userName"])
	})

	// Assertion 5
	s.Run("ActiveNotEmpty", func() {
		s.NotEmpty(entity["active"])
	})

	s.Require().NotEmpty(entity["emails"])
	emails, ok := entity["emails"].([]interface{})
	s.Require().True(ok)
	s.Require().NotEmpty(emails)
	email, ok := emails[0].(map[string]interface{})
	s.Require().True(ok)

	// Assertion 6
	s.Run("FirstEmailValueNotEmpty", func() {
		s.NotEmpty(email["value"])
	})

	// Assertion 7
	s.Run("IDsMatch", func() {
		s.Equal(id, entity["id"])
	})
}

// Required Test: Test invalid User by userName.
func (s *TestSuite) TestGetInvalidUserByUserName() {
	filter := url.Values{}
	filter.Set("filter", fmt.Sprintf("userName eq \"%s\"", s.RandomEmail()))
	resp, err := s.Get(fmt.Sprintf("/Users?%s", filter.Encode()))
	s.Require().NoError(err)

	// Assertion 0
	s.Run("StatusCode", func() {
		s.StatusOK(resp.StatusCode)
	})

	var mapData map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().NoError(json.Unmarshal(raw, &mapData))

	// Assertion 1
	s.Run("ContainsSchema", func() {
		s.Contains(mapData["schemas"], "urn:ietf:params:scim:api:messages:2.0:ListResponse")
	})

	// Assertion 2
	s.Run("NoResults", func() {
		s.Equal(0.0, mapData["totalResults"])
	})
}

// Required Test: Test invalid User by ID.
func (s *TestSuite) TestGetInvalidUser() {
	resp, err := s.Get(fmt.Sprintf("/Users/%s", s.InvalidID()))
	s.Require().NoError(err)

	// Assertion 0
	s.Run("StatusCode", func() {
		s.StatusNotFound(resp.StatusCode)
	})

	var mapData map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().NoError(json.Unmarshal(raw, &mapData))

	// Assertion 1
	s.Run("DetailNotEmpty", func() {
		s.NotEmpty(mapData["detail"])
	})

	// Assertion 2
	s.Run("ContainsSchema", func() {
		s.Contains(mapData["schemas"], "urn:ietf:params:scim:api:messages:2.0:Error")
	})
}

// Required Test: Make sure random user doesn't exist.
func (s *TestSuite) TestGetUserByRandomUserName() {
	// NOTE: UserNames are always an email in the original Okta Spec Test.
	filter := url.Values{}
	filter.Set("filter", fmt.Sprintf("userName eq \"%s\"", s.RandomEmail()))
	resp, err := s.Get(fmt.Sprintf("/Users?%s", filter.Encode()))
	s.Require().NoError(err)

	// Assertion 0
	s.Run("StatusCode", func() {
		s.StatusOK(resp.StatusCode)
	})

	var mapData map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().NoError(json.Unmarshal(raw, &mapData))

	// Assertion 1
	s.Run("TotalResultsIsNumber0", func() {
		s.Equal(float64(0), mapData["totalResults"])
	})

	// Assertion 2
	s.Run("ContainsSchema", func() {
		s.Contains(mapData["schemas"], "urn:ietf:params:scim:api:messages:2.0:ListResponse")
	})
}

// Required Test: Create Okta user with realistic values.
func (s *TestSuite) TestCreateUser() {
	randomUserName, randomEmail := s.RandomEmail(), s.RandomEmail()
	randomGivenName, randomFamilyName := s.RandomName(), s.RandomName()
	body, err := json.Marshal(map[string]interface{}{
		"schemas":  []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		"userName": randomUserName,
		"name": map[string]interface{}{
			"givenName":  randomGivenName,
			"familyName": randomFamilyName,
		},
		"emails": []map[string]interface{}{
			{
				"primary": true,
				"value":   randomEmail,
				"type":    "work",
			},
		},
		"displayName": fmt.Sprintf("%s %s", randomGivenName, randomFamilyName),
		"active":      true,
	})
	s.Require().NoError(err)
	resp, err := s.Post("/Users", bytes.NewReader(body))
	s.Require().NoError(err)

	// Assertion 0
	s.Run("StatusCode", func() {
		s.StatusCreated(resp.StatusCode)
	})

	var entity map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().NoError(json.Unmarshal(raw, &entity))

	// Assertion 1
	s.Run("ActiveTrue", func() {
		s.Equal(true, entity["active"])
	})

	// Assertion 2
	s.Run("IDNotEmpty", func() {
		s.NotEmpty(entity["id"])
	})

	s.Require().NotEmpty(entity["name"])
	name, ok := entity["name"].(map[string]interface{})
	s.Require().True(ok)

	// Assertion 3
	s.Run("FamilyNameMatches", func() {
		s.Equal(randomFamilyName, name["familyName"])
	})

	// Assertion 4
	s.Run("GivenNameMatches", func() {
		s.Equal(randomGivenName, name["givenName"])
	})

	// Assertion 5
	s.Run("ContainsSchema", func() {
		s.Contains(entity["schemas"], "urn:ietf:params:scim:schemas:core:2.0:User")
	})

	// Assertion 6
	s.Run("UserNameMatches", func() {
		s.Equal(randomUserName, entity["userName"])
	})

	s.Require().NotEmpty(entity["id"])
	id, ok := entity["id"].(string)
	s.Require().True(ok)

	// Next Tests
	s.Run("VerifyCreation", func() {
		s.testVerifyUserCreated(id, randomUserName, randomFamilyName, randomGivenName)
	})

	s.Run("CreateDuplicate", func() {
		s.testCreateDuplicate(body)
	})
}

// Required Test: Verify that user was created.
// NOTE: Gets called at the end of TestCreateUser().
func (s *TestSuite) testVerifyUserCreated(id, userName, familyName, givenName string) {
	resp, err := s.Get(fmt.Sprintf("/Users/%s", id))
	s.Require().NoError(err)

	// Assertion 0
	s.Run("StatusCode", func() {
		s.StatusOK(resp.StatusCode)
	})

	var entity map[string]interface{}
	raw, err := ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	s.Require().NoError(json.Unmarshal(raw, &entity))

	// Assertion 1
	s.Run("UserNameMatches", func() {
		s.Equal(userName, entity["userName"])
	})

	s.Require().NotEmpty(entity["name"])
	name, ok := entity["name"].(map[string]interface{})
	s.Require().True(ok)

	// Assertion 2
	s.Run("FamilyNameMatches", func() {
		s.Equal(familyName, name["familyName"])
	})

	// Assertion 3
	s.Run("GivenNameMatches", func() {
		s.Equal(givenName, name["givenName"])
	})
}

// Required Test: Expect failure when recreating user with same values
// NOTE: Gets called at the end of TestCreateUser().
func (s *TestSuite) testCreateDuplicate(body []byte) {
	resp, err := s.Post("/Users", bytes.NewReader(body))
	s.Require().NoError(err)

	// Assertion 0
	s.Run("StatusCode", func() {
		s.StatusConflict(resp.StatusCode)
	})
}

// Required Test: Username Case Sensitivity Check.
func (s *TestSuite) TestUserNameCS() {
	filter := url.Values{}
	filter.Set("filter", fmt.Sprintf("userName eq \"%s\"", strings.ToUpper(s.RandomEmail())))
	resp, err := s.Get(fmt.Sprintf("/Users?%s", filter.Encode()))
	s.Require().NoError(err)

	// Assertion 0
	s.Run("StatusCode", func() {
		s.StatusOK(resp.StatusCode)
	})
}

// Optional Test: Verify Groups endpoint.
func (s *TestSuite) TestGetGroups() {
	t := time.Now()
	resp, err := s.Get("/Groups")
	d := time.Since(t)
	s.Require().NoError(err)

	// Assertion 0
	s.Run("StatusCode", func() {
		s.StatusOK(resp.StatusCode)
	})

	// Assertion 1
	s.Run("ResponseTime", func() {
		s.LessOrEqual(d.Milliseconds(), int64(600))
	})
}
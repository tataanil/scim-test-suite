package azure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

// SOURCE: https://github.com/AzureAD/SCIMReferenceCode/wiki/Test-Your-SCIM-Endpoint
//		   https://www.postman.com/collections/3b5c4b838ec66cacd53b
// NOTE: in the postman collection, there are requests to "/users", yet url paths are case sensitive.
//		 in this test suite these paths are corrected as defined by the SCIM spec ("/users" -> "/Users").
// 		 -- August 01 2020

func (s *TestSuite) TestEndpoints() {
	s.Run("Get empty Users", func() {
		// NOTE: invalid path in source "/users"
		resp := s.Get("/Users")
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})
	})

	s.Run("Get empty Groups", func() {
		resp := s.Get("/Groups")
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})
	})

	s.Run("Get ResourceTypes", func() {
		resp := s.Get("/ResourceTypes")
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})
		var (
			mapData   = s.ReadAllToMap(resp)
			resources = s.GetSlice("Resources", mapData)
		)
		// NOTE: resourceTypes[0].endpoint == "/Users" (original test, not correct)
		//		 the first resource type is not necessarily the Users resource type.
		s.Run("User endpoint exists", func() {
			var hit bool
			for _, r := range resources {
				var (
					resource = s.IsMap(r)
					endpoint = s.GetString("endpoint", resource)
				)
				if endpoint == "/Users" {
					hit = true
					break
				}
			}
			s.Require().True(hit)
		})
	})

	s.Run("Get ServiceProviderConfig", func() {
		// NOTE: invalid path in source "/serviceConfiguration"
		resp := s.Get("/ServiceProviderConfig")
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			mapData   = s.ReadAllToMap(resp)
			patch     = s.GetMap("patch", mapData)
			supported = s.GetBool("supported", patch)
		)

		s.Run("Patch supported is true", func() { // NOTE: typo in source code: "Pach"
			s.Require().True(supported)
		})
	})

	s.Run("Get Schemas", func() {
		resp := s.Get("/Schemas")
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			mapData   = s.ReadAllToMap(resp)
			resources = s.GetSlice("Resources", mapData)
		)

		s.Run("Body contains User Account", func() { // NOTE: typo in source code: "contians"
			var hit bool
			for _, r := range resources {
				var (
					resource    = s.IsMap(r)
					description = s.GetString("description", resource)
				)
				if description == "User Account" {
					hit = true
					break
				}
			}
			s.Require().True(hit)
		})
	})
}

func (s *TestSuite) TestUsers() {
	var id1, id2 string // saved id for later use

	s.Run("Post User", func() {
		body, err := json.Marshal(s.createUserBody("di-wu", "di-wu"))
		s.Require().NoError(err)
		resp := s.Post("/Users", bytes.NewReader(body))
		s.Run("Status code is 201", func() {
			s.StatusCreated(resp.StatusCode)
		})

		userData := s.ReadAllToMap(resp)
		id1 = s.GetString("id", userData)
	})

	s.Run("Post EnterpriseUser", func() {
		body, err := json.Marshal(s.createEnterpriseUserBody("quint", "quint"))
		s.Require().NoError(err)
		resp := s.Post("/Users", bytes.NewReader(body))
		s.Run("Status code is 201", func() {
			s.StatusCreated(resp.StatusCode)
		})

		userData := s.ReadAllToMap(resp)
		id2 = s.GetString("id", userData)
	})

	s.Run("Get user1", func() {
		resp := s.Get(fmt.Sprintf("/Users/%s", id1))
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			user = s.ReadAllToMap(resp)
			id   = s.GetString("id", user)
		)
		s.Run("Id is requested", func() { // NOTE: typo in source code: "requsted"
			s.Equal(id1, id)
		})
	})

	s.Run("Get user2", func() {
		resp := s.Get(fmt.Sprintf("/Users/%s", id2))
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			user = s.ReadAllToMap(resp)
			id   = s.GetString("id", user)
		)
		s.Run("Id is requested", func() { // NOTE: typo in source code: "requsted"
			s.Equal(id2, id)
		})
	})

	s.Run("Get User Attributes", func() {
		resp := s.Get("/Users?attributes=userName,emails")
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			mapData   = s.ReadAllToMap(resp)
			resources = s.GetSlice("Resources", mapData)
		)

		s.Run("Body contains id1", func() { // NOTE: typo in source code: "contians"
			var hit bool
			for _, r := range resources {
				var (
					resource = s.IsMap(r)
					id       = s.GetString("id", resource)
				)
				if id == id1 {
					hit = true
					break
				}
			}
			s.Require().True(hit)
		})
	})

	s.Run("Get User Filters", func() {
		var (
			// NOTE: typo: "/Users/?filter=DisplayName+eq+%22BobIsAmazing%22"
			filter    = url.Values{"filter": []string{"displayName eq \"di-wu\""}}
			resp      = s.Get(fmt.Sprintf("/Users?%s", filter.Encode()))
			mapData   = s.ReadAllToMap(resp)
			resources = s.GetSlice("Resources", mapData)
		)
		s.Require().NotEmpty(resources)
		var (
			user = s.IsMap(resources[0])
			id   = s.GetString("id", user)
		)
		s.Run("Body contains id1", func() { // NOTE: typo in source code: "contians"
			s.Require().Equal(id1, id)
		})
	})

	s.Run("Patch user1", func() {
		body, err := json.Marshal(map[string]interface{}{
			"schemas": []string{
				"urn:ietf:params:scim:api:messages:2.0:PatchOp",
			},
			"Operations": []map[string]interface{}{
				{
					"op":    "replace",
					"path":  "userName",
					"value": "quint.d",
				},
			},
		})
		s.Require().NoError(err)
		resp := s.Patch(fmt.Sprintf("/Users/%s", id1), bytes.NewReader(body))
		// NOTE: no tests?
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})
	})

	s.Run("Get user1 check Patch", func() {
		resp := s.Get(fmt.Sprintf("/Users/%s", id1))
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			user     = s.ReadAllToMap(resp)
			id       = s.GetString("id", user)
			userName = s.GetString("userName", user)
		)

		s.Run("Id matches", func() {
			s.Require().Equal(id1, id)
		})

		s.Run("Username is changed", func() {
			s.Require().Equal("quint.d", userName)
		})
	})

	s.Run("Replace user2", func() {
		user2 := s.createUserBody("quint", "quint")
		user2["name"].(map[string]interface{})["formatted"] = "NewName"
		body, err := json.Marshal(user2)
		s.Require().NoError(err)
		resp := s.Put(fmt.Sprintf("/Users/%s", id2), bytes.NewReader(body))
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})
	})

	s.Run("Get user2 check Replace", func() {
		resp := s.Get(fmt.Sprintf("/Users/%s", id2))
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			user          = s.ReadAllToMap(resp)
			name          = s.GetMap("name", user)
			formattedName = s.GetString("formatted", name)
		)

		s.Run("Name is updated", func() {
			s.Require().Equal("NewName", formattedName)
		})
	})

	s.Run("Delete user1", func() {
		resp := s.Delete(fmt.Sprintf("/Users/%s", id1))
		s.Run("Status code is 204", func() {
			s.StatusNoContent(resp.StatusCode)
		})
	})

	s.Run("Delete user2", func() {
		resp := s.Delete(fmt.Sprintf("/Users/%s", id2))
		s.Run("Status code is 204", func() {
			s.StatusNoContent(resp.StatusCode)
		})
	})
}

func (s *TestSuite) TestGroups() {
	var userIDs, groupIDs []string

	s.Run("Create empty group", func() {
		body, err := json.Marshal(s.createGroup("Group1"))
		s.Require().NoError(err)
		resp := s.Post("/Groups", bytes.NewReader(body))
		s.Run("Status code is 201", func() {
			s.StatusCreated(resp.StatusCode)
		})

		groupData := s.ReadAllToMap(resp)
		groupIDs = append(groupIDs, s.GetString("id", groupData))
	})

	s.Run("Create users", func() {
		for i := 0; i < 4; i++ {
			body, err := json.Marshal(s.createUserBody(
				fmt.Sprintf("UserName%02d", i),
				fmt.Sprintf("DisplayName%02d", i),
			))
			s.Require().NoError(err)
			resp := s.Post("/Users", bytes.NewReader(body))
			s.Run("Status code is 201", func() {
				s.StatusCreated(resp.StatusCode)
			})

			userData := s.ReadAllToMap(resp)
			userIDs = append(userIDs, s.GetString("id", userData))
		}
	})

	s.Run("Create filled group2", func() {
		body, err := json.Marshal(s.createGroup("Group2", userIDs[2]))
		s.Require().NoError(err)
		resp := s.Post("/Groups", bytes.NewReader(body))
		s.Run("Status code is 201", func() {
			s.StatusCreated(resp.StatusCode)
		})

		groupData := s.ReadAllToMap(resp)
		groupIDs = append(groupIDs, s.GetString("id", groupData))
	})

	s.Run("Get Groups", func() {
		resp := s.Get("/Groups")
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})
	})

	s.Run("Create  group3", func() {
		body, err := json.Marshal(s.createGroup("Group3"))
		s.Require().NoError(err)
		resp := s.Post("/Groups", bytes.NewReader(body))
		s.Run("Status code is 201", func() {
			s.StatusCreated(resp.StatusCode)
		})

		groupData := s.ReadAllToMap(resp)
		groupIDs = append(groupIDs, s.GetString("id", groupData))
	})

	s.Run("Replace group3", func() {
		body, err := json.Marshal(s.createGroup("Group3", userIDs[2], userIDs[3]))
		s.Require().NoError(err)
		resp := s.Put(fmt.Sprintf("/Groups/%s", groupIDs[2]), bytes.NewReader(body))
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})
	})

	s.Run("Validate group3", func() {
		resp := s.Get(fmt.Sprintf("/Groups/%s", groupIDs[2]))
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			group   = s.ReadAllToMap(resp)
			id      = s.GetString("id", group)
			members = s.GetSlice("members", group)
		)

		s.Run("Id matches", func() {
			s.Require().Equal(groupIDs[2], id)
		})

		s.Run("Body contains user id3", func() { // NOTE: typo in source code: "contians"
			var hit bool
			for _, r := range members {
				var (
					resource = s.IsMap(r)
					id       = s.GetString("value", resource)
				)
				if id == userIDs[2] {
					hit = true
					break
				}
			}
			s.Require().True(hit)
		})
	})

	addUser4ToGroup := func() {
		body, err := json.Marshal(map[string]interface{}{
			"schemas": []string{
				"urn:ietf:params:scim:api:messages:2.0:PatchOp",
			},
			"Operations": []map[string]interface{}{
				{
					"op":   "add",
					"path": "members",
					"value": []map[string]interface{}{
						{
							"value": userIDs[3],
						},
					},
				},
			},
		})
		s.Require().NoError(err)
		resp := s.Patch(fmt.Sprintf("/Groups/%s", groupIDs[0]), bytes.NewReader(body))
		// NOTE: no tests?
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})
	}

	s.Run("Patch add user4 to group1", addUser4ToGroup)

	s.Run("Patch remove user4 to group1", func() {
		body, err := json.Marshal(map[string]interface{}{
			"schemas": []string{
				"urn:ietf:params:scim:api:messages:2.0:PatchOp",
			},
			"Operations": []map[string]interface{}{
				{
					"op":   "remove",
					"path": "members", // TODO: fmt.Sprintf("members[value eq \"%s\"]", userIDs[3]),
				},
			},
		})
		s.Require().NoError(err)
		resp := s.Patch(fmt.Sprintf("/Groups/%s", groupIDs[0]), bytes.NewReader(body))
		// NOTE: no tests?
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})
	})

	s.Run("Patch add user4 to group1", addUser4ToGroup)

	s.Run("Get group1 by id", func() {
		resp := s.Get(fmt.Sprintf("/Groups/%s", groupIDs[0]))
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			group   = s.ReadAllToMap(resp)
			id      = s.GetString("id", group)
			members = s.GetSlice("members", group)
		)

		s.Run("Id matches", func() {
			s.Require().Equal(groupIDs[0], id)
		})

		s.Run("Body contains user4", func() { // NOTE: typo in source code: "contians"
			var hit bool
			for _, r := range members {
				var (
					resource = s.IsMap(r)
					value    = s.GetString("value", resource)
				)
				if value == userIDs[3] {
					hit = true
					break
				}
			}
			s.Require().True(hit)
		})
	})

	s.Run("Patch remove all users", func() {
		body, err := json.Marshal(map[string]interface{}{
			"schemas": []string{
				"urn:ietf:params:scim:api:messages:2.0:PatchOp",
			},
			"Operations": []map[string]interface{}{
				{
					"op":   "remove",
					"path": "members",
				},
			},
		})
		s.Require().NoError(err)
		resp := s.Patch(fmt.Sprintf("/Groups/%s", groupIDs[0]), bytes.NewReader(body))
		// NOTE: no tests?
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})
	})

	s.Run("Get group1 by id", func() {
		resp := s.Get(fmt.Sprintf("/Groups/%s", groupIDs[0]))
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			group = s.ReadAllToMap(resp)
			id    = s.GetString("id", group)
		)

		s.Run("Id matches", func() {
			s.Require().Equal(groupIDs[0], id)
		})

		var members []interface{}
		if m, ok := group["members"]; m == nil || !ok {
			// members can be nil
			members = make([]interface{}, 0)
		} else {
			members = s.GetSlice("members", group)
		}

		s.Run("Body contains user4", func() { // NOTE: typo in source code: "contians"
			var hit bool
			for _, r := range members {
				var (
					resource = s.IsMap(r)
					value    = s.GetString("value", resource)
				)
				if value == userIDs[3] {
					hit = true
					break
				}
			}
			s.Require().False(hit)
		})
	})

	s.Run("Delete groups", func() {
		for _, id := range groupIDs {
			resp := s.Delete(fmt.Sprintf("/Groups/%s", id))
			s.Run("Status code is 204", func() {
				s.StatusNoContent(resp.StatusCode)
			})
		}
	})

	s.Run("Delete users", func() {
		for _, id := range userIDs {
			resp := s.Delete(fmt.Sprintf("/Users/%s", id))
			s.Run("Status code is 204", func() {
				s.StatusNoContent(resp.StatusCode)
			})
		}
	})
}

func (s *TestSuite) TestComplexAttributes() {
	var id1, id2 string

	s.Run("Create user1", func() {
		body, err := json.Marshal(s.createUserBody("complex1", "complex1"))
		s.Require().NoError(err)
		resp := s.Post("/Users", bytes.NewReader(body))

		s.Run("Status code is 201", func() {
			s.StatusCreated(resp.StatusCode)
		})

		userData := s.ReadAllToMap(resp)
		id1 = s.GetString("id", userData)
	})

	s.Run("Create user2", func() {
		body, err := json.Marshal(s.createUserBody("complex2", "complex2"))
		s.Require().NoError(err)
		resp := s.Post("/Users", bytes.NewReader(body))

		s.Run("Status code is 201", func() {
			s.StatusCreated(resp.StatusCode)
		})

		userData := s.ReadAllToMap(resp)
		id2 = s.GetString("id", userData)
	})

	s.Run("Get user attributes", func() {
		var (
			filter = url.Values{"attributes": []string{"emails[type eq \"other\"]"}}
			resp   = s.Get(fmt.Sprintf("/Users?%s", filter.Encode()))
		)
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			mapData   = s.ReadAllToMap(resp)
			resources = s.GetSlice("Resources", mapData)
		)

		s.Run("Body contains user1 email", func() { // NOTE: typo in source code: "contians"
			var hit bool
			for _, r := range resources {
				var (
					resource = s.IsMap(r)
					emails   = s.GetSlice("emails", resource)
				)
				for _, e := range emails {
					var (
						email = s.IsMap(e)
						value = s.GetString("value", email)
					)

					if value == "complex1@elimity.com" {
						hit = true
						break
					}
				}
			}
			s.Require().True(hit)
		})
	})

	s.Run("Get user via attribute filter", func() {
		var (
			filter = url.Values{"attributes": []string{"emails[value eq \"complex1@elimity.com\"]"}}
			resp   = s.Get(fmt.Sprintf("/Users?%s", filter.Encode()))
		)
		s.Run("Status code is 200", func() {
			s.StatusOK(resp.StatusCode)
		})

		var (
			mapData   = s.ReadAllToMap(resp)
			resources = s.GetSlice("Resources", mapData)
		)

		s.Run("Body contains user1 email", func() { // NOTE: typo in source code: "contians"
			var hit bool
			for _, r := range resources {
				var (
					resource = s.IsMap(r)
					emails   = s.GetSlice("emails", resource)
				)
				for _, e := range emails {
					var (
						email = s.IsMap(e)
						value = s.GetString("value", email)
					)

					if value == "complex1@elimity.com" {
						hit = true
						break
					}
				}
			}
			s.Require().True(hit)
		})
	})

	s.Run("Delete user1", func() {
		resp := s.Delete(fmt.Sprintf("/Users/%s", id1))
		s.Run("Status code is 204", func() {
			s.StatusNoContent(resp.StatusCode)
		})
	})

	s.Run("Delete user2", func() {
		resp := s.Delete(fmt.Sprintf("/Users/%s", id2))
		s.Run("Status code is 204", func() {
			s.StatusNoContent(resp.StatusCode)
		})
	})
}

// NOTE: not included: test w/ garbage

package test

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/di-wu/scim-filter-parser"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/errors"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
)

func Server() scim.Server {
	return scim.Server{
		Config: scim.ServiceProviderConfig{
			SupportFiltering: true,
		},
		ResourceTypes: []scim.ResourceType{
			{
				ID:          optional.NewString("User"),
				Name:        "User",
				Endpoint:    "/Users",
				Description: optional.NewString("User Account"),
				Schema:      schema.CoreUserSchema(),
				Handler:     newTestResourceHandler(),
			},
			{
				ID:          optional.NewString("EnterpriseUser"),
				Name:        "EnterpriseUser",
				Endpoint:    "/EnterpriseUsers",
				Description: optional.NewString("Enterprise User Account"),
				Schema:      schema.CoreUserSchema(),
				SchemaExtensions: []scim.SchemaExtension{
					{Schema: schema.ExtensionEnterpriseUser()},
				},
				Handler: newTestResourceHandler(),
			},
			{
				ID:          optional.NewString("Group"),
				Name:        "Group",
				Endpoint:    "/Groups",
				Description: optional.NewString("Group"),
				Schema:      schema.CoreGroupSchema(),
				Handler:     newTestResourceHandler(),
			},
		},
	}
}

func newTestResourceHandler() scim.ResourceHandler {
	data := make(map[string]Data)

	// Generate enough test data to test pagination
	for i := 1; i < 21; i++ {
		data[fmt.Sprintf("%04d", i)] = Data{
			resourceAttributes: scim.ResourceAttributes{
				"userName":   fmt.Sprintf("test%02d", i),
				"externalId": fmt.Sprintf("external%02d", i),
				"name": map[string]interface{}{
					"familyName": fmt.Sprintf("familyName%02d", i),
					"givenName":  fmt.Sprintf("givenName%02d", i),
				},
				"active": true,
				"emails": []map[string]interface{}{
					{
						"value": fmt.Sprintf("%02d@example.com", i),
					},
				},
			},
			meta: map[string]string{
				"created":      fmt.Sprintf("2020-01-%02dT15:04:05+07:00", i),
				"lastModified": fmt.Sprintf("2020-02-%02dT16:05:04+07:00", i),
				"version":      fmt.Sprintf("v%09d", i),
			},
		}
	}

	return ResourceHandler{
		data: data,
	}
}

type Data struct {
	resourceAttributes scim.ResourceAttributes
	meta               map[string]string
}

// simple in-memory resource database
type ResourceHandler struct {
	data map[string]Data
}

func (h ResourceHandler) Create(r *http.Request, attributes scim.ResourceAttributes) (scim.Resource, error) {
	// create unique identifier
	rand.Seed(time.Now().UnixNano())
	id := fmt.Sprintf("%04d", rand.Intn(9999))

	for _, entity := range h.data {
		if entity.resourceAttributes["userName"] == attributes["userName"] {
			return scim.Resource{}, errors.ScimErrorUniqueness
		}
	}

	// store resource
	h.data[id] = Data{
		resourceAttributes: attributes,
	}

	now := time.Now()

	// return stored resource
	return scim.Resource{
		ID:         id,
		ExternalID: h.externalID(attributes),
		Attributes: attributes,
		Meta: scim.Meta{
			Created:      &now,
			LastModified: &now,
			Version:      fmt.Sprintf("v%s", id),
		},
	}, nil
}

func (h ResourceHandler) Get(r *http.Request, id string) (scim.Resource, error) {
	// check if resource exists
	data, ok := h.data[id]
	if !ok {
		return scim.Resource{}, errors.ScimErrorResourceNotFound(id)
	}

	created, _ := time.ParseInLocation(time.RFC3339, fmt.Sprintf("%v", data.meta["created"]), time.UTC)
	lastModified, _ := time.Parse(time.RFC3339, fmt.Sprintf("%v", data.meta["lastModified"]))

	// return resource with given identifier
	return scim.Resource{
		ID:         id,
		ExternalID: h.externalID(data.resourceAttributes),
		Attributes: data.resourceAttributes,
		Meta: scim.Meta{
			Created:      &created,
			LastModified: &lastModified,
			Version:      fmt.Sprintf("%v", data.meta["version"]),
		},
	}, nil
}

func (h ResourceHandler) GetAll(r *http.Request, params scim.ListRequestParams) (scim.Page, error) {
	data := h.data
	if params.Filter != nil {
		filteredData := make(map[string]Data)
		for k, v := range h.data {
			// Only 'x eq userName' is implemented
			switch params.Filter.(type) {
			case filter.AttributeExpression:
				attrExp := params.Filter.(filter.AttributeExpression)
				if attrExp.CompareOperator == filter.EQ {
					switch name := attrExp.AttributePath.AttributeName; name {
					case "userName":
						if name == h.userName(v.resourceAttributes) {
							filteredData[k] = v
						}
					}
				}
			}
		}
		data = filteredData
	}

	i := 1
	resources := make([]scim.Resource, 0)
	for k, v := range data {
		if i > (params.StartIndex + params.Count - 1) {
			break
		}

		if i >= params.StartIndex {
			resources = append(resources, scim.Resource{
				ID:         k,
				ExternalID: h.externalID(v.resourceAttributes),
				Attributes: v.resourceAttributes,
			})
		}
		i++
	}

	return scim.Page{
		TotalResults: len(data),
		Resources:    resources,
	}, nil
}

func (h ResourceHandler) Replace(r *http.Request, id string, attributes scim.ResourceAttributes) (scim.Resource, error) {
	// check if resource exists
	_, ok := h.data[id]
	if !ok {
		return scim.Resource{}, errors.ScimErrorResourceNotFound(id)
	}

	// replace (all) attributes
	h.data[id] = Data{
		resourceAttributes: attributes,
	}

	// return resource with replaced attributes
	return scim.Resource{
		ID:         id,
		ExternalID: h.externalID(attributes),
		Attributes: attributes,
	}, nil
}

func (h ResourceHandler) Delete(r *http.Request, id string) error {
	// check if resource exists
	_, ok := h.data[id]
	if !ok {
		return errors.ScimErrorResourceNotFound(id)
	}

	// delete resource
	delete(h.data, id)

	return nil
}

func (h ResourceHandler) Patch(r *http.Request, id string, req scim.PatchRequest) (scim.Resource, error) {
	for _, op := range req.Operations {
		switch op.Op {
		case scim.PatchOperationAdd:
			if op.Path != "" {
				h.data[id].resourceAttributes[op.Path] = op.Value
			} else {
				valueMap := op.Value.(map[string]interface{})
				for k, v := range valueMap {
					if arr, ok := h.data[id].resourceAttributes[k].([]interface{}); ok {
						arr = append(arr, v)
						h.data[id].resourceAttributes[k] = arr
					} else {
						h.data[id].resourceAttributes[k] = v
					}
				}
			}
		case scim.PatchOperationReplace:
			if op.Path != "" {
				h.data[id].resourceAttributes[op.Path] = op.Value
			} else {
				valueMap := op.Value.(map[string]interface{})
				for k, v := range valueMap {
					h.data[id].resourceAttributes[k] = v
				}
			}
		case scim.PatchOperationRemove:
			h.data[id].resourceAttributes[op.Path] = nil
		}
	}

	created, _ := time.ParseInLocation(time.RFC3339, fmt.Sprintf("%v", h.data[id].meta["created"]), time.UTC)
	now := time.Now()

	// return resource with replaced attributes
	return scim.Resource{
		ID:         id,
		ExternalID: h.externalID(h.data[id].resourceAttributes),
		Attributes: h.data[id].resourceAttributes,
		Meta: scim.Meta{
			Created:      &created,
			LastModified: &now,
			Version:      fmt.Sprintf("%s.patch", h.data[id].meta["version"]),
		},
	}, nil
}

func (h ResourceHandler) externalID(attributes scim.ResourceAttributes) optional.String {
	if eID, ok := attributes["externalId"]; ok {
		externalID, ok := eID.(string)
		if !ok {
			return optional.String{}
		}
		return optional.NewString(externalID)
	}

	return optional.String{}
}

func (h ResourceHandler) userName(attributes scim.ResourceAttributes) string {
	if eID, ok := attributes["userName"]; ok {
		if userName, ok := eID.(string); ok {
			return userName
		}
	}

	return ""
}

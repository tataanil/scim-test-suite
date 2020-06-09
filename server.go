package suite

import (
	"fmt"
	filter "github.com/di-wu/scim-filter-parser"
	"github.com/elimity-com/scim/errors"
	"github.com/elimity-com/scim/schema"
	"math/rand"
	"net/http"
	"time"

	. "github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
)

func TestServer() Server {
	return Server{
		Config: ServiceProviderConfig{
			SupportFiltering: true,
		},
		ResourceTypes: []ResourceType{
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
				SchemaExtensions: []SchemaExtension{
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

func newTestResourceHandler() ResourceHandler {
	data := make(map[string]TestData)

	// Generate enough test data to test pagination
	for i := 1; i < 21; i++ {
		data[fmt.Sprintf("000%d", i)] = TestData{
			resourceAttributes: ResourceAttributes{
				"userName":   fmt.Sprintf("test%d", i),
				"externalId": fmt.Sprintf("external%d", i),
				"name": map[string]interface{}{
					"familyName": fmt.Sprintf("familyName%d", i),
					"givenName":  fmt.Sprintf("givenName%d", i),
				},
				"active": true,
				"emails": []map[string]interface{}{
					{
						"value": fmt.Sprintf("%d@example.com", i),
					},
				},
			},
			meta: map[string]string{
				"created":      fmt.Sprintf("2020-01-%02dT15:04:05+07:00", i),
				"lastModified": fmt.Sprintf("2020-02-%02dT16:05:04+07:00", i),
				"version":      fmt.Sprintf("v%d", i),
			},
		}
	}

	return TestResourceHandler{
		data: data,
	}
}

type TestData struct {
	resourceAttributes ResourceAttributes
	meta               map[string]string
}

// simple in-memory resource database
type TestResourceHandler struct {
	data map[string]TestData
}

func (h TestResourceHandler) Create(r *http.Request, attributes ResourceAttributes) (Resource, error) {
	// create unique identifier
	rand.Seed(time.Now().UnixNano())
	id := fmt.Sprintf("%04d", rand.Intn(9999))

	for _, entity := range h.data {
		if entity.resourceAttributes["userName"] == attributes["userName"] {
			return Resource{}, errors.ScimErrorUniqueness
		}
	}

	// store resource
	h.data[id] = TestData{
		resourceAttributes: attributes,
	}

	now := time.Now()

	// return stored resource
	return Resource{
		ID:         id,
		ExternalID: h.externalID(attributes),
		Attributes: attributes,
		Meta: Meta{
			Created:      &now,
			LastModified: &now,
			Version:      fmt.Sprintf("v%s", id),
		},
	}, nil
}

func (h TestResourceHandler) Get(r *http.Request, id string) (Resource, error) {
	// check if resource exists
	data, ok := h.data[id]
	if !ok {
		return Resource{}, errors.ScimErrorResourceNotFound(id)
	}

	created, _ := time.ParseInLocation(time.RFC3339, fmt.Sprintf("%v", data.meta["created"]), time.UTC)
	lastModified, _ := time.Parse(time.RFC3339, fmt.Sprintf("%v", data.meta["lastModified"]))

	// return resource with given identifier
	return Resource{
		ID:         id,
		ExternalID: h.externalID(data.resourceAttributes),
		Attributes: data.resourceAttributes,
		Meta: Meta{
			Created:      &created,
			LastModified: &lastModified,
			Version:      fmt.Sprintf("%v", data.meta["version"]),
		},
	}, nil
}

func (h TestResourceHandler) GetAll(r *http.Request, params ListRequestParams) (Page, error) {
	data := h.data
	if params.Filter != nil {
		filteredData := make( map[string]TestData)
		for k, v := range h.data {
			// Only 'x eq userName' is implemented
			switch params.Filter.(type) {
			case filter.AttributeExpression:
				attrExp := params.Filter.(filter.AttributeExpression)
				if attrExp.CompareOperator == filter.EQ {
					switch name := attrExp.AttributePath.AttributeName; name {
					case "userName":
						if  name == h.userName(v.resourceAttributes) {
							filteredData[k] = v
						}
					}
				}
			}
		}
		data = filteredData
	}

	i := 1
	resources := make([]Resource, 0)
	for k, v := range data {
		if i > (params.StartIndex + params.Count - 1) {
			break
		}

		if i >= params.StartIndex {
			resources = append(resources, Resource{
				ID:         k,
				ExternalID: h.externalID(v.resourceAttributes),
				Attributes: v.resourceAttributes,
			})
		}
		i++
	}

	return Page{
		TotalResults: len(data),
		Resources:    resources,
	}, nil
}

func (h TestResourceHandler) Replace(r *http.Request, id string, attributes ResourceAttributes) (Resource, error) {
	// check if resource exists
	_, ok := h.data[id]
	if !ok {
		return Resource{}, errors.ScimErrorResourceNotFound(id)
	}

	// replace (all) attributes
	h.data[id] = TestData{
		resourceAttributes: attributes,
	}

	// return resource with replaced attributes
	return Resource{
		ID:         id,
		ExternalID: h.externalID(attributes),
		Attributes: attributes,
	}, nil
}

func (h TestResourceHandler) Delete(r *http.Request, id string) error {
	// check if resource exists
	_, ok := h.data[id]
	if !ok {
		return errors.ScimErrorResourceNotFound(id)
	}

	// delete resource
	delete(h.data, id)

	return nil
}

func (h TestResourceHandler) Patch(r *http.Request, id string, req PatchRequest) (Resource, error) {
	for _, op := range req.Operations {
		switch op.Op {
		case PatchOperationAdd:
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
		case PatchOperationReplace:
			if op.Path != "" {
				h.data[id].resourceAttributes[op.Path] = op.Value
			} else {
				valueMap := op.Value.(map[string]interface{})
				for k, v := range valueMap {
					h.data[id].resourceAttributes[k] = v
				}
			}
		case PatchOperationRemove:
			h.data[id].resourceAttributes[op.Path] = nil
		}
	}

	created, _ := time.ParseInLocation(time.RFC3339, fmt.Sprintf("%v", h.data[id].meta["created"]), time.UTC)
	now := time.Now()

	// return resource with replaced attributes
	return Resource{
		ID:         id,
		ExternalID: h.externalID(h.data[id].resourceAttributes),
		Attributes: h.data[id].resourceAttributes,
		Meta: Meta{
			Created:      &created,
			LastModified: &now,
			Version:      fmt.Sprintf("%s.patch", h.data[id].meta["version"]),
		},
	}, nil
}

func (h TestResourceHandler) externalID(attributes ResourceAttributes) optional.String {
	if eID, ok := attributes["externalId"]; ok {
		externalID, ok := eID.(string)
		if !ok {
			return optional.String{}
		}
		return optional.NewString(externalID)
	}

	return optional.String{}
}

func (h TestResourceHandler) userName(attributes ResourceAttributes) string {
	if eID, ok := attributes["userName"]; ok {
		if userName, ok := eID.(string); ok {
			return userName
		}
	}

	return ""
}

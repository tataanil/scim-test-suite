package suite

import (
	"fmt"
	"github.com/elimity-com/scim/errors"
	"math/rand"
	"net/http"
	"time"

	. "github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
)

func TestServer() Server {
	userSchema := getUserSchema()
	userSchemaExtension := getUserExtensionSchema()
	return Server{
		Config: ServiceProviderConfig{},
		ResourceTypes: []ResourceType{
			{
				ID:          optional.NewString("User"),
				Name:        "User",
				Endpoint:    "/Users",
				Description: optional.NewString("User Account"),
				Schema:      userSchema,
				Handler:     newTestResourceHandler(),
			},
			{
				ID:          optional.NewString("EnterpriseUser"),
				Name:        "EnterpriseUser",
				Endpoint:    "/EnterpriseUsers",
				Description: optional.NewString("Enterprise User Account"),
				Schema:      userSchema,
				SchemaExtensions: []SchemaExtension{
					{Schema: userSchemaExtension},
				},
				Handler: newTestResourceHandler(),
			},
		},
	}
}

func getUserSchema() schema.Schema {
	return schema.Schema{
		ID:          "urn:ietf:params:scim:schemas:core:2.0:User",
		Name:        optional.NewString("User"),
		Description: optional.NewString("User Account"),
		Attributes: []schema.CoreAttribute{
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "userName",
				Required:   true,
				Uniqueness: schema.AttributeUniquenessServer(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleBooleanParams(schema.BooleanParams{
				Name:     "active",
				Required: false,
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "readonlyThing",
				Required:   false,
				Mutability: schema.AttributeMutabilityReadOnly(),
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "immutableThing",
				Required:   false,
				Mutability: schema.AttributeMutabilityImmutable(),
			})),
			schema.ComplexCoreAttribute(schema.ComplexParams{
				Name:     "Name",
				Required: false,
				SubAttributes: []schema.SimpleParams{
					schema.SimpleStringParams(schema.StringParams{
						Name: "familyName",
					}),
					schema.SimpleStringParams(schema.StringParams{
						Name: "givenName",
					}),
				},
			}),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name: "displayName",
			})),
			schema.ComplexCoreAttribute(schema.ComplexParams{
				Name:        "emails",
				MultiValued: true,
				SubAttributes: []schema.SimpleParams{
					schema.SimpleStringParams(schema.StringParams{
						Name: "value",
					}),
					schema.SimpleStringParams(schema.StringParams{
						Name: "display",
					}),
					schema.SimpleStringParams(schema.StringParams{
						Name: "type",
						CanonicalValues: []string{
							"work", "home", "other",
						},
					}),
					schema.SimpleBooleanParams(schema.BooleanParams{
						Name: "primary",
					}),
				},
			}),
		},
	}
}

func getUserExtensionSchema() schema.Schema {
	return schema.Schema{
		ID:          "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
		Name:        optional.NewString("EnterpriseUser"),
		Description: optional.NewString("Enterprise User"),
		Attributes: []schema.CoreAttribute{
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name: "employeeNumber",
			})),
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name: "organization",
			})),
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
	resources := make([]Resource, 0)
	i := 1

	for k, v := range h.data {
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
		TotalResults: len(h.data),
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

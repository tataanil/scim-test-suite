package schema

import (
	"encoding/json"
	"fmt"

	"github.com/elimity-com/scim/optional"
	. "github.com/elimity-com/scim/schema"
)

// ParseJSONSchema converts raw json data into a SCIM Schema.
// RFC: https://tools.ietf.org/html/rfc7643#section-7
func ParseJSONSchema(raw []byte) (Schema, error) {
	var jsonSchema map[string]interface{}
	if err := json.Unmarshal(raw, &jsonSchema); err != nil {
		return Schema{}, err
	}

	var schema Schema
	var jsonAttributes []interface{}
	for k, v := range jsonSchema {
		switch k {
		case "id":
			id, ok := v.(string)
			if !ok {
				return Schema{}, fmt.Errorf("id is not a string")
			}
			schema.ID = id
		case "name":
			name, ok := v.(string)
			if !ok {
				return Schema{}, fmt.Errorf("name is not a string")
			}
			schema.Name = optional.NewString(name)
		case "description":
			desc, ok := v.(string)
			if !ok {
				return Schema{}, fmt.Errorf("name is not a string")
			}
			schema.Description = optional.NewString(desc)
		case "attributes":
			attrs, ok := v.([]interface{})
			if !ok {
				return Schema{}, fmt.Errorf("attributes is not an array")
			}
			jsonAttributes = attrs
		}
	}

	if schema.ID == "" {
		return Schema{}, fmt.Errorf("id is empty")
	}

	schemaAttributes, err := parseAttributes(jsonAttributes)
	if err != nil {
		return Schema{}, err
	}
	for _, attribute := range schemaAttributes {
		switch attribute.typ {
		case "complex":
			if len(attribute.subAttributes) == 0 {
				return Schema{}, fmt.Errorf("complex attributes should have sub attributes")
			}
			var subAttributes []SimpleParams

			schema.Attributes = append(schema.Attributes, ComplexCoreAttribute(ComplexParams{
				Description:   attribute.description,
				MultiValued:   attribute.multiValued,
				Mutability:    attribute.mutability,
				Name:          attribute.name,
				Required:      attribute.required,
				Returned:      attribute.returned,
				SubAttributes: subAttributes,
				Uniqueness:    attribute.uniqueness,
			}))
		default:
			simple, err := attribute.simple()
			if err != nil {
				return Schema{}, err
			}
			schema.Attributes = append(schema.Attributes, SimpleCoreAttribute(simple))
		}
	}

	return schema, nil
}

type attribute struct {
	name, typ                        string
	description                      optional.String
	multiValued, required, caseExact bool
	canonicalValues                  []string
	mutability                       AttributeMutability
	returned                         AttributeReturned
	uniqueness                       AttributeUniqueness
	subAttributes                    []attribute
	referenceTypes                   []AttributeReferenceType
}

func (a attribute) simple() (SimpleParams, error) {
	switch a.typ {
	case "string":
		return SimpleStringParams(StringParams{
			CanonicalValues: a.canonicalValues,
			CaseExact:       a.caseExact,
			Description:     a.description,
			MultiValued:     a.multiValued,
			Mutability:      a.mutability,
			Name:            a.name,
			Required:        a.required,
			Returned:        a.returned,
			Uniqueness:      a.uniqueness,
		}), nil
	case "boolean":
		return SimpleBooleanParams(BooleanParams{
			Description: a.description,
			MultiValued: a.multiValued,
			Mutability:  a.mutability,
			Name:        a.name,
			Required:    a.required,
			Returned:    a.returned,
		}), nil
	case "decimal", "integer":
		var numberType AttributeDataType
		if a.typ == "decimal" {
			numberType = AttributeTypeDecimal()
		} else {
			numberType = AttributeTypeInteger()
		}

		return SimpleNumberParams(NumberParams{
			Description: a.description,
			MultiValued: a.multiValued,
			Mutability:  a.mutability,
			Name:        a.name,
			Required:    a.required,
			Returned:    a.returned,
			Type:        numberType,
			Uniqueness:  a.uniqueness,
		}), nil
	case "dateTime":
		return SimpleDateTimeParams(DateTimeParams{
			Description: a.description,
			MultiValued: a.multiValued,
			Mutability:  a.mutability,
			Name:        a.name,
			Required:    a.required,
			Returned:    a.returned,
		}), nil
	case "binary":
		return SimpleBinaryParams(BinaryParams{
			Description: a.description,
			MultiValued: a.multiValued,
			Mutability:  a.mutability,
			Name:        a.name,
			Required:    a.required,
			Returned:    a.returned,
		}), nil
	case "reference":
		return SimpleReferenceParams(ReferenceParams{
			Description:    a.description,
			MultiValued:    a.multiValued,
			Mutability:     a.mutability,
			Name:           a.name,
			ReferenceTypes: a.referenceTypes,
			Required:       a.required,
			Returned:       a.returned,
			Uniqueness:     a.uniqueness,
		}), nil
	default:
		return SimpleParams{}, fmt.Errorf("invalid attribute type: %s", a.typ)
	}
}

func parseAttributes(attributes []interface{}) ([]attribute, error) {
	var schemaAttributes []attribute
	for _, a := range attributes {
		jsonAttribute, ok := a.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("attribute is not an object")
		}

		var attribute attribute
		for k, v := range jsonAttribute {
			switch k {
			case "name":
				var ok bool
				attribute.name, ok = v.(string)
				if !ok {
					return nil, fmt.Errorf("name is not a string")
				}
			case "type":
				var ok bool
				attribute.typ, ok = v.(string)
				if !ok {
					return nil, fmt.Errorf("type is not a string")
				}
			case "subAttributes":
				jsonSubAttribute, ok := v.([]interface{})
				if !ok {
					return nil, fmt.Errorf("sub attribute is not an object")
				}
				var err error
				attribute.subAttributes, err = parseAttributes(jsonSubAttribute)
				if err != nil {
					return nil, err
				}
			case "multiValued":
				var ok bool
				attribute.multiValued, ok = v.(bool)
				if !ok {
					return nil, fmt.Errorf("multi valued is not a boolean")
				}
			case "description":
				desc, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("description is not a string")
				}
				attribute.description = optional.NewString(desc)
			case "required":
				var ok bool
				attribute.required, ok = v.(bool)
				if !ok {
					return nil, fmt.Errorf("required is not a boolean")
				}
			case "canonicalValues":
				cv, ok := v.([]interface{})
				if !ok {
					return nil, fmt.Errorf("canonical values is not an array of strings")
				}
				for _, s := range cv {
					vs, ok := s.(string)
					if !ok {
						return nil, fmt.Errorf("canonical value is not a string")
					}
					attribute.canonicalValues = append(attribute.canonicalValues, vs)
				}
			case "caseExact":
				var ok bool
				attribute.caseExact, ok = v.(bool)
				if !ok {
					return nil, fmt.Errorf("case exact is not a boolean")
				}
			case "mutability":
				mut, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("mutability is not a string")
				}
				switch mut {
				case "readOnly":
					attribute.mutability = AttributeMutabilityReadOnly()
				case "readWrite":
					attribute.mutability = AttributeMutabilityReadWrite()
				case "immutable":
					attribute.mutability = AttributeMutabilityImmutable()
				case "writeOnly":
					attribute.mutability = AttributeMutabilityWriteOnly()
				default:
					return nil, fmt.Errorf("invalid mutability type: %s", mut)
				}
			case "returned":
				ret, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("returned is not a string")
				}
				switch ret {
				case "always":
					attribute.returned = AttributeReturnedAlways()
				case "never":
					attribute.returned = AttributeReturnedNever()
				case "default":
					attribute.returned = AttributeReturnedDefault()
				case "request":
					attribute.returned = AttributeReturnedRequest()
				default:
					return nil, fmt.Errorf("invalid returned type: %s", ret)
				}
			case "uniqueness":
				uni, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("uniqueness is not a string")
				}
				switch uni {
				case "none":
					attribute.uniqueness = AttributeUniquenessNone()
				case "server":
					attribute.uniqueness = AttributeUniquenessServer()
				case "global":
					attribute.uniqueness = AttributeUniquenessGlobal()
				default:
					return nil, fmt.Errorf("invalid uniqueness type: %s", uni)
				}
			case "referenceTypes":
				rt, ok := v.([]interface{})
				if !ok {
					return nil, fmt.Errorf("reference types is not an array of strings")
				}
				for _, s := range rt {
					vs, ok := s.(string)
					if !ok {
						return nil, fmt.Errorf("reference type is not a string")
					}
					switch vs {
					case "external":
						attribute.referenceTypes = append(attribute.referenceTypes, AttributeReferenceTypeExternal)
					case "uri":
						attribute.referenceTypes = append(attribute.referenceTypes, AttributeReferenceTypeURI)
					default:
						attribute.referenceTypes = append(attribute.referenceTypes, AttributeReferenceType(vs))
					}
				}
			}
		}
	}
	return schemaAttributes, nil
}

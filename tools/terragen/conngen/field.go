package conngen

import (
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

const (
	FieldTypeString       = "string"
	FieldTypeSecret       = "secret"
	FieldTypeBool         = "boolean"
	FieldTypeNumber       = "number"
	FieldTypeHTTPAuth     = "httpAuth"
	FieldTypeObject       = "object"
	FieldTypeAuditFilters = "auditFilters"
)

// Generated

var UseStaticIPsField = &Field{
	Name:        "useStaticIps",
	Description: "Whether the connector should send all requests from specific static IPs.",
	Type:        FieldTypeBool,
}

// Field

type Field struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Type        string           `json:"type"`
	Required    bool             `json:"required"`
	Dynamic     bool             `json:"dynamic"`
	Initial     any              `json:"initialValue"`
	Dependency  *FieldDependency `json:"dependsOn"`

	naming *Naming
}

func (f *Field) StructName() string {
	return f.naming.GetName("field", f.Name, "struct", f.defaultStructName())
}

func (f *Field) defaultStructName() string {
	return utils.CapitalCase(f.Name)
}

func (f *Field) StructType() string {
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return `types.String`
	case FieldTypeBool:
		return `types.Bool`
	case FieldTypeNumber:
		return `types.Float64`
	case FieldTypeObject:
		return `map[string]string`
	case FieldTypeAuditFilters:
		return `[]*AuditFilterFieldModel`
	case FieldTypeHTTPAuth:
		return `*HTTPAuthFieldModel`
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) AttributeName() string {
	return f.naming.GetName("field", f.Name, "attribute", f.defaultAttributeName())
}

func (f *Field) defaultAttributeName() string {
	return utils.SnakeCase(f.Name)
}

func (f *Field) AttributeType() string {
	switch f.Type {
	case FieldTypeString:
		if f.Required {
			return `stringattr.Required()`
		}
		if v, ok := f.Initial.(string); ok {
			return fmt.Sprintf(`stringattr.Default(%q)`, v)
		}
		return `stringattr.Default("")`
	case FieldTypeSecret:
		if f.Required {
			return `stringattr.SecretRequired()`
		}
		return `stringattr.SecretOptional()`
	case FieldTypeBool:
		if f.Required {
			return `boolattr.Required()`
		}
		if f.Initial == true {
			return `boolattr.Default(true)`
		}
		return `boolattr.Default(false)`
	case FieldTypeNumber:
		if f.Required {
			return `floatattr.Required()`
		}
		if v, ok := f.Initial.(float64); ok {
			return fmt.Sprintf(`floatattr.Default(%g)`, v)
		}
		return `floatattr.Default(0)`
	case FieldTypeObject:
		return `mapattr.StringOptional()`
	case FieldTypeAuditFilters:
		return `listattr.Optional(AuditFilterFieldAttributes)`
	case FieldTypeHTTPAuth:
		if f.Required {
			return `objectattr.Required(HTTPAuthFieldAttributes, HTTPAuthFieldValidator)`
		}
		return `objectattr.Optional(HTTPAuthFieldAttributes, HTTPAuthFieldValidator)`
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) GetValueStatement() string {
	accessor := fmt.Sprintf(`m.%s`, f.StructName())
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return fmt.Sprintf(`stringattr.Get(%s, c, %q)`, accessor, f.Name)
	case FieldTypeBool:
		return fmt.Sprintf(`boolattr.Get(%s, c, %q)`, accessor, f.Name)
	case FieldTypeNumber:
		return fmt.Sprintf(`floatattr.Get(%s, c, %q)`, accessor, f.Name)
	case FieldTypeObject:
		return fmt.Sprintf(`c[%q] = %s`, f.Name, accessor)
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`listattr.Get(%s, c, %q, h)`, accessor, f.Name)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`objectattr.Get(%s, c, %q, h)`, accessor, f.Name)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) ValidateNonZero() string {
	accessor := fmt.Sprintf(`m.%s`, f.StructName())
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		initial, _ := f.Initial.(string)
		return fmt.Sprintf(`%s.ValueString() != %q`, accessor, initial)
	case FieldTypeBool:
		operator := ""
		if f.Initial == true {
			operator = "!"
		}
		return fmt.Sprintf(`%s%s.ValueBool()`, operator, accessor)
	case FieldTypeNumber:
		initial, _ := f.Initial.(float64)
		return fmt.Sprintf(`%s.ValueFloat64() != %g`, accessor, initial)
	case FieldTypeObject:
		return fmt.Sprintf(`len(%s) != 0`, accessor)
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`len(%s) != 0`, accessor)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`%s != nil`, accessor)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) IsNotNull() string {
	accessor := fmt.Sprintf(`m.%s`, f.StructName())
	switch f.Type {
	case FieldTypeString, FieldTypeSecret, FieldTypeBool, FieldTypeNumber:
		return fmt.Sprintf(`!%s.IsNull()`, accessor)
	case FieldTypeObject, FieldTypeAuditFilters:
		return fmt.Sprintf(`len(%s) != 0`, accessor)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`%s != nil`, accessor)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

// Tests

func (f *Field) GetTestAssignment() string {
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return fmt.Sprintf(`%q`, f.TestString())
	case FieldTypeBool:
		return `true`
	case FieldTypeNumber:
		return fmt.Sprintf(`%d`, f.TestNumber())
	case FieldTypeObject:
		return fmt.Sprintf(`{
    							"key" = %q
    						}`, f.TestString())
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`[{ key = "actions", operator = "includes", values = [%q] }]`, f.TestString())
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`{
    							bearer_token = %q
    						}`, f.TestString())
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) GetTestCheck(list string, index int) string {
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return fmt.Sprintf(`"connectors.%s.%d.%s": %q`, list, index, f.AttributeName(), f.TestString())
	case FieldTypeBool:
		return fmt.Sprintf(`"connectors.%s.%d.%s": true`, list, index, f.AttributeName())
	case FieldTypeNumber:
		return fmt.Sprintf(`"connectors.%s.%d.%s": %d`, list, index, f.AttributeName(), f.TestNumber())
	case FieldTypeObject:
		return fmt.Sprintf(`"connectors.%s.%d.%s.key": %q`, list, index, f.AttributeName(), f.TestString())
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`"connectors.%s.%d.%s.0.values": []string{%q}`, list, index, f.AttributeName(), f.TestString())
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`"connectors.%s.%d.%s.bearer_token": %q`, list, index, f.AttributeName(), f.TestString())
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) TestString() string {
	b := sha256.Sum256([]byte(f.Name))
	s := base32.StdEncoding.EncodeToString(b[:])
	return strings.ToLower(s[:min(len(s), len(f.Name))])
}

func (f *Field) TestNumber() int {
	return len(f.Name)
}

// Dependency

type FieldDependency struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
	*Field
}

func (d *FieldDependency) DefaultValue() bool {
	v, _ := d.Field.Initial.(bool)
	return v
}

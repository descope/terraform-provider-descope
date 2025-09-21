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
	Hidden      bool             `json:"hidden"`
	Options     []*FieldOption   `json:"options"`
	Dependency  *FieldDependency `json:"dependsOn"`

	naming *Naming
}

func (f *Field) StructName() string {
	return f.naming.GetName("field", f.Name, "struct", f.defaultStructName())
}

func (f *Field) defaultStructName() string {
	return utils.CapitalCase(f.Name)
}

func (f *Field) OptionValues() []string {
	values := []string{}
	for _, option := range f.Options {
		values = append(values, option.Value)
	}
	return values
}

func (f *Field) StructType() string {
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return `stringattr.Type`
	case FieldTypeBool:
		return `boolattr.Type`
	case FieldTypeNumber:
		return `floatattr.Type`
	case FieldTypeObject:
		return `strmapattr.Type`
	case FieldTypeAuditFilters:
		return `listattr.Type[AuditFilterFieldModel]`
	case FieldTypeHTTPAuth:
		return `objattr.Type[HTTPAuthFieldModel]`
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
		validator := ""

		if len(f.Options) > 0 {
			values := []string{}
			if !f.Required {
				values = append(values, `""`)
			}
			for _, option := range f.Options {
				values = append(values, fmt.Sprintf("%q", option.Value))
			}
			validator = fmt.Sprintf("stringvalidator.OneOf(%s)", strings.Join(values, ", "))

			if v, ok := f.Initial.(string); ok {
				return fmt.Sprintf(`stringattr.Default(%q, %s)`, v, validator)
			}
			if f.Required {
				return fmt.Sprintf(`stringattr.Required(%s)`, validator)
			}
			return fmt.Sprintf(`stringattr.Default("", %s)`, validator)
		}

		if f.Required && f.Dependency == nil {
			return fmt.Sprintf(`stringattr.Required(%s)`, validator)
		}

		if validator != "" {
			validator = ", " + validator
		}

		defValue := ""
		if v, ok := f.Initial.(string); ok {
			defValue = v
		}

		return fmt.Sprintf(`stringattr.Default(%q, %s)`, defValue, validator)
	case FieldTypeSecret:
		if f.Required && f.Dependency == nil {
			return `stringattr.SecretRequired()`
		}
		return `stringattr.SecretOptional()`
	case FieldTypeBool:
		if f.Required && f.Dependency == nil {
			return `boolattr.Required()`
		}
		if f.Initial == true {
			return `boolattr.Default(true)`
		}
		return `boolattr.Default(false)`
	case FieldTypeNumber:
		if f.Required && f.Dependency == nil {
			return `floatattr.Required()`
		}
		if v, ok := f.Initial.(float64); ok {
			return fmt.Sprintf(`floatattr.Default(%g)`, v)
		}
		return `floatattr.Default(0)`
	case FieldTypeObject:
		return `strmapattr.Default()`
	case FieldTypeAuditFilters:
		return `listattr.Default[AuditFilterFieldModel](AuditFilterFieldAttributes)`
	case FieldTypeHTTPAuth:
		if f.Required && f.Dependency == nil {
			return `objattr.Required[HTTPAuthFieldModel](HTTPAuthFieldAttributes, HTTPAuthFieldValidator)`
		}
		return `objattr.Default(HTTPAuthFieldDefault, HTTPAuthFieldAttributes, HTTPAuthFieldValidator)`
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) GetValueStatement() string {
	if f.Hidden {
		switch f.Type {
		case FieldTypeString:
			return fmt.Sprintf(`c[%q] = %q`, f.Name, f.Initial.(string)) // nolint:forcetypeassert
		case FieldTypeBool:
			return fmt.Sprintf(`c[%q] = %t`, f.Name, f.Initial.(bool)) // nolint:forcetypeassert
		default:
			panic("unexpected hidden field type: " + f.Type)
		}
	}

	accessor := fmt.Sprintf(`m.%s`, f.StructName())
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return fmt.Sprintf(`stringattr.Get(%s, c, %q)`, accessor, f.Name)
	case FieldTypeBool:
		return fmt.Sprintf(`boolattr.Get(%s, c, %q)`, accessor, f.Name)
	case FieldTypeNumber:
		return fmt.Sprintf(`floatattr.Get(%s, c, %q)`, accessor, f.Name)
	case FieldTypeObject:
		return fmt.Sprintf(`getHeaders(%s, c, %q, h)`, accessor, f.Name)
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`listattr.Get(%s, c, %q, h)`, accessor, f.Name)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`objattr.Get(%s, c, %q, h)`, accessor, f.Name)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) SetValueStatement() string {
	accessor := fmt.Sprintf(`&m.%s`, f.StructName())
	switch f.Type {
	case FieldTypeString:
		return fmt.Sprintf(`stringattr.Set(%s, c, %q)`, accessor, f.Name)
	case FieldTypeSecret:
		return fmt.Sprintf(`stringattr.Nil(%s)`, accessor)
	case FieldTypeBool:
		return fmt.Sprintf(`boolattr.Set(%s, c, %q)`, accessor, f.Name)
	case FieldTypeNumber:
		return fmt.Sprintf(`floatattr.Set(%s, c, %q)`, accessor, f.Name)
	case FieldTypeObject:
		return fmt.Sprintf(`setHeaders(%s, c, %q, h)`, accessor, f.Name)
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`listattr.Set(%s, c, %q, h)`, accessor, f.Name)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`objattr.Set(%s, c, %q, h)`, accessor, f.Name)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) IsZero() string {
	accessor := fmt.Sprintf(`m.%s`, f.StructName())
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return fmt.Sprintf(`%s.ValueString() == ""`, accessor)
	case FieldTypeBool:
		return fmt.Sprintf(`!%s.ValueBool()`, accessor)
	case FieldTypeNumber:
		return fmt.Sprintf(`%s.ValueFloat64() == 0`, accessor)
	case FieldTypeObject:
		return fmt.Sprintf(`%s.IsEmpty()`, accessor)
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`%s.IsEmpty()`, accessor)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`!%s.IsSet()`, accessor)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

func (f *Field) IsNonZero() string {
	accessor := fmt.Sprintf(`m.%s`, f.StructName())
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		return fmt.Sprintf(`%s.ValueString() != ""`, accessor)
	case FieldTypeBool:
		return fmt.Sprintf(`%s.ValueBool()`, accessor)
	case FieldTypeNumber:
		return fmt.Sprintf(`%s.ValueFloat64() != 0`, accessor)
	case FieldTypeObject:
		return fmt.Sprintf(`!%s.IsEmpty()`, accessor)
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`!%s.IsEmpty()`, accessor)
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`%s.IsSet()`, accessor)
	default:
		panic("unexpected field type: " + f.Type)
	}
}

// Tests

func (f *Field) GetTestAssignment() string {
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		if v, ok := f.Initial.(string); ok {
			return fmt.Sprintf(`%q`, v)
		}
		if d := f.Dependency; d != nil && d.Field.Type == FieldTypeString && d.Value != d.Field.Initial {
			return `""`
		}
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

func (f *Field) GetTestCheck() string {
	switch f.Type {
	case FieldTypeString, FieldTypeSecret:
		if v, ok := f.Initial.(string); ok {
			return fmt.Sprintf(`"%s": %q`, f.AttributeName(), v)
		}
		if d := f.Dependency; d != nil && d.Field.Type == FieldTypeString && d.Value != d.Field.Initial {
			return fmt.Sprintf(`"%s": ""`, f.AttributeName())
		}
		return fmt.Sprintf(`"%s": %q`, f.AttributeName(), f.TestString())
	case FieldTypeBool:
		return fmt.Sprintf(`"%s": true`, f.AttributeName())
	case FieldTypeNumber:
		return fmt.Sprintf(`"%s": %d`, f.AttributeName(), f.TestNumber())
	case FieldTypeObject:
		return fmt.Sprintf(`"%s.key": %q`, f.AttributeName(), f.TestString())
	case FieldTypeAuditFilters:
		return fmt.Sprintf(`"%s.0.values": []string{%q}`, f.AttributeName(), f.TestString())
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`"%s.bearer_token": %q`, f.AttributeName(), f.TestString())
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

// Options

type FieldOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

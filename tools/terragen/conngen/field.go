package conngen

import (
	"fmt"

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
	case FieldTypeString, FieldTypeSecret, FieldTypeAuditFilters:
		return `types.String`
	case FieldTypeBool:
		return `types.Bool`
	case FieldTypeNumber:
		return `types.Float64`
	case FieldTypeObject:
		return `map[string]string`
	case FieldTypeHTTPAuth:
		return `*HTTPAuthFieldModel`
	default:
		panic("Unexpected field type: " + f.Type)
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
	case FieldTypeString, FieldTypeAuditFilters:
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
	case FieldTypeHTTPAuth:
		if f.Required {
			return `objectattr.Required(HTTPAuthFieldAttributes, HTTPAuthFieldValidator)`
		}
		return `objectattr.Optional(HTTPAuthFieldAttributes, HTTPAuthFieldValidator)`
	default:
		panic("Unexpected field type: " + f.Type)
	}
}

func (f *Field) GetValueStatement() string {
	switch f.Type {
	case FieldTypeString, FieldTypeSecret, FieldTypeAuditFilters:
		return fmt.Sprintf(`stringattr.Get(m.%s, c, %q)`, f.StructName(), f.Name)
	case FieldTypeBool:
		return fmt.Sprintf(`boolattr.Get(m.%s, c, %q)`, f.StructName(), f.Name)
	case FieldTypeNumber:
		return fmt.Sprintf(`floatattr.Get(m.%s, c, %q)`, f.StructName(), f.Name)
	case FieldTypeObject:
		return fmt.Sprintf(`c[%q] = m.%s`, f.Name, f.StructName())
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`objectattr.Get(m.%s, c, %q, h)`, f.StructName(), f.Name)
	default:
		panic("Unexpected field type: " + f.Type)
	}
}

func (f *Field) ValidateNonZero() string {
	accessor := fmt.Sprintf(`m.%s`, f.StructName())
	switch f.Type {
	case FieldTypeString, FieldTypeSecret, FieldTypeAuditFilters:
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
	case FieldTypeHTTPAuth:
		return fmt.Sprintf(`%s != nil`, accessor)
	default:
		panic("Unexpected field type: " + f.Type)
	}
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
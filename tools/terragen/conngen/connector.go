package conngen

import (
	"log"
	"slices"
	"strings"
	"unicode"

	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

// Connector

type Connector struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	BuiltIn      bool           `json:"builtin"`
	Validator    bool           `json:"validator"`
	Extra        map[string]any `json:"extra"`
	Fields       []*Field       `json:"fields"`
	HiddenFields []*Field       `json:"allFields"`

	naming *Naming
}

func (c *Connector) IsExperimental() bool {
	return c.Extra["experimental"] == true
}

func (c *Connector) IsSkipped() bool {
	return c.ID == "smtp-v2"
}

func (c *Connector) SupportsStaticIPs() bool {
	return c.Extra["supportStaticIps"] == true
}

func (c *Connector) StructName() string {
	return c.naming.GetName("connector", c.ID, "struct", c.defaultStructName())
}

func (c *Connector) defaultStructName() string {
	return utils.CapitalCase(c.ID)
}

func (c *Connector) FileName() string {
	return c.naming.GetName("connector", c.ID, "file", c.defaultFileName())
}

func (c *Connector) defaultFileName() string {
	var b strings.Builder
	for _, char := range c.ID {
		if char == '_' || char == '-' {
			// skip
		} else {
			b.WriteRune(unicode.ToLower(char))
		}
	}
	return b.String()
}

func (c *Connector) AttributeName() string {
	return c.naming.GetName("connector", c.ID, "attribute", c.defaultAttributeName())
}

func (c *Connector) defaultAttributeName() string {
	return utils.SnakeCase(c.ID)
}

func (c *Connector) DataName() string {
	return c.ID
}

func (c *Connector) HasField(typ string) bool {
	for _, f := range c.Fields {
		if f.Type == typ {
			return true
		}
	}
	return false
}

func (c *Connector) HasEnumFields() bool {
	for _, f := range c.Fields {
		if f.Type == FieldTypeString && len(f.Options) > 0 {
			return true
		}
	}
	return false
}

func (c *Connector) HasValidator() bool {
	return c.Validator || slices.ContainsFunc(c.Fields, func(f *Field) bool {
		return f.Dependency != nil
	})
}

func (c *Connector) Prepare() {
	// remove any fields that are not actually for configuration
	c.Fields = slices.DeleteFunc(c.Fields, func(f *Field) bool {
		return f.Type == "cloudformation-link"
	})

	// split hidden fields to the fields list
	c.Fields = slices.DeleteFunc(c.Fields, func(f *Field) bool {
		if f.Hidden {
			if f.Type != FieldTypeBool && f.Type != FieldTypeString {
				log.Fatalf("Hidden field %s in connector %s has unsupported type %s", f.Name, c.ID, f.Type)
			}
			if f.Initial == nil {
				log.Fatalf("Hidden field %s in connector %s must have an initial value", f.Name, c.ID)
			}
			c.HiddenFields = append(c.HiddenFields, f)
		}
		return f.Hidden
	})

	// add the static IP field into the configuration as expected by the snapshot format
	if c.SupportsStaticIPs() {
		c.Fields = append(c.Fields, UseStaticIPsField)
	}

	for _, f := range c.Fields {
		// treat these types as regular string fields for now
		if f.Type == "readonly-string" {
			f.Type = FieldTypeString
		}

		// treat secret-json-file as a secret field, as they are essentially identical
		if f.Type == "secret-json-file" {
			f.Type = FieldTypeSecret
		}

		if d := f.Dependency; d != nil {
			// link dependencies and fields together
			if d.Field == nil {
				for _, curr := range c.Fields {
					if d.Name == curr.Name {
						d.Field = curr
					}
				}
			}

			// a few sanity checks to make sure we support what's expected
			if d.Field == nil {
				log.Fatalf("Failed to find matching field for dependency %s in connector %s", d.Name, c.ID)
			}
			if d.Field.Type != FieldTypeBool && d.Field.Type != FieldTypeString {
				log.Fatalf("Field %s has a dependency on %s of type %s which is not supported", f.Name, d.Name, d.Field.Type)
			}
			if d.Field.Type == FieldTypeBool && f.Required {
				log.Fatalf("Unexpected required field with boolean dependency %s", f.Name)
			}
			if d.Field.Type == FieldTypeBool && d.Value != true {
				log.Fatalf("Field %s has a boolean dependency whose value is not true", f.Name)
			}
			if _, ok := d.Value.(string); !ok && d.Field.Type == FieldTypeString {
				log.Fatalf("Field %s has a string dependency whose value is not a string", f.Name)
			}

			// we convert required fields with string dependencies to optional
			if d.Field.Type == FieldTypeString {
				if !f.Required {
					log.Fatalf("Field %s has a string dependency so we expect it to be required in the template", f.Name)
				}
				f.Required = false
				if d.Field.Initial == nil {
					log.Fatalf("Field %s has a string dependency on field %s with no initial value", f.Name, d.Name)
				}
			}
		}
	}
}

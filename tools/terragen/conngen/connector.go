package conngen

import (
	"log"
	"strings"
	"unicode"

	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

// Connector

type Connector struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	BuiltIn     bool           `json:"builtin"`
	Extra       map[string]any `json:"extra"`
	Fields      []*Field       `json:"fields"`

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

func (c *Connector) HasValidator() bool {
	found := false
	for _, f := range c.Fields {
		if f.Dependency != nil {
			found = true
			// use this chance to link dependencies and fields together
			if f.Dependency.Field == nil {
				for _, curr := range c.Fields {
					if f.Dependency.Name == curr.Name {
						f.Dependency.Field = curr
					}
				}
			}
			// a few sanity checks to make sure we support what's expected
			if f.Required {
				log.Fatalf("Unexpected required field with dependency %s", f.Name)
			}
			if f.Dependency.Field == nil {
				log.Fatalf("Failed to find matching field for dependency %s in connector %s", f.Dependency.Name, c.ID)
			}
			if f.Dependency.Field.Type != FieldTypeBool {
				log.Fatalf("Field %s has a dependency on %s whose type is not a boolean (other types are not currently supported)", f.Name, f.Dependency.Field.Name)
			}
			if f.Dependency.Value != true {
				log.Fatalf("Field %s has a dependency whose value is not true (other values are not currently supported)", f.Name)
			}
		}
	}
	return found
}

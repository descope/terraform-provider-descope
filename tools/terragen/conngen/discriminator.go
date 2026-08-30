package conngen

import (
	"fmt"
	"log"
	"strings"
)

// Discriminator is a computed connector configuration key whose value tells the backend which of several mutually-exclusive fields is
// in use. It is not a Terraform attribute: the value is derived from which underlying field the user populated (e.g. Twilio's
// `selectedProp` is `fromPhone` when the from-phone field is set, otherwise `messagingServiceSid`).
type Discriminator struct {
	Name    string               `json:"name"`    // the backend configuration key to emit
	Cases   []*DiscriminatorCase `json:"cases"`   // matched in order; the first whose field is set wins
	Default string               `json:"default"` // emitted when no case field is set
}

type DiscriminatorCase struct {
	Field string `json:"field"` // the name of the connector field this case checks
	Value string `json:"value"` // the discriminator value emitted when that field is set

	field *Field
}

// link resolves each case to its connector field so generated code can reference it by struct name, aborting on an unknown field.
func (d *Discriminator) link(c *Connector) {
	for _, dc := range d.Cases {
		for _, f := range c.Fields {
			if f.Name == dc.Field {
				dc.field = f
				break
			}
		}
		if dc.field == nil {
			log.Fatalf("discriminator %q in connector %s references unknown field %q", d.Name, c.ID, dc.Field)
		}
	}
}

// GetValueStatement renders the Go that emits the discriminator value, taking the first case whose field is set or else the default.
func (d *Discriminator) GetValueStatement() string {
	var b strings.Builder
	b.WriteString("switch {\n")
	for _, dc := range d.Cases {
		fmt.Fprintf(&b, "case m.%s.ValueString() != \"\":\n", dc.field.ResourceFieldName())
		fmt.Fprintf(&b, "c[%q] = %q\n", d.Name, dc.Value)
	}
	fmt.Fprintf(&b, "default:\nc[%q] = %q\n}", d.Name, d.Default)
	return b.String()
}

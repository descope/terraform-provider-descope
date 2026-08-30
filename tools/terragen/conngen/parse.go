package conngen

import (
	"github.com/descope/terraform-provider-descope/tools/terragen/schema"
	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
	"github.com/mitchellh/go-wordwrap"
)

func ParseConnectors(datadir string, mockdir string, templatesdir string) *Connectors {
	conns := &Connectors{
		Naming: &Naming{},
	}

	conns.Read(mockdir, templatesdir)

	conns.Naming.Read(datadir)

	for _, c := range conns.Connectors {
		c.naming = conns.Naming
		for _, field := range c.Fields {
			field.naming = conns.Naming
		}
	}

	return conns
}

// MergeDocs marks the standalone resource models as generated and describes every field, so no raw documentation files are expected
// for them: boilerplate fields get hardcoded descriptions, configuration fields get theirs from the connector template metadata.
func MergeDocs(conns *Connectors, sc *schema.Schema) {
	for _, c := range conns.Connectors {
		model := findOptionalModel(sc, c.ResourceStructName())
		if model == nil {
			continue // defensive: every connector should have a generated resource model
		}
		model.Generated = true
		for _, field := range model.Fields {
			if desc := connectorFieldDescription(c, field.Name); desc != "" {
				field.Description = desc
			}
		}
	}
}

func connectorFieldDescription(c *Connector, name string) string {
	switch name {
	case "project_id":
		return utils.DefaultConnectorProjectIDText
	case "name":
		return utils.DefaultConnectorNameText
	case "description":
		return utils.DefaultConnectorDescriptionText
	case "disabled":
		return utils.DefaultConnectorDisabledText
	}
	for _, f := range c.Fields {
		if f.ResourceAttributeName() == name && f.Description != "" {
			return wordwrap.WrapString(f.Description, 80)
		}
	}
	return ""
}

func findOptionalModel(sc *schema.Schema, name string) *schema.Model {
	for _, f := range sc.Files {
		for _, m := range f.Models {
			if m.Name == name {
				return m
			}
		}
	}
	return nil
}

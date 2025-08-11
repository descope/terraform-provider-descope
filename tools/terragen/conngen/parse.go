package conngen

import (
	"log"
	"slices"

	"github.com/descope/terraform-provider-descope/tools/terragen/schema"
	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
	"github.com/mitchellh/go-wordwrap"
)

func ParseConnectors(datadir string, templatesdir string) *Connectors {
	conns := &Connectors{
		Naming: &Naming{},
	}

	conns.Read(datadir, templatesdir)

	conns.Naming.Read(datadir)

	for _, c := range conns.Connectors {
		c.naming = conns.Naming
		for _, field := range c.Fields {
			field.naming = conns.Naming
		}
	}

	return conns
}

func MergeDocs(conns *Connectors, sc *schema.Schema) {
	model := findConnectorsContainer(sc)
	for _, c := range conns.Connectors {
		mergeConnectorDocs(c, sc)
		for _, f := range model.Fields {
			if f.Name == c.ID && c.Description != "" {
				f.Description = wordwrap.WrapString(c.Description, 80)
			}
		}
	}
}

func mergeConnectorDocs(c *Connector, sc *schema.Schema) {
	model := findConnectorModel(sc, c.StructName())
	for _, field := range model.Fields {
		for _, f := range c.Fields {
			if field.Name == f.AttributeName() {
				if f.Description != "" {
					field.Description = wordwrap.WrapString(f.Description, 80)
				}
				break
			}
		}
	}
}

func StripBoilerplate(conns *Connectors, sc *schema.Schema) {
	updateBoilerplate(conns, sc, true)
}

func AddBoilerplate(conns *Connectors, sc *schema.Schema) {
	updateBoilerplate(conns, sc, false)
}

func updateBoilerplate(conns *Connectors, sc *schema.Schema, strip bool) {
	name := &schema.Field{Name: "name", Description: utils.DefaultConnectorNameText}
	desc := &schema.Field{Name: "description", Description: utils.DefaultConnectorDescriptionText}
	for _, c := range conns.Connectors {
		model := findConnectorModel(sc, c.StructName())
		model.Generated = !c.BuiltIn
		if strip {
			model.Fields = slices.DeleteFunc(model.Fields, func(f *schema.Field) bool { return f.Name == name.Name || f.Name == desc.Name })
		} else {
			model.Fields = slices.Concat([]*schema.Field{name}, []*schema.Field{desc}, model.Fields)
		}
	}
}

func findConnectorModel(sc *schema.Schema, name string) *schema.Model {
	for _, f := range sc.Files {
		for _, m := range f.Models {
			if m.Name == name {
				return m
			}
		}
	}
	log.Fatalf("expected to find connector model for %s", name)
	return nil
}

func findConnectorsContainer(sc *schema.Schema) *schema.Model {
	for _, f := range sc.Files {
		if slices.Equal(f.Dirs, []string{"project", "connectors"}) && f.Name == "connectors" {
			if len(f.Models) != 1 {
				log.Fatalf("unexpected connectors container file with multiple models")
			}
			return f.Models[0]
		}
	}
	log.Fatal("expected to find connectors container model")
	return nil
}

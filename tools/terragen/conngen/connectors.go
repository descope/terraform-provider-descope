package conngen

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

type Connectors struct {
	Connectors []*Connector
	Naming     *Naming
}

func (c *Connectors) Read(datadir string, templatesdir string) {
	utils.Debug(0, "Connectors")
	utils.Debug(0, "==========")
	c.readBuiltins(datadir)
	c.readTemplates(templatesdir)
	slices.SortFunc(c.Connectors, func(a, b *Connector) int { return strings.Compare(a.ID, b.ID) })
	utils.Debug(0, "")
}

func (c *Connectors) readBuiltins(datadir string) {
	builtins := []*Connector{}

	path := filepath.Join(datadir, "builtins.json")
	if err := utils.ReadJSON(path, &builtins); err != nil {
		log.Fatalf("error reading builtin.json file: %s", err.Error())
	}

	utils.Debug(0, "+ builtins:")
	for _, conn := range builtins {
		utils.Debug(1, "- %s", conn.ID)
		conn.BuiltIn = true
	}

	c.Connectors = append(c.Connectors, builtins...)
}

func (c *Connectors) readTemplates(templatesdir string) {
	entries, err := os.ReadDir(templatesdir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("templates directory doesn't exist: %s", templatesdir)
		} else {
			log.Fatalf("failed to read files from path %s: %s", templatesdir, err.Error())
		}
	}

	paths := []string{}
	for _, v := range entries {
		if v.IsDir() && !strings.HasPrefix(v.Name(), ".") {
			paths = append(paths, filepath.Join(templatesdir, v.Name()))
		}
	}

	utils.Debug(0, "+ templates:")
	for _, path := range paths {
		c.readConnector(path)
	}
}

func (c *Connectors) readConnector(path string) {
	file := filepath.Join(path, "metadata.json")

	connector := &Connector{}
	if err := utils.ReadJSON(file, connector); err != nil {
		log.Fatalf("failed to read connector metadata from path %s: %s", file, err.Error())
	}

	if connector.IsExperimental() {
		utils.Debug(1, "- %s (experimental)", connector.ID)
		return
	}

	if connector.SupportsStaticIPs() {
		connector.Fields = append(connector.Fields, UseStaticIPsField)
	}

	c.Connectors = append(c.Connectors, connector)
	utils.Debug(1, "- %s", connector.ID)
}

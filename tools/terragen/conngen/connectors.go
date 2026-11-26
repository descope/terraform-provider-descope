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

// Read loads all connectors, in-repo mock templates first and external connector templates second. A mock template takes precedence
// over an external template with the same id (e.g. the messaging connectors we ship a flat mock template for).
func (c *Connectors) Read(mockdir string, templatesdir string) {
	utils.Debug(0, "Connectors")
	utils.Debug(0, "==========")
	utils.Debug(0, "+ mock templates:")
	c.readTemplates(mockdir)
	utils.Debug(0, "+ templates:")
	c.readTemplates(templatesdir)
	slices.SortFunc(c.Connectors, func(a, b *Connector) int { return strings.Compare(a.ID, b.ID) })
	utils.Debug(0, "")
}

func (c *Connectors) has(id string) bool {
	return slices.ContainsFunc(c.Connectors, func(conn *Connector) bool { return conn.ID == id })
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

	if connector.IsSkipped() {
		utils.Debug(1, "- %s (skipped)", connector.ID)
		return
	}

	// a mock template with the same id was already read, so it takes precedence
	if c.has(connector.ID) {
		utils.Debug(1, "- %s (overridden)", connector.ID)
		return
	}

	connector.Prepare()

	c.Connectors = append(c.Connectors, connector)
	utils.Debug(1, "- %s", connector.ID)
}

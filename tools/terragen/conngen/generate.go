package conngen

import (
	_ "embed"
	"os"
	"path/filepath"
	"slices"

	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

//go:embed resource.gotmpl
var resourceTemplateData []byte

//go:embed resourcetest.gotmpl
var resourceTestTemplateData []byte

//go:embed resources.gotmpl
var resourcesTemplateData []byte

// TrimConnectors drops connectors whose generated model doesn't exist yet, so a normal run only regenerates connectors that were
// already added. Run with --add-connectors to generate models for new connectors for the first time.
func TrimConnectors(dir string, conns *Connectors) {
	if utils.Flags.AddConnectors {
		return
	}

	conns.Connectors = slices.DeleteFunc(conns.Connectors, func(connector *Connector) bool {
		path := filepath.Join(dir, connector.FileName()+".go")
		_, err := os.Stat(path)
		return err != nil
	})
}

// Generates the standalone resource model .go source and acceptance test for each connector.
func GenerateResourceSources(dir string, conns *Connectors) {
	resourceTemplate := utils.LoadTemplate("resource", resourceTemplateData)
	resourceTestTemplate := utils.LoadTemplate("resourcetest", resourceTestTemplateData)
	for _, connector := range conns.Connectors {
		utils.WriteGoSource(filepath.Join(dir, connector.FileName()+".go"), connector, resourceTemplate, true)
		utils.WriteGoSource(filepath.Join(dir, connector.FileName()+"_test.go"), connector, resourceTestTemplate, true)
	}
}

// Generates the connectors.go registration file in the resources package, with a constructor per standalone connector resource.
func GenerateResourceRegistrations(resourcesdir string, conns *Connectors) {
	tpl := utils.LoadTemplate("resources", resourcesTemplateData)
	path := filepath.Join(resourcesdir, "connectors.go")
	utils.WriteGoSource(path, map[string]any{"Registered": conns.Connectors}, tpl, true)
}

func UpdateNaming(dir string, conns *Connectors) {
	if conns.Naming.HasChanges {
		conns.Naming.Write(dir)
	}
}

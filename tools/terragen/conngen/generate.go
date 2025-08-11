package conngen

import (
	_ "embed"
	"os"
	"path/filepath"
	"slices"

	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

//go:embed connector.gotmpl
var connectorTemplateData []byte

//go:embed connectors.gotmpl
var connectorsTemplateData []byte

//go:embed test.gotmpl
var testTemplateData []byte

func TrimConnectors(dir string, conns *Connectors) {
	if utils.Flags.AddConnectors {
		return
	}

	conns.Connectors = slices.DeleteFunc(conns.Connectors, func(connector *Connector) bool {
		if !connector.BuiltIn {
			path := filepath.Join(dir, connector.FileName()+".go")
			_, err := os.Stat(path)
			return err != nil
		}
		return false
	})
}

func GenerateSources(dir string, conns *Connectors) {
	connectorTemplate := utils.LoadTemplate("connector", connectorTemplateData)
	for _, connector := range conns.Connectors {
		if !connector.BuiltIn {
			path := filepath.Join(dir, connector.FileName()+".go")
			utils.WriteGoSource(path, connector, connectorTemplate, true)
		}
	}

	connectorsTemplate := utils.LoadTemplate("connectors", connectorsTemplateData)
	path := filepath.Join(dir, "connectors.go")
	utils.WriteGoSource(path, conns, connectorsTemplate, true)

	testTemplate := utils.LoadTemplate("test", testTemplateData)
	path = filepath.Join(dir, "connectors_test.go")
	utils.WriteGoSource(path, conns, testTemplate, true)
}

func UpdateNaming(dir string, conns *Connectors) {
	if conns.Naming.HasChanges {
		conns.Naming.Write(dir)
	}
}

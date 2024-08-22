package srcgen

import (
	_ "embed"
	"path/filepath"

	"github.com/descope/terraform-provider-descope/tools/terragen/schema"
	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

//go:embed docs.gotmpl
var docsTemplateData []byte

//go:embed models.gotmpl
var modelsTemplateData []byte

func GenerateSources(root string, schema *schema.Schema) {
	tpl := utils.LoadTemplate("docs", docsTemplateData)
	path := filepath.Join(root, "docs.go")
	utils.WriteGoSource(path, schema, tpl, false)

	tpl = utils.LoadTemplate("models", modelsTemplateData)
	path = filepath.Join(root, "models.go")
	utils.WriteGoSource(path, schema, tpl, true)
}

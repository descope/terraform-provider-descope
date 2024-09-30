package docgen

import (
	"bytes"
	_ "embed"
	"log"
	"os"
	"path/filepath"

	"github.com/descope/terraform-provider-descope/tools/terragen/schema"
	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

//go:embed docfile.gotmpl
var docfileTemplateData []byte

func GenerateDocs(root string, schema *schema.Schema) {
	tpl := utils.LoadTemplate("docfile", docfileTemplateData)
	for _, file := range schema.Files {
		if file.SkipDocs() {
			continue
		}

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, file); err != nil {
			log.Fatalf("error executing template file: %s", err.Error())
		}

		data := buf.Bytes()
		for bytes.HasSuffix(data, []byte{'\n', '\n'}) {
			data = data[:len(data)-1]
		}

		path := utils.EnsurePath(root, filepath.Join(file.Dirs...))
		file := filepath.Join(path, file.Name+".md")
		if err := os.WriteFile(file, data, 0644); err != nil {
			log.Fatalf("error writing documentation file: %s", err.Error())
		}
	}
}

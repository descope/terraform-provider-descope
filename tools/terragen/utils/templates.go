package utils

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"
)

// written to doc files when there's no description found for an attribute.
const PlaceholderDescription = `// description for`

// defaults for common connector fields
const DefaultConnectorNameText = "A custom name for your connector."
const DefaultConnectorDescriptionText = "A description of what your connector is used for."

// markdown descriptions with paragraphs will require a custom template to look good
const preserveParagraphs = false

func LoadTemplate(name string, data []byte) *template.Template {
	tpl, err := template.New(name).Funcs(templateUtils).Parse(string(data))
	if err != nil {
		log.Fatalf("error parsing %s template file: %s", name, err.Error())
	}
	return tpl
}

var templateUtils = map[string]any{
	"head": func(s string) string {
		return strings.Repeat("=", max(4, len(s)))
	},
	"subhead": func(s string) string {
		return strings.Repeat("-", max(4, len(s)))
	},
	"placeholder": func(s string) string {
		return fmt.Sprintf("%s %s", PlaceholderDescription, s)
	},
	"srcliteral": func(name string, description string) string {
		r := []string{}
		parts := strings.Split(description, "\n")
		for i, v := range parts {
			if len(v) == 0 && !preserveParagraphs {
				continue
			}
			v = strings.ReplaceAll(v, `\`, `\\`)
			v = strings.ReplaceAll(v, `"`, `\"`)
			if len(v) == 0 && preserveParagraphs && i != 0 {
				v = "\\n\\n"
			} else if len(v) > 0 && i != len(parts)-1 {
				v = v + " "
			}
			r = append(r, v)
		}
		joiner := `" +` + "\n\t" + strings.Repeat(" ", len(name)+4) + `"`
		return `"` + strings.Join(r, joiner) + `"`
	},
}

func WriteGoSource(path string, object any, tpl *template.Template, gofmt bool) {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, object); err != nil {
		log.Fatalf("error executing template file: %s", err.Error())
	}

	data := buf.Bytes()
	if gofmt {
		formatted, err := format.Source(data)
		if err != nil {
			log.Fatalf("error formatting generated file: %s", err.Error())
		}
		data = formatted
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Fatalf("error writing generated source file %s: %s", path, err.Error())
	}
}

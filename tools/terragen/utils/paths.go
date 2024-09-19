package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Paths struct {
	Root       string
	Models     string
	Docs       string
	Raw        string
	Connectors string
	Data       string
	Templates  string
}

func PreparePaths() *Paths {
	curr, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %s", err.Error())
	}

	root := filepath.Clean(curr)
	if filepath.Base(root) != "terraform-provider-descope" {
		log.Fatalf("expected to run from the project root directory")
	}

	models := filepath.Join(root, "internal", "models")
	if info, err := os.Stat(models); os.IsNotExist(err) || !info.IsDir() {
		log.Fatalf("expected to find models directory at path: %s", models)
	}

	docs := EnsurePath(root, "internal", "docs")

	raw := EnsurePath(docs, "raw")

	connectors := filepath.Join(models, "connectors")

	data := filepath.Join(root, "tools", "terragen", "conngen")

	templates := ""
	if !Flags.SkipTemplates {
		templates = strings.TrimSpace(os.Getenv("DESCOPE_TEMPLATES_PATH"))
		if templates == "" {
			log.Fatalf("expected path to templates in DESCOPE_TEMPLATES_PATH environment variable")
		}
		templates = filepath.Clean(templates)
		if info, err := os.Stat(templates); os.IsNotExist(err) || !info.IsDir() {
			log.Fatalf("expected to find templates directory at path: %s", templates)
		}
	}

	return &Paths{
		Root:       root,
		Models:     models,
		Docs:       docs,
		Raw:        raw,
		Connectors: connectors,
		Data:       data,
		Templates:  templates,
	}
}

func EnsurePath(path string, subdirs ...string) string {
	for _, d := range subdirs {
		path = filepath.Join(path, d)
		if err := os.Mkdir(path, 0755); err != nil && !os.IsExist(err) {
			log.Fatalf("failed to create subdirectory %s: %s", path, err.Error())
		}
	}
	return path
}

package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Paths struct {
	Root          string
	Models        string
	Docs          string
	Raw           string
	Connectors    string
	Resources     string
	Data          string
	MockTemplates string
	Templates     string
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

	models := requireDir(root, "internal", "models")

	docs := EnsurePath(root, "internal", "docs")

	raw := EnsurePath(root, "docs", "raw")

	connectors := EnsurePath(models, "connectors")

	resources := requireDir(root, "internal", "resources")

	data := filepath.Join(root, "tools", "terragen", "conngen")

	mockTemplates := requireDir(data, "templates")

	templates := strings.TrimSpace(os.Getenv("DESCOPE_TEMPLATES_PATH"))
	if templates == "" {
		log.Fatalf("expected path to templates in DESCOPE_TEMPLATES_PATH environment variable")
	}
	templates = requireDir(filepath.Clean(templates))

	return &Paths{
		Root:          root,
		Models:        models,
		Docs:          docs,
		Raw:           raw,
		Connectors:    connectors,
		Resources:     resources,
		Data:          data,
		MockTemplates: mockTemplates,
		Templates:     templates,
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

func requireDir(path string, subdirs ...string) string {
	path = filepath.Join(path, filepath.Join(subdirs...))
	if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
		log.Fatalf("expected to find directory at path: %s", path)
	}
	return path
}

package conngen

import (
	"log"
	"path/filepath"

	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

type Naming struct {
	Names      map[string]map[string]map[string]string
	HasChanges bool
}

func (n *Naming) Read(datadir string) {
	if err := utils.ReadJSON(n.Path(datadir), &n.Names); err != nil {
		log.Fatalf("error reading naming.json file: %s", err.Error())
	}
}

func (n *Naming) Write(datadir string) {
	if err := utils.WriteJSON(n.Path(datadir), &n.Names); err != nil {
		log.Fatalf("error writing naming.json file: %s", err.Error())
	}
}

func (n *Naming) Path(datadir string) string {
	return filepath.Join(datadir, "naming.json")
}

func (n *Naming) GetName(category, id, kind, fallback string) string {
	categoryNames, ok := n.Names[category]
	if !ok {
		categoryNames = map[string]map[string]string{}
		n.Names[category] = categoryNames
	}
	idNames, ok := categoryNames[id]
	if !ok {
		idNames = map[string]string{}
		categoryNames[id] = idNames
	}
	name, ok := idNames[kind]
	if !ok {
		name = fallback
		idNames[kind] = name
		n.HasChanges = true
	}
	return name
}

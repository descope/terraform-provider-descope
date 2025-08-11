package main

import (
	"github.com/descope/terraform-provider-descope/tools/terragen/conngen"
	"github.com/descope/terraform-provider-descope/tools/terragen/docgen"
	"github.com/descope/terraform-provider-descope/tools/terragen/schema"
	"github.com/descope/terraform-provider-descope/tools/terragen/srcgen"
	"github.com/descope/terraform-provider-descope/tools/terragen/utils"
)

func main() {
	// parses the command line flags into the Flags struct in utils
	utils.ParseFlags()

	// ensures that required paths are available and creates directories for generated files
	paths := utils.PreparePaths()

	// parses all connector template metadata
	conns := conngen.ParseConnectors(paths.Data, paths.Templates)

	// remove any connectors from the templates that don't already exist (unless --add-connectors flag was set)
	conngen.TrimConnectors(paths.Connectors, conns)

	// generates .go sources and tests for all connector models
	conngen.GenerateSources(paths.Connectors, conns)

	// creates a simple schema representation by parsing attributes in all model .go source files
	schema := schema.ParseSources(paths.Models)

	// copies model descriptions from the connector template metadata files into the schema
	conngen.MergeDocs(conns, schema)

	// strip repetitive boilerplate fields from generated docs
	conngen.StripBoilerplate(conns, schema)

	// copies existing model descriptions from the raw .md documentation files into the schema
	docgen.MergeDocs(paths.Raw, schema)

	// checks that nothing went wrong and all docs are available, aborts if not (unless --skip-validate flag was set)
	schema.ValidateIfNeeded()

	// generates updated raw .md documentation files
	docgen.GenerateDocs(paths.Raw, schema)

	// stop after generating .md files if needed
	schema.AbortIfNeeded()

	// add back boilerplate fields with hardcoded descriptions
	conngen.AddBoilerplate(conns, schema)

	// generates model documentation injection .go source files that are actually compiled into the provider binary
	srcgen.GenerateSources(paths.Docs, schema)

	// updates the naming.json file if needed
	conngen.UpdateNaming(paths.Data, conns)
}

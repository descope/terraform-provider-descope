package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/discover"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/export"
)

func main() {
	projectID := flag.String("project", "", "the ID of the Descope project to export")
	outDir := flag.String("out", "", "the directory to write generated files into")
	only := flag.String("only", "", "limit the export to resource types containing this substring")
	dumpIndex := flag.Bool("dump-index", false, "print a structural summary of the project's snapshot index and exit")
	dumpFile := flag.String("dump-file", "", "print the raw JSON of a single snapshot file and exit")
	printState := flag.Bool("print-state", false, "print the attribute values of all read resources and exit")
	flag.Parse()

	if *projectID == "" {
		fail("The -project flag is required")
	}
	managementKey := os.Getenv("DESCOPE_MANAGEMENT_KEY")
	if managementKey == "" {
		fail("The DESCOPE_MANAGEMENT_KEY environment variable is required")
	}

	ctx := context.Background()
	client := infra.NewClient("tfexport", managementKey, os.Getenv("DESCOPE_BASE_URL"))

	switch {
	case *dumpIndex || *dumpFile != "":
		dumpSnapshot(ctx, client, *projectID, *dumpFile)
	case *printState:
		printResources(ctx, client, *projectID, *only)
	default:
		if *outDir == "" {
			fail("The -out flag is required")
		}
		count, warnings, err := export.Run(ctx, client, *projectID, *outDir, *only)
		if err != nil {
			fail("%s", err.Error())
		}
		printWarnings(warnings)
		fmt.Printf("Generated %d resources in %s\n", count, *outDir)
	}
}

func dumpSnapshot(ctx context.Context, client *infra.Client, projectID, file string) {
	files, err := discover.Snapshot(ctx, client, projectID)
	if err != nil {
		fail("Failed to export project snapshot: %s", err.Error())
	}
	if file == "" {
		discover.DumpIndex(os.Stdout, files)
		return
	}
	value, ok := files[file]
	if !ok {
		fail("No file %q in snapshot", file)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		fail("Failed to encode snapshot file: %s", err.Error())
	}
}

func printResources(ctx context.Context, client *infra.Client, projectID, only string) {
	results, warnings, err := export.ReadProject(ctx, client, projectID, only)
	if err != nil {
		fail("%s", err.Error())
	}
	for _, result := range results {
		fmt.Printf("### %s %q (id %s)\n%s\n\n", result.Instance.Resource, result.Instance.Name, result.Instance.ID, result.Object.String())
	}
	printWarnings(warnings)
}

func printWarnings(warnings []string) {
	for _, warning := range warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

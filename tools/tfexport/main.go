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
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/warn"
)

func main() {
	projectID := flag.String("project", "", "the ID of the Descope project to export")
	outDir := flag.String("out", "", "the directory to write generated files into")
	only := flag.String("only", "", "limit the export to resource types containing this substring")
	force := flag.Bool("force", false, "write into the output directory even when it already has files in it")
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
		if entries, err := os.ReadDir(*outDir); err == nil && len(entries) > 0 && !*force {
			fail("The output directory %s is not empty: files this export doesn't overwrite would survive into it, pass -force to write anyway", *outDir)
		}
		count, warnings, err := export.Run(ctx, client, *projectID, *outDir, *only)
		if err != nil {
			fail("%s", err.Error())
		}
		printWarnings(warnings)
		fmt.Printf("Generated %d resources in %s\n", count, *outDir)
		if warn.Incomplete(warnings) {
			os.Exit(2)
		}
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
	export.PrintState(ctx, os.Stdout, results)
	printWarnings(warnings)
	if warn.Incomplete(warnings) {
		os.Exit(2)
	}
}

func printWarnings(warnings []warn.Warning) {
	for _, warning := range warnings {
		label := "warning"
		if warning.Lossy {
			label = "incomplete"
		}
		fmt.Fprintf(os.Stderr, "%s: %s\n", label, warning.Text)
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

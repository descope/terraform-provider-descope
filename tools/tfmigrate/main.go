// Command tfmigrate migrates the ownership recorded in a legacy descope_project resource onto the standalone
// resources it was split into. It is read only: it never applies anything, never touches Terraform state and never
// makes a mutating Descope API call.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/descope/terraform-provider-descope/tools/tfmigrate/migrate"
)

// Exit codes: 0 success, 1 blockers or a failed verification, 2 usage or I/O error.
const (
	exitOK      = 0
	exitBlocked = 1
	exitUsage   = 2
)

var exactProviderVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func writef(out io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(out, format, args...)
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return exitUsage
	}
	switch args[0] {
	case "plan":
		return planCommand(args[1:], stdout, stderr)
	case "verify-plan":
		return verifyCommand(args[1:], stdout, stderr)
	case "-h", "--help", "help":
		usage(stdout)
		return exitOK
	default:
		writef(stderr, "unknown command %q\n\n", args[0])
		usage(stderr)
		return exitUsage
	}
}

func usage(out io.Writer) {
	writef(out, "%s", `tfmigrate migrates legacy descope_project ownership onto the standalone descope resources.

Usage:
  tfmigrate plan         -state FILE -address ADDR -source-provider-version V -target-provider-version V -out DIR
  tfmigrate verify-plan  -plan FILE -manifest FILE

plan builds the ownership manifest and the two stage migration artifacts from a Terraform state snapshot.
verify-plan checks that the adoption plan only imports and reconciles, and never creates, replaces or destroys.

Exit codes: 0 success, 1 blockers or failed verification, 2 usage or I/O error.
`)
}

func planCommand(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("plan", flag.ContinueOnError)
	flags.SetOutput(stderr)
	state := flags.String("state", "", "path to a Terraform state file or the output of `terraform show -json`")
	address := flags.String("address", "", "address of the legacy descope_project instance, module qualified if needed")
	sourceVersion := flags.String("source-provider-version", "", "exact provider version the state was written with")
	targetVersion := flags.String("target-provider-version", "", "exact provider version being adopted")
	out := flags.String("out", "", "new output directory for the generated migration; the path must not exist")
	remote := flags.Bool("remote", false, "allow one GET-only read of the project to backfill ids missing from state")
	baseURL := flags.String("base-url", "", "HTTPS Descope base URL for -remote, using a DNS host and default TLS port; defaults to $DESCOPE_BASE_URL")
	verbose := flags.Bool("verbose", false, "log every destination as it is mapped")
	if err := flags.Parse(args); err != nil {
		return exitUsage
	}
	missing := []string{}
	for _, required := range []struct {
		name  string
		value string
	}{
		{"-state", *state},
		{"-address", *address},
		{"-source-provider-version", *sourceVersion},
		{"-target-provider-version", *targetVersion},
		{"-out", *out},
	} {
		if required.value == "" {
			missing = append(missing, required.name)
		}
	}
	if len(missing) > 0 {
		writef(stderr, "missing required flags: %s\n", strings.Join(missing, ", "))
		return exitUsage
	}
	for _, version := range []struct {
		name  string
		value string
	}{{"-source-provider-version", *sourceVersion}, {"-target-provider-version", *targetVersion}} {
		if !exactProviderVersion.MatchString(version.value) {
			writef(stderr, "%s must be an exact semantic version such as 1.8.3, got %q\n", version.name, version.value)
			return exitUsage
		}
	}

	snapshot, err := os.ReadFile(*state)
	if err != nil {
		writef(stderr, "reading state: %s\n", err)
		return exitUsage
	}
	stateFile, err := migrate.LoadStateFile(snapshot)
	if err != nil {
		writef(stderr, "%s\n", err)
		return exitUsage
	}
	legacy, err := stateFile.Find(*address)
	if err != nil {
		writef(stderr, "%s\n", err)
		return exitUsage
	}

	ctx := context.Background()
	registry, err := migrate.NewRegistry(ctx)
	if err != nil {
		writef(stderr, "discovering destination resources: %s\n", err)
		return exitUsage
	}

	if *remote {
		reader, readerErr := migrate.NewProjectReader(projectIDOf(legacy), *baseURL)
		if readerErr != nil {
			writef(stderr, "%s\n", readerErr)
			return exitUsage
		}
		filled, backfillErr := migrate.Backfill(ctx, legacy, reader)
		if backfillErr != nil {
			writef(stderr, "backfilling ids: %s\n", backfillErr)
			return exitUsage
		}
		writef(stdout, "backfilled %d id(s) with a GET-only read of project %s\n", len(filled), projectIDOf(legacy))
		if *verbose {
			for _, entry := range filled {
				writef(stdout, "  id recovered for %s\n", entry)
			}
		}
	}

	opts := migrate.Options{
		ProviderSource:        migrate.ProviderSource,
		SourceProviderVersion: *sourceVersion,
		TargetProviderVersion: *targetVersion,
		Remote:                *remote,
	}
	manifest, blocks, payloads, err := migrate.Build(ctx, registry, legacy, opts)
	if err != nil {
		writef(stderr, "%s\n", err)
		return exitUsage
	}
	artifacts, err := migrate.Generate(manifest, blocks, payloads, snapshot, opts)
	if err != nil {
		writef(stderr, "%s\n", err)
		return exitUsage
	}
	if err := artifacts.Write(*out); err != nil {
		writef(stderr, "%s\n", err)
		return exitUsage
	}

	writef(stdout, "%s\n", manifest.Summary())
	if *verbose {
		for _, entry := range manifest.Entries {
			writef(stdout, "  %s -> %s (import id %q)\n", entry.LegacyPath, entry.Address, entry.ImportID)
		}
	}
	if !manifest.Ready {
		writef(stderr, "\n%d blocker(s): see %s, no migration configuration was generated\n", len(manifest.Blockers), filepath.Join(*out, "BLOCKERS.md"))
		return exitBlocked
	}
	writef(stdout, "review %s before running either stage\n", filepath.Join(*out, "README.md"))
	return exitOK
}

func verifyCommand(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("verify-plan", flag.ContinueOnError)
	flags.SetOutput(stderr)
	plan := flags.String("plan", "", "path to the output of `terraform show -json` for the adoption plan")
	manifestPath := flags.String("manifest", "", "path to the manifest.json produced by the plan command")
	if err := flags.Parse(args); err != nil {
		return exitUsage
	}
	if *plan == "" || *manifestPath == "" {
		writef(stderr, "%s", "both -plan and -manifest are required\n")
		return exitUsage
	}
	planJSON, err := os.ReadFile(*plan)
	if err != nil {
		writef(stderr, "reading plan: %s\n", err)
		return exitUsage
	}
	manifestJSON, err := os.ReadFile(*manifestPath)
	if err != nil {
		writef(stderr, "reading manifest: %s\n", err)
		return exitUsage
	}
	manifest, err := migrate.DecodeManifest(manifestJSON)
	if err != nil {
		writef(stderr, "%s\n", err)
		return exitUsage
	}
	result, err := migrate.VerifyPlan(planJSON, manifest)
	if err != nil {
		writef(stderr, "%s\n", err)
		return exitUsage
	}
	writef(stdout, "%s", result.Report())
	if !result.OK {
		return exitBlocked
	}
	return exitOK
}

func projectIDOf(legacy *migrate.LegacyResource) string {
	id, _ := legacy.Attributes["id"].(string)
	return id
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/harness"
)

func TestConnectorSweep(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("DESCOPE_TFEXPORT_ROUNDTRIP") == "" {
		t.Skip("set TF_ACC and DESCOPE_TFEXPORT_ROUNDTRIP to run this suite")
	}
	if os.Getenv("DESCOPE_MANAGEMENT_KEY") == "" {
		t.Skip("set DESCOPE_MANAGEMENT_KEY to run this suite")
	}

	configs := connectorSeedConfigs(t)
	t.Logf("assembled seed configurations for %d connector types", len(configs))
	if len(configs) < 50 {
		t.Fatalf("expected seed configurations for most connector types, got %d", len(configs))
	}

	// connector types that require a company license the test environment lacks (creation fails with E122009, status 403)
	licensed := []string{"fingerprint_descope_connector"}

	types := make([]string, 0, len(configs))
	for connectorType := range configs {
		if !slices.Contains(licensed, connectorType) {
			types = append(types, connectorType)
		}
	}
	sort.Strings(types)

	var blocks []string
	for _, connectorType := range types {
		blocks = append(blocks, fmt.Sprintf("resource \"descope_%s\" \"sweep\" {\n  name = \"sweep-%s\"\n%s\n}\n",
			connectorType, strings.ReplaceAll(connectorType, "_", "-"), configs[connectorType]))
	}

	harness.RunScenario(t, harness.Scenario{
		Name:  "connector-sweep",
		Steps: []harness.Step{{Config: strings.Join(blocks, "\n")}},
	})
}

var (
	connectorTypePattern = regexp.MustCompile(`testacc\.NewResource\(t, "([a-z0-9_]+)"\)`)
	projectIDLinePattern = regexp.MustCompile(`(?m)^\s*project_id = .*$`)
)

func connectorSeedConfigs(t *testing.T) map[string]string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join("..", "..", "internal", "models", "connectors", "*_test.go"))
	if err != nil {
		t.Fatal(err)
	}

	configs := map[string]string{}
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		source := string(content)

		typeMatch := connectorTypePattern.FindStringSubmatch(source)
		if typeMatch == nil {
			continue
		}
		start := strings.Index(source, "c.Config(`")
		if start < 0 {
			continue
		}
		body := source[start+len("c.Config(`"):]
		end := strings.Index(body, "`)")
		if end < 0 {
			continue
		}
		body = projectIDLinePattern.ReplaceAllString(body[:end], "  project_id = descope_project.test.id")
		configs[typeMatch[1]] = strings.TrimRight(body, "\n\t ")
	}
	return configs
}

package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
)

var (
	terraformOnce sync.Once
	terraformRC   string
	terraformErr  error
)

func Env(t *testing.T) []string {
	t.Helper()
	terraformOnce.Do(func() {
		dir, err := os.MkdirTemp("", "tfexport-provider-*")
		if err != nil {
			terraformErr = err
			return
		}
		build := exec.Command("go", "build", "-o", filepath.Join(dir, "terraform-provider-descope"), "github.com/descope/terraform-provider-descope")
		if output, err := build.CombinedOutput(); err != nil {
			terraformErr = fmt.Errorf("provider build failed: %w\n%s", err, output)
			return
		}
		terraformRC = filepath.Join(dir, "terraform.tfrc")
		overrides := fmt.Sprintf("provider_installation {\n  dev_overrides {\n    %q = %q\n  }\n  direct {}\n}\n", "descope/descope", dir)
		terraformErr = os.WriteFile(terraformRC, []byte(overrides), 0o644)
	})
	if terraformErr != nil {
		t.Fatal(terraformErr)
	}
	return append(os.Environ(), "TF_CLI_CONFIG_FILE="+terraformRC, "TF_IN_AUTOMATION=1")
}

func Try(dir string, env []string, args ...string) ([]byte, error) {
	terraform, err := exec.LookPath("terraform")
	if err != nil {
		return nil, err
	}
	command := exec.Command(terraform, args...)
	command.Dir = dir
	command.Env = env
	return command.CombinedOutput()
}

func Run(t *testing.T, dir string, env []string, args ...string) []byte {
	t.Helper()
	terraform, err := exec.LookPath("terraform")
	if err != nil {
		t.Skip("terraform is not installed")
	}
	command := exec.Command(terraform, args...)
	command.Dir = dir
	command.Env = env
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("terraform %s failed: %s\n%s", strings.Join(args, " "), err, output)
	}
	return output
}

type Change struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Change  struct {
		Actions        []string        `json:"actions"`
		Before         map[string]any  `json:"before"`
		After          map[string]any  `json:"after"`
		AfterSensitive map[string]any  `json:"after_sensitive"`
		Importing      json.RawMessage `json:"importing"`
	} `json:"change"`
}

func PlanJSON(t *testing.T, dir string, env []string) []Change {
	t.Helper()
	planFile := filepath.Join(dir, "plan.bin")
	Run(t, dir, env, "plan", "-input=false", "-out="+planFile)

	terraform, err := exec.LookPath("terraform")
	if err != nil {
		t.Skip("terraform is not installed")
	}
	show := exec.Command(terraform, "show", "-json", planFile)
	show.Dir = dir
	show.Env = env
	output, err := show.Output()
	if err != nil {
		t.Fatalf("terraform show failed: %s", err)
	}

	var parsed struct {
		ResourceChanges []Change `json:"resource_changes"`
	}
	if err := json.Unmarshal(output, &parsed); err != nil {
		t.Fatalf("parsing plan JSON: %s", err)
	}
	return parsed.ResourceChanges
}

func AssertBenignPlan(t *testing.T, changes []Change) int {
	t.Helper()
	imported := 0
	for _, change := range changes {
		actions := strings.Join(change.Change.Actions, ",")
		if len(change.Change.Importing) > 0 {
			imported++
		}
		switch actions {
		case "no-op":
			continue
		case "update":
			for _, attribute := range ChangedAttributes(change.Change.Before, change.Change.After) {
				if attribute == "data" && (change.Type == "descope_flow" || change.Type == "descope_widget" || change.Type == "descope_styles") {
					continue // pretty-printed file vs compact imported state, settles on first apply
				}
				if onlySensitiveChanges(change.Change.Before[attribute], change.Change.After[attribute], change.Change.AfterSensitive[attribute]) {
					continue // dummy variable values for secrets the API never returns
				}
				if change.Change.Before[attribute] == nil {
					continue // the server never returned the attribute: the plan fills in the schema default
				}
				if compareConsistency(map[string]bool{}, change.Type, attribute, change.Change.After[attribute], change.Change.Before[attribute]) == consistencyZeroOmitted {
					continue // nested zero values the server omits, filled back in by the plan
				}
				t.Errorf("%s: unexpected change to attribute %q: %v -> %v", change.Address, attribute, change.Change.Before[attribute], change.Change.After[attribute])
			}
		default:
			t.Errorf("%s: unexpected plan actions %q", change.Address, actions)
		}
	}
	return imported
}

func AssertNoopPlan(t *testing.T, changes []Change) {
	t.Helper()
	for _, change := range changes {
		if actions := strings.Join(change.Change.Actions, ","); actions != "no-op" {
			t.Errorf("%s: expected no-op after apply, got %q", change.Address, actions)
		}
	}
}

func WriteVariables(t *testing.T, outDir, projectID string) {
	t.Helper()
	// the file is absent when the export declares no variables at all, which happens once the project supplies the id
	variables, err := os.ReadFile(filepath.Join(outDir, "variables.tf"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	// the project id is only a variable when the export didn't include the project resource itself
	var lines []string
	for _, match := range regexp.MustCompile(`variable "([^"]+)"`).FindAllStringSubmatch(string(variables), -1) {
		if match[1] == "project_id" {
			lines = append(lines, fmt.Sprintf("project_id = %q", projectID))
		} else {
			lines = append(lines, fmt.Sprintf("%s = %q", match[1], "dummy-secret-value"))
		}
	}
	if len(lines) == 0 {
		return
	}
	if err := os.WriteFile(filepath.Join(outDir, "terraform.tfvars"), []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func onlySensitiveChanges(before, after, sensitive any) bool {
	if isSensitive, _ := sensitive.(bool); isSensitive {
		return true
	}
	beforeJSON, _ := json.Marshal(before)
	afterJSON, _ := json.Marshal(after)
	if string(beforeJSON) == string(afterJSON) {
		return true
	}
	beforeMap, beforeOK := before.(map[string]any)
	afterMap, afterOK := after.(map[string]any)
	sensitiveMap, sensitiveOK := sensitive.(map[string]any)
	if !beforeOK || !afterOK || !sensitiveOK {
		return false
	}
	keys := map[string]bool{}
	for key := range beforeMap {
		keys[key] = true
	}
	for key := range afterMap {
		keys[key] = true
	}
	for key := range keys {
		if !onlySensitiveChanges(beforeMap[key], afterMap[key], sensitiveMap[key]) {
			return false
		}
	}
	return true
}

func ChangedAttributes(before, after map[string]any) []string {
	var changed []string
	for name, beforeValue := range before {
		afterValue, ok := after[name]
		if !ok {
			changed = append(changed, name)
			continue
		}
		beforeJSON, _ := json.Marshal(beforeValue)
		afterJSON, _ := json.Marshal(afterValue)
		if string(beforeJSON) != string(afterJSON) {
			changed = append(changed, name)
		}
	}
	for name := range after {
		if _, ok := before[name]; !ok {
			changed = append(changed, name)
		}
	}
	return changed
}

package harness

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func LoadScenarios(t *testing.T, dir string) []Scenario {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading scenarios: %s", err)
	}

	var scenarios []Scenario
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		steps := loadSteps(t, filepath.Join(dir, entry.Name()))
		if len(steps) == 0 {
			t.Fatalf("scenario %q has no steps", entry.Name())
		}
		scenarios = append(scenarios, Scenario{Name: entry.Name(), Steps: steps})
	}
	if len(scenarios) == 0 {
		t.Fatalf("no scenarios found in %s", dir)
	}
	return scenarios
}

func loadSteps(t *testing.T, dir string) []Step {
	t.Helper()

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading scenario %s: %s", dir, err)
	}

	numbers := make([]int, 0, len(files))
	configs := map[int]string{}
	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".tf") {
			continue
		}
		number, err := strconv.Atoi(strings.TrimSuffix(name, ".tf"))
		if err != nil {
			t.Fatalf("scenario step %s must be named <number>.tf", filepath.Join(dir, name))
		}
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("reading scenario step: %s", err)
		}
		numbers = append(numbers, number)
		configs[number] = string(content)
	}
	sort.Ints(numbers)

	steps := make([]Step, 0, len(numbers))
	for i, number := range numbers {
		if number != i+1 {
			t.Fatalf("scenario %s steps must be numbered consecutively from 1, found %d", dir, number)
		}
		steps = append(steps, Step{Config: configs[number]})
	}
	return steps
}

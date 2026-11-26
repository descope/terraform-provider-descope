package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/harness"
)

func TestRoundTripFuzz(t *testing.T) {
	if os.Getenv("TF_ACC") == "" || os.Getenv("DESCOPE_TFEXPORT_ROUNDTRIP") == "" {
		t.Skip("set TF_ACC and DESCOPE_TFEXPORT_ROUNDTRIP to run this suite")
	}
	if os.Getenv("DESCOPE_MANAGEMENT_KEY") == "" {
		t.Skip("set DESCOPE_MANAGEMENT_KEY to run this suite")
	}

	seed := int64(1)
	if value := os.Getenv("DESCOPE_TFEXPORT_FUZZ_SEED"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			t.Fatalf("invalid DESCOPE_TFEXPORT_FUZZ_SEED: %s", err)
		}
		seed = parsed
	}
	t.Logf("fuzzing with seed %d", seed)
	harness.RunScenario(t, fuzzScenario(seed))
}

var (
	fuzzWords = []string{
		"alpha", "béta", "गामा", "デルタ", "epsilon", "Ζήτα", "эта", "θήτα", "iota-🎯",
		"kappa", "λάμδα", "mu", "nu", "csi", "ómicron", "pi", "rho", "sígma",
	}
	fuzzPunctuation = []string{" ", " — ", ": ", "; ", " / ", "'", "\"", "\\", "\n", "\t", " & ", " <b> ", " $5 ", " 100% "}
	fuzzUnits       = []string{"minutes", "hours", "days", "weeks"}
)

func fuzzScenario(seed int64) harness.Scenario {
	generator := rand.New(rand.NewSource(seed))
	return harness.Scenario{
		Name:  fmt.Sprintf("fuzz-%d", seed),
		Steps: []harness.Step{{Config: fuzzConfig(generator)}, {Config: fuzzConfig(generator)}},
	}
}

func fuzzConfig(generator *rand.Rand) string {
	word := func() string { return fuzzWords[generator.Intn(len(fuzzWords))] }
	phrase := func(words int) string {
		parts := make([]string, 0, words*2)
		for i := 0; i < words; i++ {
			if i > 0 {
				parts = append(parts, fuzzPunctuation[generator.Intn(len(fuzzPunctuation))])
			}
			parts = append(parts, word())
		}
		return strings.Join(parts, "")
	}
	duration := func() string {
		return fmt.Sprintf("%d %s", 1+generator.Intn(50), fuzzUnits[generator.Intn(len(fuzzUnits))])
	}

	var blocks []string
	block := func(format string, args ...any) {
		blocks = append(blocks, fmt.Sprintf(format, args...))
	}

	roles := 2 + generator.Intn(4)
	for i := 0; i < roles; i++ {
		block(`
resource "descope_permission" "fuzz_%d" {
  project_id  = descope_project.test.id
  name        = "fuzz.perm.%d"
  description = %s
}

resource "descope_role" "fuzz_%d" {
  project_id  = descope_project.test.id
  name        = "fuzz-role-%d"
  description = %s
  permissions = [descope_permission.fuzz_%d.name]
}`, i, i, hclQuote(phrase(1+generator.Intn(4))), i, i, hclQuote(phrase(1+generator.Intn(4))), i)
	}

	lists := 1 + generator.Intn(3)
	for i := 0; i < lists; i++ {
		items := make([]string, 1+generator.Intn(4))
		seen := map[string]bool{}
		for j := range items {
			items[j] = fmt.Sprintf("%s-%d-%d", word(), i, j)
			for seen[items[j]] {
				items[j] += "x"
			}
			seen[items[j]] = true
		}
		quoted := make([]string, len(items))
		for j, item := range items {
			quoted[j] = hclQuote(item)
		}
		block(`
resource "descope_list" "fuzz_%d" {
  project_id = descope_project.test.id
  name       = "fuzz-list-%d"
  texts      = [%s]
}`, i, i, strings.Join(quoted, ", "))
	}

	block(`
resource "descope_magiclink_settings" "fuzz" {
  project_id      = descope_project.test.id
  expiration_time = %q
  redirect_url    = "https://example.com/fuzz/%d"
}

resource "descope_otp_settings" "fuzz" {
  project_id      = descope_project.test.id
  expiration_time = %q
  domain          = "fuzz-%d.example.com"
}

resource "descope_jwt_template" "fuzz" {
  project_id = descope_project.test.id
  type       = "user"
  name       = "fuzz-jwt"
  template = jsonencode({
    phrase = %s
    number = %d
  })
}`, duration(), generator.Intn(1000), duration(), generator.Intn(1000), hclQuote(phrase(2)), generator.Intn(100000))

	return strings.Join(blocks, "\n")
}

func hclQuote(s string) string {
	quoted, _ := json.Marshal(s)
	escaped := strings.ReplaceAll(string(quoted), "${", "$${")
	escaped = strings.ReplaceAll(escaped, "%{", "%%{")
	return escaped
}

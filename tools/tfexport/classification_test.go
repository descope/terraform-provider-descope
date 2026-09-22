package main

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var secretLikeName = regexp.MustCompile(`(^|_)(key|token|secret|password|passphrase|credential|credentials|cert|certificate)$`)

var publicByDesign = map[string]bool{
	"public_key":       true,
	"public_api_key":   true,
	"site_key":         true,
	"project_token":    true, // mixpanel project tokens are client-side identifiers
	"public_write_key": true, // the segment browser library needs this one in the page, unlike write_key
	"key":              true, // generic map key fields, not credentials
}

// knownMisclassified carry real secret material but are not marked sensitive - each entry is an open bug that must go once its template is fixed.
var knownMisclassified = map[string]string{
	"descope_eight_by_eight_viber_connector.api_key":    "8x8 api keys are bearer credentials",
	"descope_eight_by_eight_whatsapp_connector.api_key": "8x8 api keys are bearer credentials",
}

func TestSecretClassification(t *testing.T) {
	ctx := context.Background()
	stillMisclassified := map[string]bool{}

	for resourceType, exportable := range registry.Load(ctx) {
		for name, attribute := range exportable.ExportSchema().Attributes {
			if _, isString := attribute.(schema.StringAttribute); !isString {
				continue
			}
			if !secretLikeName.MatchString(name) || attribute.IsSensitive() {
				continue
			}
			if publicByDesign[name] {
				continue
			}
			key := resourceType + "." + name
			if reason, known := knownMisclassified[key]; known {
				stillMisclassified[key] = true
				t.Logf("KNOWN misclassification %s (%s)", key, reason)
				continue
			}
			t.Errorf("attribute %s has a secret-like name but is not marked sensitive; classify it or add it to an allowlist with a reason", key)
		}
	}

	for key := range knownMisclassified {
		if !stillMisclassified[key] && !strings.Contains(key, "|skip|") {
			t.Errorf("knownMisclassified entry %s no longer applies - remove it and close the finding", key)
		}
	}
}

package registry_test

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/descope/terraform-provider-descope/internal/resources"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/registry"
)

func TestAllResourcesExportable(t *testing.T) {
	ctx := context.Background()

	total := 0
	notExportable := 0
	for _, constructor := range provider.NewDescopeProvider("test")().Resources(ctx) {
		total++
		if _, ok := constructor().(resources.ExportableResource); !ok {
			notExportable++
		}
	}
	if notExportable != 0 {
		t.Errorf("expected every resource to be exportable, got %d of %d that aren't", notExportable, total)
	}

	loaded := registry.Load(ctx)
	if len(loaded) != total-notExportable {
		t.Errorf("registry has %d resources, expected %d", len(loaded), total-notExportable)
	}

	for _, name := range []string{
		"descope_project", "descope_flow", "descope_widget", "descope_styles", "descope_magiclink_settings",
		"descope_email_template", "descope_text_template", "descope_smtp_connector",
		"descope_role", "descope_oidc_app", "descope_app_role", "descope_jwt_template",
	} {
		exportable, ok := loaded[name]
		if !ok {
			t.Errorf("expected %s in the registry", name)
			continue
		}
		if len(exportable.ExportSchema().Attributes) == 0 {
			t.Errorf("expected %s to have a non-empty schema", name)
		}
	}

	for name, exportable := range loaded {
		if name != "descope_"+exportable.ExportName() {
			t.Errorf("registry key %q doesn't match resource name %q", name, exportable.ExportName())
		}
	}
}

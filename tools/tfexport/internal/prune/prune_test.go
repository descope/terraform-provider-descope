package prune_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func object(ctx context.Context, t *testing.T, sc schema.Schema, values map[string]any) types.Object {
	t.Helper()
	objectType, ok := sc.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatal("unexpected schema type")
	}
	value, err := tfValue(objectType, values)
	if err != nil {
		t.Fatal(err)
	}
	converted, err := sc.Type().ValueFromTerraform(ctx, value)
	if err != nil {
		t.Fatal(err)
	}
	result, ok := converted.(types.Object)
	if !ok {
		t.Fatal("unexpected value type")
	}
	return result
}

func tfValue(typ tftypes.Type, value any) (tftypes.Value, error) {
	if value == nil {
		return tftypes.NewValue(typ, nil), nil
	}
	if typ.Is(tftypes.String) || typ.Is(tftypes.Bool) || typ.Is(tftypes.Number) {
		if n, ok := value.(int); ok {
			value = float64(n)
		}
		return tftypes.NewValue(typ, value), nil
	}
	switch t := typ.(type) {
	case tftypes.Object:
		fields := map[string]tftypes.Value{}
		object, _ := value.(map[string]any)
		for name, fieldType := range t.AttributeTypes {
			field, err := tfValue(fieldType, object[name])
			if err != nil {
				return tftypes.Value{}, err
			}
			fields[name] = field
		}
		return tftypes.NewValue(typ, fields), nil
	case tftypes.List, tftypes.Set:
		elementType := elementTypeOf(typ)
		items, _ := value.([]any)
		var elements []tftypes.Value
		for _, item := range items {
			element, err := tfValue(elementType, item)
			if err != nil {
				return tftypes.Value{}, err
			}
			elements = append(elements, element)
		}
		return tftypes.NewValue(typ, elements), nil
	case tftypes.Map:
		items, _ := value.(map[string]any)
		elements := map[string]tftypes.Value{}
		for key, item := range items {
			element, err := tfValue(t.ElementType, item)
			if err != nil {
				return tftypes.Value{}, err
			}
			elements[key] = element
		}
		return tftypes.NewValue(typ, elements), nil
	}
	return tftypes.Value{}, fmt.Errorf("unsupported type %s for value %v", typ, value)
}

func elementTypeOf(typ tftypes.Type) tftypes.Type {
	switch t := typ.(type) {
	case tftypes.List:
		return t.ElementType
	case tftypes.Set:
		return t.ElementType
	}
	return nil
}

func names(attrs []prune.Attr) []string {
	result := make([]string, 0, len(attrs))
	for _, a := range attrs {
		result = append(result, a.Name)
	}
	return result
}

func find(attrs []prune.Attr, name string) *prune.Attr {
	for i := range attrs {
		if attrs[i].Name == name {
			return &attrs[i]
		}
	}
	return nil
}

func TestPrune(t *testing.T) {
	ctx := context.Background()
	loaded := registry.Load(ctx)

	magiclink := loaded["descope_magiclink_settings"].ExportSchema()
	amplitude := loaded["descope_amplitude_connector"].ExportSchema()
	adminportal := loaded["descope_admin_portal"].ExportSchema()

	t.Run("NestedObjectsPruneToTheirNonDefaultFields", func(t *testing.T) {
		sso := loaded["descope_sso_settings"].ExportSchema()
		obj := object(ctx, t, sso, map[string]any{
			"id":                 "P123",
			"project_id":         "P123",
			"sso_suite_settings": map[string]any{"hide_sso": true},
		})
		attrs, _ := prune.Object(ctx, sso.Attributes, obj, true)
		suite := find(attrs, "sso_suite_settings")
		if suite == nil || !suite.IsNested {
			t.Fatalf("expected sso_suite_settings to survive as a nested object, got %+v", suite)
		}
		if len(suite.Nested) != 1 || suite.Nested[0].Name != "hide_sso" {
			t.Errorf("expected only hide_sso to survive inside it, got %+v", suite.Nested)
		}
	})

	t.Run("AllDefaultsPruneToNothing", func(t *testing.T) {
		obj := object(ctx, t, magiclink, map[string]any{
			"id":                 "P123",
			"project_id":         "P123",
			"disabled":           false,
			"expiration_time":    "3 minutes",
			"redirect_url":       "",
			"email_template_id":  "",
			"text_template_id":   "",
			"email_connector_id": "",
			"text_connector_id":  "",
		})
		attrs, secrets := prune.Object(ctx, magiclink.Attributes, obj, true)
		if prune.Meaningful(attrs) {
			t.Errorf("expected no meaningful attributes, got %v", names(attrs))
		}
		for _, a := range attrs {
			if !a.Verbatim {
				t.Errorf("expected only verbatim attributes to survive, got %+v", a)
			}
		}
		if len(secrets) != 0 {
			t.Errorf("expected no secrets, got %v", secrets)
		}
	})

	t.Run("NonDefaultsSurvive", func(t *testing.T) {
		obj := object(ctx, t, magiclink, map[string]any{
			"id":                 "P123",
			"project_id":         "P123",
			"disabled":           false,
			"expiration_time":    "5 minutes",
			"redirect_url":       "https://example.com",
			"email_connector_id": "CI123",
		})
		attrs, _ := prune.Object(ctx, magiclink.Attributes, obj, true)
		if len(attrs) != 3 {
			t.Fatalf("expected 3 surviving attributes, got %v", names(attrs))
		}
		for _, name := range []string{"expiration_time", "redirect_url", "email_connector_id"} {
			if find(attrs, name) == nil {
				t.Errorf("expected %s to survive, got %v", name, names(attrs))
			}
		}
	})

	t.Run("NullsPrune", func(t *testing.T) {
		obj := object(ctx, t, magiclink, map[string]any{"id": "P123", "project_id": "P123"})
		attrs, _ := prune.Object(ctx, magiclink.Attributes, obj, true)
		if len(attrs) != 0 {
			t.Errorf("expected all null attributes to prune, got %v", names(attrs))
		}
	})

	t.Run("RequiredSensitiveBecomesPlaceholder", func(t *testing.T) {
		obj := object(ctx, t, amplitude, map[string]any{
			"id":         "CI123",
			"project_id": "P123",
			"name":       "my connector",
			"server_url": "https://example.com",
		})
		attrs, secrets := prune.Object(ctx, amplitude.Attributes, obj, true)
		key := find(attrs, "api_key")
		if key == nil || !key.Placeholder {
			t.Fatalf("expected api_key placeholder, got %v", names(attrs))
		}
		if len(secrets) != 1 || secrets[0].Path != "api_key" || !secrets[0].Required {
			t.Errorf("expected a required api_key secret, got %v", secrets)
		}
		if find(attrs, "name") == nil || find(attrs, "server_url") == nil {
			t.Errorf("expected name and server_url to survive, got %v", names(attrs))
		}
	})

	t.Run("NestedListElementFieldsSurvive", func(t *testing.T) {
		obj := object(ctx, t, adminportal, map[string]any{
			"id":         "P123",
			"project_id": "P123",
			"widgets": []any{
				map[string]any{"widget_id": "user-management", "type": "user-management"},
				map[string]any{"widget_id": "audit-management", "type": "audit-management"},
			},
		})
		attrs, _ := prune.Object(ctx, adminportal.Attributes, obj, true)
		widgets := find(attrs, "widgets")
		if widgets == nil || !widgets.IsElements || len(widgets.Elements) != 2 {
			t.Fatalf("expected widgets with 2 elements, got %+v", attrs)
		}
		for i, element := range widgets.Elements {
			if fields := names(element); len(fields) != 2 { // both fields are required and always survive
				t.Errorf("expected widget %d to keep widget_id and type, got %v", i, fields)
			}
		}
	})

	t.Run("EmptyNestedListPrunesToDefault", func(t *testing.T) {
		obj := object(ctx, t, adminportal, map[string]any{
			"id":         "P123",
			"project_id": "P123",
			"widgets":    []any{},
		})
		attrs, _ := prune.Object(ctx, adminportal.Attributes, obj, true)
		if len(attrs) != 0 {
			t.Errorf("expected empty widgets list to prune as the default, got %v", names(attrs))
		}
	})
}

func TestPruneGeneratedSecrets(t *testing.T) {
	ctx := context.Background()
	loaded := registry.Load(ctx)
	oidc := loaded["descope_oidc_app"].ExportSchema()

	secretFor := func(secrets []prune.Secret, path string) (prune.Secret, bool) {
		for _, secret := range secrets {
			if secret.Path == path {
				return secret, true
			}
		}
		return prune.Secret{}, false
	}

	t.Run("EmptyGeneratedSecretIsNotDropped", func(t *testing.T) {
		// the read blanks a generated client_secret, so nothing is lost by omitting it and no warning is warranted
		obj := object(ctx, t, oidc, map[string]any{
			"id": "APP1", "project_id": "P123", "name": "app", "client_secret": "",
		})
		_, secrets := prune.Object(ctx, oidc.Attributes, obj, true)
		secret, ok := secretFor(secrets, "client_secret")
		if !ok {
			t.Fatalf("expected client_secret to be recorded, got %v", secrets)
		}
		if secret.Dropped {
			t.Error("an empty generated secret should not be reported as dropped")
		}
		if secret.Required {
			t.Error("an empty generated secret should not be promoted on its own")
		}
	})

	t.Run("RealGeneratedSecretIsDropped", func(t *testing.T) {
		obj := object(ctx, t, oidc, map[string]any{
			"id": "APP1", "project_id": "P123", "name": "app", "client_secret": "actual-value",
		})
		_, secrets := prune.Object(ctx, oidc.Attributes, obj, true)
		secret, ok := secretFor(secrets, "client_secret")
		if !ok {
			t.Fatalf("expected client_secret to be recorded, got %v", secrets)
		}
		if !secret.Dropped {
			t.Error("a secret the configuration cannot carry should be reported as dropped")
		}
	})
}

func TestPruneDurations(t *testing.T) {
	ctx := context.Background()
	session := registry.Load(ctx)["descope_session_settings"].ExportSchema()

	// an equal duration spelled differently still has to be emitted: an import keeps the stored spelling, so omitting it would plan a change
	for _, test := range []struct {
		value  string
		pruned bool
	}{
		{value: "4 weeks", pruned: true},
		{value: "4 week", pruned: false},
		{value: "28 days", pruned: false},
		{value: "3 weeks", pruned: false},
	} {
		t.Run(test.value, func(t *testing.T) {
			obj := object(ctx, t, session, map[string]any{
				"id": "P123", "project_id": "P123", "refresh_token_expiration": test.value,
			})
			attrs, _ := prune.Object(ctx, session.Attributes, obj, true)
			_, survived := findAttr(attrs, "refresh_token_expiration")
			if survived == test.pruned {
				t.Errorf("value %q: survived=%v, expected pruned=%v", test.value, survived, test.pruned)
			}
		})
	}
}

func findAttr(attrs []prune.Attr, name string) (prune.Attr, bool) {
	for _, attr := range attrs {
		if attr.Name == name {
			return attr, true
		}
	}
	return prune.Attr{}, false
}

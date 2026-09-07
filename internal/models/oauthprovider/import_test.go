package oauthprovider

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSetValuesLeavesIDUntouched(t *testing.T) {
	var diags diag.Diagnostics
	handler := helpers.NewHandler(context.Background(), &diags)

	model := &OAuthProviderModel{ID: stringattr.Value("google")}
	model.SetValues(handler, map[string]any{"id": "google", "enabled": true, "clientId": "cid"})

	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags.Errors())
	}
	if actual := model.ID.ValueString(); actual != "google" {
		t.Errorf("id was overwritten with %q", actual)
	}
	if actual := model.ClientID.ValueString(); actual != "cid" {
		t.Errorf("client_id is %q, expected %q", actual, "cid")
	}
}

func TestValuesRunsResolvedValidation(t *testing.T) {
	var diags diag.Diagnostics
	handler := helpers.NewHandler(context.Background(), &diags)

	model := &OAuthProviderModel{ID: stringattr.Value("google"), Scopes: strlistattr.Value([]string{"email"}),
		Prompts: strlistattr.Value(nil), AllowedGrantTypes: strlistattr.Value(nil)}
	model.Values(handler)
	if !diags.HasError() {
		t.Fatalf("expected scopes without credentials to fail on a system provider")
	}

	diags = diag.Diagnostics{}
	handler = helpers.NewHandler(context.Background(), &diags)
	model.ClientID = stringattr.Value("cid")
	model.ClientSecret = stringattr.Value("secret")
	model.Values(handler)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags.Errors())
	}
}

func TestSystemCasingValidator(t *testing.T) {
	for _, test := range []struct {
		value string
		valid bool
	}{
		{"google", true},
		{"apple", true},
		{"my-custom-provider", true},
		{"My Custom Provider", true}, // custom names keep their casing
		{"Google", false},
		{"GOOGLE", false},
		{"Apple", false},
	} {
		req := validator.StringRequest{
			Path:        path.Root("id"),
			ConfigValue: types.StringValue(test.value),
		}
		resp := &validator.StringResponse{}
		systemCasingValidator{}.ValidateString(context.Background(), req, resp)

		if valid := !resp.Diagnostics.HasError(); valid != test.valid {
			t.Errorf("id %q: valid is %v, expected %v", test.value, valid, test.valid)
		}
	}
}

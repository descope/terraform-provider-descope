package export

import (
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/discover"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/read"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func provider(t *testing.T, id, redirectURL string, builtin bool) read.Result {
	t.Helper()
	object, diags := types.ObjectValue(map[string]attr.Type{"redirect_url": types.StringType}, map[string]attr.Value{"redirect_url": types.StringValue(redirectURL)})
	if diags.HasError() {
		t.Fatal(diags)
	}
	return read.Result{Instance: discover.Instance{Resource: "descope_oauth_provider", ID: id, Builtin: builtin}, Object: object}
}

func TestBuiltinRedirectURL(t *testing.T) {
	const standard, custom = "https://api.descope.com/v1/auth/oauth/callback", "https://auth.example.com/v1/auth/oauth/callback"

	for _, test := range []struct {
		name    string
		results []read.Result
		want    string
	}{
		{"SharedByMost", []read.Result{provider(t, "github", standard, true), provider(t, "google", custom, true), provider(t, "slack", standard, true)}, standard},
		{"CustomProvidersDontCount", []read.Result{provider(t, "github", standard, true), provider(t, "a", custom, false), provider(t, "b", custom, false), provider(t, "slack", standard, true)}, standard},
		{"TieHasNoDefault", []read.Result{provider(t, "github", standard, true), provider(t, "google", custom, true)}, ""},
		{"SingleProviderHasNoDefault", []read.Result{provider(t, "github", standard, true)}, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := builtinRedirectURL(test.results); got != test.want {
				t.Errorf("expected %q, got %q", test.want, got)
			}
		})
	}
}

func TestBuiltinCustomization(t *testing.T) {
	const standard = "https://api.descope.com/v1/auth/oauth/callback"
	id := prune.Attr{Name: "id", Value: types.StringValue("github")}

	for _, test := range []struct {
		name         string
		attrs        []prune.Attr
		redirectURL  string
		uncustomized bool
	}{
		{"OnlyIdentity", []prune.Attr{id}, standard, true},
		{"DefaultRedirect", []prune.Attr{id, {Name: "redirect_url", Value: types.StringValue(standard)}}, standard, true},
		{"UnknownDefaultRedirect", []prune.Attr{id, {Name: "redirect_url", Value: types.StringValue(standard)}}, "", false},
		{"CustomRedirect", []prune.Attr{id, {Name: "redirect_url", Value: types.StringValue("https://auth.example.com/cb")}}, standard, false},
		{"Disabled", []prune.Attr{id, {Name: "disabled", Value: types.BoolValue(true)}}, standard, false},
		{"OwnCredentials", []prune.Attr{id, {Name: "client_id", Value: types.StringValue("client")}}, standard, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := uncustomized(withoutBuiltinRedirectURL(test.attrs, test.redirectURL)); got != test.uncustomized {
				t.Errorf("expected uncustomized to be %t, got %t", test.uncustomized, got)
			}
		})
	}

	attrs := []prune.Attr{id, {Name: "client_id", Value: types.StringValue("client")}, {Name: "redirect_url", Value: types.StringValue(standard)}}
	if kept := withoutBuiltinRedirectURL(attrs, standard); len(kept) != 2 || len(attrs) != 3 {
		t.Errorf("expected only the default redirect_url to be dropped from a copy, got %+v from %+v", kept, attrs)
	}
}

package inboundapp_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/inboundapp"
	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestInboundAppClientTypeReplacement(t *testing.T) {
	testCases := []struct {
		name    string
		state   string
		plan    string
		replace bool
	}{
		{"legacy to confidential updates in place", "", "confidential", false},
		{"confidential to legacy replaces", "confidential", "", true},
		{"legacy to public replaces", "", "public", true},
		{"confidential to public replaces", "confidential", "public", true},
		{"public to confidential replaces", "public", "confidential", true},
		{"public to legacy replaces", "public", "", true},
		{"unchanged public is kept", "public", "public", false},
	}
	attr, ok := inboundapp.InboundAppAttributes["client_type"].(schema.StringAttribute)
	if !ok {
		t.Fatalf("expected client_type to be a string attribute")
	}
	raw := tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{})
	modifyPlan := func(state, config, planned types.String) (types.String, bool) {
		req := planmodifier.StringRequest{State: tfsdk.State{Raw: raw}, Plan: tfsdk.Plan{Raw: raw}, StateValue: state, ConfigValue: config, PlanValue: planned}
		replace := false
		for _, modifier := range attr.PlanModifiers {
			resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
			modifier.PlanModifyString(context.Background(), req, resp)
			req.PlanValue = resp.PlanValue
			replace = replace || resp.RequiresReplace
		}
		return req.PlanValue, replace
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, got := modifyPlan(types.StringValue(tc.state), types.StringValue(tc.plan), types.StringValue(tc.plan)); got != tc.replace {
				t.Errorf("expected requires replace to be %t, got %t", tc.replace, got)
			}
		})
	}
	t.Run("unknown config value replaces", func(t *testing.T) {
		if _, replace := modifyPlan(types.StringValue("confidential"), types.StringUnknown(), types.StringUnknown()); !replace {
			t.Errorf("expected requires replace when the planned client type is unknown")
		}
	})
	for name, state := range map[string]string{"legacy": "", "confidential": "confidential", "public": "public"} {
		t.Run("omitted value keeps "+name+" in place", func(t *testing.T) {
			planned, replace := modifyPlan(types.StringValue(state), types.StringNull(), types.StringUnknown())
			if replace || !planned.Equal(types.StringValue(state)) {
				t.Errorf("expected the stored client type %q to be kept in place, got %s with requires replace %t", state, planned, replace)
			}
		})
	}
	t.Run("omitted value on create is left to the server", func(t *testing.T) {
		if attr.Default != nil || !attr.Computed {
			t.Errorf("expected client_type to be computed with no default")
		}
	})
}

func TestInboundAppClientTypeValues(t *testing.T) {
	h := helpers.NewHandler(context.Background(), &diag.Diagnostics{})
	for _, value := range []tftypes.Value{tftypes.NewValue(tftypes.String, tftypes.UnknownValue), tftypes.NewValue(tftypes.String, nil)} {
		m := newInboundAppModel(t, value)
		if v, ok := m.Values(h)["clientType"]; ok {
			t.Errorf("expected a %s client type not to be sent, got %v", value, v)
		}
	}
	m := newInboundAppModel(t, tftypes.NewValue(tftypes.String, "public"))
	if v := m.Values(h)["clientType"]; v != "public" {
		t.Errorf("expected the client type to be sent as public, got %v", v)
	}
}

func TestInboundAppClientTypeSetValues(t *testing.T) {
	h := helpers.NewHandler(context.Background(), &diag.Diagnostics{})
	testCases := []struct {
		name     string
		data     map[string]any
		expected string
	}{
		{"server type is adopted", map[string]any{"clientType": "confidential"}, "confidential"},
		{"omitted server type is legacy", map[string]any{}, ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := newInboundAppModel(t, tftypes.NewValue(tftypes.String, tftypes.UnknownValue))
			m.SetValues(h, tc.data)
			if !m.ClientType.Equal(types.StringValue(tc.expected)) {
				t.Errorf("expected client type %q, got %s", tc.expected, m.ClientType)
			}
		})
	}
}

func TestInboundAppGrantTypesModifyPlan(t *testing.T) {
	str := func(v any) tftypes.Value { return tftypes.NewValue(tftypes.String, v) }
	testCases := []struct {
		name       string
		configType tftypes.Value
		planType   tftypes.Value
		configured types.Bool
		expected   types.Bool
		conflict   bool
	}{
		{"public defaults to disabled", str("public"), str("public"), types.BoolNull(), types.BoolValue(true), false},
		{"public keeps the configured value", str("public"), str("public"), types.BoolValue(true), types.BoolValue(true), false},
		{"confidential keeps the default", str("confidential"), str("confidential"), types.BoolNull(), types.BoolValue(false), false},
		{"legacy keeps the default", str(""), str(""), types.BoolNull(), types.BoolValue(false), false},
		{"omitted type on create is known after apply", str(nil), str(tftypes.UnknownValue), types.BoolNull(), types.BoolUnknown(), false},
		{"unknown type is known after apply", str(tftypes.UnknownValue), str(tftypes.UnknownValue), types.BoolNull(), types.BoolUnknown(), false},
		{"unknown type keeps the configured value", str(tftypes.UnknownValue), str(tftypes.UnknownValue), types.BoolValue(false), types.BoolValue(false), false},
		{"stored public defaults to disabled", str(nil), str("public"), types.BoolNull(), types.BoolValue(true), false},
		{"stored public cannot enable the grant", str(nil), str("public"), types.BoolValue(false), types.BoolValue(false), true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := newInboundAppModel(t, tc.configType)
			config.ClientCredentialsDisabled = tc.configured
			config.JWTBearerDisabled = tc.configured
			planned := types.BoolValue(false) // the attribute default when not configured
			if !tc.configured.IsNull() {
				planned = tc.configured
			}
			plan := newInboundAppModel(t, tc.planType)
			plan.ClientCredentialsDisabled = planned
			plan.JWTBearerDisabled = planned
			plan.RefreshTokenDisabled = types.BoolValue(false)
			diags := diag.Diagnostics{}
			plan.ModifyPlan(helpers.NewHandler(context.Background(), &diags), config, newInboundAppModel(t, str(nil)))
			if diags.HasError() != tc.conflict {
				t.Errorf("expected a conflict to be %t, got %v", tc.conflict, diags)
			}
			if !plan.ClientCredentialsDisabled.Equal(tc.expected) || !plan.JWTBearerDisabled.Equal(tc.expected) {
				t.Errorf("expected client_credentials_disabled and jwt_bearer_disabled to be planned as %s, got %s and %s", tc.expected, plan.ClientCredentialsDisabled, plan.JWTBearerDisabled)
			}
			if !plan.RefreshTokenDisabled.Equal(types.BoolValue(false)) {
				t.Errorf("expected refresh_token_disabled to be left alone, got %s", plan.RefreshTokenDisabled)
			}
		})
	}
}

func TestInboundAppGrantTypesValidate(t *testing.T) {
	testCases := []struct {
		name       string
		clientType string
		configured types.Bool
		conflict   bool
	}{
		{"public cannot enable the grant", "public", types.BoolValue(false), true},
		{"public can disable the grant", "public", types.BoolValue(true), false},
		{"public can omit the grant", "public", types.BoolNull(), false},
		{"public allows an unknown value", "public", types.BoolUnknown(), false},
		{"confidential can enable the grant", "confidential", types.BoolValue(false), false},
	}
	for _, tc := range testCases {
		for _, field := range []string{"client_credentials_disabled", "jwt_bearer_disabled"} {
			t.Run(tc.name+" "+field, func(t *testing.T) {
				m := newInboundAppModel(t, tftypes.NewValue(tftypes.String, tc.clientType))
				if field == "client_credentials_disabled" {
					m.ClientCredentialsDisabled = tc.configured
				} else {
					m.JWTBearerDisabled = tc.configured
				}
				diags := diag.Diagnostics{}
				m.Validate(helpers.NewHandler(context.Background(), &diags))
				if diags.HasError() != tc.conflict {
					t.Errorf("expected a conflict to be %t, got %v", tc.conflict, diags)
				}
			})
		}
	}
}

func TestInboundAppGrantTypesValues(t *testing.T) {
	h := helpers.NewHandler(context.Background(), &diag.Diagnostics{})
	keys := []string{"authorizationCodeDisabled", "clientCredentialsDisabled", "refreshTokenDisabled", "jwtBearerDisabled", "deviceCodeDisabled", "tokenExchangeDisabled"}
	m := newInboundAppModel(t, tftypes.NewValue(tftypes.String, "confidential"))
	data := map[string]any{}
	for _, key := range keys {
		data[key] = true
	}
	m.SetValues(h, data)
	values := m.Values(h)
	for _, key := range keys {
		if values[key] != true {
			t.Errorf("expected %s to round trip as true, got %v", key, values[key])
		}
	}
}

func newInboundAppModel(t *testing.T, clientType tftypes.Value) *inboundapp.InboundAppModel {
	ctx := context.Background()
	objType, ok := inboundapp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("expected the inbound app schema to be an object")
	}
	values := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(typ, nil)
	}
	values["client_type"] = clientType
	m := &inboundapp.InboundAppModel{}
	if diags := (tfsdk.Plan{Schema: inboundapp.Schema, Raw: tftypes.NewValue(objType, values)}).Get(ctx, m); diags.HasError() {
		t.Fatalf("failed to build the inbound app model: %v", diags)
	}
	return m
}

func TestInboundAppDeletionProtection(t *testing.T) {
	p := testacc.Project(t)
	a := testacc.InboundApp(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
			`),
			Check: a.Check(map[string]any{
				"deletion_protection": testacc.AttributeIsNotSet,
			}),
		},
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
			`),
			Destroy:     true,
			ExpectError: regexp.MustCompile(`Deletion Protection Enabled`),
		},
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "public"
			`),
			ExpectError: regexp.MustCompile(`Deletion Protection Enabled`),
		},
		// The flag change must be applied on its own before the replacement is allowed
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
			`),
		},
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				deletion_protection = false
				client_type = "public"
			`),
			Check: a.Check(map[string]any{
				"client_type":         "public",
				"deletion_protection": "false",
			}),
		},
	)
}

func TestInboundApp(t *testing.T) {
	p := testacc.Project(t)
	a := testacc.InboundApp(t)
	j := testacc.JWTTemplate(t)
	j2 := testacc.JWTTemplate(t)
	j2.ID = "other"
	testacc.Run(t,
		// Test basic creation with the required fields and a client type
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
			`),
			Check: a.Check(map[string]any{
				"id":                               testacc.AttributeIsSet,
				"project_id":                       testacc.AttributeIsSet,
				"name":                             a.Name,
				"description":                      "",
				"client_type":                      "confidential",
				"client_id":                        testacc.AttributeIsSet,
				"client_secret":                    testacc.AttributeIsSet,
				"force_pkce":                       "false",
				"approved_callback_urls.#":         "0",
				"permissions_scopes.#":             "0",
				"attributes_scopes.#":              "0",
				"connections_scopes.#":             "0",
				"audience_whitelist.#":             "0",
				"force_add_all_authorization_info": "false",
				"force_dpop":                       "false",
				"default_audience":                 "",
			}),
		},
		// Test update with description and callback URLs
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				description = "Updated description"
				force_pkce = true
				approved_callback_urls = ["https://example.com/callback"]
				permissions_scopes = [
					{
						name = "openid"
						description = "Foo"
					},
					{
						name = "email"
						description = "Access email"
						optional = true
					},
				]
			`),
			Check: a.Check(map[string]any{
				"description":                      "Updated description",
				"force_pkce":                       "true",
				"approved_callback_urls.#":         "1",
				"approved_callback_urls.0":         "https://example.com/callback",
				"permissions_scopes.#":             "2",
				"permissions_scopes.0.name":        "openid",
				"permissions_scopes.1.name":        "email",
				"permissions_scopes.1.description": "Access email",
				"permissions_scopes.1.optional":    "true",
			}),
		},
		// Test update with session settings
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				session_settings = {
					enabled = true
					refresh_token_expiration = "4 weeks"
					session_token_expiration = "10 minutes"
					key_session_token_expiration = "30 minutes"
				}
			`),
			Check: a.Check(map[string]any{
				"session_settings.enabled":                      "true",
				"session_settings.refresh_token_expiration":     "4 weeks",
				"session_settings.session_token_expiration":     "10 minutes",
				"session_settings.key_session_token_expiration": "30 minutes",
			}),
		},
		// Test session settings with a user JWT template
		resource.TestStep{
			Config: p.Config() + j.Config(`
				project_id = `+p.Path()+`.id
				type = "user"
				template = "{}"
			`) + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				session_settings = {
					enabled = true
					refresh_token_expiration = "4 weeks"
					session_token_expiration = "10 minutes"
					key_session_token_expiration = "30 minutes"
					user_template_id = `+j.Path()+`.id
				}
			`),
			Check: a.Check(map[string]any{
				"session_settings.enabled":                      "true",
				"session_settings.refresh_token_expiration":     "4 weeks",
				"session_settings.session_token_expiration":     "10 minutes",
				"session_settings.key_session_token_expiration": "30 minutes",
				"session_settings.user_template_id":             testacc.AttributeHasPrefix("JT"),
			}),
		},
		resource.TestStep{
			Config: p.Config() + j2.Config(`
				project_id = `+p.Path()+`.id
				type = "user"
				template = "{}"
			`) + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				session_settings = {
					enabled = true
					refresh_token_expiration = "4 weeks"
					session_token_expiration = "10 minutes"
					key_session_token_expiration = "30 minutes"
					user_template_id = `+j2.Path()+`.id
				}
			`),
			Check: a.Check(map[string]any{
				"session_settings.enabled":                      "true",
				"session_settings.refresh_token_expiration":     "4 weeks",
				"session_settings.session_token_expiration":     "10 minutes",
				"session_settings.key_session_token_expiration": "30 minutes",
				"session_settings.user_template_id":             testacc.AttributeHasPrefix("JT"),
			}),
		},
		// Test import with composite ID
		resource.TestStep{
			ResourceName:      a.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
		},
		// Disable the default deletion protection so the resource can be destroyed
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
			`),
		},
		// force_dpop needs an explicit client type
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				deletion_protection = false
				client_type = "confidential"
				force_dpop = true
			`),
			Check: a.Check(map[string]any{
				"client_type": "confidential",
				"force_dpop":  "true",
			}),
		},
		// Destroy resource
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
			`),
			Destroy: true,
		},
	)
}

func TestInboundAppGrantTypes(t *testing.T) {
	p := testacc.Project(t)
	a := testacc.InboundApp(t)
	pub := testacc.InboundApp(t)
	pub.ID = "public"
	testacc.Run(t,
		// Every grant type is enabled by default
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
			`),
			Check: a.Check(map[string]any{
				"authorization_code_disabled": "false",
				"client_credentials_disabled": "false",
				"refresh_token_disabled":      "false",
				"jwt_bearer_disabled":         "false",
				"device_code_disabled":        "false",
				"token_exchange_disabled":     "false",
			}),
		},
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
				authorization_code_disabled = true
				client_credentials_disabled = true
				refresh_token_disabled = true
				jwt_bearer_disabled = true
				device_code_disabled = true
				token_exchange_disabled = true
			`),
			Check: a.Check(map[string]any{
				"authorization_code_disabled": "true",
				"client_credentials_disabled": "true",
				"refresh_token_disabled":      "true",
				"jwt_bearer_disabled":         "true",
				"device_code_disabled":        "true",
				"token_exchange_disabled":     "true",
			}),
		},
		// Import reads the stored grant types back
		resource.TestStep{
			ResourceName:      a.Path(),
			ImportState:       true,
			ImportStateIdFunc: testacc.GenerateImportStateID(a.Path(), "project_id", "id"),
			ImportStateCheck: func(states []*terraform.InstanceState) error {
				if len(states) != 1 {
					return fmt.Errorf("expected one imported inbound app, got %d", len(states))
				}
				for _, attr := range []string{"authorization_code_disabled", "client_credentials_disabled", "refresh_token_disabled", "jwt_bearer_disabled", "device_code_disabled", "token_exchange_disabled"} {
					if v := states[0].Attributes[attr]; v != "true" {
						return fmt.Errorf("expected the imported %s to be true, got %q", attr, v)
					}
				}
				return nil
			},
		},
		// Removing the attributes enables the grant types again
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
			`),
			Check: a.Check(map[string]any{
				"authorization_code_disabled": "false",
				"client_credentials_disabled": "false",
				"refresh_token_disabled":      "false",
				"jwt_bearer_disabled":         "false",
				"device_code_disabled":        "false",
				"token_exchange_disabled":     "false",
			}),
		},
		// A public client cannot use the client_credentials or jwt-bearer grant types
		resource.TestStep{
			Config: p.Config() + pub.Config(`
				project_id = `+p.Path()+`.id
				client_type = "public"
				deletion_protection = false
				client_credentials_disabled = false
			`),
			ExpectError: regexp.MustCompile(`client_credentials_disabled field cannot be false`),
		},
		resource.TestStep{
			Config: p.Config() + pub.Config(`
				project_id = `+p.Path()+`.id
				client_type = "public"
				deletion_protection = false
			`),
			Check: pub.Check(map[string]any{
				"client_credentials_disabled": "true",
				"jwt_bearer_disabled":         "true",
				"refresh_token_disabled":      "false",
			}),
		},
		// A public client is updated in place with its grant types
		resource.TestStep{
			Config: p.Config() + pub.Config(`
				project_id = `+p.Path()+`.id
				client_type = "public"
				deletion_protection = false
				description = "Public client"
				refresh_token_disabled = true
			`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(pub.Path(), plancheck.ResourceActionUpdate)},
			},
			Check: pub.Check(map[string]any{
				"description":                 "Public client",
				"client_credentials_disabled": "true",
				"jwt_bearer_disabled":         "true",
				"refresh_token_disabled":      "true",
			}),
		},
		// The stored public client type still rejects the grant types when client_type is removed from the configuration
		resource.TestStep{
			Config: p.Config() + pub.Config(`
				project_id = `+p.Path()+`.id
				deletion_protection = false
				description = "Public client"
				refresh_token_disabled = true
				jwt_bearer_disabled = false
			`),
			ExpectError: regexp.MustCompile(`jwt_bearer_disabled field cannot be false`),
		},
	)
}

func TestInboundAppLegacyClientType(t *testing.T) {
	p := testacc.Project(t)
	a := testacc.InboundApp(t)

	var clientID, clientSecret string
	captureClientID := func(s string) error {
		clientID = s
		return nil
	}
	captureClientSecret := func(s string) error {
		if s == "" {
			return fmt.Errorf("expected a client secret for a legacy app")
		}
		clientSecret = s
		return nil
	}
	sameClientID := func(s string) error {
		if s != clientID {
			return fmt.Errorf("expected client_id to be preserved as %s, got %s", clientID, s)
		}
		return nil
	}
	sameClientSecret := func(s string) error {
		if s != clientSecret {
			return fmt.Errorf("expected client_secret to be preserved")
		}
		return nil
	}

	testacc.Run(t,
		// A legacy app is created with a client ID and secret
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = ""
				deletion_protection = false
			`),
			Check: a.Check(map[string]any{
				"client_type":   "",
				"client_id":     captureClientID,
				"client_secret": captureClientSecret,
			}),
		},
		// Moving from the legacy type to confidential updates the app in place and keeps its credentials
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = "confidential"
				deletion_protection = false
			`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(a.Path(), plancheck.ResourceActionUpdate)},
			},
			Check: a.Check(map[string]any{
				"client_type":   "confidential",
				"client_id":     sameClientID,
				"client_secret": sameClientSecret,
			}),
		},
		// Removing client_type from the configuration keeps the stored type
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				deletion_protection = false
			`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
			},
			Check: a.Check(map[string]any{
				"client_type": "confidential",
			}),
		},
		// Moving back to the legacy type replaces the app
		resource.TestStep{
			Config: p.Config() + a.Config(`
				project_id = `+p.Path()+`.id
				client_type = ""
				deletion_protection = false
			`),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(a.Path(), plancheck.ResourceActionReplace)},
			},
			Check: a.Check(map[string]any{
				"client_type": "",
			}),
		},
	)
}

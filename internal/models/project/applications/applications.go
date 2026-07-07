package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ApplicationsValidator = objattr.NewValidator[ApplicationsModel]("must have a valid applications configuration")

var ApplicationsModifier = objattr.NewModifier[ApplicationsModel]("rejects changes to immutable application fields and maintains application identifiers between plan changes")

var ApplicationsAttributes = map[string]schema.Attribute{
	"oidc_applications":  listattr.Default[OIDCModel](OIDCAttributes),
	"saml_applications":  listattr.Default[SAMLModel](SAMLAttributes),
	"wsfed_applications": listattr.Default[WSFedModel](WSFedAttributes),
}

type ApplicationsModel struct {
	OIDCApplications  listattr.Type[OIDCModel]  `tfsdk:"oidc_applications"`
	SAMLApplications  listattr.Type[SAMLModel]  `tfsdk:"saml_applications"`
	WSFedApplications listattr.Type[WSFedModel] `tfsdk:"wsfed_applications"`
}

func (m *ApplicationsModel) Values(h *helpers.Handler) map[string]any {
	m.Check(h)
	data := map[string]any{}
	listattr.Get(m.OIDCApplications, data, "oidc", h)
	listattr.Get(m.SAMLApplications, data, "saml", h)
	listattr.Get(m.WSFedApplications, data, "wsfed", h)
	return data
}

func (m *ApplicationsModel) SetValues(h *helpers.Handler, data map[string]any) {
	listattr.SetMatchingNames(&m.OIDCApplications, data, "oidc", "name", h)
	listattr.SetMatchingNames(&m.SAMLApplications, data, "saml", "name", h)
	listattr.SetMatchingNames(&m.WSFedApplications, data, "wsfed", "name", h)
}

func (m *ApplicationsModel) Check(h *helpers.Handler) {
	for app := range listattr.Iterator(m.SAMLApplications, h) {
		if !app.DynamicConfiguration.IsSet() && !app.ManualConfiguration.IsSet() {
			h.Missing("Either the dynamic_configuration or manual_configuration attribute must be set in the '%s' saml application", app.Name.ValueString())
		} else if app.DynamicConfiguration.IsSet() && app.ManualConfiguration.IsSet() {
			h.Warn("Both dynamic_configuration and manual_configuration supplied - dynamic configuration will take precedence", "dynamic_configuration and manual_configuration are mutually exclusive. If both given - dynamic takes precedence")
		}
	}
}

func (m *ApplicationsModel) Modify(h *helpers.Handler, state *ApplicationsModel) {
	listattr.ModifyMatchingNames(h, &m.OIDCApplications, state.OIDCApplications)
	listattr.ModifyMatchingNames(h, &m.SAMLApplications, state.SAMLApplications)
	listattr.ModifyMatchingNames(h, &m.WSFedApplications, state.WSFedApplications)

	for p := range listattr.Iterator(m.OIDCApplications, h) {
		for s := range listattr.Iterator(state.OIDCApplications, h) {
			if p.ID.ValueString() != "" && s.ID.Equal(p.ID) && s.ClientID.ValueString() != "" && p.ClientID.ValueString() != "" && !p.ClientID.Equal(s.ClientID) {
				h.Error("Immutable Field Changed", "The 'client_id' of OIDC application '%s' cannot be changed after creation, delete the application first to recreate it", p.Name.ValueString())
			}
		}
	}
}

func (m *ApplicationsModel) Validate(h *helpers.Handler) {
	// XXX move Check here eventually
	for app := range listattr.Iterator(m.OIDCApplications, h) {
		validateSSOAppRoles(h, app.Name, app.Permissions, app.Roles)
	}
	for app := range listattr.Iterator(m.SAMLApplications, h) {
		validateSSOAppRoles(h, app.Name, app.Permissions, app.Roles)
	}
	for app := range listattr.Iterator(m.WSFedApplications, h) {
		validateSSOAppRoles(h, app.Name, app.Permissions, app.Roles)
	}
}

func validateSSOAppRoles(h *helpers.Handler, appName stringattr.Type, perms listattr.Type[SSOAppPermissionModel], roles listattr.Type[SSOAppRoleModel]) {
	if helpers.HasUnknownValues(perms, roles) {
		return
	}
	defined := map[string]bool{}
	for p := range listattr.Iterator(perms, h) {
		defined[p.Name.ValueString()] = true
	}
	for r := range listattr.Iterator(roles, h) {
		for name := range strsetattr.Iterator(r.Permissions, h) {
			if !defined[name] {
				h.Error("Missing SSO App Permission", "The role '%s' in app '%s' references a permission '%s' that doesn't exist on the same app", r.Name.ValueString(), appName.ValueString(), name)
			}
		}
	}
}

package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ApplicationsValidator = objattr.NewValidator[ApplicationsModel]("must have a valid SAML configuration")

var ApplicationsAttributes = map[string]schema.Attribute{
	"oidc_applications": listattr.Default[OIDCModel](OIDCAttributes),
	"saml_applications": listattr.Default[SAMLModel](SAMLAttributes),
}

type ApplicationsModel struct {
	OIDCApplications listattr.Type[OIDCModel] `tfsdk:"oidc_applications"`
	SAMLApplications listattr.Type[SAMLModel] `tfsdk:"saml_applications"`
}

var ApplicationsDefault = &ApplicationsModel{
	OIDCApplications: listattr.Empty[OIDCModel](),
	SAMLApplications: listattr.Empty[SAMLModel](),
}

func (m *ApplicationsModel) Values(h *helpers.Handler) map[string]any {
	m.Check(h)
	data := map[string]any{}
	listattr.Get(m.OIDCApplications, data, "oidc", h)
	listattr.Get(m.SAMLApplications, data, "saml", h)
	return data
}

func (m *ApplicationsModel) SetValues(h *helpers.Handler, data map[string]any) {
	if m.OIDCApplications.IsUnknown() {
		listattr.Set(&m.OIDCApplications, data, "oidc", h)
	} else {
		for app := range listattr.MutatingIterator(&m.OIDCApplications, h) {
			RequireID(h, data, "oidc", app.Name, &app.ID)
		}
	}
	if m.SAMLApplications.IsUnknown() {
		listattr.Set(&m.SAMLApplications, data, "saml", h)
	} else {
		for app := range listattr.MutatingIterator(&m.SAMLApplications, h) {
			RequireID(h, data, "saml", app.Name, &app.ID)
		}
	}
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

func (m *ApplicationsModel) Validate(h *helpers.Handler) {
	// XXX move Check here eventually
}

package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
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

func (m *ApplicationsModel) Values(h *helpers.Handler) map[string]any {
	m.Check(h)
	data := map[string]any{}
	listattr.Get(m.OIDCApplications, data, "oidc", h)
	listattr.Get(m.SAMLApplications, data, "saml", h)
	return data
}

func (m *ApplicationsModel) SetValues(h *helpers.Handler, data map[string]any) {
	listattr.SetMatching(&m.OIDCApplications, data, "oidc", "name", h)
	listattr.SetMatching(&m.SAMLApplications, data, "saml", "name", h)
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

package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ApplicationValidator = objectattr.NewValidator[ApplicationModel]("must have a valid SAML configuration")

var ApplicationAttributes = map[string]schema.Attribute{
	"oidc_applications": listattr.Optional(OIDCAttributes),
	"saml_applications": listattr.Optional(SAMLAttributes),
}

type ApplicationModel struct {
	OIDCApplications []*OIDCModel `tfsdk:"oidc_applications"`
	SAMLApplications []*SAMLModel `tfsdk:"saml_applications"`
}

func (m *ApplicationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.OIDCApplications, data, "oidc", h)
	listattr.Get(m.SAMLApplications, data, "saml", h)
	return data
}

func (m *ApplicationModel) SetValues(h *helpers.Handler, data map[string]any) {
	for _, app := range m.OIDCApplications {
		RequireID(h, data, "oidc", app.Name, &app.ID)
	}
	for _, app := range m.SAMLApplications {
		RequireID(h, data, "saml", app.Name, &app.ID)
	}
	if m.OIDCApplications == nil {
		m.OIDCApplications = []*OIDCModel{}
		listattr.Set(&m.OIDCApplications, data, "oidc", h)
	}
	if m.SAMLApplications == nil {
		m.SAMLApplications = []*SAMLModel{}
		listattr.Set(&m.SAMLApplications, data, "saml", h)
	}
}

func (m *ApplicationModel) Validate(h *helpers.Handler) {
	for _, app := range m.SAMLApplications {
		if app.DynamicConfiguration == nil && app.ManualConfiguration == nil {
			h.Error("Either dynamic_configuration or manual_configuration must be set", "no configuration found for application")
		} else if app.DynamicConfiguration != nil && app.ManualConfiguration != nil {
			h.Warn("Both dynamic_configuration and manual_configuration supplied - dynamic configuration will take precedence", "dynamic_configuration and manual_configuration are mutually exclusive. If both given - dynamic takes precedence")
		}
	}
}

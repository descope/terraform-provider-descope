package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ApplicationValidator = objectattr.NewValidator[ApplicationModel]("must have a valid SAML configuration")

var ApplicationAttributes = map[string]schema.Attribute{
	"oidc": listattr.Optional(OIDCAttributes),
	"saml": listattr.Optional(SAMLAttributes),
}

type ApplicationModel struct {
	OIDC []*OIDCModel `tfsdk:"oidc"`
	SAML []*SAMLModel `tfsdk:"saml"`
}

func (m *ApplicationModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.OIDC, data, "oidc", h)
	listattr.Get(m.SAML, data, "saml", h)
	return data
}

func (m *ApplicationModel) SetValues(h *helpers.Handler, data map[string]any) {
	for _, app := range m.OIDC {
		RequireID(h, data, "oidc", app.Name, &app.ID)
	}
	for _, app := range m.SAML {
		RequireID(h, data, "saml", app.Name, &app.ID)
	}
}

func (m *ApplicationModel) Validate(h *helpers.Handler) {
	for _, app := range m.SAML {
		if app.DynamicConfiguration == nil && app.ManualConfiguration == nil {
			h.Error("Either dynamic_configuration or manual_configuration must be set", "no configuration found for application")
		} else if app.DynamicConfiguration != nil && app.ManualConfiguration != nil {
			h.Warn("Both dynamic_configuration and manual_configuration supplied - dynamic configuration will take precedence", "dynamic_configuration and manual_configuration are mutually exclusive. If both given - dynamic takes precedence")
		}
	}
}

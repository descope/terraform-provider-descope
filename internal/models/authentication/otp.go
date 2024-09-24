package authentication

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/templates"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var OTPAttributes = map[string]schema.Attribute{
	"enabled":         boolattr.Optional(),
	"domain":          stringattr.Optional(),
	"expiration_time": durationattr.Optional(durationattr.MinimumValue("1 minute")),
	"email_service":   objectattr.Optional(templates.EmailServiceAttributes, templates.EmailServiceValidator),
	"text_service":    objectattr.Optional(templates.TextServiceAttributes, templates.TextServiceValidator),
	"voice_service":   objectattr.Optional(templates.VoiceServiceAttributes, templates.VoiceServiceValidator),
}

type OTPModel struct {
	Enabled        types.Bool                   `tfsdk:"enabled"`
	Domain         types.String                 `tfsdk:"domain"`
	ExpirationTime types.String                 `tfsdk:"expiration_time"`
	EmailService   *templates.EmailServiceModel `tfsdk:"email_service"`
	TextService    *templates.TextServiceModel  `tfsdk:"text_service"`
	VoiceService   *templates.VoiceServiceModel `tfsdk:"voice_service"`
}

func (m *OTPModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Enabled, data, "enabled")
	stringattr.Get(m.Domain, data, "domain")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	if v := m.EmailService; v != nil {
		maps.Copy(data, v.Values(h))
	}
	if v := m.TextService; v != nil {
		maps.Copy(data, v.Values(h))
	}
	if v := m.VoiceService; v != nil {
		maps.Copy(data, v.Values(h))
	}
	return data
}

func (m *OTPModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Enabled, data, "enabled")
	stringattr.Set(&m.Domain, data, "domain")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
	if v := m.EmailService; v != nil {
		v.SetValues(h, data)
	}
	if v := m.TextService; v != nil {
		v.SetValues(h, data)
	}
	if v := m.VoiceService; v != nil {
		v.SetValues(h, data)
	}
}

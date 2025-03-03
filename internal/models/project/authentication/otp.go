package authentication

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/project/templates"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var OTPAttributes = map[string]schema.Attribute{
	"disabled":        boolattr.Default(false),
	"domain":          stringattr.Optional(),
	"expiration_time": durationattr.Optional(durationattr.MinimumValue("1 minute")),
	"email_service":   objectattr.Optional(templates.EmailServiceAttributes, templates.EmailServiceValidator),
	"text_service":    objectattr.Optional(templates.TextServiceAttributes, templates.TextServiceValidator),
	"voice_service":   objectattr.Optional(templates.VoiceServiceAttributes, templates.VoiceServiceValidator),
}

type OTPModel struct {
	Disabled       types.Bool                   `tfsdk:"disabled"`
	Domain         types.String                 `tfsdk:"domain"`
	ExpirationTime types.String                 `tfsdk:"expiration_time"`
	EmailService   *templates.EmailServiceModel `tfsdk:"email_service"`
	TextService    *templates.TextServiceModel  `tfsdk:"text_service"`
	VoiceService   *templates.VoiceServiceModel `tfsdk:"voice_service"`
}

func (m *OTPModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
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
	boolattr.SetNot(&m.Disabled, data, "enabled")
	stringattr.Set(&m.Domain, data, "domain")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
	if m.EmailService = helpers.InitIfImport(h.Ctx, m.EmailService); m.EmailService != nil {
		m.EmailService.SetValues(h, data)
	}
	if m.TextService = helpers.InitIfImport(h.Ctx, m.TextService); m.TextService != nil {
		m.TextService.SetValues(h, data)
	}
	if m.VoiceService = helpers.InitIfImport(h.Ctx, m.VoiceService); m.VoiceService != nil {
		m.VoiceService.SetValues(h, data)
	}
}

func (m *OTPModel) SetReferences(h *helpers.Handler) {
	if m.EmailService != nil {
		m.EmailService.SetReferences(h)
	}
	if m.TextService != nil {
		m.TextService.SetReferences(h)
	}
	if m.VoiceService != nil {
		m.VoiceService.SetReferences(h)
	}
}

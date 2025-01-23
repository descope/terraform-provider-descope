package authentication

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/templates"
	"github.com/descope/terraform-provider-descope/internal/utils"
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
	emailService := utils.ZVL(m.EmailService)
	emailService.SetValues(h, data)
	if emailService.Connector.ValueString() != "" {
		m.EmailService = emailService
	}
	textService := utils.ZVL(m.TextService)
	textService.SetValues(h, data)
	if textService.Connector.ValueString() != "" {
		m.TextService = textService
	}
	voiceService := utils.ZVL(m.VoiceService)
	voiceService.SetValues(h, data)
	if voiceService.Connector.ValueString() != "" {
		m.VoiceService = voiceService
	}
	}
}

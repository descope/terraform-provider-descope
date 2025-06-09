package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/project/templates"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var OTPAttributes = map[string]schema.Attribute{
	"disabled":        boolattr.Default(false),
	"domain":          stringattr.Optional(),
	"expiration_time": durationattr.Optional(durationattr.MinimumValue("1 minute")),
	"email_service":   objattr.Optional[templates.EmailServiceModel](templates.EmailServiceAttributes, templates.EmailServiceValidator),
	"text_service":    objattr.Optional[templates.TextServiceModel](templates.TextServiceAttributes, templates.TextServiceValidator),
	"voice_service":   objattr.Optional[templates.VoiceServiceModel](templates.VoiceServiceAttributes, templates.VoiceServiceValidator),
}

type OTPModel struct {
	Disabled       boolattr.Type                             `tfsdk:"disabled"`
	Domain         stringattr.Type                           `tfsdk:"domain"`
	ExpirationTime stringattr.Type                           `tfsdk:"expiration_time"`
	EmailService   objattr.Type[templates.EmailServiceModel] `tfsdk:"email_service"`
	TextService    objattr.Type[templates.TextServiceModel]  `tfsdk:"text_service"`
	VoiceService   objattr.Type[templates.VoiceServiceModel] `tfsdk:"voice_service"`
}

func (m *OTPModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	stringattr.Get(m.Domain, data, "domain")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	objattr.Get(m.EmailService, data, helpers.RootKey, h)
	objattr.Get(m.TextService, data, helpers.RootKey, h)
	objattr.Get(m.VoiceService, data, helpers.RootKey, h)
	return data
}

func (m *OTPModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	stringattr.Set(&m.Domain, data, "domain")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
	objattr.Set(&m.EmailService, data, helpers.RootKey, h)
	objattr.Set(&m.TextService, data, helpers.RootKey, h)
	objattr.Set(&m.VoiceService, data, helpers.RootKey, h)
}

func (m *OTPModel) UpdateReferences(h *helpers.Handler) {
	objattr.UpdateReferences(&m.EmailService, h)
	objattr.UpdateReferences(&m.TextService, h)
	objattr.UpdateReferences(&m.VoiceService, h)
}

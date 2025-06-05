package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/project/templates"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var MagicLinkAttributes = map[string]schema.Attribute{
	"disabled":        boolattr.Default(false),
	"expiration_time": durationattr.Optional(durationattr.MinimumValue("1 minute")),
	"redirect_url":    stringattr.Optional(),
	"email_service":   objattr.Optional[templates.EmailServiceModel](templates.EmailServiceAttributes, templates.EmailServiceValidator),
	"text_service":    objattr.Optional[templates.TextServiceModel](templates.TextServiceAttributes, templates.TextServiceValidator),
}

type MagicLinkModel struct {
	Disabled       boolattr.Type                             `tfsdk:"disabled"`
	ExpirationTime stringattr.Type                           `tfsdk:"expiration_time"`
	RedirectURL    stringattr.Type                           `tfsdk:"redirect_url"`
	EmailService   objattr.Type[templates.EmailServiceModel] `tfsdk:"email_service"`
	TextService    objattr.Type[templates.TextServiceModel]  `tfsdk:"text_service"`
}

func (m *MagicLinkModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	objattr.Get(m.EmailService, data, helpers.RootKey, h)
	objattr.Get(m.TextService, data, helpers.RootKey, h)
	return data
}

func (m *MagicLinkModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	objattr.Set(&m.EmailService, data, helpers.RootKey, h)
	objattr.Set(&m.TextService, data, helpers.RootKey, h)
}

func (m *MagicLinkModel) UpdateReferences(h *helpers.Handler) {
	objattr.UpdateReferences(&m.EmailService, h)
	objattr.UpdateReferences(&m.TextService, h)
}

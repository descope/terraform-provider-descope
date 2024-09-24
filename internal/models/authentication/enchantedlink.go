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

var EnchantedLinkAttributes = map[string]schema.Attribute{
	"disabled":        boolattr.Default(false),
	"expiration_time": durationattr.Optional(durationattr.MinimumValue("1 minute")),
	"redirect_url":    stringattr.Optional(),
	"email_service":   objectattr.Optional(templates.EmailServiceAttributes, templates.EmailServiceValidator),
}

type EnchantedLinkModel struct {
	Disabled       types.Bool                   `tfsdk:"disabled"`
	ExpirationTime types.String                 `tfsdk:"expiration_time"`
	RedirectURL    types.String                 `tfsdk:"redirect_url"`
	EmailService   *templates.EmailServiceModel `tfsdk:"email_service"`
}

func (m *EnchantedLinkModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	if v := m.EmailService; v != nil {
		maps.Copy(data, v.Values(h))
	}
	return data
}

func (m *EnchantedLinkModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	if v := m.EmailService; v != nil {
		v.SetValues(h, data)
	}
}

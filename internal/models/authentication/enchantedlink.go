package authentication

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/templates"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var EnchantedLinkAttributes = map[string]schema.Attribute{
	"enabled":              boolattr.Optional(),
	"expiration_time":      intattr.Optional(),
	"expiration_time_unit": stringattr.Optional(stringattr.TimeUnitValidator),
	"redirect_url":         stringattr.Optional(),
	"email_service":        objectattr.Optional(templates.EmailServiceAttributes, templates.EmailServiceValidator),
}

type EnchantedLinkModel struct {
	Enabled            types.Bool                   `tfsdk:"enabled"`
	ExpirationTime     types.Int64                  `tfsdk:"expiration_time"`
	ExpirationTimeUnit types.String                 `tfsdk:"expiration_time_unit"`
	RedirectURL        types.String                 `tfsdk:"redirect_url"`
	EmailService       *templates.EmailServiceModel `tfsdk:"email_service"`
}

func (m *EnchantedLinkModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Enabled, data, "enabled")
	intattr.Get(m.ExpirationTime, data, "expirationTime")
	stringattr.Get(m.ExpirationTimeUnit, data, "expirationTimeUnit")
	stringattr.Get(m.RedirectURL, data, "redirectUrl")
	if v := m.EmailService; v != nil {
		maps.Copy(data, v.Values(h))
	}
	return data
}

func (m *EnchantedLinkModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Enabled, data, "enabled")
	intattr.Set(&m.ExpirationTime, data, "expirationTime")
	stringattr.Set(&m.ExpirationTimeUnit, data, "expirationTimeUnit")
	stringattr.Set(&m.RedirectURL, data, "redirectUrl")
	if v := m.EmailService; v != nil {
		v.SetValues(h, data)
	}
}

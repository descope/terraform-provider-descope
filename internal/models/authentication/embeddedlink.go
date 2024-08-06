package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var EmbeddedLinkAttributes = map[string]schema.Attribute{
	"enabled":              boolattr.Optional(),
	"expiration_time":      intattr.Optional(),
	"expiration_time_unit": stringattr.Optional(stringattr.TimeUnitValidator),
}

type EmbeddedLinkModel struct {
	Enabled            types.Bool   `tfsdk:"enabled"`
	ExpirationTime     types.Int64  `tfsdk:"expiration_time"`
	ExpirationTimeUnit types.String `tfsdk:"expiration_time_unit"`
}

func (m *EmbeddedLinkModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Enabled, data, "enabled")
	intattr.Get(m.ExpirationTime, data, "expirationTime")
	stringattr.Get(m.ExpirationTimeUnit, data, "expirationTimeUnit")
	return data
}

func (m *EmbeddedLinkModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Enabled, data, "enabled")
	intattr.Set(&m.ExpirationTime, data, "expirationTime")
	stringattr.Set(&m.ExpirationTimeUnit, data, "expirationTimeUnit")
}

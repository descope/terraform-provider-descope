package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/durationattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var EmbeddedLinkAttributes = map[string]schema.Attribute{
	"enabled":         boolattr.Optional(),
	"expiration_time": durationattr.Optional(durationattr.MinimumValue("1 minute")),
}

type EmbeddedLinkModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	ExpirationTime types.String `tfsdk:"expiration_time"`
}

func (m *EmbeddedLinkModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Enabled, data, "enabled")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	return data
}

func (m *EmbeddedLinkModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Enabled, data, "enabled")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
}

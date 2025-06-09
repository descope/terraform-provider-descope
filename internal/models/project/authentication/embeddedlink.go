package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var EmbeddedLinkAttributes = map[string]schema.Attribute{
	"disabled":        boolattr.Default(false),
	"expiration_time": durationattr.Optional(durationattr.MinimumValue("1 minute")),
}

type EmbeddedLinkModel struct {
	Disabled       boolattr.Type     `tfsdk:"disabled"`
	ExpirationTime durationattr.Type `tfsdk:"expiration_time"`
}

func (m *EmbeddedLinkModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	return data
}

func (m *EmbeddedLinkModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
}

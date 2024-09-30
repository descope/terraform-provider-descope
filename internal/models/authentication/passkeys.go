package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var PasskeysAttributes = map[string]schema.Attribute{
	"disabled":         boolattr.Default(false),
	"top_level_domain": stringattr.Optional(),
}

type PasskeysModel struct {
	Disabled       types.Bool   `tfsdk:"disabled"`
	TopLevelDomain types.String `tfsdk:"top_level_domain"`
}

func (m *PasskeysModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	stringattr.Get(m.TopLevelDomain, data, "relyingPartyId")
	return data
}

func (m *PasskeysModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	stringattr.Set(&m.TopLevelDomain, data, "relyingPartyId")
}

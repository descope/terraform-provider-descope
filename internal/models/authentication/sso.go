package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SSOAttributes = map[string]schema.Attribute{
	"enabled":     boolattr.Optional(),
	"merge_users": boolattr.Optional(),
}

type SSOModel struct {
	Enabled    types.Bool `tfsdk:"enabled"`
	MergeUsers types.Bool `tfsdk:"merge_users"`
}

func (m *SSOModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Enabled, data, "enabled")
	boolattr.Get(m.MergeUsers, data, "mergeUsers")
	return data
}

func (m *SSOModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Enabled, data, "enabled")
	boolattr.Set(&m.MergeUsers, data, "mergeUsers")
}

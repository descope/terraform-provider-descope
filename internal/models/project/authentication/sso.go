package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SSOAttributes = map[string]schema.Attribute{
	"disabled":    boolattr.Default(false),
	"merge_users": boolattr.Default(false),
	// "redirect_url": stringattr.Default(""), // XXX not yet
}

type SSOModel struct {
	Disabled   boolattr.Type `tfsdk:"disabled"`
	MergeUsers boolattr.Type `tfsdk:"merge_users"`
	// RedirectURL stringattr.Type `tfsdk:"redirect_url"` // XXX not yet
}

func (m *SSOModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	boolattr.Get(m.MergeUsers, data, "mergeUsers")
	// stringattr.Get(m.RedirectURL, data, "redirectUrl") // XXX not yet
	return data
}

func (m *SSOModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	boolattr.Set(&m.MergeUsers, data, "mergeUsers")
	// stringattr.Set(&m.RedirectURL, data, "redirectUrl") // XXX not yet
}

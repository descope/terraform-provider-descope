package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var TOTPAttributes = map[string]schema.Attribute{
	"disabled": boolattr.Default(false),
}

type TOTPModel struct {
	Disabled boolattr.Type `tfsdk:"disabled"`
}

func (m *TOTPModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	return data
}

func (m *TOTPModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
}

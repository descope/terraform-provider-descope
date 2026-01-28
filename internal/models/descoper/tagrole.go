package descoper

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var DescoperTagRoleAttributes = map[string]schema.Attribute{
	"tags": strsetattr.Required(),
	"role": stringattr.Required(stringvalidator.OneOf("admin", "developer", "support", "auditor")),
}

type DescoperTagRoleModel struct {
	Tags strsetattr.Type `tfsdk:"tags"`
	Role stringattr.Type `tfsdk:"role"`
}

func (m *DescoperTagRoleModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.Tags, data, "tags", h)
	stringattr.Get(m.Role, data, "role")
	return data
}

func (m *DescoperTagRoleModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.Tags, data, "tags", h)
	stringattr.Set(&m.Role, data, "role")
}

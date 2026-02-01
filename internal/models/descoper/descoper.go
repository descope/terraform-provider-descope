package descoper

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var DescoperAttributes = map[string]schema.Attribute{
	"id":    stringattr.Identifier(),
	"email": stringattr.Required(),
	"phone": stringattr.Default(""),
	"name":  stringattr.Default(""),
	"rbac":  objattr.Required[RBacModel](RBacAttributes, RBacValidator),
}

type DescoperModel struct {
	ID    stringattr.Type         `tfsdk:"id"`
	Email stringattr.Type         `tfsdk:"email"`
	Phone stringattr.Type         `tfsdk:"phone"`
	Name  stringattr.Type         `tfsdk:"name"`
	RBac  objattr.Type[RBacModel] `tfsdk:"rbac"`
}

func (m *DescoperModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Email, data, "email")
	stringattr.Get(m.Phone, data, "phone")
	stringattr.Get(m.Name, data, "name")
	objattr.Get(m.RBac, data, "rbac", h)
	return data
}

func (m *DescoperModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Email, data, "email", stringattr.SkipIfAlreadySet)
	stringattr.Set(&m.Phone, data, "phone", stringattr.SkipIfAlreadySet)
	stringattr.Set(&m.Name, data, "name", stringattr.SkipIfAlreadySet)
	objattr.Set(&m.RBac, data, "rbac", h)
}

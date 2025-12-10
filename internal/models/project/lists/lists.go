package lists

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ListsAttributes = map[string]schema.Attribute{
	"list": listattr.Default[ListModel](ListAttributes),
}

type ListsModel struct {
	List listattr.Type[ListModel] `tfsdk:"list"`
}

func (m *ListsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.List, data, "lists", h)
	return data
}

func (m *ListsModel) SetValues(h *helpers.Handler, data map[string]any) {
	listattr.SetMatchingNames(&m.List, data, "lists", "name", h)
}

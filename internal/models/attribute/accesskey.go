package attribute

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var AccessKeyAttributeSchema = schema.Schema{
	MarkdownDescription: "Manages a single access key custom attribute definition in a Descope project.",
	Attributes:          AccessKeyAttributeAttributes,
}

var AccessKeyAttributeAttributes = map[string]schema.Attribute{
	"id":                   stringattr.Required(stringattr.MachineIDValidator, stringvalidator.LengthAtMost(20), stringplanmodifier.RequiresReplace()),
	"project_id":           stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":                 stringattr.Required(stringattr.StandardLenValidator),
	"type":                 stringattr.Required(typeValidator, stringplanmodifier.RequiresReplace()),
	"select_options":       strsetattr.Default(setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1))),
	"widget_authorization": objattr.Default(widgetAuthorizationDefault, WidgetAuthorizationAttributes),
}

type AccessKeyAttributeModel struct {
	AttributeModel
	WidgetAuthorization objattr.Type[widgetAuthorizationModel] `tfsdk:"widget_authorization"`
}

func (m *AccessKeyAttributeModel) Values(h *helpers.Handler) map[string]any {
	data := m.AttributeModel.Values(h)
	objattr.Get(m.WidgetAuthorization, data, helpers.RootKey, h)
	return data
}

func (m *AccessKeyAttributeModel) SetValues(h *helpers.Handler, data map[string]any) {
	m.AttributeModel.SetValues(h, data)
	objattr.Set(&m.WidgetAuthorization, data, helpers.RootKey, h)
}

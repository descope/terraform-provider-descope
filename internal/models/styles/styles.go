package styles

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/jsonattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages the styles (theme) of the project's flow pages via the exported theme JSON representation.",
	Attributes:          StylesAttributes,
}

var StylesAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"data":       jsonattr.Required("styles"),
}

type StylesModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	Data      jsonattr.Type   `tfsdk:"data"`
}

func (m *StylesModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	jsonattr.Get(m.Data, data, helpers.RootKey)
	return data
}

func (m *StylesModel) SetValues(h *helpers.Handler, data map[string]any) {
	jsonattr.Set(&m.Data, data, helpers.RootKey, jsonattr.SkipIfAlreadySet)
}

func (m *StylesModel) GetID() stringattr.Type        { return m.ID }
func (m *StylesModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *StylesModel) GetProjectID() stringattr.Type { return m.ProjectID }

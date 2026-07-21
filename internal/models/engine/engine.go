package engine

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var EngineAttributes = map[string]schema.Attribute{
	"id":           stringattr.Identifier(),
	"project_id":   stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":         stringattr.Required(),
	"created_time": intattr.Generated(),
	"secret":       stringattr.SecretGenerated(false), // returned only on create; kept in state afterwards
}

var Schema = schema.Schema{
	Attributes: EngineAttributes,
}

type EngineModel struct {
	ID          stringattr.Type `tfsdk:"id"`
	ProjectID   stringattr.Type `tfsdk:"project_id"`
	Name        stringattr.Type `tfsdk:"name"`
	CreatedTime intattr.Type    `tfsdk:"created_time"`
	Secret      stringattr.Type `tfsdk:"secret"`
}

func (m *EngineModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	return data
}

func (m *EngineModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	intattr.Set(&m.CreatedTime, data, "createdTime")
	stringattr.Set(&m.Secret, data, "secret") // absent on read/update, so the create-time value is kept
}

func (m *EngineModel) GetID() stringattr.Type {
	return m.ID
}

func (m *EngineModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *EngineModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

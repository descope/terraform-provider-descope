package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_oauth_settings is the project-level OAuth settings singleton (id = project_id).
// The general OAuth fields this resource doesn't own are preserved by the backend on write.

var OAuthSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level OAuth authentication settings. This is a singleton resource, and its id is always the project ID. Individual OAuth providers are managed with the descope_oauth_provider resource.",
	Attributes:          OAuthSettingsAttributes,
}

var OAuthSettingsAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"disabled":   boolattr.Default(false),
}

type OAuthSettingsModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	Disabled  boolattr.Type   `tfsdk:"disabled"`
}

func (m *OAuthSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	return data
}

func (m *OAuthSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
}

func (m *OAuthSettingsModel) GetID() stringattr.Type {
	return m.ID
}

func (m *OAuthSettingsModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *OAuthSettingsModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

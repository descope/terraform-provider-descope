package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_totp_settings is the project-level TOTP (authenticator app) settings singleton (id = project_id).

var TOTPSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level TOTP (authenticator app) authentication settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          TOTPSettingsAttributes,
}

var TOTPSettingsAttributes = map[string]schema.Attribute{
	"id":            stringattr.Identifier(),
	"project_id":    stringattr.Required(stringplanmodifier.RequiresReplace()),
	"disabled":      boolattr.Default(false),
	"service_label": stringattr.Default(""),
}

type TOTPSettingsModel struct {
	ID           stringattr.Type `tfsdk:"id"`
	ProjectID    stringattr.Type `tfsdk:"project_id"`
	Disabled     boolattr.Type   `tfsdk:"disabled"`
	ServiceLabel stringattr.Type `tfsdk:"service_label"`
}

func (m *TOTPSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	stringattr.Get(m.ServiceLabel, data, "issuerLabelTemplate")
	return data
}

func (m *TOTPSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	stringattr.Set(&m.ServiceLabel, data, "issuerLabelTemplate")
}

func (m *TOTPSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *TOTPSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *TOTPSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }

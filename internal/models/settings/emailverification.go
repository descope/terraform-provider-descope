package settings

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var EmailVerificationSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level email verification authentication settings (beta). This is a singleton resource, and its id is always the project ID.",
	Attributes:          EmailVerificationSettingsAttributes,
}

var EmailVerificationSettingsAttributes = map[string]schema.Attribute{
	"id":         stringattr.Identifier(),
	"project_id": stringattr.Required(stringplanmodifier.RequiresReplace()),
	"enabled":    boolattr.Default(false),
}

type EmailVerificationSettingsModel struct {
	ID        stringattr.Type `tfsdk:"id"`
	ProjectID stringattr.Type `tfsdk:"project_id"`
	Enabled   boolattr.Type   `tfsdk:"enabled"`
}

func (m *EmailVerificationSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Enabled, data, "enabled")
	return data
}

func (m *EmailVerificationSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Enabled, data, "enabled")
}

func (m *EmailVerificationSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *EmailVerificationSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *EmailVerificationSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }

package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_embeddedlink_settings is the project-level embedded link settings singleton (id = project_id).
// Embedded links are generated via the management API rather than delivered to users, so there are no message templates or delivery services.

var EmbeddedLinkSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level embedded link authentication settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          EmbeddedLinkSettingsAttributes,
}

var EmbeddedLinkSettingsAttributes = map[string]schema.Attribute{
	"id":              stringattr.Identifier(),
	"project_id":      stringattr.Required(stringplanmodifier.RequiresReplace()),
	"disabled":        boolattr.Default(false),
	"expiration_time": durationattr.Default("3 minutes", durationattr.MinimumValue("1 minute")),
}

type EmbeddedLinkSettingsModel struct {
	ID             stringattr.Type `tfsdk:"id"`
	ProjectID      stringattr.Type `tfsdk:"project_id"`
	Disabled       boolattr.Type   `tfsdk:"disabled"`
	ExpirationTime stringattr.Type `tfsdk:"expiration_time"`
}

func (m *EmbeddedLinkSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	durationattr.Get(m.ExpirationTime, data, "expirationTime")
	return data
}

func (m *EmbeddedLinkSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	durationattr.Set(&m.ExpirationTime, data, "expirationTime")
}

func (m *EmbeddedLinkSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *EmbeddedLinkSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *EmbeddedLinkSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }

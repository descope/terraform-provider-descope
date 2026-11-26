package settings

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/setattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var WellKnownSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project's customization of its OIDC discovery documents. This is a singleton resource, and its id is always the project ID.",
	Attributes:          WellKnownSettingsAttributes,
}

var WellKnownSettingsAttributes = map[string]schema.Attribute{
	"id":               stringattr.Identifier(),
	"project_id":       stringattr.Required(stringplanmodifier.RequiresReplace()),
	"enabled":          boolattr.Default(false),
	"hidden_fields":    strsetattr.Default(),
	"removed_values":   setattr.Default[WellKnownRemovedValuesModel](WellKnownRemovedValuesAttributes),
	"scopes_supported": objattr.Default(WellKnownValueDeltaDefault, WellKnownValueDeltaAttributes),
	"claims_supported": objattr.Default(WellKnownValueDeltaDefault, WellKnownValueDeltaAttributes),
}

type WellKnownSettingsModel struct {
	ID              stringattr.Type                           `tfsdk:"id"`
	ProjectID       stringattr.Type                           `tfsdk:"project_id"`
	Enabled         boolattr.Type                             `tfsdk:"enabled"`
	HiddenFields    strsetattr.Type                           `tfsdk:"hidden_fields"`
	RemovedValues   setattr.Type[WellKnownRemovedValuesModel] `tfsdk:"removed_values"`
	ScopesSupported objattr.Type[WellKnownValueDeltaModel]    `tfsdk:"scopes_supported"`
	ClaimsSupported objattr.Type[WellKnownValueDeltaModel]    `tfsdk:"claims_supported"`
}

func (m *WellKnownSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Enabled, data, "enabled")
	strsetattr.Get(m.HiddenFields, data, "hiddenFields", h)
	setattr.Get(m.RemovedValues, data, "removedValues", h)
	objattr.Get(m.ScopesSupported, data, "scopesSupported", h)
	objattr.Get(m.ClaimsSupported, data, "claimsSupported", h)
	return data
}

func (m *WellKnownSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Enabled, data, "enabled")
	strsetattr.Set(&m.HiddenFields, data, "hiddenFields", h)
	setattr.Set(&m.RemovedValues, data, "removedValues", h)
	objattr.Set(&m.ScopesSupported, data, "scopesSupported", h)
	objattr.Set(&m.ClaimsSupported, data, "claimsSupported", h)
}

func (m *WellKnownSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *WellKnownSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *WellKnownSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }

var WellKnownRemovedValuesAttributes = map[string]schema.Attribute{
	"field":  stringattr.Required(),
	"values": strsetattr.Required(),
}

type WellKnownRemovedValuesModel struct {
	Field       stringattr.Type `tfsdk:"field"`
	FieldValues strsetattr.Type `tfsdk:"values"`
}

func (m *WellKnownRemovedValuesModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Field, data, "field")
	strsetattr.Get(m.FieldValues, data, "values", h)
	return data
}

func (m *WellKnownRemovedValuesModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Field, data, "field")
	strsetattr.Set(&m.FieldValues, data, "values", h)
}

var WellKnownValueDeltaAttributes = map[string]schema.Attribute{
	"extra":    strsetattr.Default(),
	"excluded": strsetattr.Default(),
}

var WellKnownValueDeltaDefault = &WellKnownValueDeltaModel{
	Extra:    strsetattr.Empty(),
	Excluded: strsetattr.Empty(),
}

type WellKnownValueDeltaModel struct {
	Extra    strsetattr.Type `tfsdk:"extra"`
	Excluded strsetattr.Type `tfsdk:"excluded"`
}

func (m *WellKnownValueDeltaModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	strsetattr.Get(m.Extra, data, "extra", h)
	strsetattr.Get(m.Excluded, data, "excluded", h)
	return data
}

func (m *WellKnownValueDeltaModel) SetValues(h *helpers.Handler, data map[string]any) {
	strsetattr.Set(&m.Extra, data, "extra", h)
	strsetattr.Set(&m.Excluded, data, "excluded", h)
}

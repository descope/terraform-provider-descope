package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/templates"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Delivery service blocks for settings whose message templates are managed by the standalone template resources.

// Email Service

var EmailServiceRefAttributes = map[string]schema.Attribute{
	"connector_id": stringattr.Default(helpers.DescopeConnector),
}

type EmailServiceRefModel struct {
	ConnectorID stringattr.Type `tfsdk:"connector_id"`
}

func (m *EmailServiceRefModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ConnectorID, data, "emailServiceProvider")
	return data
}

func (m *EmailServiceRefModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ConnectorID, data, "emailServiceProvider")
	templates.SetServiceConnectorID(&m.ConnectorID)
}

// Text Service

var TextServiceRefAttributes = map[string]schema.Attribute{
	"connector_id": stringattr.Default(helpers.DescopeConnector),
}

type TextServiceRefModel struct {
	ConnectorID stringattr.Type `tfsdk:"connector_id"`
}

func (m *TextServiceRefModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ConnectorID, data, "textServiceProvider")
	return data
}

func (m *TextServiceRefModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ConnectorID, data, "textServiceProvider")
	templates.SetServiceConnectorID(&m.ConnectorID)
}

// Voice Service

var VoiceServiceRefAttributes = map[string]schema.Attribute{
	"connector_id": stringattr.Default(helpers.DescopeConnector),
}

type VoiceServiceRefModel struct {
	ConnectorID stringattr.Type `tfsdk:"connector_id"`
}

func (m *VoiceServiceRefModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ConnectorID, data, "voiceServiceProvider")
	return data
}

func (m *VoiceServiceRefModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ConnectorID, data, "voiceServiceProvider")
	templates.SetServiceConnectorID(&m.ConnectorID)
}

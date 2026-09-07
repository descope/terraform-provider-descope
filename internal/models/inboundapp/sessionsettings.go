package inboundapp

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SessionSettingsAttributes = map[string]schema.Attribute{
	"enabled":                      boolattr.Default(false),
	"refresh_token_expiration":     durationattr.Default("520 weeks"),
	"session_token_expiration":     durationattr.Default("10 minutes"),
	"key_session_token_expiration": durationattr.Default("10 minutes"),
	"user_template_id":             stringattr.Default(""),
	"key_template_id":              stringattr.Default(""),
}

type SessionSettingsModel struct {
	Enabled                   boolattr.Type     `tfsdk:"enabled"`
	RefreshTokenExpiration    durationattr.Type `tfsdk:"refresh_token_expiration"`
	SessionTokenExpiration    durationattr.Type `tfsdk:"session_token_expiration"`
	KeySessionTokenExpiration durationattr.Type `tfsdk:"key_session_token_expiration"`
	UserTemplateId            stringattr.Type   `tfsdk:"user_template_id"`
	KeyTemplateId             stringattr.Type   `tfsdk:"key_template_id"`
}

func (m *SessionSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Enabled, data, "enabled")
	durationattr.Get(m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	durationattr.Get(m.SessionTokenExpiration, data, "sessionTokenExpiration")
	durationattr.Get(m.KeySessionTokenExpiration, data, "keySessionTokenExpiration")
	stringattr.Get(m.UserTemplateId, data, "userTemplateId")
	stringattr.Get(m.KeyTemplateId, data, "keyTemplateId")
	return data
}

func (m *SessionSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Enabled, data, "enabled")
	durationattr.Set(&m.RefreshTokenExpiration, data, "refreshTokenExpiration")
	durationattr.Set(&m.SessionTokenExpiration, data, "sessionTokenExpiration")
	durationattr.Set(&m.KeySessionTokenExpiration, data, "keySessionTokenExpiration")
	stringattr.Set(&m.UserTemplateId, data, "userTemplateId")
	stringattr.Set(&m.KeyTemplateId, data, "keyTemplateId")
}

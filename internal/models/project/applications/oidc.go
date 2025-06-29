package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var OIDCAttributes = map[string]schema.Attribute{
	"id":          stringattr.Optional(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),
	"logo":        stringattr.Default(""),
	"disabled":    boolattr.Default(false),

	"login_page_url":       stringattr.Default(""),
	"claims":               strlistattr.Default(),
	"force_authentication": boolattr.Default(false),
}

// Model

type OIDCModel struct {
	ID                  stringattr.Type  `tfsdk:"id"`
	Name                stringattr.Type  `tfsdk:"name"`
	Description         stringattr.Type  `tfsdk:"description"`
	Logo                stringattr.Type  `tfsdk:"logo"`
	Disabled            boolattr.Type    `tfsdk:"disabled"`
	LoginPageURL        stringattr.Type  `tfsdk:"login_page_url"`
	Claims              strlistattr.Type `tfsdk:"claims"`
	ForceAuthentication boolattr.Type    `tfsdk:"force_authentication"`
}

func (m *OIDCModel) Values(h *helpers.Handler) map[string]any {
	settings := map[string]any{}
	stringattr.Get(m.LoginPageURL, settings, "loginPageUrl")
	strlistattr.Get(m.Claims, settings, "claims", h)
	boolattr.Get(m.ForceAuthentication, settings, "forceAuthentication")

	data := sharedApplicationData(h, m.ID, m.Name, m.Description, m.Logo, m.Disabled)
	data["oidc"] = settings
	return data
}

func (m *OIDCModel) SetValues(h *helpers.Handler, data map[string]any) {
	setSharedApplicationData(h, data, &m.ID, &m.Name, &m.Description, &m.Logo, &m.Disabled)
	if settings, ok := data["oidc"].(map[string]any); ok {
		stringattr.Set(&m.LoginPageURL, settings, "loginPageUrl")
		strlistattr.Set(&m.Claims, settings, "claims", h)
		boolattr.Set(&m.ForceAuthentication, settings, "forceAuthentication")
	}
}

// Matching

func (m *OIDCModel) GetName() stringattr.Type {
	return m.Name
}

func (m *OIDCModel) GetID() stringattr.Type {
	return m.ID
}

func (m *OIDCModel) SetID(id stringattr.Type) {
	m.ID = id
}

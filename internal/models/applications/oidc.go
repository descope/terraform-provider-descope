package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strlistattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var OIDCAttributes = map[string]schema.Attribute{
	"id":          stringattr.Optional(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),
	"logo":        stringattr.Default(""),
	"disabled":    boolattr.Default(false),

	"login_page_url":       stringattr.Default(""),
	"claims":               strlistattr.Optional(),
	"force_authentication": boolattr.Default(false),
}

// Model

type OIDCModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Logo                types.String `tfsdk:"logo"`
	Disabled            types.Bool   `tfsdk:"disabled"`
	LoginPageURL        types.String `tfsdk:"login_page_url"`
	Claims              []string     `tfsdk:"claims"`
	ForceAuthentication types.Bool   `tfsdk:"force_authentication"`
}

func (m *OIDCModel) Values(h *Handler) map[string]any {
	data := sharedApplicationData(h, m.ID, m.Name, m.Description, m.Logo, m.Disabled)
	settings := map[string]any{}
	stringattr.Get(m.LoginPageURL, settings, "loginPageUrl")
	strlistattr.Get(m.Claims, settings, "claims", h)
	boolattr.Get(m.ForceAuthentication, settings, "forceAuthentication")
	data["oidc"] = settings
	return data
}

func (m *OIDCModel) SetValues(h *Handler, data map[string]any) {
	setSharedApplicationData(h, data, &m.ID, &m.Name, &m.Description, &m.Logo, &m.Disabled)
	if settings, ok := data["oidc"].(map[string]any); ok {
		stringattr.Set(&m.LoginPageURL, settings, "loginPageUrl")
		strlistattr.Set(&m.Claims, settings, "claims", h)
	}
}

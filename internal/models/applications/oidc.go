package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var OIDCAttributes = map[string]schema.Attribute{
	"id":          stringattr.Optional(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),
	"logo":        stringattr.Default(""),
	"disabled":    boolattr.Default(false),
	// oidc
	"login_page_url": stringattr.Default(""),
	"claims":         listattr.StringOptional(),
}

// Model

type OIDCModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Logo         types.String `tfsdk:"logo"`
	Disabled     types.Bool   `tfsdk:"disabled"`
	LoginPageURL types.String `tfsdk:"login_page_url"`
	Claims       []string     `tfsdk:"claims"`
}

func (m *OIDCModel) Values(h *Handler) map[string]any {
	data := sharedApplicationData(h, m.ID, m.Name, m.Description, m.Logo, m.Disabled)
	settings := map[string]any{}
	stringattr.Get(m.LoginPageURL, settings, "loginPageUrl")
	settings["claims"] = m.Claims
	data["oidc"] = settings
	return data
}

func (m *OIDCModel) SetValues(h *Handler, data map[string]any) {
	// all oidc application values are specified in the configuration
}
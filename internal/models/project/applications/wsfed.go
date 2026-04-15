package applications

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var WSFedAttributes = map[string]schema.Attribute{
	"id":          stringattr.Optional(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),
	"logo":        stringattr.Default(""),
	"disabled":    boolattr.Default(false),

	"realm":                stringattr.Default(""),
	"reply_url":            stringattr.Default(""),
	"login_page_url":       stringattr.Default(""),
	"attribute_mapping":    listattr.Default[AttributeMappingModel](AttributeMappingAttributes),
	"groups_mapping":       listattr.Default[GroupsMappingModel](GroupsMappingAttributes),
	"force_authentication": boolattr.Default(false),
	"logout_redirect_url":  stringattr.Default(""),
	"error_redirect_url":   stringattr.Default(""),
}

// Model

type WSFedModel struct {
	ID                  stringattr.Type                      `tfsdk:"id"`
	Name                stringattr.Type                      `tfsdk:"name"`
	Description         stringattr.Type                      `tfsdk:"description"`
	Logo                stringattr.Type                      `tfsdk:"logo"`
	Disabled            boolattr.Type                        `tfsdk:"disabled"`
	Realm               stringattr.Type                      `tfsdk:"realm"`
	ReplyURL            stringattr.Type                      `tfsdk:"reply_url"`
	LoginPageURL        stringattr.Type                      `tfsdk:"login_page_url"`
	AttributeMapping    listattr.Type[AttributeMappingModel] `tfsdk:"attribute_mapping"`
	GroupsMapping       listattr.Type[GroupsMappingModel]    `tfsdk:"groups_mapping"`
	ForceAuthentication boolattr.Type                        `tfsdk:"force_authentication"`
	LogoutRedirectURL   stringattr.Type                      `tfsdk:"logout_redirect_url"`
	ErrorRedirectURL    stringattr.Type                      `tfsdk:"error_redirect_url"`
}

func (m *WSFedModel) Values(h *helpers.Handler) map[string]any {
	settings := map[string]any{}
	stringattr.Get(m.Realm, settings, "realm")
	stringattr.Get(m.ReplyURL, settings, "replyUrl")
	stringattr.Get(m.LoginPageURL, settings, "loginPageUrl")
	listattr.Get(m.AttributeMapping, settings, "attributeMapping", h)
	listattr.Get(m.GroupsMapping, settings, "groupsMapping", h)
	boolattr.Get(m.ForceAuthentication, settings, "forceAuthentication")
	stringattr.Get(m.LogoutRedirectURL, settings, "logoutRedirectUrl")
	stringattr.Get(m.ErrorRedirectURL, settings, "errorRedirectUrl")

	data := sharedApplicationData(h, m.ID, m.Name, m.Description, m.Logo, m.Disabled)
	data["wsfed"] = settings
	return data
}

func (m *WSFedModel) SetValues(h *helpers.Handler, data map[string]any) {
	setSharedApplicationData(h, data, &m.ID, &m.Name, &m.Description, &m.Logo, &m.Disabled)
	if settings, ok := data["wsfed"].(map[string]any); ok {
		stringattr.Set(&m.Realm, settings, "realm")
		stringattr.Set(&m.ReplyURL, settings, "replyUrl")
		stringattr.Nil(&m.LoginPageURL) // XXX reset by the backend on response for now
		listattr.Set(&m.AttributeMapping, settings, "attributeMapping", h)
		listattr.Set(&m.GroupsMapping, settings, "groupsMapping", h)
		boolattr.Set(&m.ForceAuthentication, settings, "forceAuthentication")
		stringattr.Set(&m.LogoutRedirectURL, settings, "logoutRedirectUrl")
		stringattr.Set(&m.ErrorRedirectURL, settings, "errorRedirectUrl")
	}
}

// Matching

func (m *WSFedModel) GetName() stringattr.Type {
	return m.Name
}

func (m *WSFedModel) GetID() stringattr.Type {
	return m.ID
}

func (m *WSFedModel) SetID(id stringattr.Type) {
	m.ID = id
}

// Groups Mapping

var GroupsMappingAttributes = map[string]schema.Attribute{
	"name":        stringattr.Required(),
	"type":        stringattr.Required(),
	"filter_type": stringattr.Required(),
	"value":       stringattr.Required(),
	"roles":       listattr.Default[RoleGroupMappingModel](RoleGroupMappingAttributes),
}

type GroupsMappingModel struct {
	Name       stringattr.Type                      `tfsdk:"name"`
	Type       stringattr.Type                      `tfsdk:"type"`
	FilterType stringattr.Type                      `tfsdk:"filter_type"`
	Value      stringattr.Type                      `tfsdk:"value"`
	Roles      listattr.Type[RoleGroupMappingModel] `tfsdk:"roles"`
}

func (m *GroupsMappingModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Type, data, "type")
	stringattr.Get(m.FilterType, data, "filterType")
	stringattr.Get(m.Value, data, "value")
	listattr.Get(m.Roles, data, "roles", h)
	return data
}

func (m *GroupsMappingModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Type, data, "type")
	stringattr.Set(&m.FilterType, data, "filterType")
	stringattr.Set(&m.Value, data, "value")
	listattr.Set(&m.Roles, data, "roles", h)
}

// Role Group Mapping

var RoleGroupMappingAttributes = map[string]schema.Attribute{
	"id":   stringattr.Required(),
	"name": stringattr.Required(),
}

type RoleGroupMappingModel struct {
	ID   stringattr.Type `tfsdk:"id"`
	Name stringattr.Type `tfsdk:"name"`
}

func (m *RoleGroupMappingModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ID, data, "id")
	stringattr.Get(m.Name, data, "name")
	return data
}

func (m *RoleGroupMappingModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ID, data, "id")
	stringattr.Set(&m.Name, data, "name")
}

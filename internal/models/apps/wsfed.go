package apps

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var WSFedAppSchema = schema.Schema{
	MarkdownDescription: "Manages a single WS-Federation federated application where Descope acts as the identity provider.",
	Attributes:          WSFedAppAttributes,
}

var WSFedAppAttributes = map[string]schema.Attribute{
	"id":                  stringattr.Optional(stringplanmodifier.RequiresReplace(), appIDValidator),
	"project_id":          stringattr.Required(stringplanmodifier.RequiresReplace()),
	"deletion_protection": boolattr.Tristate(),
	"name":                stringattr.Required(stringattr.StandardLenValidator),
	"description":         stringattr.Default("", stringvalidator.LengthAtMost(100000)),
	"logo":                stringattr.Default(""),
	"disabled":            boolattr.Default(false),

	"realm":                       stringattr.Default(""),
	"reply_url":                   stringattr.Default(""),
	"reply_allowed_callback_urls": strsetattr.Default(),
	"login_page_url":              stringattr.Optional(),
	"attribute_mapping":           listattr.Default[AttributeMappingModel](AttributeMappingAttributes),
	"groups_mapping":              listattr.Default[GroupsMappingModel](GroupsMappingAttributes),
	"force_authentication":        boolattr.Default(false),
	"logout_redirect_url":         stringattr.Default(""),
	"error_redirect_url":          stringattr.Default(""),
}

// Model

type WSFedAppModel struct {
	ID                 stringattr.Type `tfsdk:"id"`
	ProjectID          stringattr.Type `tfsdk:"project_id"`
	DeletionProtection boolattr.Type   `tfsdk:"deletion_protection"`
	Name               stringattr.Type `tfsdk:"name"`
	Description        stringattr.Type `tfsdk:"description"`
	Logo               stringattr.Type `tfsdk:"logo"`
	Disabled           boolattr.Type   `tfsdk:"disabled"`

	Realm                    stringattr.Type                      `tfsdk:"realm"`
	ReplyURL                 stringattr.Type                      `tfsdk:"reply_url"`
	ReplyAllowedCallbackURLs strsetattr.Type                      `tfsdk:"reply_allowed_callback_urls"`
	LoginPageURL             stringattr.Type                      `tfsdk:"login_page_url"`
	AttributeMapping         listattr.Type[AttributeMappingModel] `tfsdk:"attribute_mapping"`
	GroupsMapping            listattr.Type[GroupsMappingModel]    `tfsdk:"groups_mapping"`
	ForceAuthentication      boolattr.Type                        `tfsdk:"force_authentication"`
	LogoutRedirectURL        stringattr.Type                      `tfsdk:"logout_redirect_url"`
	ErrorRedirectURL         stringattr.Type                      `tfsdk:"error_redirect_url"`
}

func (m *WSFedAppModel) Values(h *helpers.Handler) map[string]any {
	data := sharedAppValues(h, m.ID, m.Name, m.Description, m.Logo, m.Disabled)
	stringattr.Get(m.Realm, data, "realm")
	stringattr.Get(m.ReplyURL, data, "replyUrl")
	strsetattr.Get(m.ReplyAllowedCallbackURLs, data, "replyAllowedCallbacks", h)
	stringattr.Get(m.LoginPageURL, data, "loginPageUrl")
	listattr.Get(m.AttributeMapping, data, "attributeMapping", h)
	listattr.Get(m.GroupsMapping, data, "groupsMapping", h)
	boolattr.Get(m.ForceAuthentication, data, "forceAuthentication")
	stringattr.Get(m.LogoutRedirectURL, data, "logoutRedirectUrl")
	stringattr.Get(m.ErrorRedirectURL, data, "errorRedirectUrl")
	return data
}

func (m *WSFedAppModel) SetValues(h *helpers.Handler, data map[string]any) {
	setSharedAppValues(h, data, &m.ID, &m.Name, &m.Description, &m.Logo, &m.Disabled)
	if settings, ok := data["wsfedSettings"].(map[string]any); ok {
		stringattr.Set(&m.Realm, settings, "realm")
		stringattr.Set(&m.ReplyURL, settings, "replyUrl")
		strsetattr.Set(&m.ReplyAllowedCallbackURLs, settings, "replyAllowedCallbacks", h)
		stringattr.Set(&m.LoginPageURL, settings, "loginPageUrl")
		listattr.Set(&m.AttributeMapping, settings, "attributeMapping", h)
		listattr.Set(&m.GroupsMapping, settings, "groupsMapping", h)
		boolattr.Set(&m.ForceAuthentication, settings, "forceAuthentication")
		stringattr.Set(&m.LogoutRedirectURL, settings, "logoutRedirectUrl")
		stringattr.Set(&m.ErrorRedirectURL, settings, "errorRedirectUrl")
	}
}

func (m *WSFedAppModel) DeletionProtectionDefault(_ context.Context) bool {
	return true
}

func (m *WSFedAppModel) GetID() stringattr.Type {
	return m.ID
}

func (m *WSFedAppModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *WSFedAppModel) GetProjectID() stringattr.Type {
	return m.ProjectID
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

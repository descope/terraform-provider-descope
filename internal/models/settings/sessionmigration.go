package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_session_migration is the project-level session migration settings singleton (id = project_id).
// An empty vendor (the default) means session migration is disabled.

var SessionMigrationSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level session migration settings for migrating user sessions from another vendor. This is a singleton resource, and its id is always the project ID.",
	Attributes:          SessionMigrationSettingsAttributes,
}

var SessionMigrationSettingsAttributes = map[string]schema.Attribute{
	"id":                         stringattr.Identifier(),
	"project_id":                 stringattr.Required(stringplanmodifier.RequiresReplace()),
	"vendor":                     stringattr.Default("", stringattr.StandardLenValidator),
	"client_id":                  stringattr.Default("", stringattr.StandardLenValidator),
	"domain":                     stringattr.Default("", stringattr.StandardLenValidator),
	"audience":                   stringattr.Default("", stringattr.StandardLenValidator),
	"issuer":                     stringattr.Default("", stringattr.StandardLenValidator),
	"api_token":                  stringattr.SecretOptional(),
	"loginid_matched_attributes": strsetattr.Default(stringattr.StandardLenValidator),
	"user_sync_type":             stringattr.Default("", stringvalidator.OneOf("", "matchOnly", "jit")),
	"user_mapping":               listattr.Default[UserMappingItemModel](UserMappingItemAttributes),
}

var UserMappingItemAttributes = map[string]schema.Attribute{
	"external_key": stringattr.Required(stringattr.StandardLenValidator),
	"descope_key":  stringattr.Required(stringattr.StandardLenValidator),
}

type UserMappingItemModel struct {
	ExternalKey stringattr.Type `tfsdk:"external_key"`
	DescopeKey  stringattr.Type `tfsdk:"descope_key"`
}

func (m *UserMappingItemModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ExternalKey, data, "externalKey")
	stringattr.Get(m.DescopeKey, data, "descopeKey")
	return data
}

func (m *UserMappingItemModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ExternalKey, data, "externalKey")
	stringattr.Set(&m.DescopeKey, data, "descopeKey")
}

type SessionMigrationSettingsModel struct {
	ID                       stringattr.Type                     `tfsdk:"id"`
	ProjectID                stringattr.Type                     `tfsdk:"project_id"`
	Vendor                   stringattr.Type                     `tfsdk:"vendor"`
	ClientID                 stringattr.Type                     `tfsdk:"client_id"`
	Domain                   stringattr.Type                     `tfsdk:"domain"`
	Audience                 stringattr.Type                     `tfsdk:"audience"`
	Issuer                   stringattr.Type                     `tfsdk:"issuer"`
	APIToken                 stringattr.Type                     `tfsdk:"api_token"`
	LoginIDMatchedAttributes strsetattr.Type                     `tfsdk:"loginid_matched_attributes"`
	UserSyncType             stringattr.Type                     `tfsdk:"user_sync_type"`
	UserMapping              listattr.Type[UserMappingItemModel] `tfsdk:"user_mapping"`
}

func (m *SessionMigrationSettingsModel) Values(h *helpers.Handler) map[string]any {
	migration := map[string]any{}
	stringattr.Get(m.Vendor, migration, "vendor")
	switch m.Vendor.ValueString() {
	case "auth0":
		stringattr.Get(m.ClientID, migration, "clientId")
		stringattr.Get(m.Domain, migration, "domain")
		stringattr.Get(m.Audience, migration, "audience")
	case "okta":
		stringattr.Get(m.ClientID, migration, "clientId")
		stringattr.Get(m.Issuer, migration, "issuer")
		stringattr.Get(m.APIToken, migration, "apiToken")
	}
	loginIDs := map[string]any{}
	strsetattr.Get(m.LoginIDMatchedAttributes, loginIDs, "loginIdExternalUserSources", h)
	migration["loginIdExternalUserSources"] = loginIDs["loginIdExternalUserSources"]
	stringattr.Get(m.UserSyncType, migration, "userSyncType")
	listattr.Get(m.UserMapping, migration, "userMapping", h)
	return map[string]any{"sessionMigration": migration}
}

func (m *SessionMigrationSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	migration, _ := data["sessionMigration"].(map[string]any)
	if migration == nil {
		migration = map[string]any{}
	}
	stringattr.Set(&m.Vendor, migration, "vendor")
	stringattr.Set(&m.ClientID, migration, "clientId")
	stringattr.Set(&m.Domain, migration, "domain")
	stringattr.Set(&m.Audience, migration, "audience")
	stringattr.Set(&m.Issuer, migration, "issuer")
	// api_token is write-only: never read back from the server, the config value is authoritative.
	strsetattr.Set(&m.LoginIDMatchedAttributes, migration, "loginIdExternalUserSources", h)
	stringattr.Set(&m.UserSyncType, migration, "userSyncType")
	listattr.Set(&m.UserMapping, migration, "userMapping", h)
}

func (m *SessionMigrationSettingsModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Vendor, m.ClientID, m.Domain, m.Audience, m.Issuer, m.APIToken, m.LoginIDMatchedAttributes, m.UserSyncType, m.UserMapping) {
		return
	}

	vendor := m.Vendor.ValueString()

	switch vendor {
	case "":
		if m.ClientID.ValueString() != "" || m.Domain.ValueString() != "" || m.Audience.ValueString() != "" || m.Issuer.ValueString() != "" || m.APIToken.ValueString() != "" || !m.LoginIDMatchedAttributes.IsEmpty() || m.UserSyncType.ValueString() != "" || !m.UserMapping.IsEmpty() {
			h.Invalid("The other session migration attributes must not be set when vendor is not specified")
		}
		return
	case "auth0":
		if m.Domain.ValueString() == "" {
			h.Missing("The domain attribute is required for %s session migration", vendor)
		}
		if m.Issuer.ValueString() != "" {
			h.Invalid("The issuer attribute should not be set for %s session migration", vendor)
		}
	case "okta":
		if m.Domain.ValueString() != "" {
			h.Invalid("The domain attribute should not be set for %s session migration", vendor)
		}
		if m.Issuer.ValueString() == "" {
			h.Missing("The issuer attribute is required for %s session migration", vendor)
		}
		if m.Audience.ValueString() != "" {
			h.Invalid("The audience attribute should not be set for %s session migration", vendor)
		}
	default:
		h.Invalid("Unsupported session migration vendor: %s", vendor)
	}

	if m.ClientID.ValueString() == "" {
		h.Missing("The client_id attribute is required for %s session migration", vendor)
	}
	if m.LoginIDMatchedAttributes.IsEmpty() {
		h.Missing("The loginid_matched_attributes attribute is expected to be a non-empty list for %s session migration", vendor)
	}
}

func (m *SessionMigrationSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *SessionMigrationSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *SessionMigrationSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }

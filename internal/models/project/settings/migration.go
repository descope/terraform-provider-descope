package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SessionMigrationValidator = objattr.NewValidator[SessionMigrationModel]("must have a valid configuration")

var SessionMigrationAttributes = map[string]schema.Attribute{
	"vendor":                     stringattr.Default("", stringattr.StandardLenValidator),
	"client_id":                  stringattr.Default("", stringattr.StandardLenValidator),
	"domain":                     stringattr.Default("", stringattr.StandardLenValidator),
	"audience":                   stringattr.Default("", stringattr.StandardLenValidator),
	"issuer":                     stringattr.Default("", stringattr.StandardLenValidator),
	"api_token":                  stringattr.SecretOptional(),
	"loginid_matched_attributes": strsetattr.Default(stringattr.StandardLenValidator),
	"user_sync_type":             stringattr.Optional(stringvalidator.OneOf("matchOnly", "jit")),
	"user_mapping":               listattr.Default[ExternalAuthUserMappingItemModel](ExternalAuthUserMappingItemAttributes),
}

var ExternalAuthUserMappingItemAttributes = map[string]schema.Attribute{
	"external_key": stringattr.Required(stringattr.StandardLenValidator),
	"descope_key":  stringattr.Required(stringattr.StandardLenValidator),
}

type ExternalAuthUserMappingItemModel struct {
	ExternalKey stringattr.Type `tfsdk:"external_key"`
	DescopeKey  stringattr.Type `tfsdk:"descope_key"`
}

func (m *ExternalAuthUserMappingItemModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.ExternalKey, data, "externalKey")
	stringattr.Get(m.DescopeKey, data, "descopeKey")
	return data
}

func (m *ExternalAuthUserMappingItemModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.ExternalKey, data, "externalKey")
	stringattr.Set(&m.DescopeKey, data, "descopeKey")
}

type SessionMigrationModel struct {
	Vendor                   stringattr.Type                                 `tfsdk:"vendor"`
	ClientID                 stringattr.Type                                 `tfsdk:"client_id"`
	Domain                   stringattr.Type                                 `tfsdk:"domain"`
	Audience                 stringattr.Type                                 `tfsdk:"audience"`
	Issuer                   stringattr.Type                                 `tfsdk:"issuer"`
	ApiToken                 stringattr.Type                                 `tfsdk:"api_token"`
	LoginIDMatchedAttributes strsetattr.Type                                 `tfsdk:"loginid_matched_attributes"`
	UserSyncType             stringattr.Type                                 `tfsdk:"user_sync_type"`
	UserMapping              listattr.Type[ExternalAuthUserMappingItemModel] `tfsdk:"user_mapping"`
}

var SessionMigrationDefault = &SessionMigrationModel{
	Vendor:                   stringattr.Value(""),
	ClientID:                 stringattr.Value(""),
	Domain:                   stringattr.Value(""),
	Audience:                 stringattr.Value(""),
	Issuer:                   stringattr.Value(""),
	LoginIDMatchedAttributes: strsetattr.Empty(),
	UserMapping:              listattr.Empty[ExternalAuthUserMappingItemModel](),
}

func (m *SessionMigrationModel) Values(h *helpers.Handler) map[string]any {
	vendor := m.Vendor.ValueString()

	c := map[string]any{}
	stringattr.Get(m.ClientID, c, "clientId")
	switch vendor {
	case "auth0":
		stringattr.Get(m.Domain, c, "domain")
		stringattr.Get(m.Audience, c, "audience")
	case "okta":
		stringattr.Get(m.Issuer, c, "issuer")
		stringattr.Get(m.ApiToken, c, "apiToken")
	}

	data := map[string]any{}
	data[vendor] = c
	strsetattr.Get(m.LoginIDMatchedAttributes, data, "loginIdExternalUserSources", h)
	stringattr.Get(m.UserSyncType, data, "userSyncType")
	listattr.Get(m.UserMapping, data, "userMapping", h)
	return data
}

func (m *SessionMigrationModel) SetValues(h *helpers.Handler, data map[string]any) {
	if v, ok := data["auth0"].(map[string]any); ok {
		m.Vendor = stringattr.Value("auth0")
		stringattr.Set(&m.ClientID, v, "clientId")
		stringattr.Set(&m.Domain, v, "domain")
		stringattr.Set(&m.Audience, v, "audience")
		stringattr.Nil(&m.Issuer)
		stringattr.Nil(&m.ApiToken)
	} else if v, ok := data["okta"].(map[string]any); ok {
		m.Vendor = stringattr.Value("okta")
		stringattr.Set(&m.ClientID, v, "clientId")
		stringattr.Nil(&m.Domain)
		stringattr.Nil(&m.Audience)
		stringattr.Set(&m.Issuer, v, "issuer")
		// server masks the api_token with a placeholder on export; preserve prior state value
		// see: managementservice/common/pkg/common/domain/domain.go ManagementServicePlaceholderValueIndicator
		if t, ok := v["apiToken"].(string); ok && t != "" && t != "PLACEHOLDER_VALUE" {
			m.ApiToken = stringattr.Value(t)
		}
	} else {
		m.Vendor = stringattr.Value("")
		stringattr.Nil(&m.ClientID)
		stringattr.Nil(&m.Domain)
		stringattr.Nil(&m.Audience)
		stringattr.Nil(&m.Issuer)
		stringattr.Nil(&m.ApiToken)
	}
	strsetattr.Set(&m.LoginIDMatchedAttributes, data, "loginIdExternalUserSources", h)
	if v, ok := data["userSyncType"].(string); ok && v != "" {
		m.UserSyncType = stringattr.Value(v)
	} else {
		stringattr.Nil(&m.UserSyncType)
	}
	listattr.Set(&m.UserMapping, data, "userMapping", h)
}

func (m *SessionMigrationModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Vendor, m.ClientID, m.Domain, m.Audience, m.Issuer, m.ApiToken, m.UserSyncType, m.LoginIDMatchedAttributes) {
		return
	}

	vendor := m.Vendor.ValueString()

	switch vendor {
	case "":
		if m.ClientID.ValueString() != "" || m.Domain.ValueString() != "" || m.Audience.ValueString() != "" || m.Issuer.ValueString() != "" || m.ApiToken.ValueString() != "" || m.UserSyncType.ValueString() != "" || !m.LoginIDMatchedAttributes.IsEmpty() || !m.UserMapping.IsEmpty() {
			h.Invalid("The other session_migration attributes must not be set when vendor is not specified")
		}
		return
	case "auth0":
		if m.Domain.ValueString() == "" {
			h.Missing("The domain attribute is required for %s session migration", vendor)
		}
		if m.Issuer.ValueString() != "" {
			h.Invalid("The issuer attribute should not be set for %s session migration", vendor)
		}
		if m.ApiToken.ValueString() != "" {
			h.Invalid("The api_token attribute should not be set for %s session migration", vendor)
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
		if m.ApiToken.ValueString() == "" {
			h.Missing("The api_token attribute is required for %s session migration", vendor)
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

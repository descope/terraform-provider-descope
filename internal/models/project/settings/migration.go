package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SessionMigrationValidator = objattr.NewValidator[SessionMigrationModel]("must have a valid configuration")

var SessionMigrationAttributes = map[string]schema.Attribute{
	"vendor":                     stringattr.Default("", stringattr.StandardLenValidator),
	"client_id":                  stringattr.Default("", stringattr.StandardLenValidator),
	"domain":                     stringattr.Default("", stringattr.StandardLenValidator),
	"audience":                   stringattr.Default("", stringattr.StandardLenValidator),
	"issuer":                     stringattr.Default("", stringattr.StandardLenValidator),
	"loginid_matched_attributes": strsetattr.Default(stringattr.StandardLenValidator),
}

type SessionMigrationModel struct {
	Vendor                   stringattr.Type `tfsdk:"vendor"`
	ClientID                 stringattr.Type `tfsdk:"client_id"`
	Domain                   stringattr.Type `tfsdk:"domain"`
	Audience                 stringattr.Type `tfsdk:"audience"`
	Issuer                   stringattr.Type `tfsdk:"issuer"`
	LoginIDMatchedAttributes strsetattr.Type `tfsdk:"loginid_matched_attributes"`
}

var SessionMigrationDefaujlt = &SessionMigrationModel{
	Vendor:                   stringattr.Value(""),
	ClientID:                 stringattr.Value(""),
	Domain:                   stringattr.Value(""),
	Audience:                 stringattr.Value(""),
	Issuer:                   stringattr.Value(""),
	LoginIDMatchedAttributes: strsetattr.Empty(),
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
	}

	data := map[string]any{}
	data[vendor] = c
	strsetattr.Get(m.LoginIDMatchedAttributes, data, "loginIdExternalUserSources", h)
	return data
}

func (m *SessionMigrationModel) SetValues(h *helpers.Handler, data map[string]any) {
	if v, ok := data["auth0"].(map[string]any); ok {
		m.Vendor = stringattr.Value("auth0")
		stringattr.Set(&m.ClientID, v, "clientId")
		stringattr.Set(&m.Domain, v, "domain")
		stringattr.Set(&m.Audience, v, "audience")
		stringattr.Nil(&m.Issuer)
	} else if v, ok := data["okta"].(map[string]any); ok {
		m.Vendor = stringattr.Value("okta")
		stringattr.Set(&m.ClientID, v, "clientId")
		stringattr.Nil(&m.Domain)
		stringattr.Nil(&m.Audience)
		stringattr.Set(&m.Issuer, v, "issuer")
	} else {
		m.Vendor = stringattr.Value("")
		stringattr.Nil(&m.ClientID)
		stringattr.Nil(&m.Domain)
		stringattr.Nil(&m.Audience)
		stringattr.Nil(&m.Issuer)
	}
	strsetattr.Set(&m.LoginIDMatchedAttributes, data, "loginIdExternalUserSources", h)
}

func (m *SessionMigrationModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Vendor, m.ClientID, m.Domain, m.Audience, m.Issuer, m.LoginIDMatchedAttributes) {
		return
	}

	vendor := m.Vendor.ValueString()

	switch vendor {
	case "":
		if m.ClientID.ValueString() != "" || m.Domain.ValueString() != "" || m.Audience.ValueString() != "" || m.Issuer.ValueString() != "" || !m.LoginIDMatchedAttributes.IsEmpty() {
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
	case "okta":
		if m.Domain.ValueString() != "" {
			h.Missing("The domain attribute should not be set for %s session migration", vendor)
		}
		if m.Issuer.ValueString() == "" {
			h.Invalid("The issuer attribute is required for %s session migration", vendor)
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

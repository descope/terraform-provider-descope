package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/floatattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var SmtpStaticIpAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"host":           stringattr.Required(),
	"port":           floatattr.Default(25),
	"auth_method":    stringattr.Default("PLAIN"),
	"username":       stringattr.Required(),
	"password":       stringattr.SecretRequired(),
	"from_email":     stringattr.Required(),
	"from_name":      stringattr.Default(""),
	"use_static_ips": boolattr.Default(false),
}

// Model

type SmtpStaticIpModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	Host         types.String  `tfsdk:"host"`
	Port         types.Float64 `tfsdk:"port"`
	AuthMethod   types.String  `tfsdk:"auth_method"`
	Username     types.String  `tfsdk:"username"`
	Password     types.String  `tfsdk:"password"`
	FromEmail    types.String  `tfsdk:"from_email"`
	FromName     types.String  `tfsdk:"from_name"`
	UseStaticIPs types.Bool    `tfsdk:"use_static_ips"`
}

func (m *SmtpStaticIpModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "smtp-static-ip"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *SmtpStaticIpModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *SmtpStaticIpModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.Host, c, "host")
	floatattr.Get(m.Port, c, "port")
	stringattr.Get(m.AuthMethod, c, "authMethod")
	stringattr.Get(m.Username, c, "username")
	stringattr.Get(m.Password, c, "password")
	stringattr.Get(m.FromEmail, c, "fromEmail")
	stringattr.Get(m.FromName, c, "fromName")
	boolattr.Get(m.UseStaticIPs, c, "useStaticIps")
	return c
}

// Matching

func (m *SmtpStaticIpModel) GetName() types.String {
	return m.Name
}

func (m *SmtpStaticIpModel) GetID() types.String {
	return m.ID
}

func (m *SmtpStaticIpModel) SetID(id types.String) {
	m.ID = id
}

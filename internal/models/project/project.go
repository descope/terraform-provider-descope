package project

import (
	"context"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/mapattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/project/applications"
	"github.com/descope/terraform-provider-descope/internal/models/project/attributes"
	"github.com/descope/terraform-provider-descope/internal/models/project/authentication"
	"github.com/descope/terraform-provider-descope/internal/models/project/authorization"
	"github.com/descope/terraform-provider-descope/internal/models/project/connectors"
	"github.com/descope/terraform-provider-descope/internal/models/project/flows"
	"github.com/descope/terraform-provider-descope/internal/models/project/jwttemplates"
	"github.com/descope/terraform-provider-descope/internal/models/project/settings"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ProjectAttributes = map[string]schema.Attribute{
	"id":               stringattr.Identifier(),
	"name":             stringattr.Required(),
	"environment":      stringattr.Optional(stringvalidator.OneOf("", "production")),
	"tags":             strlistattr.Optional(listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 50))),
	"project_settings": objectattr.Optional(settings.SettingsAttributes, settings.SettingsValidator),
	"invite_settings":  objectattr.Optional(settings.InviteSettingsAttributes),
	"authentication":   objectattr.Optional(authentication.AuthenticationAttributes),
	"authorization":    objectattr.Optional(authorization.AuthorizationAttributes, authorization.AuthorizationValidator),
	"attributes":       objectattr.Optional(attributes.AttributesAttributes),
	"connectors":       objectattr.Optional(connectors.ConnectorsAttributes, connectors.ConnectorsModifier, connectors.ConnectorsValidator),
	"applications":     objectattr.Optional(applications.ApplicationAttributes, applications.ApplicationValidator),
	"jwt_templates":    objectattr.Optional(jwttemplates.JWTTemplatesAttributes, jwttemplates.JWTTemplatesValidator),
	"styles":           objectattr.Optional(flows.StylesAttributes, flows.StylesValidator),
	"flows":            mapattr.Optional(flows.FlowAttributes, flows.FlowsValidator),
}

type ProjectModel struct {
	ID             types.String                        `tfsdk:"id"`
	Name           types.String                        `tfsdk:"name"`
	Environment    types.String                        `tfsdk:"environment"`
	Tags           []string                            `tfsdk:"tags"`
	Settings       *settings.SettingsModel             `tfsdk:"project_settings"`
	Invite         *settings.InviteSettingsModel       `tfsdk:"invite_settings"`
	Authentication *authentication.AuthenticationModel `tfsdk:"authentication"`
	Authorization  *authorization.AuthorizationModel   `tfsdk:"authorization"`
	Attributes     *attributes.AttributesModel         `tfsdk:"attributes"`
	Connectors     *connectors.ConnectorsModel         `tfsdk:"connectors"`
	Applications   *applications.ApplicationModel      `tfsdk:"applications"`
	JWTTemplates   *jwttemplates.JWTTemplatesModel     `tfsdk:"jwt_templates"`
	Styles         *flows.StylesModel                  `tfsdk:"styles"`
	Flows          *flows.FlowsModel                   `tfsdk:"flows"` // this is just a map but use pointer to stay consistent with other models
}

func (m *ProjectModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	data["version"] = helpers.ModelVersion
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Environment, data, "environment")
	strlistattr.Get(m.Tags, data, "tags", h)
	objectattr.Get(m.Settings, data, "settings", h)
	objectattr.Get(m.Invite, data, "settings", h)
	objectattr.Get(m.Authentication, data, "authentication", h)
	objectattr.Get(m.Connectors, data, "connectors", h)
	objectattr.Get(m.Applications, data, "applications", h)
	objectattr.Get(m.Authorization, data, "authorization", h)
	objectattr.Get(m.Attributes, data, "attributes", h)
	objectattr.Get(m.JWTTemplates, data, "jwtTemplates", h)
	objectattr.Get(m.Styles, data, "styles", h)
	mapattr.Get(m.Flows, data, "flows", h)
	return data
}

func (m *ProjectModel) SetValues(h *helpers.Handler, data map[string]any) {
	if v, ok := data["version"].(float64); ok {
		helpers.EnsureModelVersion(v, h.Diagnostics)
	}

	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Environment, data, "environment")
	strlistattr.Set(&m.Tags, data, "tags", h)
	objectattr.Set(&m.Settings, data, "settings", h)
	objectattr.Set(&m.Invite, data, "settings", h)
	objectattr.Set(&m.Authentication, data, "authentication", h)
	objectattr.Set(&m.Connectors, data, "connectors", h)
	objectattr.Set(&m.Applications, data, "applications", h)
	objectattr.Set(&m.Authorization, data, "authorization", h)
	objectattr.Set(&m.Attributes, data, "attributes", h)
	objectattr.Set(&m.JWTTemplates, data, "jwtTemplates", h)
	objectattr.Set(&m.Styles, data, "styles", h)
	mapattr.Set(&m.Flows, data, "flows", h)
}

func (m *ProjectModel) References(ctx context.Context) helpers.ReferencesMap {
	refs := helpers.ReferencesMap{}
	if m.Connectors != nil {
		maps.Copy(refs, m.Connectors.References(ctx))
	}
	if m.Authorization != nil {
		maps.Copy(refs, m.Authorization.References(ctx))
	}
	if m.JWTTemplates != nil {
		maps.Copy(refs, m.JWTTemplates.References(ctx))
	}
	return refs
}

func (m *ProjectModel) SetReferences(h *helpers.Handler) {
	if m.Authentication != nil {
		m.Authentication.SetReferences(h)
	}
	if m.Settings != nil {
		m.Settings.SetReferences(h)
	}
}

package project

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/mapattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/project/applications"
	"github.com/descope/terraform-provider-descope/internal/models/project/attributes"
	"github.com/descope/terraform-provider-descope/internal/models/project/authentication"
	"github.com/descope/terraform-provider-descope/internal/models/project/authorization"
	"github.com/descope/terraform-provider-descope/internal/models/project/connectors"
	"github.com/descope/terraform-provider-descope/internal/models/project/flows"
	"github.com/descope/terraform-provider-descope/internal/models/project/jwttemplates"
	"github.com/descope/terraform-provider-descope/internal/models/project/settings"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ProjectAttributes = map[string]schema.Attribute{
	"id":               stringattr.Identifier(),
	"name":             stringattr.Required(),
	"environment":      stringattr.Optional(stringvalidator.OneOf("", "production")),
	"tags":             strsetattr.Optional(stringvalidator.LengthBetween(1, 50)),
	"project_settings": objattr.Optional[settings.SettingsModel](settings.SettingsAttributes, settings.SettingsValidator),
	"invite_settings":  objattr.Default(settings.InviteSettingsDefault, settings.InviteSettingsAttributes),
	"authentication":   objattr.Optional[authentication.AuthenticationModel](authentication.AuthenticationAttributes),
	"authorization":    objattr.Optional[authorization.AuthorizationModel](authorization.AuthorizationAttributes, authorization.AuthorizationValidator),
	"attributes":       objattr.Default(attributes.AttributesDefault, attributes.AttributesAttributes),
	"connectors":       objattr.Optional[connectors.ConnectorsModel](connectors.ConnectorsAttributes, connectors.ConnectorsModifier, connectors.ConnectorsValidator),
	"applications":     objattr.Default(applications.ApplicationsDefault, applications.ApplicationsAttributes, applications.ApplicationsValidator),
	"jwt_templates":    objattr.Default(jwttemplates.JWTTemplatesDefault, jwttemplates.JWTTemplatesAttributes, jwttemplates.JWTTemplatesValidator),
	"styles":           objattr.Optional[flows.StylesModel](flows.StylesAttributes, flows.StylesValidator),
	"flows":            mapattr.Optional[flows.FlowModel](flows.FlowAttributes, flows.FlowsValidator),
}

type ProjectModel struct {
	ID             stringattr.Type                                  `tfsdk:"id"`
	Name           stringattr.Type                                  `tfsdk:"name"`
	Environment    stringattr.Type                                  `tfsdk:"environment"`
	Tags           strsetattr.Type                                  `tfsdk:"tags"`
	Settings       objattr.Type[settings.SettingsModel]             `tfsdk:"project_settings"`
	Invite         objattr.Type[settings.InviteSettingsModel]       `tfsdk:"invite_settings"`
	Authentication objattr.Type[authentication.AuthenticationModel] `tfsdk:"authentication"`
	Authorization  objattr.Type[authorization.AuthorizationModel]   `tfsdk:"authorization"`
	Attributes     objattr.Type[attributes.AttributesModel]         `tfsdk:"attributes"`
	Connectors     objattr.Type[connectors.ConnectorsModel]         `tfsdk:"connectors"`
	Applications   objattr.Type[applications.ApplicationsModel]     `tfsdk:"applications"`
	JWTTemplates   objattr.Type[jwttemplates.JWTTemplatesModel]     `tfsdk:"jwt_templates"`
	Styles         objattr.Type[flows.StylesModel]                  `tfsdk:"styles"`
	Flows          mapattr.Type[flows.FlowModel]                    `tfsdk:"flows"`
}

func (m *ProjectModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	data["version"] = helpers.ModelVersion
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Environment, data, "environment")
	strsetattr.Get(m.Tags, data, "tags", h)
	objattr.Get(m.Settings, data, "settings", h)
	objattr.Get(m.Invite, data, "settings", h)
	objattr.Get(m.Authentication, data, "authentication", h)
	objattr.Get(m.Connectors, data, "connectors", h)
	objattr.Get(m.Applications, data, "applications", h)
	objattr.Get(m.Authorization, data, "authorization", h)
	objattr.Get(m.Attributes, data, "attributes", h)
	objattr.Get(m.JWTTemplates, data, "jwtTemplates", h)
	objattr.Get(m.Styles, data, "styles", h)
	if !m.Flows.IsNull() && !m.Flows.IsUnknown() {
		data["flows"] = flows.Get(m.Flows, h)
	}
	return data
}

func (m *ProjectModel) SetValues(h *helpers.Handler, data map[string]any) {
	if v, ok := data["version"].(float64); ok {
		helpers.EnsureModelVersion(v, h.Diagnostics)
	}

	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Environment, data, "environment")
	strsetattr.Set(&m.Tags, data, "tags", h)
	objattr.Set(&m.Settings, data, "settings", h)
	objattr.Set(&m.Invite, data, "settings", h)
	objattr.Set(&m.Authentication, data, "authentication", h)
	objattr.Set(&m.Connectors, data, "connectors", h)
	objattr.Set(&m.Applications, data, "applications", h)
	objattr.Set(&m.Authorization, data, "authorization", h)
	objattr.Set(&m.Attributes, data, "attributes", h)
	objattr.Set(&m.JWTTemplates, data, "jwtTemplates", h)
	objattr.Set(&m.Styles, data, "styles", h)
	if m.Flows.IsEmpty() {
		flows.Set(&m.Flows, data, "flows", h)
	}
}

func (m *ProjectModel) CollectReferences(h *helpers.Handler) {
	objattr.CollectReferences(m.Connectors, h)
	objattr.CollectReferences(m.Authorization, h)
	objattr.CollectReferences(m.JWTTemplates, h)
}

func (m *ProjectModel) UpdateReferences(h *helpers.Handler) {
	objattr.UpdateReferences(&m.Authentication, h)
	objattr.UpdateReferences(&m.Settings, h)
}

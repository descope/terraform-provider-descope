package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/floatattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var RecaptchaEnterpriseValidator = objectattr.NewValidator[RecaptchaEnterpriseModel]("must have a valid configuration")

var RecaptchaEnterpriseAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"project_id":          stringattr.Required(),
	"site_key":            stringattr.Required(),
	"api_key":             stringattr.SecretRequired(),
	"override_assessment": boolattr.Default(false),
	"assessment_score":    floatattr.Default(0.5),
}

// Model

type RecaptchaEnterpriseModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	ProjectID          types.String  `tfsdk:"project_id"`
	SiteKey            types.String  `tfsdk:"site_key"`
	APIKey             types.String  `tfsdk:"api_key"`
	OverrideAssessment types.Bool    `tfsdk:"override_assessment"`
	AssessmentScore    types.Float64 `tfsdk:"assessment_score"`
}

func (m *RecaptchaEnterpriseModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "recaptcha-enterprise"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *RecaptchaEnterpriseModel) SetValues(h *helpers.Handler, data map[string]any) {
	setConnectorValues(&m.ID, &m.Name, &m.Description, data, h)
	if c, ok := data["configuration"].(map[string]any); ok {
		stringattr.Set(&m.ProjectID, c, "projectId")
		stringattr.Set(&m.SiteKey, c, "siteKey")
		stringattr.Set(&m.APIKey, c, "apiKey")
		boolattr.Set(&m.OverrideAssessment, c, "overrideAssessment")
		floatattr.Set(&m.AssessmentScore, c, "assessmentScore")
	}
}

func (m *RecaptchaEnterpriseModel) Validate(h *helpers.Handler) {
	if !m.AssessmentScore.IsNull() && !m.OverrideAssessment.ValueBool() {
		h.Error("Invalid connector configuration", "The assessment_score field cannot be used unless override_assessment is set to true")
	}
}

// Configuration

func (m *RecaptchaEnterpriseModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.ProjectID, c, "projectId")
	stringattr.Get(m.SiteKey, c, "siteKey")
	stringattr.Get(m.APIKey, c, "apiKey")
	boolattr.Get(m.OverrideAssessment, c, "overrideAssessment")
	floatattr.Get(m.AssessmentScore, c, "assessmentScore")
	return c
}

// Matching

func (m *RecaptchaEnterpriseModel) GetName() types.String {
	return m.Name
}

func (m *RecaptchaEnterpriseModel) GetID() types.String {
	return m.ID
}

func (m *RecaptchaEnterpriseModel) SetID(id types.String) {
	m.ID = id
}

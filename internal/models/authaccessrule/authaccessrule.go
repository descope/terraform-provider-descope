package authaccessrule

import (
	"math"

	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages a single rule of the project's auth access policy. Rules are evaluated in the order set by the `descope_auth_access_policy` resource, and the first rule whose conditions all match decides whether the authentication is allowed or blocked.",
	Attributes:          AuthAccessRuleAttributes,
}

var AuthAccessRuleAttributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"project_id":  stringattr.Required(stringplanmodifier.RequiresReplace()),
	"name":        stringattr.Default("", stringattr.StandardLenValidator),
	"description": stringattr.Default(""),
	"enabled":     boolattr.Default(true),
	"effect":      stringattr.Required(stringvalidator.OneOf("allow", "block")),
	"conditions":  listattr.Default[AuthAccessConditionModel](AuthAccessConditionAttributes, AuthAccessConditionValidator),
	"expire_time": intattr.Default(0, int64validator.Between(0, math.MaxInt32)),
}

type AuthAccessRuleModel struct {
	ID          stringattr.Type                         `tfsdk:"id"`
	ProjectID   stringattr.Type                         `tfsdk:"project_id"`
	Name        stringattr.Type                         `tfsdk:"name"`
	Description stringattr.Type                         `tfsdk:"description"`
	Enabled     boolattr.Type                           `tfsdk:"enabled"`
	Effect      stringattr.Type                         `tfsdk:"effect"`
	Conditions  listattr.Type[AuthAccessConditionModel] `tfsdk:"conditions"`
	ExpireTime  intattr.Type                            `tfsdk:"expire_time"`
}

func (m *AuthAccessRuleModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Name, data, "name")
	stringattr.Get(m.Description, data, "description")
	boolattr.Get(m.Enabled, data, "enabled")
	stringattr.Get(m.Effect, data, "effect")
	listattr.Get(m.Conditions, data, "conditions", h)
	intattr.Get(m.ExpireTime, data, "expireTime")
	return data
}

func (m *AuthAccessRuleModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Name, data, "name")
	stringattr.Set(&m.Description, data, "description")
	boolattr.Set(&m.Enabled, data, "enabled")
	stringattr.Set(&m.Effect, data, "effect")
	listattr.Set(&m.Conditions, data, "conditions", h)
	intattr.Set(&m.ExpireTime, data, "expireTime")
}

func (m *AuthAccessRuleModel) GetID() stringattr.Type {
	return m.ID
}

func (m *AuthAccessRuleModel) SetID(id stringattr.Type) {
	m.ID = id
}

func (m *AuthAccessRuleModel) GetProjectID() stringattr.Type {
	return m.ProjectID
}

// AuthAccessCondition

var AuthAccessConditionValidator = objattr.NewValidator[AuthAccessConditionModel]("must have a valid operator and values")

var AuthAccessConditionAttributes = map[string]schema.Attribute{
	"key":      stringattr.Required(stringvalidator.OneOf("request.ip", "request.country", "user.email", "user.emailDomain", "user.loginId", "user.phone")),
	"operator": stringattr.Required(stringvalidator.OneOf("equal", "notEqual", "in", "notIn", "inCidr", "notInCidr")),
	"values":   strlistattr.Required(listvalidator.SizeAtLeast(1), stringattr.NonEmptyValidator),
}

type AuthAccessConditionModel struct {
	Key      stringattr.Type  `tfsdk:"key"`
	Operator stringattr.Type  `tfsdk:"operator"`
	Vals     strlistattr.Type `tfsdk:"values"`
}

func (m *AuthAccessConditionModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.Key, data, "key")
	stringattr.Get(m.Operator, data, "operator")
	strlistattr.Get(m.Vals, data, "values", h)
	return data
}

func (m *AuthAccessConditionModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.Key, data, "key")
	stringattr.Set(&m.Operator, data, "operator")
	strlistattr.Set(&m.Vals, data, "values", h)
}

func (m *AuthAccessConditionModel) Validate(h *helpers.Handler) {
	if helpers.HasUnknownValues(m.Key, m.Operator, m.Vals) {
		return
	}
	switch operator := m.Operator.ValueString(); operator {
	case "inCidr", "notInCidr":
		if m.Key.ValueString() != "request.ip" {
			h.Invalid("The %s operator is only supported with the request.ip key", operator)
		}
		fallthrough
	case "equal", "notEqual":
		if len(m.Vals.Elements()) != 1 {
			h.Invalid("The %s operator requires exactly one value", operator)
		}
	}
}

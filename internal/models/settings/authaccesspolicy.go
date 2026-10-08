package settings

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_auth_access_policy is the project-level auth access policy singleton (id = project_id).

var AuthAccessPolicySchema = schema.Schema{
	MarkdownDescription: "Manages the project's auth access policy: the evaluation order of its `descope_auth_access_rule` resources and the action taken when no rule matches. This is a singleton resource, and its id is always the project ID. Declare only one per project. Removing this resource stops Terraform from managing the policy but leaves it in place, so set `default_action` to `allow` first if the policy should stop blocking.",
	Attributes:          AuthAccessPolicyAttributes,
}

var AuthAccessPolicyAttributes = map[string]schema.Attribute{
	"id":             stringattr.Identifier(),
	"project_id":     stringattr.Required(stringplanmodifier.RequiresReplace()),
	"default_action": stringattr.Default("allow", stringvalidator.OneOf("allow", "block")),
	"rule_ids":       strlistattr.Default(listvalidator.UniqueValues()),
}

type AuthAccessPolicyModel struct {
	ID            stringattr.Type  `tfsdk:"id"`
	ProjectID     stringattr.Type  `tfsdk:"project_id"`
	DefaultAction stringattr.Type  `tfsdk:"default_action"`
	RuleIDs       strlistattr.Type `tfsdk:"rule_ids"`
}

func (m *AuthAccessPolicyModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.DefaultAction, data, "defaultAction")
	strlistattr.Get(m.RuleIDs, data, "ruleIds", h)
	return data
}

func (m *AuthAccessPolicyModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.DefaultAction, data, "defaultAction")
	strlistattr.Set(&m.RuleIDs, data, "ruleIds", h)
}

func (m *AuthAccessPolicyModel) GetID() stringattr.Type        { return m.ID }
func (m *AuthAccessPolicyModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *AuthAccessPolicyModel) GetProjectID() stringattr.Type { return m.ProjectID }

package resourcepolicy

import (
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Schema = schema.Schema{
	MarkdownDescription: "Manages the access policy connecting one Descope inbound application to one API resource.",
	Attributes:          resourcePolicyAttributes(),
}

type Model struct {
	ID                 stringattr.Type `tfsdk:"id"`
	ProjectID          stringattr.Type `tfsdk:"project_id"`
	ApplicationID      stringattr.Type `tfsdk:"application_id"`
	ResourceID         stringattr.Type `tfsdk:"resource_id"`
	UserAccessScopes   strsetattr.Type `tfsdk:"user_access_scopes"`
	ClientAccessScopes strsetattr.Type `tfsdk:"client_access_scopes"`
	AllUserScopes      boolattr.Type   `tfsdk:"all_user_scopes"`
	AllClientScopes    boolattr.Type   `tfsdk:"all_client_scopes"`
}

func (m *Model) Identity() infra.ResourcePolicyIdentity {
	return infra.ResourcePolicyIdentity{
		ApplicationID: m.ApplicationID.ValueString(),
		ResourceID:    m.ResourceID.ValueString(),
	}
}

func (m *Model) Policy(h *helpers.Handler) infra.ResourcePolicy {
	return infra.ResourcePolicy{
		ApplicationID:      m.ApplicationID.ValueString(),
		ResourceID:         m.ResourceID.ValueString(),
		UserAccessScopes:   scopes(m.UserAccessScopes, h),
		ClientAccessScopes: scopes(m.ClientAccessScopes, h),
		AllUserScopes:      m.AllUserScopes.ValueBool(),
		AllClientScopes:    m.AllClientScopes.ValueBool(),
	}
}

func (m *Model) SetPolicy(policy *infra.ResourcePolicy) {
	m.UserAccessScopes = strsetattr.Value(policy.UserAccessScopes)
	m.ClientAccessScopes = strsetattr.Value(policy.ClientAccessScopes)
	m.ApplicationID = stringattr.Value(policy.ApplicationID)
	m.ResourceID = stringattr.Value(policy.ResourceID)
	m.AllUserScopes = boolattr.Value(policy.AllUserScopes)
	m.AllClientScopes = boolattr.Value(policy.AllClientScopes)
}

func scopes(value strsetattr.Type, h *helpers.Handler) []string {
	data := map[string]any{}
	strsetattr.Get(value, data, "scopes", h)
	result, _ := data["scopes"].([]string)
	if result == nil {
		return []string{}
	}
	return result
}

func resourcePolicyAttributes() map[string]schema.Attribute {
	id := stringattr.Identifier()
	id.MarkdownDescription = "The provider identity, composed as `application_id/resource_id`."

	projectID := stringattr.Required(stringattr.StandardLenValidator, stringattr.MachineIDValidator, stringplanmodifier.RequiresReplace())
	projectID.MarkdownDescription = "The Descope project ID containing the application and resource."

	applicationID := stringattr.Required(stringattr.StandardLenValidator, stringattr.MachineIDValidator, stringplanmodifier.RequiresReplace())
	applicationID.MarkdownDescription = "The inbound application ID referenced by this policy."

	resourceID := stringattr.Required(stringattr.StandardLenValidator, stringattr.MachineIDValidator, stringplanmodifier.RequiresReplace())
	resourceID.MarkdownDescription = "The API resource ID referenced by this policy."

	userAccessScopes := strsetattr.Default(stringattr.NonEmptyValidator, stringattr.StandardLenValidator)
	userAccessScopes.MarkdownDescription = "The resource scopes granted to authenticated users."

	clientAccessScopes := strsetattr.Default(stringattr.NonEmptyValidator, stringattr.StandardLenValidator)
	clientAccessScopes.MarkdownDescription = "The resource scopes granted to clients."

	allUserScopes := boolattr.Default(false)
	allUserScopes.MarkdownDescription = "Whether authenticated users receive all scopes exposed by the resource."

	allClientScopes := boolattr.Default(false)
	allClientScopes.MarkdownDescription = "Whether clients receive all scopes exposed by the resource."

	return map[string]schema.Attribute{
		"id":                   id,
		"project_id":           projectID,
		"application_id":       applicationID,
		"resource_id":          resourceID,
		"user_access_scopes":   userAccessScopes,
		"client_access_scopes": clientAccessScopes,
		"all_user_scopes":      allUserScopes,
		"all_client_scopes":    allClientScopes,
	}
}

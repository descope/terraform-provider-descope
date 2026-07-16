package tenant

import (
	"context"
	"regexp"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var tenantIDValidators = []validator.String{
	stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9_-]{1,30}$`), "must contain 1 to 30 alphanumeric, underscore, or hyphen characters"),
}

var TenantAttributes = map[string]schema.Attribute{
	"project_id":                stringattr.ReplaceRequired(),
	"id":                        stringattr.ReplaceOptional(tenantIDValidators...),
	"name":                      stringattr.Required(stringattr.StandardLenValidator),
	"self_provisioning_domains": strlistattr.Default(stringattr.StandardLenValidator),
	"disabled":                  boolattr.Default(false),
	"enforce_sso":               boolattr.Default(false),
	"enforce_sso_exclusions":    strlistattr.Default(stringattr.StandardLenValidator),
	"federated_application_ids": strlistattr.Default(stringattr.MachineIDValidator, stringattr.StandardLenValidator),
	"parent":                    stringattr.ReplaceOptional(tenantIDValidators...),
	"role_inheritance":          stringattr.Optional(stringvalidator.OneOf("none", "userOnly")),
}

type Model struct {
	ProjectID               stringattr.Type  `tfsdk:"project_id"`
	ID                      stringattr.Type  `tfsdk:"id"`
	Name                    stringattr.Type  `tfsdk:"name"`
	SelfProvisioningDomains strlistattr.Type `tfsdk:"self_provisioning_domains"`
	Disabled                boolattr.Type    `tfsdk:"disabled"`
	EnforceSSO              boolattr.Type    `tfsdk:"enforce_sso"`
	EnforceSSOExclusions    strlistattr.Type `tfsdk:"enforce_sso_exclusions"`
	FederatedApplicationIDs strlistattr.Type `tfsdk:"federated_application_ids"`
	Parent                  stringattr.Type  `tfsdk:"parent"`
	RoleInheritance         stringattr.Type  `tfsdk:"role_inheritance"`
}

type Values struct {
	ProjectID               string
	ID                      string
	Name                    string
	SelfProvisioningDomains []string
	Disabled                bool
	EnforceSSO              bool
	EnforceSSOExclusions    []string
	FederatedApplicationIDs []string
	Parent                  string
	RoleInheritance         string
}

func (m Model) Values(ctx context.Context, diagnostics *diag.Diagnostics) Values {
	handler := helpers.NewHandler(ctx, diagnostics)
	return Values{
		ProjectID:               m.ProjectID.ValueString(),
		ID:                      m.ID.ValueString(),
		Name:                    m.Name.ValueString(),
		SelfProvisioningDomains: collect(strlistattr.Iterator(m.SelfProvisioningDomains, handler)),
		Disabled:                m.Disabled.ValueBool(),
		EnforceSSO:              m.EnforceSSO.ValueBool(),
		EnforceSSOExclusions:    collect(strlistattr.Iterator(m.EnforceSSOExclusions, handler)),
		FederatedApplicationIDs: collect(strlistattr.Iterator(m.FederatedApplicationIDs, handler)),
		Parent:                  m.Parent.ValueString(),
		RoleInheritance:         m.RoleInheritance.ValueString(),
	}
}

func (m *Model) SetTenant(ctx context.Context, value *infra.Tenant) {
	m.ID = stringattr.Value(value.ID)
	m.Name = stringattr.Value(value.Name)
	m.SelfProvisioningDomains = strlistattr.ValueContext(ctx, value.SelfProvisioningDomains)
	m.Disabled = boolattr.Value(value.Disabled)
	m.EnforceSSO = boolattr.Value(value.EnforceSSO)
	m.EnforceSSOExclusions = strlistattr.ValueContext(ctx, value.EnforceSSOExclusions)
	m.FederatedApplicationIDs = strlistattr.ValueContext(ctx, value.FederatedApplicationIDs)
	m.Parent = stringattr.Value(value.Parent)
	m.RoleInheritance = stringattr.Value(value.RoleInheritance)
}

func collect(values func(func(string) bool)) []string {
	result := []string{}
	for value := range values {
		result = append(result, value)
	}
	return result
}

package settings

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_governance is the project-level Agentic Governance Suite configuration singleton (id = project_id).

var GovernanceSchema = schema.Schema{
	MarkdownDescription: "Manages the project's Agentic Governance Suite configuration. This is a singleton resource, and its id is always the project ID.",
	Attributes:          GovernanceAttributes,
}

var GovernanceAttributes = map[string]schema.Attribute{
	"id":             stringattr.Identifier(),
	"project_id":     stringattr.Required(stringplanmodifier.RequiresReplace()),
	"configured":     boolattr.Default(false),
	"auto_approval":  boolattr.Default(false),
	"suite_disabled": boolattr.Default(false),
	"logo":           stringattr.Default(""),
	// The write is version checked, so the value read has to travel back out.
	// A string because int64 crosses protojson as one.
	"version": stringattr.Generated(),
}

type GovernanceModel struct {
	ID            stringattr.Type `tfsdk:"id"`
	ProjectID     stringattr.Type `tfsdk:"project_id"`
	Configured    boolattr.Type   `tfsdk:"configured"`
	AutoApproval  boolattr.Type   `tfsdk:"auto_approval"`
	SuiteDisabled boolattr.Type   `tfsdk:"suite_disabled"`
	Logo          stringattr.Type `tfsdk:"logo"`
	Version       stringattr.Type `tfsdk:"version"`
}

func (m *GovernanceModel) Values(_ *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Configured, data, "configured")
	boolattr.Get(m.AutoApproval, data, "autoApproval")
	boolattr.Get(m.SuiteDisabled, data, "suiteDisabled")
	stringattr.Get(m.Logo, data, "logo")
	stringattr.Get(m.Version, data, "version")
	return data
}

func (m *GovernanceModel) SetValues(_ *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Configured, data, "configured")
	boolattr.Set(&m.AutoApproval, data, "autoApproval")
	boolattr.Set(&m.SuiteDisabled, data, "suiteDisabled")
	stringattr.Set(&m.Logo, data, "logo")
	stringattr.Set(&m.Version, data, "version")
}

func (m *GovernanceModel) GetID() stringattr.Type        { return m.ID }
func (m *GovernanceModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *GovernanceModel) GetProjectID() stringattr.Type { return m.ProjectID }

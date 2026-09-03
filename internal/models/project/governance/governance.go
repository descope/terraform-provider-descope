package governance

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Agentic Governance Suite settings. Unlike admin_portal these are flat -- the
// snapshot file's whole body is the object, with no "config" wrapper.
var GovernanceAttributes = map[string]schema.Attribute{
	"configured":     boolattr.Default(false),
	"auto_approval":  boolattr.Default(false),
	"suite_disabled": boolattr.Default(false),
	// Every field the suite has is modelled here on purpose. The infra layer
	// replaces governance.json whole rather than merging, so a field left out of
	// the model is written back as its zero value -- omitting the logo would blank
	// a configured one on the next apply.
	"logo": stringattr.Default(""),
}

type GovernanceModel struct {
	Configured    boolattr.Type   `tfsdk:"configured"`
	AutoApproval  boolattr.Type   `tfsdk:"auto_approval"`
	SuiteDisabled boolattr.Type   `tfsdk:"suite_disabled"`
	Logo          stringattr.Type `tfsdk:"logo"`
}

func (m *GovernanceModel) Values(_ *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Configured, data, "configured")
	boolattr.Get(m.AutoApproval, data, "autoApproval")
	boolattr.Get(m.SuiteDisabled, data, "suiteDisabled")
	stringattr.Get(m.Logo, data, "logo")
	return data
}

func (m *GovernanceModel) SetValues(_ *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Configured, data, "configured")
	boolattr.Set(&m.AutoApproval, data, "autoApproval")
	boolattr.Set(&m.SuiteDisabled, data, "suiteDisabled")
	stringattr.Set(&m.Logo, data, "logo")
}

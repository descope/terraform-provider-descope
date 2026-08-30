package settings

import (
	"regexp"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_passkey_settings is the project-level passkey (WebAuthn) settings singleton (id = project_id).

var androidFingerprintValidator = stringvalidator.RegexMatches(
	regexp.MustCompile(`^([0-9A-Fa-f]{2}:){31}[0-9A-Fa-f]{2}$`), "must be a colon-separated SHA-256 hex fingerprint (e.g. AB:CD:EF:...)",
)

var PasskeySettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level passkey (WebAuthn) authentication settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          PasskeySettingsAttributes,
}

var PasskeySettingsAttributes = map[string]schema.Attribute{
	"id":                   stringattr.Identifier(),
	"project_id":           stringattr.Required(stringplanmodifier.RequiresReplace()),
	"disabled":             boolattr.Default(false),
	"display_name":         stringattr.Default("", stringvalidator.LengthAtMost(256)),
	"top_level_domain":     stringattr.Optional(),
	"android_fingerprints": strsetattr.Default(androidFingerprintValidator),
}

type PasskeySettingsModel struct {
	ID                  stringattr.Type `tfsdk:"id"`
	ProjectID           stringattr.Type `tfsdk:"project_id"`
	Disabled            boolattr.Type   `tfsdk:"disabled"`
	DisplayName         stringattr.Type `tfsdk:"display_name"`
	TopLevelDomain      stringattr.Type `tfsdk:"top_level_domain"`
	AndroidFingerprints strsetattr.Type `tfsdk:"android_fingerprints"`
}

func (m *PasskeySettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	stringattr.Get(m.DisplayName, data, "name")
	stringattr.Get(m.TopLevelDomain, data, "relyingPartyId")
	strsetattr.Get(m.AndroidFingerprints, data, "androidFingerprints", h)
	return data
}

func (m *PasskeySettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	stringattr.Set(&m.DisplayName, data, "name")
	stringattr.Set(&m.TopLevelDomain, data, "relyingPartyId")
	strsetattr.Set(&m.AndroidFingerprints, data, "androidFingerprints", h)
}

func (m *PasskeySettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *PasskeySettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *PasskeySettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }

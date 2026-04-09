package authentication

import (
	"regexp"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var androidFingerprintValidator = stringvalidator.RegexMatches(
	regexp.MustCompile(`^([0-9A-Fa-f]{2}:){31}[0-9A-Fa-f]{2}$`), "must be a colon-separated SHA-256 hex fingerprint (e.g. AB:CD:EF:...)",
)

var PasskeysAttributes = map[string]schema.Attribute{
	"disabled":             boolattr.Default(false),
	"top_level_domain":     stringattr.Optional(),
	"android_fingerprints": strsetattr.Default(androidFingerprintValidator),
}

type PasskeysModel struct {
	Disabled            boolattr.Type   `tfsdk:"disabled"`
	TopLevelDomain      stringattr.Type `tfsdk:"top_level_domain"`
	AndroidFingerprints strsetattr.Type `tfsdk:"android_fingerprints"`
}

func (m *PasskeysModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	stringattr.Get(m.TopLevelDomain, data, "relyingPartyId")
	strsetattr.Get(m.AndroidFingerprints, data, "androidFingerprints", h)
	return data
}

func (m *PasskeysModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	stringattr.Set(&m.TopLevelDomain, data, "relyingPartyId")
	strsetattr.Set(&m.AndroidFingerprints, data, "androidFingerprints", h)
}

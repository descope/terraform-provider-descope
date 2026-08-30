package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/durationattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// descope_password_settings is the project-level password settings singleton (id = project_id).
// The password reset email templates are managed by the descope_email_template resource and selected here by id.

var PasswordSettingsSchema = schema.Schema{
	MarkdownDescription: "Manages the project-level password authentication settings. This is a singleton resource, and its id is always the project ID.",
	Attributes:          PasswordSettingsAttributes,
}

var PasswordSettingsAttributes = map[string]schema.Attribute{
	"id":                      stringattr.Identifier(),
	"project_id":              stringattr.Required(stringplanmodifier.RequiresReplace()),
	"disabled":                boolattr.Default(false),
	"min_length":              intattr.Default(8, int64validator.Between(4, 64)),
	"lowercase":               boolattr.Default(true),
	"uppercase":               boolattr.Default(true),
	"number":                  boolattr.Default(true),
	"non_alphanumeric":        boolattr.Default(true),
	"any_letter":              boolattr.Default(false),
	"disallowed_characters":   stringattr.Default(""),
	"disallow_email_match":    boolattr.Default(false),
	"expiration":              boolattr.Default(false),
	"expiration_weeks":        intattr.Default(20, int64validator.Between(1, 999)),
	"reuse":                   boolattr.Default(false),
	"reuse_amount":            intattr.Default(10, int64validator.Between(1, 50)),
	"lock":                    boolattr.Default(false),
	"lock_attempts":           intattr.Default(5, int64validator.Between(2, 10)),
	"temporary_lock":          boolattr.Default(false),
	"temporary_lock_attempts": intattr.Default(3, int64validator.Between(1, 10)),
	"temporary_lock_duration": durationattr.Default("5 minutes", durationattr.MinimumValue("1 minute"), durationattr.MaximumValue("24 hours")),
	"enforce_strength":        stringattr.Default("none", stringvalidator.OneOf("none", "very_weak", "weak", "average", "strong", "very_strong")),
	"mask_errors":             boolattr.Default(true),
	"email_service":           objattr.Default[EmailServiceRefModel](nil, EmailServiceRefAttributes),
	"email_template_id":       stringattr.Default(""),
}

type PasswordSettingsModel struct {
	ID                    stringattr.Type                    `tfsdk:"id"`
	ProjectID             stringattr.Type                    `tfsdk:"project_id"`
	Disabled              boolattr.Type                      `tfsdk:"disabled"`
	MinLength             intattr.Type                       `tfsdk:"min_length"`
	Lowercase             boolattr.Type                      `tfsdk:"lowercase"`
	Uppercase             boolattr.Type                      `tfsdk:"uppercase"`
	Number                boolattr.Type                      `tfsdk:"number"`
	NonAlphanumeric       boolattr.Type                      `tfsdk:"non_alphanumeric"`
	AnyLetter             boolattr.Type                      `tfsdk:"any_letter"`
	DisallowedCharacters  stringattr.Type                    `tfsdk:"disallowed_characters"`
	DisallowEmailMatch    boolattr.Type                      `tfsdk:"disallow_email_match"`
	Expiration            boolattr.Type                      `tfsdk:"expiration"`
	ExpirationWeeks       intattr.Type                       `tfsdk:"expiration_weeks"`
	Reuse                 boolattr.Type                      `tfsdk:"reuse"`
	ReuseAmount           intattr.Type                       `tfsdk:"reuse_amount"`
	Lock                  boolattr.Type                      `tfsdk:"lock"`
	LockAttempts          intattr.Type                       `tfsdk:"lock_attempts"`
	TemporaryLock         boolattr.Type                      `tfsdk:"temporary_lock"`
	TemporaryLockAttempts intattr.Type                       `tfsdk:"temporary_lock_attempts"`
	TemporaryLockDuration durationattr.Type                  `tfsdk:"temporary_lock_duration"`
	EnforceStrength       stringattr.Type                    `tfsdk:"enforce_strength"`
	MaskErrors            boolattr.Type                      `tfsdk:"mask_errors"`
	EmailService          objattr.Type[EmailServiceRefModel] `tfsdk:"email_service"`
	EmailTemplateID       stringattr.Type                    `tfsdk:"email_template_id"`
}

func (m *PasswordSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	intattr.Get(m.MinLength, data, "minLength")
	boolattr.Get(m.Lowercase, data, "lowercase")
	boolattr.Get(m.Uppercase, data, "uppercase")
	boolattr.Get(m.Number, data, "number")
	boolattr.Get(m.NonAlphanumeric, data, "nonAlphanumeric")
	boolattr.Get(m.AnyLetter, data, "anyLetter")
	stringattr.Get(m.DisallowedCharacters, data, "disallowedCharacters")
	boolattr.Get(m.DisallowEmailMatch, data, "disallowEmailMatch")
	boolattr.Get(m.Expiration, data, "expiration")
	intattr.Get(m.ExpirationWeeks, data, "expirationWeeks")
	boolattr.Get(m.Reuse, data, "reuse")
	intattr.Get(m.ReuseAmount, data, "reuseAmount")
	boolattr.Get(m.Lock, data, "lock")
	intattr.Get(m.LockAttempts, data, "lockAttempts")
	boolattr.Get(m.TemporaryLock, data, "tempLock")
	intattr.Get(m.TemporaryLockAttempts, data, "tempLockAttempts")
	durationattr.GetMinutes(m.TemporaryLockDuration, data, "tempLockDuration")
	if m.EnforceStrength.ValueString() == "none" {
		data["enablePasswordStrength"] = false
		data["passwordStrengthScore"] = 0
	} else {
		data["enablePasswordStrength"] = true
		data["passwordStrengthScore"] = strengthScoreFromString(m.EnforceStrength.ValueString())
	}
	boolattr.Get(m.MaskErrors, data, "maskError")
	objattr.Get(m.EmailService, data, helpers.RootKey, h)
	stringattr.Get(m.EmailTemplateID, data, "emailTemplateId")

	useDescopeService(m.EmailService, data, "emailServiceProvider")

	return data
}

func (m *PasswordSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	intattr.Set(&m.MinLength, data, "minLength")
	boolattr.Set(&m.Lowercase, data, "lowercase")
	boolattr.Set(&m.Uppercase, data, "uppercase")
	boolattr.Set(&m.Number, data, "number")
	boolattr.Set(&m.NonAlphanumeric, data, "nonAlphanumeric")
	boolattr.Set(&m.AnyLetter, data, "anyLetter")
	stringattr.Set(&m.DisallowedCharacters, data, "disallowedCharacters")
	boolattr.Set(&m.DisallowEmailMatch, data, "disallowEmailMatch")
	boolattr.Set(&m.Expiration, data, "expiration")
	intattr.Set(&m.ExpirationWeeks, data, "expirationWeeks")
	boolattr.Set(&m.Reuse, data, "reuse")
	intattr.Set(&m.ReuseAmount, data, "reuseAmount")
	boolattr.Set(&m.Lock, data, "lock")
	intattr.Set(&m.LockAttempts, data, "lockAttempts")
	boolattr.Set(&m.TemporaryLock, data, "tempLock")
	intattr.Set(&m.TemporaryLockAttempts, data, "tempLockAttempts")
	durationattr.SetMinutes(&m.TemporaryLockDuration, data, "tempLockDuration")
	if enabled, _ := data["enablePasswordStrength"].(bool); !enabled {
		m.EnforceStrength = stringattr.Value("none")
	} else {
		score, _ := data["passwordStrengthScore"].(float64)
		m.EnforceStrength = stringattr.Value(strengthStringFromScore(int(score)))
	}
	boolattr.Set(&m.MaskErrors, data, "maskError")
	objattr.Set(&m.EmailService, data, helpers.RootKey, h)
	stringattr.Set(&m.EmailTemplateID, data, "emailTemplateId")
}

func (m *PasswordSettingsModel) GetID() stringattr.Type        { return m.ID }
func (m *PasswordSettingsModel) SetID(id stringattr.Type)      { m.ID = id }
func (m *PasswordSettingsModel) GetProjectID() stringattr.Type { return m.ProjectID }

var strengthLevels = map[string]int{
	"very_weak":   0,
	"weak":        1,
	"average":     2,
	"strong":      3,
	"very_strong": 4,
}

func strengthScoreFromString(s string) int {
	return strengthLevels[s]
}

func strengthStringFromScore(n int) string {
	for name, score := range strengthLevels {
		if score == n {
			return name
		}
	}
	return "none"
}

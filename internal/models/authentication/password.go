package authentication

import (
	"maps"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/descope/terraform-provider-descope/internal/models/templates"
	"github.com/descope/terraform-provider-descope/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var PasswordAttributes = map[string]schema.Attribute{
	"disabled":         boolattr.Default(false),
	"expiration":       boolattr.Optional(),
	"expiration_weeks": intattr.Optional(int64validator.Between(1, 999)),
	"lock":             boolattr.Optional(),
	"lock_attempts":    intattr.Optional(int64validator.Between(2, 10)),
	"lowercase":        boolattr.Optional(),
	"min_length":       intattr.Optional(int64validator.Between(4, 64)),
	"non_alphanumeric": boolattr.Optional(),
	"number":           boolattr.Optional(),
	"reuse":            boolattr.Optional(),
	"reuse_amount":     intattr.Optional(int64validator.Between(1, 50)),
	"uppercase":        boolattr.Optional(),
	"email_service":    objectattr.Optional(templates.EmailServiceAttributes, templates.EmailServiceValidator),
}

type PasswordModel struct {
	Disabled        types.Bool                   `tfsdk:"disabled"`
	Expiration      types.Bool                   `tfsdk:"expiration"`
	ExpirationWeeks types.Int64                  `tfsdk:"expiration_weeks"`
	Lock            types.Bool                   `tfsdk:"lock"`
	LockAttempts    types.Int64                  `tfsdk:"lock_attempts"`
	Lowercase       types.Bool                   `tfsdk:"lowercase"`
	MinLength       types.Int64                  `tfsdk:"min_length"`
	NonAlphanumeric types.Bool                   `tfsdk:"non_alphanumeric"`
	Number          types.Bool                   `tfsdk:"number"`
	Reuse           types.Bool                   `tfsdk:"reuse"`
	ReuseAmount     types.Int64                  `tfsdk:"reuse_amount"`
	Uppercase       types.Bool                   `tfsdk:"uppercase"`
	EmailService    *templates.EmailServiceModel `tfsdk:"email_service"`
}

func (m *PasswordModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.Disabled, data, "enabled")
	boolattr.Get(m.Expiration, data, "expiration")
	intattr.Get(m.ExpirationWeeks, data, "expirationWeeks")
	boolattr.Get(m.Lock, data, "lock")
	intattr.Get(m.LockAttempts, data, "lockAttempts")
	boolattr.Get(m.Lowercase, data, "lowercase")
	intattr.Get(m.MinLength, data, "minLength")
	boolattr.Get(m.NonAlphanumeric, data, "nonAlphanumeric")
	boolattr.Get(m.Number, data, "number")
	boolattr.Get(m.Reuse, data, "reuse")
	intattr.Get(m.ReuseAmount, data, "reuseAmount")
	boolattr.Get(m.Uppercase, data, "uppercase")
	if v := m.EmailService; v != nil {
		maps.Copy(data, v.Values(h))
	}
	return data
}

func (m *PasswordModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.Disabled, data, "enabled")
	boolattr.Set(&m.Expiration, data, "expiration")
	intattr.Set(&m.ExpirationWeeks, data, "expirationWeeks")
	boolattr.Set(&m.Lock, data, "lock")
	intattr.Set(&m.LockAttempts, data, "lockAttempts")
	boolattr.Set(&m.Lowercase, data, "lowercase")
	intattr.Set(&m.MinLength, data, "minLength")
	boolattr.Set(&m.NonAlphanumeric, data, "nonAlphanumeric")
	boolattr.Set(&m.Number, data, "number")
	boolattr.Set(&m.Reuse, data, "reuse")
	intattr.Set(&m.ReuseAmount, data, "reuseAmount")
	boolattr.Set(&m.Uppercase, data, "uppercase")
	emailService := utils.ZVL(m.EmailService)
	emailService.SetValues(h, data)
	if emailService.Connector.ValueString() != "" {
		m.EmailService = emailService
	}
}

func (m *PasswordModel) UseReferences(h *helpers.Handler) {
	if m.EmailService != nil {
		m.EmailService.UseReferences(h)
	}
}

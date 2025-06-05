package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/project/templates"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	"email_service":    objattr.Optional[templates.EmailServiceModel](templates.EmailServiceAttributes, templates.EmailServiceValidator),
}

type PasswordModel struct {
	Disabled        boolattr.Type                             `tfsdk:"disabled"`
	Expiration      boolattr.Type                             `tfsdk:"expiration"`
	ExpirationWeeks intattr.Type                              `tfsdk:"expiration_weeks"`
	Lock            boolattr.Type                             `tfsdk:"lock"`
	LockAttempts    intattr.Type                              `tfsdk:"lock_attempts"`
	Lowercase       boolattr.Type                             `tfsdk:"lowercase"`
	MinLength       intattr.Type                              `tfsdk:"min_length"`
	NonAlphanumeric boolattr.Type                             `tfsdk:"non_alphanumeric"`
	Number          boolattr.Type                             `tfsdk:"number"`
	Reuse           boolattr.Type                             `tfsdk:"reuse"`
	ReuseAmount     intattr.Type                              `tfsdk:"reuse_amount"`
	Uppercase       boolattr.Type                             `tfsdk:"uppercase"`
	EmailService    objattr.Type[templates.EmailServiceModel] `tfsdk:"email_service"`
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
	objattr.Get(m.EmailService, data, helpers.RootKey, h)
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
	objattr.Set(&m.EmailService, data, helpers.RootKey, h)
}

func (m *PasswordModel) UpdateReferences(h *helpers.Handler) {
	objattr.UpdateReferences(&m.EmailService, h)
}

package settings

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var TestUsersAttributes = map[string]schema.Attribute{
	"loginid_regular_expression":             stringattr.Default(""),
	"static_otp_enabled":                     boolattr.Default(false),
	"static_otp_code":                        stringattr.Default(""),
	"static_otp_verifier_regular_expression": stringattr.Default(""),
}

type TestUsersModel struct {
	LoginIDRegularExpression           types.String `tfsdk:"loginid_regular_expression"`
	StaticOTPEnabled                   types.Bool   `tfsdk:"static_otp_enabled"`
	StaticOTPCode                      types.String `tfsdk:"static_otp_code"`
	StaticOTPVerifierRegularExpression types.String `tfsdk:"verifier_regular_expression"`
}

func (m *TestUsersModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	stringattr.Get(m.LoginIDRegularExpression, data, "testUserRegex")
	boolattr.Get(m.StaticOTPEnabled, data, "testUserAllowFixedAuth")
	stringattr.Get(m.StaticOTPCode, data, "testUserFixedAuthToken")
	stringattr.Get(m.StaticOTPVerifierRegularExpression, data, "testUserFixedAuthVerifierRegex")
	return data
}

func (m *TestUsersModel) SetValues(h *helpers.Handler, data map[string]any) {
	stringattr.Set(&m.LoginIDRegularExpression, data, "testUserRegex")
	boolattr.Set(&m.StaticOTPEnabled, data, "testUserAllowFixedAuth")
	stringattr.Set(&m.StaticOTPCode, data, "testUserFixedAuthToken")
	stringattr.Set(&m.StaticOTPVerifierRegularExpression, data, "testUserFixedAuthVerifierRegex")
}

package settings

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var InviteSettingsAttributes = map[string]schema.Attribute{
	"require_invitation":  boolattr.Default(false),
	"invite_url":          stringattr.Default(""),
	"add_magiclink_token": boolattr.Default(false),
	"send_email":          boolattr.Default(true),
	"send_text":           boolattr.Default(false),
	"foo": schema.SingleNestedAttribute{
		Optional:   true,
		Computed:   true,
		CustomType: NewObjectTypeOf[FooType](context.Background()),
		Attributes: map[string]schema.Attribute{},
	},
}

type InviteSettingsModel struct {
	RequireInvitation types.Bool             `tfsdk:"require_invitation"`
	InviteURL         types.String           `tfsdk:"invite_url"`
	AddMagicLinkToken types.Bool             `tfsdk:"add_magiclink_token"`
	SendEmail         types.Bool             `tfsdk:"send_email"`
	SendText          types.Bool             `tfsdk:"send_text"`
	Foo               ObjectValueOf[FooType] `tfsdk:"foo"`
}

func (m *InviteSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.RequireInvitation, data, "projectSelfProvisioning")
	stringattr.Get(m.InviteURL, data, "inviteUrl")
	boolattr.Get(m.AddMagicLinkToken, data, "inviteMagicLink")
	boolattr.Get(m.SendEmail, data, "inviteSendEmail")
	boolattr.Get(m.SendText, data, "inviteSendSms")
	v, _ := m.Foo.ToPtr(h.Ctx)
	fmt.Println(v)
	return data
}

func (m *InviteSettingsModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.SetNot(&m.RequireInvitation, data, "projectSelfProvisioning")
	stringattr.Set(&m.InviteURL, data, "inviteUrl")
	boolattr.Set(&m.AddMagicLinkToken, data, "inviteMagicLink")
	boolattr.Set(&m.SendEmail, data, "inviteSendEmail")
	boolattr.Set(&m.SendText, data, "inviteSendSms")
}

// Type

type FooType struct {
}

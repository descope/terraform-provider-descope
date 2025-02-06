package settings

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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
		CustomType: FooType{},
		Attributes: map[string]schema.Attribute{},
	},
}

type InviteSettingsModel struct {
	RequireInvitation types.Bool   `tfsdk:"require_invitation"`
	InviteURL         types.String `tfsdk:"invite_url"`
	AddMagicLinkToken types.Bool   `tfsdk:"add_magiclink_token"`
	SendEmail         types.Bool   `tfsdk:"send_email"`
	SendText          types.Bool   `tfsdk:"send_text"`
	Foo               FooValue     `tfsdk:"foo"`
}

func (m *InviteSettingsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.GetNot(m.RequireInvitation, data, "projectSelfProvisioning")
	stringattr.Get(m.InviteURL, data, "inviteUrl")
	boolattr.Get(m.AddMagicLinkToken, data, "inviteMagicLink")
	boolattr.Get(m.SendEmail, data, "inviteSendEmail")
	boolattr.Get(m.SendText, data, "inviteSendSms")
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

var _ basetypes.ObjectTypable = FooType{}

type FooType struct {
	basetypes.ObjectType
}

func (f FooType) Equal(t attr.Type) bool {
	other, ok := t.(FooType)
	if !ok {
		return false
	}
	return f.ObjectType.Equal(other.ObjectType)
}

func (f FooType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := f.ObjectType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	objectValue, ok := attrValue.(basetypes.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	objectValuable, diags := f.ValueFromObject(ctx, objectValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting FooValue to ObjectValuable: %v", diags)
	}

	return objectValuable, nil
}

func (f FooType) ValueFromObject(_ context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	value := FooValue{
		ObjectValue: in,
	}
	return value, nil
}

func (f FooType) ValueType(_ context.Context) attr.Value {
	return FooValue{}
}

func (f FooType) String() string {
	return "FooType"
}

// Value

var _ basetypes.ObjectValuable = FooValue{}

type FooValue struct {
	basetypes.ObjectValue
}

func (c FooValue) Equal(o attr.Value) bool {
	other, ok := o.(FooValue)

	if !ok {
		return false
	}

	return c.ObjectValue.Equal(other.ObjectValue)
}

func (c FooValue) Type(ctx context.Context) attr.Type {
	return FooType{
		basetypes.ObjectType{
			AttrTypes: c.AttributeTypes(ctx),
		},
	}
}

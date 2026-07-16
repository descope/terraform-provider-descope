package outboundapp

import (
	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var URLParamAttributes = map[string]schema.Attribute{
	"key":   describedString(stringattr.Required(), "The request parameter name."),
	"value": describedString(stringattr.Required(), "The request parameter value."),
}

var OutboundAppAttributes = map[string]schema.Attribute{
	"id":          describedString(stringattr.Identifier(), "The stable ID of the outbound app."),
	"project_id":  describedString(stringattr.Required(stringplanmodifier.RequiresReplace()), "The ID of the Descope project this outbound app belongs to."),
	"name":        describedString(stringattr.Required(), "A name for the outbound app."),
	"description": describedString(stringattr.Default(""), "A description for the outbound app."),
	"client_id":   describedString(stringattr.Optional(), "The OAuth client ID. The backend may populate this value when it is omitted."),
	"client_secret": schema.StringAttribute{
		Optional:            true,
		Sensitive:           true,
		WriteOnly:           true,
		Validators:          []validator.String{stringattr.NonEmptyValidator},
		MarkdownDescription: "The OAuth client secret. This write-only value is never returned by Descope or stored in Terraform state. Omitting it during an update preserves the existing secret; imported resources leave it unmanaged until configured.",
	},
	"logo":                     describedString(stringattr.Default(""), "A URL or data URI for the outbound app logo."),
	"discovery_url":            describedString(stringattr.Default(""), "The OpenID Connect discovery endpoint."),
	"authorization_url":        describedString(stringattr.Default(""), "The OAuth authorization endpoint."),
	"authorization_url_params": describedList(listattr.Default[URLParamModel](URLParamAttributes), "Additional parameters sent to the authorization endpoint."),
	"token_url":                describedString(stringattr.Default(""), "The OAuth token endpoint."),
	"token_url_params":         describedList(listattr.Default[URLParamModel](URLParamAttributes), "Additional parameters sent to the token endpoint."),
	"revocation_url":           describedString(stringattr.Default(""), "The OAuth token revocation endpoint."),
	"default_scopes":           describedSet(strsetattr.Default(), "The scopes requested by default."),
	"default_redirect_url":     describedString(stringattr.Default(""), "The default redirect URL used after authorization."),
	"callback_domain":          describedString(stringattr.Default(""), "The callback domain used for the OAuth redirect."),
	"pkce":                     describedBool(boolattr.Default(false), "Whether the outbound app uses PKCE."),
	"access_type":              describedString(stringattr.Default("", stringvalidator.OneOf("", string(descope.AccessTypeOffline), string(descope.AccessTypeOnline))), "The OAuth access type. Supported values are `offline` and `online`."),
	"prompt": describedSet(strsetattr.Default(setvalidator.ValueStringsAre(stringvalidator.OneOf(
		string(descope.PromptTypeNone),
		string(descope.PromptTypeLogin),
		string(descope.PromptTypeConsent),
		string(descope.PromptTypeSelectAccount),
	))), "The OAuth prompt values sent during authorization."),
}

var Schema = schema.Schema{
	MarkdownDescription: "Manages a custom OAuth outbound application in a Descope project.",
	Attributes:          OutboundAppAttributes,
}

func describedString(attribute schema.StringAttribute, description string) schema.StringAttribute {
	attribute.MarkdownDescription = description
	return attribute
}

func describedBool(attribute schema.BoolAttribute, description string) schema.BoolAttribute {
	attribute.MarkdownDescription = description
	return attribute
}

func describedSet(attribute schema.SetAttribute, description string) schema.SetAttribute {
	attribute.MarkdownDescription = description
	return attribute
}

func describedList(attribute schema.ListNestedAttribute, description string) schema.ListNestedAttribute {
	attribute.MarkdownDescription = description
	return attribute
}

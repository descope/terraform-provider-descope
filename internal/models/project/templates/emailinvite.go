package templates

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/objattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var EmailInviteValidator = objattr.NewValidator[EmailInviteModel]("must have unique template names and a valid configuration")

var EmailInviteAttributes = map[string]schema.Attribute{
	"connector": stringattr.Required(),
	"templates": listattr.Default[EmailTemplateModel](EmailTemplateAttributes, EmailTemplateValidator),
}

type EmailInviteModel struct {
	Connector stringattr.Type                   `tfsdk:"connector"`
	Templates listattr.Type[EmailTemplateModel] `tfsdk:"templates"`
}

func (m *EmailInviteModel) Values(h *helpers.Handler) map[string]any {
	return getEmailValues(h, m.Connector, "inviteEmailProviderId", m.Templates, "inviteEmailTemplates")
}

func (m *EmailInviteModel) SetValues(h *helpers.Handler, data map[string]any) {
	setEmailValues(h, data, &m.Connector, "inviteEmailProviderId", &m.Templates, "inviteEmailTemplates")
}

func (m *EmailInviteModel) Validate(h *helpers.Handler) {
	validateEmailValues(h, m.Connector, m.Templates)
}

func (m *EmailInviteModel) UpdateReferences(h *helpers.Handler) {
	replaceConnectorIDWithReference(&m.Connector, h)
}

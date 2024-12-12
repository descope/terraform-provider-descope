// Code generated by terragen. DO NOT EDIT.

package docs

import (
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models"
	"github.com/descope/terraform-provider-descope/internal/models/applications"
	"github.com/descope/terraform-provider-descope/internal/models/attributes"
	"github.com/descope/terraform-provider-descope/internal/models/authentication"
	"github.com/descope/terraform-provider-descope/internal/models/authorization"
	"github.com/descope/terraform-provider-descope/internal/models/connectors"
	"github.com/descope/terraform-provider-descope/internal/models/flows"
	"github.com/descope/terraform-provider-descope/internal/models/jwttemplates"
	"github.com/descope/terraform-provider-descope/internal/models/settings"
	"github.com/descope/terraform-provider-descope/internal/models/templates"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func InjectModels() {
	inject(models.ProjectAttributes, docsProject)
	inject(applications.ApplicationAttributes, docsApplication)
	inject(applications.OIDCAttributes, docsOIDC)
	inject(applications.SAMLAttributes, docsSAML)
	inject(applications.AttributeMappingAttributes, docsAttributeMapping)
	inject(applications.DynamicConfigurationAttributes, docsDynamicConfiguration)
	inject(applications.ManualConfigurationAttributes, docsManualConfiguration)
	inject(attributes.AttributesAttributes, docsAttributes)
	inject(attributes.TenantAttributeAttributes, docsTenantAttribute)
	inject(attributes.UserAttributeAttributes, docsUserAttribute)
	inject(attributes.UserAttributeWidgetAuthorizationAttributes, docsUserAttributeWidgetAuthorization)
	inject(authentication.AuthenticationAttributes, docsAuthentication)
	inject(authentication.EmbeddedLinkAttributes, docsEmbeddedLink)
	inject(authentication.EnchantedLinkAttributes, docsEnchantedLink)
	inject(authentication.MagicLinkAttributes, docsMagicLink)
	inject(authentication.OAuthAttributes, docsOAuth)
	inject(authentication.OAuthSystemProviderAttributes, docsOAuthSystemProvider)
	inject(authentication.OAuthProviderAttributes, docsOAuthProvider)
	inject(authentication.OTPAttributes, docsOTP)
	inject(authentication.PasskeysAttributes, docsPasskeys)
	inject(authentication.PasswordAttributes, docsPassword)
	inject(authentication.SSOAttributes, docsSSO)
	inject(authentication.TOTPAttributes, docsTOTP)
	inject(authorization.AuthorizationAttributes, docsAuthorization)
	inject(authorization.PermissionAttributes, docsPermission)
	inject(authorization.RoleAttributes, docsRole)
	inject(connectors.AbuseIPDBAttributes, docsAbuseIPDB)
	inject(connectors.AmplitudeAttributes, docsAmplitude)
	inject(connectors.AuditWebhookAttributes, docsAuditWebhook)
	inject(connectors.AWSS3Attributes, docsAWSS3)
	inject(connectors.AWSTranslateAttributes, docsAWSTranslate)
	inject(connectors.ClearAttributes, docsClear)
	inject(connectors.ConnectorsAttributes, docsConnectors)
	inject(connectors.DatadogAttributes, docsDatadog)
	inject(connectors.DevRevGrowAttributes, docsDevRevGrow)
	inject(connectors.DoceboAttributes, docsDocebo)
	inject(connectors.FingerprintAttributes, docsFingerprint)
	inject(connectors.FingerprintDescopeAttributes, docsFingerprintDescope)
	inject(connectors.ForterAttributes, docsForter)
	inject(connectors.GoogleCloudTranslationAttributes, docsGoogleCloudTranslation)
	inject(connectors.HIBPAttributes, docsHIBP)
	inject(connectors.HTTPAttributes, docsHTTP)
	inject(connectors.HttpStaticIPAttributes, docsHttpStaticIP)
	inject(connectors.HubSpotAttributes, docsHubSpot)
	inject(connectors.IntercomAttributes, docsIntercom)
	inject(connectors.LokaliseAttributes, docsLokalise)
	inject(connectors.MParticleAttributes, docsMParticle)
	inject(connectors.NewRelicAttributes, docsNewRelic)
	inject(connectors.RecaptchaAttributes, docsRecaptcha)
	inject(connectors.RecaptchaEnterpriseAttributes, docsRecaptchaEnterprise)
	inject(connectors.RekognitionAttributes, docsRekognition)
	inject(connectors.SalesforceAttributes, docsSalesforce)
	inject(connectors.SegmentAttributes, docsSegment)
	inject(connectors.SendGridAttributes, docsSendGrid)
	inject(connectors.SendGridAuthFieldAttributes, docsSendGridAuthField)
	inject(connectors.SenderFieldAttributes, docsSenderField)
	inject(connectors.ServerFieldAttributes, docsServerField)
	inject(connectors.HTTPAuthFieldAttributes, docsHTTPAuthField)
	inject(connectors.HTTPAuthBasicFieldAttributes, docsHTTPAuthBasicField)
	inject(connectors.HTTPAuthAPIKeyFieldAttributes, docsHTTPAuthAPIKeyField)
	inject(connectors.SmartlingAttributes, docsSmartling)
	inject(connectors.SMTPAttributes, docsSMTP)
	inject(connectors.SMTPAuthFieldAttributes, docsSMTPAuthField)
	inject(connectors.SumoLogicAttributes, docsSumoLogic)
	inject(connectors.TelesignAttributes, docsTelesign)
	inject(connectors.TraceableAttributes, docsTraceable)
	inject(connectors.TwilioCoreAttributes, docsTwilioCore)
	inject(connectors.TwilioCoreSendersFieldAttributes, docsTwilioCoreSendersField)
	inject(connectors.TwilioCoreSendersSMSFieldAttributes, docsTwilioCoreSendersSMSField)
	inject(connectors.TwilioCoreSendersVoiceFieldAttributes, docsTwilioCoreSendersVoiceField)
	inject(connectors.TwilioAuthFieldAttributes, docsTwilioAuthField)
	inject(connectors.TwilioVerifyAttributes, docsTwilioVerify)
	inject(connectors.VeriffAttributes, docsVeriff)
	inject(flows.FlowAttributes, docsFlow)
	inject(flows.StylesAttributes, docsStyles)
	inject(jwttemplates.JWTTemplateAttributes, docsJWTTemplate)
	inject(jwttemplates.JWTTemplatesAttributes, docsJWTTemplates)
	inject(settings.SettingsAttributes, docsSettings)
	inject(templates.EmailServiceAttributes, docsEmailService)
	inject(templates.EmailTemplateAttributes, docsEmailTemplate)
	inject(templates.TextServiceAttributes, docsTextService)
	inject(templates.TextTemplateAttributes, docsTextTemplate)
	inject(templates.VoiceServiceAttributes, docsVoiceService)
	inject(templates.VoiceTemplateAttributes, docsVoiceTemplate)
}

func inject(model map[string]schema.Attribute, docs map[string]string) {
	for field, desc := range docs {
		if _, found := model[field]; !found {
			panic(fmt.Sprintf("missing schema attribute for documentation key %s", field))
		}
		switch attr := model[field].(type) {
		case schema.StringAttribute:
			attr.MarkdownDescription = desc
			model[field] = attr
		case schema.BoolAttribute:
			attr.MarkdownDescription = desc
			model[field] = attr
		case schema.Int64Attribute:
			attr.MarkdownDescription = desc
			model[field] = attr
		case schema.Float64Attribute:
			attr.MarkdownDescription = desc
			model[field] = attr
		case schema.MapAttribute:
			attr.MarkdownDescription = desc
			model[field] = attr
		case schema.ListAttribute:
			attr.MarkdownDescription = desc
			model[field] = attr
		case schema.MapNestedAttribute:
			attr.MarkdownDescription = desc
			model[field] = attr
		case schema.SingleNestedAttribute:
			attr.MarkdownDescription = desc
			model[field] = attr
		case schema.ListNestedAttribute:
			attr.MarkdownDescription = desc
			model[field] = attr
		default:
			panic(fmt.Sprintf("unexpected schema type for documentation key %s: %T", field, attr))
		}
	}
}

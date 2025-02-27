package connectors

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/objectattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var ConnectorsValidator = objectattr.NewValidator[ConnectorsModel]("must have unique connector names")

var ConnectorsModifier = objectattr.NewModifier[ConnectorsModel]("maintains connector identifiers between plan changes")

var ConnectorsAttributes = map[string]schema.Attribute{
	"abuseipdb":                  listattr.Optional(AbuseIPDBAttributes),
	"amplitude":                  listattr.Optional(AmplitudeAttributes),
	"audit_webhook":              listattr.Optional(AuditWebhookAttributes),
	"aws_s3":                     listattr.Optional(AWSS3Attributes, AWSS3Validator),
	"aws_translate":              listattr.Optional(AWSTranslateAttributes),
	"clear":                      listattr.Optional(ClearAttributes),
	"datadog":                    listattr.Optional(DatadogAttributes, DatadogValidator),
	"devrev_grow":                listattr.Optional(DevRevGrowAttributes),
	"docebo":                     listattr.Optional(DoceboAttributes),
	"fingerprint":                listattr.Optional(FingerprintAttributes, FingerprintValidator),
	"fingerprint_descope":        listattr.Optional(FingerprintDescopeAttributes),
	"forter":                     listattr.Optional(ForterAttributes, ForterValidator),
	"generic_sms_gateway":        listattr.Optional(GenericSMSGatewayAttributes),
	"google_cloud_translation":   listattr.Optional(GoogleCloudTranslationAttributes),
	"hibp":                       listattr.Optional(HIBPAttributes),
	"http":                       listattr.Optional(HTTPAttributes),
	"hubspot":                    listattr.Optional(HubSpotAttributes),
	"incode":                     listattr.Optional(IncodeAttributes),
	"intercom":                   listattr.Optional(IntercomAttributes),
	"lokalise":                   listattr.Optional(LokaliseAttributes),
	"mparticle":                  listattr.Optional(MParticleAttributes),
	"newrelic":                   listattr.Optional(NewRelicAttributes, NewRelicValidator),
	"recaptcha":                  listattr.Optional(RecaptchaAttributes),
	"recaptcha_enterprise":       listattr.Optional(RecaptchaEnterpriseAttributes, RecaptchaEnterpriseValidator),
	"rekognition":                listattr.Optional(RekognitionAttributes),
	"salesforce":                 listattr.Optional(SalesforceAttributes),
	"salesforce_marketing_cloud": listattr.Optional(SalesforceMarketingCloudAttributes),
	"segment":                    listattr.Optional(SegmentAttributes),
	"sendgrid":                   listattr.Optional(SendGridAttributes),
	"ses":                        listattr.Optional(SESAttributes),
	"slack":                      listattr.Optional(SlackAttributes),
	"smartling":                  listattr.Optional(SmartlingAttributes),
	"smtp":                       listattr.Optional(SMTPAttributes),
	"sns":                        listattr.Optional(SNSAttributes),
	"sumologic":                  listattr.Optional(SumoLogicAttributes, SumoLogicValidator),
	"telesign":                   listattr.Optional(TelesignAttributes),
	"traceable":                  listattr.Optional(TraceableAttributes),
	"twilio_core":                listattr.Optional(TwilioCoreAttributes),
	"twilio_verify":              listattr.Optional(TwilioVerifyAttributes),
}

type ConnectorsModel struct {
	AbuseIPDB                []*AbuseIPDBModel                `tfsdk:"abuseipdb"`
	Amplitude                []*AmplitudeModel                `tfsdk:"amplitude"`
	AuditWebhook             []*AuditWebhookModel             `tfsdk:"audit_webhook"`
	AWSS3                    []*AWSS3Model                    `tfsdk:"aws_s3"`
	AWSTranslate             []*AWSTranslateModel             `tfsdk:"aws_translate"`
	Clear                    []*ClearModel                    `tfsdk:"clear"`
	Datadog                  []*DatadogModel                  `tfsdk:"datadog"`
	DevRevGrow               []*DevRevGrowModel               `tfsdk:"devrev_grow"`
	Docebo                   []*DoceboModel                   `tfsdk:"docebo"`
	Fingerprint              []*FingerprintModel              `tfsdk:"fingerprint"`
	FingerprintDescope       []*FingerprintDescopeModel       `tfsdk:"fingerprint_descope"`
	Forter                   []*ForterModel                   `tfsdk:"forter"`
	GenericSMSGateway        []*GenericSMSGatewayModel        `tfsdk:"generic_sms_gateway"`
	GoogleCloudTranslation   []*GoogleCloudTranslationModel   `tfsdk:"google_cloud_translation"`
	HIBP                     []*HIBPModel                     `tfsdk:"hibp"`
	HTTP                     []*HTTPModel                     `tfsdk:"http"`
	HubSpot                  []*HubSpotModel                  `tfsdk:"hubspot"`
	Incode                   []*IncodeModel                   `tfsdk:"incode"`
	Intercom                 []*IntercomModel                 `tfsdk:"intercom"`
	Lokalise                 []*LokaliseModel                 `tfsdk:"lokalise"`
	MParticle                []*MParticleModel                `tfsdk:"mparticle"`
	NewRelic                 []*NewRelicModel                 `tfsdk:"newrelic"`
	Recaptcha                []*RecaptchaModel                `tfsdk:"recaptcha"`
	RecaptchaEnterprise      []*RecaptchaEnterpriseModel      `tfsdk:"recaptcha_enterprise"`
	Rekognition              []*RekognitionModel              `tfsdk:"rekognition"`
	Salesforce               []*SalesforceModel               `tfsdk:"salesforce"`
	SalesforceMarketingCloud []*SalesforceMarketingCloudModel `tfsdk:"salesforce_marketing_cloud"`
	Segment                  []*SegmentModel                  `tfsdk:"segment"`
	SendGrid                 []*SendGridModel                 `tfsdk:"sendgrid"`
	SES                      []*SESModel                      `tfsdk:"ses"`
	Slack                    []*SlackModel                    `tfsdk:"slack"`
	Smartling                []*SmartlingModel                `tfsdk:"smartling"`
	SMTP                     []*SMTPModel                     `tfsdk:"smtp"`
	SNS                      []*SNSModel                      `tfsdk:"sns"`
	SumoLogic                []*SumoLogicModel                `tfsdk:"sumologic"`
	Telesign                 []*TelesignModel                 `tfsdk:"telesign"`
	Traceable                []*TraceableModel                `tfsdk:"traceable"`
	TwilioCore               []*TwilioCoreModel               `tfsdk:"twilio_core"`
	TwilioVerify             []*TwilioVerifyModel             `tfsdk:"twilio_verify"`
}

func (m *ConnectorsModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	listattr.Get(m.AbuseIPDB, data, "abuseipdb", h)
	listattr.Get(m.Amplitude, data, "amplitude", h)
	listattr.Get(m.AuditWebhook, data, "audit-webhook", h)
	listattr.Get(m.AWSS3, data, "aws-s3", h)
	listattr.Get(m.AWSTranslate, data, "aws-translate", h)
	listattr.Get(m.Clear, data, "clear", h)
	listattr.Get(m.Datadog, data, "datadog", h)
	listattr.Get(m.DevRevGrow, data, "devrev-grow", h)
	listattr.Get(m.Docebo, data, "docebo", h)
	listattr.Get(m.Fingerprint, data, "fingerprint", h)
	listattr.Get(m.FingerprintDescope, data, "fingerprint-descope", h)
	listattr.Get(m.Forter, data, "forter", h)
	listattr.Get(m.GenericSMSGateway, data, "generic-sms-gateway", h)
	listattr.Get(m.GoogleCloudTranslation, data, "google-cloud-translation", h)
	listattr.Get(m.HIBP, data, "hibp", h)
	listattr.Get(m.HTTP, data, "http", h)
	listattr.Get(m.HubSpot, data, "hubspot", h)
	listattr.Get(m.Incode, data, "incode", h)
	listattr.Get(m.Intercom, data, "intercom", h)
	listattr.Get(m.Lokalise, data, "lokalise", h)
	listattr.Get(m.MParticle, data, "mparticle", h)
	listattr.Get(m.NewRelic, data, "newrelic", h)
	listattr.Get(m.Recaptcha, data, "recaptcha", h)
	listattr.Get(m.RecaptchaEnterprise, data, "recaptcha-enterprise", h)
	listattr.Get(m.Rekognition, data, "rekognition", h)
	listattr.Get(m.Salesforce, data, "salesforce", h)
	listattr.Get(m.SalesforceMarketingCloud, data, "salesforce-marketing-cloud", h)
	listattr.Get(m.Segment, data, "segment", h)
	listattr.Get(m.SendGrid, data, "sendgrid", h)
	listattr.Get(m.SES, data, "ses", h)
	listattr.Get(m.Slack, data, "slack", h)
	listattr.Get(m.Smartling, data, "smartling", h)
	listattr.Get(m.SMTP, data, "smtp", h)
	listattr.Get(m.SNS, data, "sns", h)
	listattr.Get(m.SumoLogic, data, "sumologic", h)
	listattr.Get(m.Telesign, data, "telesign", h)
	listattr.Get(m.Traceable, data, "traceable", h)
	listattr.Get(m.TwilioCore, data, "twilio-core", h)
	listattr.Get(m.TwilioVerify, data, "twilio-verify", h)
	return data
}

func (m *ConnectorsModel) SetValues(h *helpers.Handler, data map[string]any) {
	SetConnectorIDs(h, data, "abuseipdb", &m.AbuseIPDB)
	SetConnectorIDs(h, data, "amplitude", &m.Amplitude)
	SetConnectorIDs(h, data, "audit-webhook", &m.AuditWebhook)
	SetConnectorIDs(h, data, "aws-s3", &m.AWSS3)
	SetConnectorIDs(h, data, "aws-translate", &m.AWSTranslate)
	SetConnectorIDs(h, data, "clear", &m.Clear)
	SetConnectorIDs(h, data, "datadog", &m.Datadog)
	SetConnectorIDs(h, data, "devrev-grow", &m.DevRevGrow)
	SetConnectorIDs(h, data, "docebo", &m.Docebo)
	SetConnectorIDs(h, data, "fingerprint", &m.Fingerprint)
	SetConnectorIDs(h, data, "fingerprint-descope", &m.FingerprintDescope)
	SetConnectorIDs(h, data, "forter", &m.Forter)
	SetConnectorIDs(h, data, "generic-sms-gateway", &m.GenericSMSGateway)
	SetConnectorIDs(h, data, "google-cloud-translation", &m.GoogleCloudTranslation)
	SetConnectorIDs(h, data, "hibp", &m.HIBP)
	SetConnectorIDs(h, data, "http", &m.HTTP)
	SetConnectorIDs(h, data, "hubspot", &m.HubSpot)
	SetConnectorIDs(h, data, "incode", &m.Incode)
	SetConnectorIDs(h, data, "intercom", &m.Intercom)
	SetConnectorIDs(h, data, "lokalise", &m.Lokalise)
	SetConnectorIDs(h, data, "mparticle", &m.MParticle)
	SetConnectorIDs(h, data, "newrelic", &m.NewRelic)
	SetConnectorIDs(h, data, "recaptcha", &m.Recaptcha)
	SetConnectorIDs(h, data, "recaptcha-enterprise", &m.RecaptchaEnterprise)
	SetConnectorIDs(h, data, "rekognition", &m.Rekognition)
	SetConnectorIDs(h, data, "salesforce", &m.Salesforce)
	SetConnectorIDs(h, data, "salesforce-marketing-cloud", &m.SalesforceMarketingCloud)
	SetConnectorIDs(h, data, "segment", &m.Segment)
	SetConnectorIDs(h, data, "sendgrid", &m.SendGrid)
	SetConnectorIDs(h, data, "ses", &m.SES)
	SetConnectorIDs(h, data, "slack", &m.Slack)
	SetConnectorIDs(h, data, "smartling", &m.Smartling)
	SetConnectorIDs(h, data, "smtp", &m.SMTP)
	SetConnectorIDs(h, data, "sns", &m.SNS)
	SetConnectorIDs(h, data, "sumologic", &m.SumoLogic)
	SetConnectorIDs(h, data, "telesign", &m.Telesign)
	SetConnectorIDs(h, data, "traceable", &m.Traceable)
	SetConnectorIDs(h, data, "twilio-core", &m.TwilioCore)
	SetConnectorIDs(h, data, "twilio-verify", &m.TwilioVerify)
}

func (m *ConnectorsModel) References(ctx context.Context) helpers.ReferencesMap {
	refs := helpers.ReferencesMap{}
	addConnectorReferences(refs, "abuseipdb", m.AbuseIPDB)
	addConnectorReferences(refs, "amplitude", m.Amplitude)
	addConnectorReferences(refs, "audit-webhook", m.AuditWebhook)
	addConnectorReferences(refs, "aws-s3", m.AWSS3)
	addConnectorReferences(refs, "aws-translate", m.AWSTranslate)
	addConnectorReferences(refs, "clear", m.Clear)
	addConnectorReferences(refs, "datadog", m.Datadog)
	addConnectorReferences(refs, "devrev-grow", m.DevRevGrow)
	addConnectorReferences(refs, "docebo", m.Docebo)
	addConnectorReferences(refs, "fingerprint", m.Fingerprint)
	addConnectorReferences(refs, "fingerprint-descope", m.FingerprintDescope)
	addConnectorReferences(refs, "forter", m.Forter)
	addConnectorReferences(refs, "generic-sms-gateway", m.GenericSMSGateway)
	addConnectorReferences(refs, "google-cloud-translation", m.GoogleCloudTranslation)
	addConnectorReferences(refs, "hibp", m.HIBP)
	addConnectorReferences(refs, "http", m.HTTP)
	addConnectorReferences(refs, "hubspot", m.HubSpot)
	addConnectorReferences(refs, "incode", m.Incode)
	addConnectorReferences(refs, "intercom", m.Intercom)
	addConnectorReferences(refs, "lokalise", m.Lokalise)
	addConnectorReferences(refs, "mparticle", m.MParticle)
	addConnectorReferences(refs, "newrelic", m.NewRelic)
	addConnectorReferences(refs, "recaptcha", m.Recaptcha)
	addConnectorReferences(refs, "recaptcha-enterprise", m.RecaptchaEnterprise)
	addConnectorReferences(refs, "rekognition", m.Rekognition)
	addConnectorReferences(refs, "salesforce", m.Salesforce)
	addConnectorReferences(refs, "salesforce-marketing-cloud", m.SalesforceMarketingCloud)
	addConnectorReferences(refs, "segment", m.Segment)
	addConnectorReferences(refs, "sendgrid", m.SendGrid)
	addConnectorReferences(refs, "ses", m.SES)
	addConnectorReferences(refs, "slack", m.Slack)
	addConnectorReferences(refs, "smartling", m.Smartling)
	addConnectorReferences(refs, "smtp", m.SMTP)
	addConnectorReferences(refs, "sns", m.SNS)
	addConnectorReferences(refs, "sumologic", m.SumoLogic)
	addConnectorReferences(refs, "telesign", m.Telesign)
	addConnectorReferences(refs, "traceable", m.Traceable)
	addConnectorReferences(refs, "twilio-core", m.TwilioCore)
	addConnectorReferences(refs, "twilio-verify", m.TwilioVerify)
	return refs
}

func (m *ConnectorsModel) Validate(h *helpers.Handler) {
	names := map[string]int{}
	addConnectorNames(names, m.AbuseIPDB)
	addConnectorNames(names, m.Amplitude)
	addConnectorNames(names, m.AuditWebhook)
	addConnectorNames(names, m.AWSS3)
	addConnectorNames(names, m.AWSTranslate)
	addConnectorNames(names, m.Clear)
	addConnectorNames(names, m.Datadog)
	addConnectorNames(names, m.DevRevGrow)
	addConnectorNames(names, m.Docebo)
	addConnectorNames(names, m.Fingerprint)
	addConnectorNames(names, m.FingerprintDescope)
	addConnectorNames(names, m.Forter)
	addConnectorNames(names, m.GenericSMSGateway)
	addConnectorNames(names, m.GoogleCloudTranslation)
	addConnectorNames(names, m.HIBP)
	addConnectorNames(names, m.HTTP)
	addConnectorNames(names, m.HubSpot)
	addConnectorNames(names, m.Incode)
	addConnectorNames(names, m.Intercom)
	addConnectorNames(names, m.Lokalise)
	addConnectorNames(names, m.MParticle)
	addConnectorNames(names, m.NewRelic)
	addConnectorNames(names, m.Recaptcha)
	addConnectorNames(names, m.RecaptchaEnterprise)
	addConnectorNames(names, m.Rekognition)
	addConnectorNames(names, m.Salesforce)
	addConnectorNames(names, m.SalesforceMarketingCloud)
	addConnectorNames(names, m.Segment)
	addConnectorNames(names, m.SendGrid)
	addConnectorNames(names, m.SES)
	addConnectorNames(names, m.Slack)
	addConnectorNames(names, m.Smartling)
	addConnectorNames(names, m.SMTP)
	addConnectorNames(names, m.SNS)
	addConnectorNames(names, m.SumoLogic)
	addConnectorNames(names, m.Telesign)
	addConnectorNames(names, m.Traceable)
	addConnectorNames(names, m.TwilioCore)
	addConnectorNames(names, m.TwilioVerify)
	for k, v := range names {
		if v > 1 {
			h.Error("Connector names must be unique", "The connector name '%s' is used %d times", k, v)
		}
	}
}

func (m *ConnectorsModel) Modify(ctx context.Context, state *ConnectorsModel, diags *diag.Diagnostics) {
	helpers.MatchModels(ctx, m.AbuseIPDB, state.AbuseIPDB)
	helpers.MatchModels(ctx, m.Amplitude, state.Amplitude)
	helpers.MatchModels(ctx, m.AuditWebhook, state.AuditWebhook)
	helpers.MatchModels(ctx, m.AWSS3, state.AWSS3)
	helpers.MatchModels(ctx, m.AWSTranslate, state.AWSTranslate)
	helpers.MatchModels(ctx, m.Clear, state.Clear)
	helpers.MatchModels(ctx, m.Datadog, state.Datadog)
	helpers.MatchModels(ctx, m.DevRevGrow, state.DevRevGrow)
	helpers.MatchModels(ctx, m.Docebo, state.Docebo)
	helpers.MatchModels(ctx, m.Fingerprint, state.Fingerprint)
	helpers.MatchModels(ctx, m.FingerprintDescope, state.FingerprintDescope)
	helpers.MatchModels(ctx, m.Forter, state.Forter)
	helpers.MatchModels(ctx, m.GenericSMSGateway, state.GenericSMSGateway)
	helpers.MatchModels(ctx, m.GoogleCloudTranslation, state.GoogleCloudTranslation)
	helpers.MatchModels(ctx, m.HIBP, state.HIBP)
	helpers.MatchModels(ctx, m.HTTP, state.HTTP)
	helpers.MatchModels(ctx, m.HubSpot, state.HubSpot)
	helpers.MatchModels(ctx, m.Incode, state.Incode)
	helpers.MatchModels(ctx, m.Intercom, state.Intercom)
	helpers.MatchModels(ctx, m.Lokalise, state.Lokalise)
	helpers.MatchModels(ctx, m.MParticle, state.MParticle)
	helpers.MatchModels(ctx, m.NewRelic, state.NewRelic)
	helpers.MatchModels(ctx, m.Recaptcha, state.Recaptcha)
	helpers.MatchModels(ctx, m.RecaptchaEnterprise, state.RecaptchaEnterprise)
	helpers.MatchModels(ctx, m.Rekognition, state.Rekognition)
	helpers.MatchModels(ctx, m.Salesforce, state.Salesforce)
	helpers.MatchModels(ctx, m.SalesforceMarketingCloud, state.SalesforceMarketingCloud)
	helpers.MatchModels(ctx, m.Segment, state.Segment)
	helpers.MatchModels(ctx, m.SendGrid, state.SendGrid)
	helpers.MatchModels(ctx, m.SES, state.SES)
	helpers.MatchModels(ctx, m.Slack, state.Slack)
	helpers.MatchModels(ctx, m.Smartling, state.Smartling)
	helpers.MatchModels(ctx, m.SMTP, state.SMTP)
	helpers.MatchModels(ctx, m.SNS, state.SNS)
	helpers.MatchModels(ctx, m.SumoLogic, state.SumoLogic)
	helpers.MatchModels(ctx, m.Telesign, state.Telesign)
	helpers.MatchModels(ctx, m.Traceable, state.Traceable)
	helpers.MatchModels(ctx, m.TwilioCore, state.TwilioCore)
	helpers.MatchModels(ctx, m.TwilioVerify, state.TwilioVerify)
}

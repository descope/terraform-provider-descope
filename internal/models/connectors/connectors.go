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
	"abuseipdb":                listattr.Optional(AbuseIPDBAttributes),
	"amplitude":                listattr.Optional(AmplitudeAttributes),
	"audit_webhook":            listattr.Optional(AuditWebhookAttributes),
	"aws_s3":                   listattr.Optional(AWSS3Attributes, AWSS3Validator),
	"aws_translate":            listattr.Optional(AWSTranslateAttributes),
	"clear":                    listattr.Optional(ClearAttributes),
	"datadog":                  listattr.Optional(DatadogAttributes, DatadogValidator),
	"devrev_grow":              listattr.Optional(DevRevGrowAttributes),
	"docebo":                   listattr.Optional(DoceboAttributes),
	"fingerprint":              listattr.Optional(FingerprintAttributes, FingerprintValidator),
	"fingerprint_descope":      listattr.Optional(FingerprintDescopeAttributes),
	"forter":                   listattr.Optional(ForterAttributes, ForterValidator),
	"google_cloud_translation": listattr.Optional(GoogleCloudTranslationAttributes),
	"hibp":                     listattr.Optional(HIBPAttributes),
	"http":                     listattr.Optional(HTTPAttributes),
	"http_static_ip":           listattr.Optional(HttpStaticIPAttributes),
	"hubspot":                  listattr.Optional(HubSpotAttributes),
	"intercom":                 listattr.Optional(IntercomAttributes),
	"newrelic":                 listattr.Optional(NewRelicAttributes, NewRelicValidator),
	"recaptcha":                listattr.Optional(RecaptchaAttributes),
	"recaptcha_enterprise":     listattr.Optional(RecaptchaEnterpriseAttributes, RecaptchaEnterpriseValidator),
	"rekognition":              listattr.Optional(RekognitionAttributes),
	"salesforce":               listattr.Optional(SalesforceAttributes),
	"segment":                  listattr.Optional(SegmentAttributes),
	"sendgrid":                 listattr.Optional(SendGridAttributes),
	"smtp":                     listattr.Optional(SMTPAttributes),
	"sumologic":                listattr.Optional(SumoLogicAttributes, SumoLogicValidator),
	"telesign":                 listattr.Optional(TelesignAttributes),
	"traceable":                listattr.Optional(TraceableAttributes),
	"twilio_core":              listattr.Optional(TwilioCoreAttributes),
	"twilio_verify":            listattr.Optional(TwilioVerifyAttributes),
	"veriff":                   listattr.Optional(VeriffAttributes),
}

type ConnectorsModel struct {
	AbuseIPDB              []*AbuseIPDBModel              `tfsdk:"abuseipdb"`
	Amplitude              []*AmplitudeModel              `tfsdk:"amplitude"`
	AuditWebhook           []*AuditWebhookModel           `tfsdk:"audit_webhook"`
	AWSS3                  []*AWSS3Model                  `tfsdk:"aws_s3"`
	AWSTranslate           []*AWSTranslateModel           `tfsdk:"aws_translate"`
	Clear                  []*ClearModel                  `tfsdk:"clear"`
	Datadog                []*DatadogModel                `tfsdk:"datadog"`
	DevRevGrow             []*DevRevGrowModel             `tfsdk:"devrev_grow"`
	Docebo                 []*DoceboModel                 `tfsdk:"docebo"`
	Fingerprint            []*FingerprintModel            `tfsdk:"fingerprint"`
	FingerprintDescope     []*FingerprintDescopeModel     `tfsdk:"fingerprint_descope"`
	Forter                 []*ForterModel                 `tfsdk:"forter"`
	GoogleCloudTranslation []*GoogleCloudTranslationModel `tfsdk:"google_cloud_translation"`
	HIBP                   []*HIBPModel                   `tfsdk:"hibp"`
	HTTP                   []*HTTPModel                   `tfsdk:"http"`
	HttpStaticIP           []*HttpStaticIPModel           `tfsdk:"http_static_ip"`
	HubSpot                []*HubSpotModel                `tfsdk:"hubspot"`
	Intercom               []*IntercomModel               `tfsdk:"intercom"`
	NewRelic               []*NewRelicModel               `tfsdk:"newrelic"`
	Recaptcha              []*RecaptchaModel              `tfsdk:"recaptcha"`
	RecaptchaEnterprise    []*RecaptchaEnterpriseModel    `tfsdk:"recaptcha_enterprise"`
	Rekognition            []*RekognitionModel            `tfsdk:"rekognition"`
	Salesforce             []*SalesforceModel             `tfsdk:"salesforce"`
	Segment                []*SegmentModel                `tfsdk:"segment"`
	SendGrid               []*SendGridModel               `tfsdk:"sendgrid"`
	SMTP                   []*SMTPModel                   `tfsdk:"smtp"`
	SumoLogic              []*SumoLogicModel              `tfsdk:"sumologic"`
	Telesign               []*TelesignModel               `tfsdk:"telesign"`
	Traceable              []*TraceableModel              `tfsdk:"traceable"`
	TwilioCore             []*TwilioCoreModel             `tfsdk:"twilio_core"`
	TwilioVerify           []*TwilioVerifyModel           `tfsdk:"twilio_verify"`
	Veriff                 []*VeriffModel                 `tfsdk:"veriff"`
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
	listattr.Get(m.GoogleCloudTranslation, data, "google-cloud-translation", h)
	listattr.Get(m.HIBP, data, "hibp", h)
	listattr.Get(m.HTTP, data, "http", h)
	listattr.Get(m.HttpStaticIP, data, "http-static-ip", h)
	listattr.Get(m.HubSpot, data, "hubspot", h)
	listattr.Get(m.Intercom, data, "intercom", h)
	listattr.Get(m.NewRelic, data, "newrelic", h)
	listattr.Get(m.Recaptcha, data, "recaptcha", h)
	listattr.Get(m.RecaptchaEnterprise, data, "recaptcha-enterprise", h)
	listattr.Get(m.Rekognition, data, "rekognition", h)
	listattr.Get(m.Salesforce, data, "salesforce", h)
	listattr.Get(m.Segment, data, "segment", h)
	listattr.Get(m.SendGrid, data, "sendgrid", h)
	listattr.Get(m.SMTP, data, "smtp", h)
	listattr.Get(m.SumoLogic, data, "sumologic", h)
	listattr.Get(m.Telesign, data, "telesign", h)
	listattr.Get(m.Traceable, data, "traceable", h)
	listattr.Get(m.TwilioCore, data, "twilio_core", h)
	listattr.Get(m.TwilioVerify, data, "twilio_verify", h)
	listattr.Get(m.Veriff, data, "veriff", h)
	return data
}

func (m *ConnectorsModel) SetValues(h *helpers.Handler, data map[string]any) {
	SetConnectorIDs(h, data, "abuseipdb", m.AbuseIPDB)
	SetConnectorIDs(h, data, "amplitude", m.Amplitude)
	SetConnectorIDs(h, data, "audit-webhook", m.AuditWebhook)
	SetConnectorIDs(h, data, "aws-s3", m.AWSS3)
	SetConnectorIDs(h, data, "aws-translate", m.AWSTranslate)
	SetConnectorIDs(h, data, "clear", m.Clear)
	SetConnectorIDs(h, data, "datadog", m.Datadog)
	SetConnectorIDs(h, data, "devrev-grow", m.DevRevGrow)
	SetConnectorIDs(h, data, "docebo", m.Docebo)
	SetConnectorIDs(h, data, "fingerprint", m.Fingerprint)
	SetConnectorIDs(h, data, "fingerprint-descope", m.FingerprintDescope)
	SetConnectorIDs(h, data, "forter", m.Forter)
	SetConnectorIDs(h, data, "google-cloud-translation", m.GoogleCloudTranslation)
	SetConnectorIDs(h, data, "hibp", m.HIBP)
	SetConnectorIDs(h, data, "http", m.HTTP)
	SetConnectorIDs(h, data, "http-static-ip", m.HttpStaticIP)
	SetConnectorIDs(h, data, "hubspot", m.HubSpot)
	SetConnectorIDs(h, data, "intercom", m.Intercom)
	SetConnectorIDs(h, data, "newrelic", m.NewRelic)
	SetConnectorIDs(h, data, "recaptcha", m.Recaptcha)
	SetConnectorIDs(h, data, "recaptcha-enterprise", m.RecaptchaEnterprise)
	SetConnectorIDs(h, data, "rekognition", m.Rekognition)
	SetConnectorIDs(h, data, "salesforce", m.Salesforce)
	SetConnectorIDs(h, data, "segment", m.Segment)
	SetConnectorIDs(h, data, "sendgrid", m.SendGrid)
	SetConnectorIDs(h, data, "smtp", m.SMTP)
	SetConnectorIDs(h, data, "sumologic", m.SumoLogic)
	SetConnectorIDs(h, data, "telesign", m.Telesign)
	SetConnectorIDs(h, data, "traceable", m.Traceable)
	SetConnectorIDs(h, data, "twilio_core", m.TwilioCore)
	SetConnectorIDs(h, data, "twilio_verify", m.TwilioVerify)
	SetConnectorIDs(h, data, "veriff", m.Veriff)
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
	addConnectorReferences(refs, "google-cloud-translation", m.GoogleCloudTranslation)
	addConnectorReferences(refs, "hibp", m.HIBP)
	addConnectorReferences(refs, "http", m.HTTP)
	addConnectorReferences(refs, "http-static-ip", m.HttpStaticIP)
	addConnectorReferences(refs, "hubspot", m.HubSpot)
	addConnectorReferences(refs, "intercom", m.Intercom)
	addConnectorReferences(refs, "newrelic", m.NewRelic)
	addConnectorReferences(refs, "recaptcha", m.Recaptcha)
	addConnectorReferences(refs, "recaptcha-enterprise", m.RecaptchaEnterprise)
	addConnectorReferences(refs, "rekognition", m.Rekognition)
	addConnectorReferences(refs, "salesforce", m.Salesforce)
	addConnectorReferences(refs, "segment", m.Segment)
	addConnectorReferences(refs, "sendgrid", m.SendGrid)
	addConnectorReferences(refs, "smtp", m.SMTP)
	addConnectorReferences(refs, "sumologic", m.SumoLogic)
	addConnectorReferences(refs, "telesign", m.Telesign)
	addConnectorReferences(refs, "traceable", m.Traceable)
	addConnectorReferences(refs, "twilio_core", m.TwilioCore)
	addConnectorReferences(refs, "twilio_verify", m.TwilioVerify)
	addConnectorReferences(refs, "veriff", m.Veriff)
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
	addConnectorNames(names, m.GoogleCloudTranslation)
	addConnectorNames(names, m.HIBP)
	addConnectorNames(names, m.HTTP)
	addConnectorNames(names, m.HttpStaticIP)
	addConnectorNames(names, m.HubSpot)
	addConnectorNames(names, m.Intercom)
	addConnectorNames(names, m.NewRelic)
	addConnectorNames(names, m.Recaptcha)
	addConnectorNames(names, m.RecaptchaEnterprise)
	addConnectorNames(names, m.Rekognition)
	addConnectorNames(names, m.Salesforce)
	addConnectorNames(names, m.Segment)
	addConnectorNames(names, m.SendGrid)
	addConnectorNames(names, m.SMTP)
	addConnectorNames(names, m.SumoLogic)
	addConnectorNames(names, m.Telesign)
	addConnectorNames(names, m.Traceable)
	addConnectorNames(names, m.TwilioCore)
	addConnectorNames(names, m.TwilioVerify)
	addConnectorNames(names, m.Veriff)
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
	helpers.MatchModels(ctx, m.GoogleCloudTranslation, state.GoogleCloudTranslation)
	helpers.MatchModels(ctx, m.HIBP, state.HIBP)
	helpers.MatchModels(ctx, m.HTTP, state.HTTP)
	helpers.MatchModels(ctx, m.HttpStaticIP, state.HttpStaticIP)
	helpers.MatchModels(ctx, m.HubSpot, state.HubSpot)
	helpers.MatchModels(ctx, m.Intercom, state.Intercom)
	helpers.MatchModels(ctx, m.NewRelic, state.NewRelic)
	helpers.MatchModels(ctx, m.Recaptcha, state.Recaptcha)
	helpers.MatchModels(ctx, m.RecaptchaEnterprise, state.RecaptchaEnterprise)
	helpers.MatchModels(ctx, m.Rekognition, state.Rekognition)
	helpers.MatchModels(ctx, m.Salesforce, state.Salesforce)
	helpers.MatchModels(ctx, m.Segment, state.Segment)
	helpers.MatchModels(ctx, m.SendGrid, state.SendGrid)
	helpers.MatchModels(ctx, m.SMTP, state.SMTP)
	helpers.MatchModels(ctx, m.SumoLogic, state.SumoLogic)
	helpers.MatchModels(ctx, m.Telesign, state.Telesign)
	helpers.MatchModels(ctx, m.Traceable, state.Traceable)
	helpers.MatchModels(ctx, m.TwilioCore, state.TwilioCore)
	helpers.MatchModels(ctx, m.TwilioVerify, state.TwilioVerify)
	helpers.MatchModels(ctx, m.Veriff, state.Veriff)
}

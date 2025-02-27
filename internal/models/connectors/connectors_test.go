package connectors_test

import (
	"regexp"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConnectors(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"connectors.%": 0,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"abuseipdb": [
						{
							name = "Test abuseipdb Connector"
							description = "A description for the abuseipdb connector"
    						api_key = "mhvece"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.abuseipdb.#":             1,
				"connectors.abuseipdb.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.abuseipdb.0.name":        "Test abuseipdb Connector",
				"connectors.abuseipdb.0.description": "A description for the abuseipdb connector",
				"connectors.abuseipdb.0.api_key":     "mhvece",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"amplitude": [
						{
							name = "Test amplitude Connector"
							description = "A description for the amplitude connector"
    						api_key = "mhvece"
    						server_url = "wluvduhqc"
    						server_zone = "lr32nx7xfo"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.amplitude.#":             1,
				"connectors.amplitude.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.amplitude.0.name":        "Test amplitude Connector",
				"connectors.amplitude.0.description": "A description for the amplitude connector",
				"connectors.amplitude.0.api_key":     "mhvece",
				"connectors.amplitude.0.server_url":  "wluvduhqc",
				"connectors.amplitude.0.server_zone": "lr32nx7xfo",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"audit_webhook": [
						{
							name = "Test audit-webhook Connector"
							description = "A description for the audit-webhook connector"
    						base_url = "bceszn6"
    						authentication = {
    							bearer_token = "xhmqmkcfhe4mk6"
    						}
    						headers = {
    							"key" = "g6htpmp"
    						}
    						hmac_secret = "ooxzct5yxz"
    						insecure = true
    						audit_filters = [{ key = "actions", operator = "includes", values = ["kekpon4oj34w"] }]
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.audit_webhook.#":                             1,
				"connectors.audit_webhook.0.id":                          testacc.AttributeHasPrefix("CI"),
				"connectors.audit_webhook.0.name":                        "Test audit-webhook Connector",
				"connectors.audit_webhook.0.description":                 "A description for the audit-webhook connector",
				"connectors.audit_webhook.0.base_url":                    "bceszn6",
				"connectors.audit_webhook.0.authentication.bearer_token": "xhmqmkcfhe4mk6",
				"connectors.audit_webhook.0.headers.key":                 "g6htpmp",
				"connectors.audit_webhook.0.hmac_secret":                 "ooxzct5yxz",
				"connectors.audit_webhook.0.insecure":                    true,
				"connectors.audit_webhook.0.audit_filters.0.values":      []string{"kekpon4oj34w"},
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"aws_s3": [
						{
							name = "Test aws-s3 Connector"
							description = "A description for the aws-s3 connector"
    						access_key_id = "ezzrllbqu22"
    						secret_access_key = "xiyuadzk4w64hog"
    						region = "y2l5fg"
    						bucket = "ywdzxd"
    						audit_enabled = true
    						audit_filters = [{ key = "actions", operator = "includes", values = ["kekpon4oj34w"] }]
    						troubleshoot_log_enabled = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.aws_s3.#":                          1,
				"connectors.aws_s3.0.id":                       testacc.AttributeHasPrefix("CI"),
				"connectors.aws_s3.0.name":                     "Test aws-s3 Connector",
				"connectors.aws_s3.0.description":              "A description for the aws-s3 connector",
				"connectors.aws_s3.0.access_key_id":            "ezzrllbqu22",
				"connectors.aws_s3.0.secret_access_key":        "xiyuadzk4w64hog",
				"connectors.aws_s3.0.region":                   "y2l5fg",
				"connectors.aws_s3.0.bucket":                   "ywdzxd",
				"connectors.aws_s3.0.audit_enabled":            true,
				"connectors.aws_s3.0.audit_filters.0.values":   []string{"kekpon4oj34w"},
				"connectors.aws_s3.0.troubleshoot_log_enabled": true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"aws_translate": [
						{
							name = "Test aws-translate Connector"
							description = "A description for the aws-translate connector"
    						access_key_id = "ezzrllbqu22"
    						secret_access_key = "xiyuadzk4w64hog"
    						session_token = "wnx4upgg3mft"
    						region = "y2l5fg"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.aws_translate.#":                   1,
				"connectors.aws_translate.0.id":                testacc.AttributeHasPrefix("CI"),
				"connectors.aws_translate.0.name":              "Test aws-translate Connector",
				"connectors.aws_translate.0.description":       "A description for the aws-translate connector",
				"connectors.aws_translate.0.access_key_id":     "ezzrllbqu22",
				"connectors.aws_translate.0.secret_access_key": "xiyuadzk4w64hog",
				"connectors.aws_translate.0.session_token":     "wnx4upgg3mft",
				"connectors.aws_translate.0.region":            "y2l5fg",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"clear": [
						{
							name = "Test clear Connector"
							description = "A description for the clear connector"
    						project_id = "yhw7b6yel"
    						api_key = "mhvece"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.clear.#":             1,
				"connectors.clear.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.clear.0.name":        "Test clear Connector",
				"connectors.clear.0.description": "A description for the clear connector",
				"connectors.clear.0.project_id":  "yhw7b6yel",
				"connectors.clear.0.api_key":     "mhvece",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"datadog": [
						{
							name = "Test datadog Connector"
							description = "A description for the datadog connector"
    						api_key = "mhvece"
    						site = "7oxa"
    						audit_enabled = true
    						audit_filters = [{ key = "actions", operator = "includes", values = ["kekpon4oj34w"] }]
    						troubleshoot_log_enabled = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.datadog.#":                          1,
				"connectors.datadog.0.id":                       testacc.AttributeHasPrefix("CI"),
				"connectors.datadog.0.name":                     "Test datadog Connector",
				"connectors.datadog.0.description":              "A description for the datadog connector",
				"connectors.datadog.0.api_key":                  "mhvece",
				"connectors.datadog.0.site":                     "7oxa",
				"connectors.datadog.0.audit_enabled":            true,
				"connectors.datadog.0.audit_filters.0.values":   []string{"kekpon4oj34w"},
				"connectors.datadog.0.troubleshoot_log_enabled": true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"devrev_grow": [
						{
							name = "Test devrev-grow Connector"
							description = "A description for the devrev-grow connector"
    						api_key = "mhvece"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.devrev_grow.#":             1,
				"connectors.devrev_grow.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.devrev_grow.0.name":        "Test devrev-grow Connector",
				"connectors.devrev_grow.0.description": "A description for the devrev-grow connector",
				"connectors.devrev_grow.0.api_key":     "mhvece",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"docebo": [
						{
							name = "Test docebo Connector"
							description = "A description for the docebo connector"
    						base_url = "bceszn6"
    						client_id = "sgetgyvq"
    						client_secret = "tjpxl7uy4wbb"
    						username = "c33yu7ld"
    						password = "l2eergg2"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.docebo.#":               1,
				"connectors.docebo.0.id":            testacc.AttributeHasPrefix("CI"),
				"connectors.docebo.0.name":          "Test docebo Connector",
				"connectors.docebo.0.description":   "A description for the docebo connector",
				"connectors.docebo.0.base_url":      "bceszn6",
				"connectors.docebo.0.client_id":     "sgetgyvq",
				"connectors.docebo.0.client_secret": "tjpxl7uy4wbb",
				"connectors.docebo.0.username":      "c33yu7ld",
				"connectors.docebo.0.password":      "l2eergg2",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"fingerprint": [
						{
							name = "Test fingerprint Connector"
							description = "A description for the fingerprint connector"
    						public_api_key = "htt624yz4z6i"
    						secret_api_key = "qxt75gbg4234"
    						use_cloudflare_integration = true
    						cloudflare_script_url = "p5sop7bd2jskwpzwdm6"
    						cloudflare_endpoint_url = "ad7li7hhec3doqaf33abq"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.fingerprint.#":                            1,
				"connectors.fingerprint.0.id":                         testacc.AttributeHasPrefix("CI"),
				"connectors.fingerprint.0.name":                       "Test fingerprint Connector",
				"connectors.fingerprint.0.description":                "A description for the fingerprint connector",
				"connectors.fingerprint.0.public_api_key":             "htt624yz4z6i",
				"connectors.fingerprint.0.secret_api_key":             "qxt75gbg4234",
				"connectors.fingerprint.0.use_cloudflare_integration": true,
				"connectors.fingerprint.0.cloudflare_script_url":      "p5sop7bd2jskwpzwdm6",
				"connectors.fingerprint.0.cloudflare_endpoint_url":    "ad7li7hhec3doqaf33abq",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"fingerprint_descope": [
						{
							name = "Test fingerprint-descope Connector"
							description = "A description for the fingerprint-descope connector"
    						custom_domain = "chk55vpucvwg"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.fingerprint_descope.#":               1,
				"connectors.fingerprint_descope.0.id":            testacc.AttributeHasPrefix("CI"),
				"connectors.fingerprint_descope.0.name":          "Test fingerprint-descope Connector",
				"connectors.fingerprint_descope.0.description":   "A description for the fingerprint-descope connector",
				"connectors.fingerprint_descope.0.custom_domain": "chk55vpucvwg",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"forter": [
						{
							name = "Test forter Connector"
							description = "A description for the forter connector"
    						site_id = "raavxa"
    						secret_key = "wi4bhwt7a"
    						overrides = true
    						override_ip_address = "urtsgzb7hjbuj5k3y"
    						override_user_email = "kqyebsfvqy6w6y6wn"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.forter.#":                     1,
				"connectors.forter.0.id":                  testacc.AttributeHasPrefix("CI"),
				"connectors.forter.0.name":                "Test forter Connector",
				"connectors.forter.0.description":         "A description for the forter connector",
				"connectors.forter.0.site_id":             "raavxa",
				"connectors.forter.0.secret_key":          "wi4bhwt7a",
				"connectors.forter.0.overrides":           true,
				"connectors.forter.0.override_ip_address": "urtsgzb7hjbuj5k3y",
				"connectors.forter.0.override_user_email": "kqyebsfvqy6w6y6wn",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"generic_sms_gateway": [
						{
							name = "Test generic-sms-gateway Connector"
							description = "A description for the generic-sms-gateway connector"
    						post_url = "efnv6ac"
    						sender = "bi3hxe"
    						authentication = {
    							bearer_token = "xhmqmkcfhe4mk6"
    						}
    						headers = {
    							"key" = "g6htpmp"
    						}
    						hmac_secret = "ooxzct5yxz"
    						insecure = true
    						use_static_ips = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.generic_sms_gateway.#":                             1,
				"connectors.generic_sms_gateway.0.id":                          testacc.AttributeHasPrefix("CI"),
				"connectors.generic_sms_gateway.0.name":                        "Test generic-sms-gateway Connector",
				"connectors.generic_sms_gateway.0.description":                 "A description for the generic-sms-gateway connector",
				"connectors.generic_sms_gateway.0.post_url":                    "efnv6ac",
				"connectors.generic_sms_gateway.0.sender":                      "bi3hxe",
				"connectors.generic_sms_gateway.0.authentication.bearer_token": "xhmqmkcfhe4mk6",
				"connectors.generic_sms_gateway.0.headers.key":                 "g6htpmp",
				"connectors.generic_sms_gateway.0.hmac_secret":                 "ooxzct5yxz",
				"connectors.generic_sms_gateway.0.insecure":                    true,
				"connectors.generic_sms_gateway.0.use_static_ips":              true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"google_cloud_translation": [
						{
							name = "Test google-cloud-translation Connector"
							description = "A description for the google-cloud-translation connector"
    						project_id = "yhw7b6yel"
    						service_account_json = "4wrifr235ikphcluwt"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.google_cloud_translation.#":                      1,
				"connectors.google_cloud_translation.0.id":                   testacc.AttributeHasPrefix("CI"),
				"connectors.google_cloud_translation.0.name":                 "Test google-cloud-translation Connector",
				"connectors.google_cloud_translation.0.description":          "A description for the google-cloud-translation connector",
				"connectors.google_cloud_translation.0.project_id":           "yhw7b6yel",
				"connectors.google_cloud_translation.0.service_account_json": "4wrifr235ikphcluwt",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"hibp": [
						{
							name = "Test hibp Connector"
							description = "A description for the hibp connector"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.hibp.#":             1,
				"connectors.hibp.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.hibp.0.name":        "Test hibp Connector",
				"connectors.hibp.0.description": "A description for the hibp connector",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"http": [
						{
							name = "Test http Connector"
							description = "A description for the http connector"
    						base_url = "bceszn6"
    						authentication = {
    							bearer_token = "xhmqmkcfhe4mk6"
    						}
    						headers = {
    							"key" = "g6htpmp"
    						}
    						hmac_secret = "ooxzct5yxz"
    						insecure = true
    						include_headers_in_context = true
    						use_static_ips = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.http.#":                             1,
				"connectors.http.0.id":                          testacc.AttributeHasPrefix("CI"),
				"connectors.http.0.name":                        "Test http Connector",
				"connectors.http.0.description":                 "A description for the http connector",
				"connectors.http.0.base_url":                    "bceszn6",
				"connectors.http.0.authentication.bearer_token": "xhmqmkcfhe4mk6",
				"connectors.http.0.headers.key":                 "g6htpmp",
				"connectors.http.0.hmac_secret":                 "ooxzct5yxz",
				"connectors.http.0.insecure":                    true,
				"connectors.http.0.include_headers_in_context":  true,
				"connectors.http.0.use_static_ips":              true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"hubspot": [
						{
							name = "Test hubspot Connector"
							description = "A description for the hubspot connector"
    						access_token = "ssrho3t3233"
    						base_url = "bceszn6"
    						use_static_ips = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.hubspot.#":                1,
				"connectors.hubspot.0.id":             testacc.AttributeHasPrefix("CI"),
				"connectors.hubspot.0.name":           "Test hubspot Connector",
				"connectors.hubspot.0.description":    "A description for the hubspot connector",
				"connectors.hubspot.0.access_token":   "ssrho3t3233",
				"connectors.hubspot.0.base_url":       "bceszn6",
				"connectors.hubspot.0.use_static_ips": true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"incode": [
						{
							name = "Test incode Connector"
							description = "A description for the incode connector"
    						api_key = "mhvece"
    						api_url = "dg2kp4"
    						flow_id = "xkg6re"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.incode.#":             1,
				"connectors.incode.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.incode.0.name":        "Test incode Connector",
				"connectors.incode.0.description": "A description for the incode connector",
				"connectors.incode.0.api_key":     "mhvece",
				"connectors.incode.0.api_url":     "dg2kp4",
				"connectors.incode.0.flow_id":     "xkg6re",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"intercom": [
						{
							name = "Test intercom Connector"
							description = "A description for the intercom connector"
    						token = "hrdj5"
    						region = "y2l5fg"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.intercom.#":             1,
				"connectors.intercom.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.intercom.0.name":        "Test intercom Connector",
				"connectors.intercom.0.description": "A description for the intercom connector",
				"connectors.intercom.0.token":       "hrdj5",
				"connectors.intercom.0.region":      "y2l5fg",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"lokalise": [
						{
							name = "Test lokalise Connector"
							description = "A description for the lokalise connector"
    						api_token = "mybopddv"
    						project_id = "yhw7b6yel"
    						team_id = "ontsaz"
    						card_id = "uo4way"
    						translation_provider = "zdmwgn7cvt7zfpsmrww"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.lokalise.#":                      1,
				"connectors.lokalise.0.id":                   testacc.AttributeHasPrefix("CI"),
				"connectors.lokalise.0.name":                 "Test lokalise Connector",
				"connectors.lokalise.0.description":          "A description for the lokalise connector",
				"connectors.lokalise.0.api_token":            "mybopddv",
				"connectors.lokalise.0.project_id":           "yhw7b6yel",
				"connectors.lokalise.0.team_id":              "ontsaz",
				"connectors.lokalise.0.card_id":              "uo4way",
				"connectors.lokalise.0.translation_provider": "zdmwgn7cvt7zfpsmrww",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"mparticle": [
						{
							name = "Test mparticle Connector"
							description = "A description for the mparticle connector"
    						api_key = "mhvece"
    						api_secret = "hgg666mus"
    						base_url = "bceszn6"
    						default_environment = "mekqliza6drwrn7azt"
    						use_static_ips = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.mparticle.#":                     1,
				"connectors.mparticle.0.id":                  testacc.AttributeHasPrefix("CI"),
				"connectors.mparticle.0.name":                "Test mparticle Connector",
				"connectors.mparticle.0.description":         "A description for the mparticle connector",
				"connectors.mparticle.0.api_key":             "mhvece",
				"connectors.mparticle.0.api_secret":          "hgg666mus",
				"connectors.mparticle.0.base_url":            "bceszn6",
				"connectors.mparticle.0.default_environment": "mekqliza6drwrn7azt",
				"connectors.mparticle.0.use_static_ips":      true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"newrelic": [
						{
							name = "Test newrelic Connector"
							description = "A description for the newrelic connector"
    						api_key = "mhvece"
    						data_center = "wjih4vgzrt"
    						audit_enabled = true
    						audit_filters = [{ key = "actions", operator = "includes", values = ["kekpon4oj34w"] }]
    						troubleshoot_log_enabled = true
    						override_logs_prefix = true
    						logs_prefix = "2zcsbfwbhp"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.newrelic.#":                          1,
				"connectors.newrelic.0.id":                       testacc.AttributeHasPrefix("CI"),
				"connectors.newrelic.0.name":                     "Test newrelic Connector",
				"connectors.newrelic.0.description":              "A description for the newrelic connector",
				"connectors.newrelic.0.api_key":                  "mhvece",
				"connectors.newrelic.0.data_center":              "wjih4vgzrt",
				"connectors.newrelic.0.audit_enabled":            true,
				"connectors.newrelic.0.audit_filters.0.values":   []string{"kekpon4oj34w"},
				"connectors.newrelic.0.troubleshoot_log_enabled": true,
				"connectors.newrelic.0.override_logs_prefix":     true,
				"connectors.newrelic.0.logs_prefix":              "2zcsbfwbhp",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"recaptcha": [
						{
							name = "Test recaptcha Connector"
							description = "A description for the recaptcha connector"
    						site_key = "ikzbbly"
    						secret_key = "wi4bhwt7a"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.recaptcha.#":             1,
				"connectors.recaptcha.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.recaptcha.0.name":        "Test recaptcha Connector",
				"connectors.recaptcha.0.description": "A description for the recaptcha connector",
				"connectors.recaptcha.0.site_key":    "ikzbbly",
				"connectors.recaptcha.0.secret_key":  "wi4bhwt7a",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"recaptcha_enterprise": [
						{
							name = "Test recaptcha-enterprise Connector"
							description = "A description for the recaptcha-enterprise connector"
    						project_id = "yhw7b6yel"
    						site_key = "ikzbbly"
    						api_key = "mhvece"
    						override_assessment = true
    						assessment_score = 15
    						enterprise = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.recaptcha_enterprise.#":                     1,
				"connectors.recaptcha_enterprise.0.id":                  testacc.AttributeHasPrefix("CI"),
				"connectors.recaptcha_enterprise.0.name":                "Test recaptcha-enterprise Connector",
				"connectors.recaptcha_enterprise.0.description":         "A description for the recaptcha-enterprise connector",
				"connectors.recaptcha_enterprise.0.project_id":          "yhw7b6yel",
				"connectors.recaptcha_enterprise.0.site_key":            "ikzbbly",
				"connectors.recaptcha_enterprise.0.api_key":             "mhvece",
				"connectors.recaptcha_enterprise.0.override_assessment": true,
				"connectors.recaptcha_enterprise.0.assessment_score":    15,
				"connectors.recaptcha_enterprise.0.enterprise":          true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"rekognition": [
						{
							name = "Test rekognition Connector"
							description = "A description for the rekognition connector"
    						access_key_id = "ezzrllbqu22"
    						secret_access_key = "xiyuadzk4w64hog"
    						collection_id = "y3z2olrexe5d"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.rekognition.#":                   1,
				"connectors.rekognition.0.id":                testacc.AttributeHasPrefix("CI"),
				"connectors.rekognition.0.name":              "Test rekognition Connector",
				"connectors.rekognition.0.description":       "A description for the rekognition connector",
				"connectors.rekognition.0.access_key_id":     "ezzrllbqu22",
				"connectors.rekognition.0.secret_access_key": "xiyuadzk4w64hog",
				"connectors.rekognition.0.collection_id":     "y3z2olrexe5d",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"salesforce": [
						{
							name = "Test salesforce Connector"
							description = "A description for the salesforce connector"
    						base_url = "bceszn6"
    						client_id = "sgetgyvq"
    						client_secret = "tjpxl7uy4wbb"
    						version = "lssphbi"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.salesforce.#":               1,
				"connectors.salesforce.0.id":            testacc.AttributeHasPrefix("CI"),
				"connectors.salesforce.0.name":          "Test salesforce Connector",
				"connectors.salesforce.0.description":   "A description for the salesforce connector",
				"connectors.salesforce.0.base_url":      "bceszn6",
				"connectors.salesforce.0.client_id":     "sgetgyvq",
				"connectors.salesforce.0.client_secret": "tjpxl7uy4wbb",
				"connectors.salesforce.0.version":       "lssphbi",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"salesforce_marketing_cloud": [
						{
							name = "Test salesforce-marketing-cloud Connector"
							description = "A description for the salesforce-marketing-cloud connector"
    						subdomain = "ko54m53in"
    						client_id = "sgetgyvq"
    						client_secret = "tjpxl7uy4wbb"
    						scope = "l4lbz"
    						account_id = "4bpveggea"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.salesforce_marketing_cloud.#":               1,
				"connectors.salesforce_marketing_cloud.0.id":            testacc.AttributeHasPrefix("CI"),
				"connectors.salesforce_marketing_cloud.0.name":          "Test salesforce-marketing-cloud Connector",
				"connectors.salesforce_marketing_cloud.0.description":   "A description for the salesforce-marketing-cloud connector",
				"connectors.salesforce_marketing_cloud.0.subdomain":     "ko54m53in",
				"connectors.salesforce_marketing_cloud.0.client_id":     "sgetgyvq",
				"connectors.salesforce_marketing_cloud.0.client_secret": "tjpxl7uy4wbb",
				"connectors.salesforce_marketing_cloud.0.scope":         "l4lbz",
				"connectors.salesforce_marketing_cloud.0.account_id":    "4bpveggea",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"segment": [
						{
							name = "Test segment Connector"
							description = "A description for the segment connector"
    						write_key = "vs2jik2d"
    						host = "i5ak"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.segment.#":             1,
				"connectors.segment.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.segment.0.name":        "Test segment Connector",
				"connectors.segment.0.description": "A description for the segment connector",
				"connectors.segment.0.write_key":   "vs2jik2d",
				"connectors.segment.0.host":        "i5ak",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"sendgrid": [
						{
							name = "Test sendgrid Connector"
							description = "A description for the sendgrid connector"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"ses": [
						{
							name = "Test ses Connector"
							description = "A description for the ses connector"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"slack": [
						{
							name = "Test slack Connector"
							description = "A description for the slack connector"
    						token = "hrdj5"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.slack.#":             1,
				"connectors.slack.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.slack.0.name":        "Test slack Connector",
				"connectors.slack.0.description": "A description for the slack connector",
				"connectors.slack.0.token":       "hrdj5",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"smartling": [
						{
							name = "Test smartling Connector"
							description = "A description for the smartling connector"
    						user_identifier = "h2sy3jtqq4rxwi"
    						user_secret = "gchtcl2tno"
    						account_uid = "7qxonan5tu"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.smartling.#":                 1,
				"connectors.smartling.0.id":              testacc.AttributeHasPrefix("CI"),
				"connectors.smartling.0.name":            "Test smartling Connector",
				"connectors.smartling.0.description":     "A description for the smartling connector",
				"connectors.smartling.0.user_identifier": "h2sy3jtqq4rxwi",
				"connectors.smartling.0.user_secret":     "gchtcl2tno",
				"connectors.smartling.0.account_uid":     "7qxonan5tu",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"smtp": [
						{
							name = "Test smtp Connector"
							description = "A description for the smtp connector"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"sns": [
						{
							name = "Test sns Connector"
							description = "A description for the sns connector"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"sumologic": [
						{
							name = "Test sumologic Connector"
							description = "A description for the sumologic connector"
    						http_source_url = "shhxfygq6qspm"
    						audit_enabled = true
    						audit_filters = [{ key = "actions", operator = "includes", values = ["kekpon4oj34w"] }]
    						troubleshoot_log_enabled = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.sumologic.#":                          1,
				"connectors.sumologic.0.id":                       testacc.AttributeHasPrefix("CI"),
				"connectors.sumologic.0.name":                     "Test sumologic Connector",
				"connectors.sumologic.0.description":              "A description for the sumologic connector",
				"connectors.sumologic.0.http_source_url":          "shhxfygq6qspm",
				"connectors.sumologic.0.audit_enabled":            true,
				"connectors.sumologic.0.audit_filters.0.values":   []string{"kekpon4oj34w"},
				"connectors.sumologic.0.troubleshoot_log_enabled": true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"telesign": [
						{
							name = "Test telesign Connector"
							description = "A description for the telesign connector"
    						customer_id = "yn5scsxze2"
    						api_key = "mhvece"
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.telesign.#":             1,
				"connectors.telesign.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.telesign.0.name":        "Test telesign Connector",
				"connectors.telesign.0.description": "A description for the telesign connector",
				"connectors.telesign.0.customer_id": "yn5scsxze2",
				"connectors.telesign.0.api_key":     "mhvece",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"traceable": [
						{
							name = "Test traceable Connector"
							description = "A description for the traceable connector"
    						secret_key = "wi4bhwt7a"
    						eu_region = true
						}
					]
				}
			`),
			Check: p.Check(map[string]any{
				"connectors.traceable.#":             1,
				"connectors.traceable.0.id":          testacc.AttributeHasPrefix("CI"),
				"connectors.traceable.0.name":        "Test traceable Connector",
				"connectors.traceable.0.description": "A description for the traceable connector",
				"connectors.traceable.0.secret_key":  "wi4bhwt7a",
				"connectors.traceable.0.eu_region":   true,
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"twilio_core": [
						{
							name = "Test twilio-core Connector"
							description = "A description for the twilio-core connector"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
		},
		resource.TestStep{
			Config: p.Config(`
				connectors = {
					"twilio_verify": [
						{
							name = "Test twilio-verify Connector"
							description = "A description for the twilio-verify connector"
						}
					]
				}
			`),
			ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
		},
		resource.TestStep{
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"connectors.%": 0,
			}),
		},
	)
}

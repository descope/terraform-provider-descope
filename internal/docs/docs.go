// Code generated by terragen. DO NOT EDIT.

package docs

var docsProject = map[string]string{
	"name": "The name of the Descope project.",
	"environment": "This can be set to `production` to mark production projects, otherwise this should be " +
	               "left unset for development or staging projects.",
	"project_settings": "General settings for the Descope project.",
	"authentication": "Settings for each authentication method.",
	"authorization": "Define Role-Based Access Control (RBAC) for your users by creating roles and permissions.",
	"attributes": "Custom attributes that can be attached to users and tenants.",
	"connectors": "Enrich your flows by interacting with third party services.",
	"applications": "Applications that are registered with the project.",
	"jwt_templates": "Defines templates for JSON Web Tokens (JWT) used for authentication.",
	"styles": "Custom styles that can be applied to the project's authentication flows.",
	"flows": "",
}

var docsApplication = map[string]string{
	"oidc_applications": "Applications using OpenID Connect (OIDC) for authentication.",
	"saml_applications": "Applications using SAML for authentication.",
}

var docsOIDC = map[string]string{
	"id": "An optional identifier for the application.",
	"name": "The name of the application.",
	"description": "A brief description of the application.",
	"logo": "The URL of the logo associated with the application.",
	"disabled": "Indicates whether the resource or functionality is disabled.",
	"login_page_url": "The URL of the custom login page for the application.",
	"claims": "Claims associated with JWT tokens, typically used for user information.",
}

var docsSAML = map[string]string{
	"id": "An optional identifier for the application.",
	"name": "The name of the application.",
	"description": "A brief description of the application.",
	"logo": "The URL of the logo associated with the application.",
	"disabled": "Indicates whether the application is disabled.",
	"login_page_url": "The URL of the custom login page for the application.",
	"dynamic_configuration": "",
	"manual_configuration": "",
	"acs_allowed_callback_urls": "",
	"subject_name_id_type": "",
	"subject_name_id_format": "",
	"default_relay_state": "",
	"attribute_mapping": "Map user attributes from the third party identity provider to custom attributes in Descope.",
}

var docsAttributeMapping = map[string]string{
	"name": "The name of the user attribute in the third party identity provider.",
	"value": "The name of the user custom attribute in Descope.",
}

var docsDynamicConfiguration = map[string]string{
	"metadata_url": "",
}

var docsManualConfiguration = map[string]string{
	"acs_url": "",
	"entity_id": "",
	"certificate": "",
}

var docsAttributes = map[string]string{
	"tenant": "Custom attributes to store additional details about your tenants.",
	"user": "Custom attributes to store additional details about your users.",
}

var docsTenantAttribute = map[string]string{
	"name": "The name of the tenant attribute.",
	"type": "The type of the tenant attribute. Valid valus are `string`, `number`, `boolean`, " +
	        "`date`, `singleselect`, and `multiselect`.",
	"select_options": "A list of strings to define the set of options for select attributes.",
}

var docsUserAttribute = map[string]string{
	"name": "The name of the user attribute.",
	"type": "The type of the user attribute. Valid valus are `string`, `number`, `boolean`, " +
	        "`date`, `singleselect`, and `multiselect`.",
	"select_options": "A list of strings to define the set of options for select attributes.",
	"widget_authorization": "When provided, viewing and editing the attribute values in widgets will " +
	                        "be restricted to users with the specified permissions.",
}

var docsUserAttributeWidgetAuthorization = map[string]string{
	"view_permissions": "Viewing the attribute value in widgets will be restricted to users with " +
	                    "the specified permissions.",
	"edit_permissions": "Editing the attribute value in widgets will be restricted to users with " +
	                    "the specified permissions.",
}

var docsAuthentication = map[string]string{
	"otp": "A dynamically generated set of numbers, granting the user one-time access.",
	"magic_link": "An authentication method where a user receives a unique link via email to log in.",
	"enchanted_link": "An enhanced and more secure version of Magic Link, enabling users to start the authentication " +
	                  "process on one device and execute the verification on another.",
	"embedded_link": "Make the authentication experience smoother for the user by generating their initial token in a " +
	                 "way that does not require the end user to initiate the process, requiring only verification.",
	"password": "The classic username and password combination used for authentication.",
	"oauth": "Authentication using Open Authorization, which allows users to authenticate with various external " +
	         "services.",
	"sso": "Single Sign-On (SSO) authentication method that enables users to access multiple applications with " +
	       "a single set of credentials.",
	"totp": "A one-time code generated for the user using a shared secret and time.",
	"passkeys": "Device-based passwordless authentication, using fingerprint, face scan, and more.",
}

var docsEmbeddedLink = map[string]string{
	"enabled": "",
	"expiration_time": "The amount of time that the embedded link will be valid for.",
	"expiration_time_unit": "",
}

var docsEnchantedLink = map[string]string{
	"enabled": "",
	"expiration_time": "",
	"expiration_time_unit": "",
	"redirect_url": "The URL to redirect users to after they log in using the enchanted link.",
	"email_service": "Settings related to sending emails as part of the enchanted link authentication.",
}

var docsMagicLink = map[string]string{
	"enabled": "",
	"expiration_time": "",
	"expiration_time_unit": "",
	"redirect_url": "The URL to redirect users to after they log in using the magic link.",
	"email_service": "Settings related to sending emails as part of the magic link authentication.",
	"text_service": "Settings related to sending SMS messages as part of the magic link authentication.",
}

var docsOAuth = map[string]string{
	"disabled": "",
	"system": "Custom configurations for builtin OAuth providers such as Apple, Google, GitHub, Facebook, etc.",
	"custom": "Custom OAuth providers configured for this project.",
}

var docsOAuthSystemProvider = map[string]string{
	"apple": "Apple's OAuth provider, allowing users to authenticate with their Apple Account.",
	"discord": "Discord's OAuth provider, allowing users to authenticate with their Discord account.",
	"facebook": "Facebook's OAuth provider, allowing users to authenticate with their Facebook account.",
	"github": "GitHub's OAuth provider, allowing users to authenticate with their GitHub account.",
	"gitlab": "GitLab's OAuth provider, allowing users to authenticate with their GitLab account.",
	"google": "Google's OAuth provider, allowing users to authenticate with their Google account.",
	"linkedin": "LinkedIn's OAuth provider, allowing users to authenticate with their LinkedIn account.",
	"microsoft": "Microsoft's OAuth provider, allowing users to authenticate with their Microsoft account.",
	"slack": "Slack's OAuth provider, allowing users to authenticate with their Slack account.",
}

var docsOAuthProvider = map[string]string{
	"disabled": "",
	"client_id": "The client ID for the OAuth provider, used to identify the application to the provider.",
	"client_secret": "The client secret for the OAuth provider, used to authenticate the application with the provider.",
	"provider_token_management": "Settings related to token management for the OAuth provider.",
	"prompts": "Custom prompts or consent screens that users may see during OAuth authentication.",
	"scopes": "Scopes of access that the application requests from the user's account on the OAuth provider.",
	"merge_user_accounts": "Whether to merge existing user accounts with new ones created through OAuth authentication.",
	"description": "A brief description of the OAuth provider.",
	"logo": "The URL of the logo associated with the OAuth provider.",
	"grant_type": "The type of grant (`authorization_code` or `implicit`) to use when requesting access tokens " +
	              "from the OAuth provider.",
	"issuer": "",
	"authorization_endpoint": "The URL that users are redirected to for authorization with the OAuth provider.",
	"token_endpoint": "The URL where the application requests an access token from the OAuth provider.",
	"user_info_endpoint": "The URL where the application retrieves user information from the OAuth provider.",
	"jwks_endpoint": "The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.",
	"claim_mapping": "",
}

var docsOTP = map[string]string{
	"enabled": "",
	"domain": "The domain to embed in OTP messages.",
	"expiration_time": "The amount of time that an OTP code will be valid for.",
	"expiration_time_unit": "",
	"email_service": "Settings related to sending emails with OTP codes.",
	"text_service": "Settings related to sending SMS messages with OTP codes.",
	"voice_service": "Settings related to voice calls with OTP codes.",
}

var docsPasskeys = map[string]string{
	"enabled": "",
	"top_level_domain": "Passkeys will be usable in the following domain and all its subdomains.",
}

var docsPassword = map[string]string{
	"enabled": "",
	"expiration": "Whether users are required to change their password periodically.",
	"expiration_weeks": "The number of weeks after which a user's password expires and they need to replace it.",
	"lock": "Whether the user account should be locked after a specified number of failed login attempts.",
	"lock_attempts": "The number of failed login attempts allowed before an account is locked.",
	"lowercase": "Whether passwords must contain at least one lowercase letter.",
	"min_length": "The minimum length of the password that users are required to use. The maximum length is always `64`.",
	"non_alphanumeric": "Whether passwords must contain at least one non-alphanumeric character (e.g. `!`, `@`, `#`).",
	"number": "Whether passwords must contain at least one number.",
	"reuse": "Whether to forbid password reuse when users change their password.",
	"reuse_amount": "",
	"uppercase": "Whether passwords must contain at least one uppercase letter.",
	"email_service": "Settings related to sending password reset emails as part of the password feature.",
}

var docsSSO = map[string]string{
	"enabled": "",
	"merge_users": "Whether to merge existing user accounts with new ones created through SSO authentication.",
}

var docsTOTP = map[string]string{
	"enabled": "",
}

var docsAuthorization = map[string]string{
	"roles": "",
	"permissions": "",
}

var docsPermission = map[string]string{
	"name": "",
	"description": "",
}

var docsRole = map[string]string{
	"name": "",
	"description": "",
	"permissions": "",
}

var docsAbuseIPDB = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"api_key": "The unique AbuseIPDB API key.",
}

var docsAmplitude = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"api_key": "The Amplitude API Key generated for the Descope service.",
	"server_url": "The server URL of the Amplitude API, when using different api or a custom domain " +
	              "in Amplitude.",
	"server_zone": "`EU` or `US`. Sets the Amplitude server zone. Set this to `EU` for Amplitude " +
	               "projects created in `EU` data center. Default is `US`.",
}

var docsAuditWebhook = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"base_url": "The base URL to fetch",
	"authentication": "Authentication Information",
	"headers": "The headers to send with the request",
	"hmac_secret": "HMAC is a method for message signing with a symmetrical key. This secret will be " +
	               "used to sign the payload, and the resulting signature will be sent in the " +
	               "`x-descope-webhook-s256` header. The receiving service should use this secret to " +
	               "verify the integrity and authenticity of the payload by checking the provided " +
	               "signature",
	"insecure": "Will ignore certificate errors raised by the client",
	"audit_filters": "Specify which events will be sent to the external audit service (including " +
	                 "tenant selection).",
}

var docsAWSS3 = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"access_key_id": "The unique AWS access key ID.",
	"secret_access_key": "The secret AWS access key.",
	"region": "The AWS S3 region, e.g. `us-east-1`.",
	"bucket": "The AWS S3 bucket. This bucket should already exist for the connector to work.",
	"audit_enabled": "",
	"audit_filters": "",
	"troubleshoot_log_enabled": "",
}

var docsAWSTranslate = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"access_key_id": "AWS access key ID.",
	"secret_access_key": "AWS secret access key.",
	"session_token": "(Optional) A security or session token to use with these credentials. Usually " +
	                 "present for temporary credentials.",
	"region": "The AWS region to which this client will send requests. (e.g. us-east-1.)",
}

var docsClear = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"project_id": "Your CLEAR project ID.",
	"api_key": "Your CLEAR API key.",
}

var docsConnectors = map[string]string{
	"abuseipdb": "AbuseIPDB provides an API to identify if an IP address has been associated with " +
	             "malicious activities online.",
	"amplitude": "Amplitude, an analytics product that allows you to collects events from web and " +
	             "mobile apps, unify those and use those to better understand your customers " +
	             "needs.",
	"audit_webhook": "",
	"aws_s3": "",
	"aws_translate": "",
	"clear": "An identity verification platform that allow customers to digitally verify their " +
	         "identity from anywhere.",
	"datadog": "Datadog, an observability service for cloud-scale applications, providing " +
	           "monitoring of servers, databases, tools, and services, through a SaaS-based data " +
	           "analytics platform.",
	"devrev_grow": "",
	"docebo": "Docebo is a cloud-based Learning Management System (LMS) designed to increase " +
	          "performance and learning engagement.",
	"fingerprint": "Use the Fingerprint (formerly FingerprintJS) connector to add device " +
	               "intelligence and prevent fraud.",
	"fingerprint_descope": "",
	"forter": "Use the Forter connector for account fraud prevention.",
	"google_cloud_translation": "",
	"hibp": "API to check if password appeared previously exposed in data breaches.",
	"http": "A general purpose HTTP client",
	"http_static_ip": "",
	"hubspot": "HubSpot is a CRM platform with software, integrations, and resources needed to " +
	           "connect marketing, sales, content management, and customer service.",
	"intercom": "Intercom is a Conversational Relationship Platform (CRP).",
	"newrelic": "Use this connector to send audit events and troubleshooting logs to New Relic.",
	"recaptcha": "reCAPTCHA is a free google service that protects your site from spam and abuse. " +
	             "It uses advanced risk analysis techniques to tell humans and bots apart.",
	"recaptcha_enterprise": "",
	"rekognition": "AWS Rekognition, cloud-based AI service that offers computer vision capabilities " +
	               "for analyzing and processing images. Useful for registration and verification " +
	               "processes, and can be used to detect fraud and prevent identity theft.",
	"salesforce": "Salesforce is a leading cloud-based Customer Relationship Management (CRM) " +
	              "platform that helps businesses streamline their sales, service, and marketing " +
	              "operations.",
	"segment": "Segment, an analytics product that allows you to collects events from web and " +
	           "mobile apps, unify those and use those to better understand your customers " +
	           "needs.",
	"sendgrid": "",
	"smtp": "",
	"sumologic": "Sumo Logic, fast troubleshooting and investigation with AI/ML-powered log " +
	             "analytics",
	"telesign": "Telesign Phone number intelligence API provides risk score for phone numbers.",
	"traceable": "API security for a cloud-first, API-driven world.",
	"twilio_core": "",
	"twilio_verify": "",
	"veriff": "AI-powered identity verification solution for identity fraud prevention, Know " +
	          "Your Customer compliance, and fast conversions of valuable customers.",
}

var docsDatadog = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"api_key": "The unique Datadog organization key.",
	"site": "The Datadog site to send logs to. Default is `datadoghq.com`. European, free " +
	        "tier and other customers should set their site accordingly.",
	"audit_enabled": "",
	"audit_filters": "",
	"troubleshoot_log_enabled": "",
}

var docsDevRevGrow = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"api_key": "Authentication to DevRev APIs requires a personal access token (PAT).",
}

var docsDocebo = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"base_url": "The Docebo api base url.",
	"client_id": "The Docebo OAuth 2.0 app client ID.",
	"client_secret": "The Docebo OAuth 2.0 app client secret.",
	"username": "The Docebo username.",
	"password": "The Docebo user's password.",
}

var docsFingerprint = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"public_api_key": "The Fingerprint public API key.",
	"secret_api_key": "The Fingerprint secret API key.",
	"use_cloudflare_integration": "Enable to configure the relevant Cloudflare integration parameters if Cloudflare " +
	                              "integration is set in your Fingerprint account.",
	"cloudflare_script_url": "The Cloudflare integration Script URL.",
	"cloudflare_endpoint_url": "The Cloudflare integration Endpoint URL.",
}

var docsFingerprintDescope = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"custom_domain": "The custom domain to fetch",
}

var docsForter = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"site_id": "The Forter site ID.",
	"secret_key": "The Forter secret key.",
	"overrides": "Override the user's IP address or email so that Forter can provide a specific " +
	             "decision or recommendation. Contact the Forter team for further details. Note: " +
	             "Overriding the user IP address or email is intended for testing purpose and " +
	             "should not be utilized in production environments.",
	"override_ip_address": "Override the user IP address.",
	"override_user_email": "Override the user email.",
}

var docsGoogleCloudTranslation = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"project_id": "The Google Cloud project ID where the Google Cloud Translation is managed.",
	"service_account_json": "Service Account JSON associated with the current project.",
}

var docsHIBP = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
}

var docsHTTP = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"base_url": "The base URL to fetch",
	"authentication": "Authentication Information",
	"headers": "The headers to send with the request",
	"hmac_secret": "HMAC is a method for message signing with a symmetrical key. This secret will be " +
	               "used to sign the base64 encoded payload, and the resulting signature will be " +
	               "sent in the `x-descope-webhook-s256` header. The receiving service should use " +
	               "this secret to verify the integrity and authenticity of the payload by checking " +
	               "the provided signature",
	"insecure": "Will ignore certificate errors raised by the client",
	"include_headers_in_context": "The connector response context will also include the headers. The context will " +
	                              "have a \"body\" attribute and a \"headers\" attribute. See more details in the help " +
	                              "guide",
}

var docsHttpStaticIP = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"base_url": "The base URL to fetch",
	"authentication": "Authentication Information",
	"headers": "The headers to send with the request",
	"hmac_secret": "HMAC is a method for message signing with a symmetrical key. This secret will be " +
	               "used to sign the base64 encoded payload, and the resulting signature will be " +
	               "sent in the `x-descope-webhook-s256` header. The receiving service should use " +
	               "this secret to verify the integrity and authenticity of the payload by checking " +
	               "the provided signature",
	"insecure": "Will ignore certificate errors raised by the client",
	"include_headers_in_context": "The connector response context will also include the headers. The context will " +
	                              "have a \"body\" attribute and a \"headers\" attribute. See more details in the help " +
	                              "guide",
}

var docsHubSpot = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"access_token": "The HubSpot private API access token generated for the Descope service.",
	"base_url": "The base URL of the HubSpot API, when using a custom domain in HubSpot, default " +
	            "value is https://api.hubapi.com .",
}

var docsIntercom = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"token": "The Intercom access token.",
	"region": "Regional Hosting - US, EU, or AU. default: US",
}

var docsNewRelic = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"api_key": "Ingest License Key of the account you want to report data to.",
	"data_center": "The New Relic data center the account belongs to. Possible values are: `US`, " +
	               "`EU`, `FedRAMP`. Default is `US`.",
	"audit_enabled": "",
	"audit_filters": "",
	"troubleshoot_log_enabled": "",
	"override_logs_prefix": "Enable this option to use a custom prefix for log fields.",
	"logs_prefix": "Specify a custom prefix for all log fields. The default prefix is `descope.`.",
}

var docsRecaptcha = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"site_key": "The site key is used to invoke reCAPTCHA service on your site or mobile " +
	            "application.",
	"secret_key": "The secret key authorizes communication between Descope backend and the " +
	              "reCAPTCHA server to verify the user's response.",
}

var docsRecaptchaEnterprise = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"project_id": "The Google Cloud project ID where the reCAPTCHA Enterprise is managed.",
	"site_key": "The site key is used to invoke reCAPTCHA Enterprise service on your site or " +
	            "mobile application.",
	"api_key": "API key associated with the current project.",
	"override_assessment": "Override the default assessment model. Note: Overriding assessment is intended " +
	                       "for automated testing and should not be utilized in production environments.",
	"assessment_score": "When configured, the Recaptcha action will return the score without assessing " +
	                    "the request. The score ranges between 0 and 1, where 1 is a human interaction " +
	                    "and 0 is a bot.",
}

var docsRekognition = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"access_key_id": "The AWS access key ID",
	"secret_access_key": "The AWS secret access key",
	"collection_id": "The collection to store registered users in. Should match `[a-zA-Z0-9_.-]+` " +
	                 "pattern. Changing this will cause losing existing users.",
}

var docsSalesforce = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"base_url": "The Salesforce API base URL.",
	"client_id": "The consumer key of the connected app.",
	"client_secret": "The consumer secret of the connected app.",
	"version": "REST API Version.",
}

var docsSegment = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"write_key": "The Segment Write Key generated for the Descope service.",
	"host": "The base URL of the Segment API, when using a custom domain in Segment.",
}

var docsSendGrid = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"sender": "",
	"authentication": "",
}

var docsSendGridAuthField = map[string]string{
	"api_key": "",
}

var docsSenderField = map[string]string{
	"email": "",
	"name": "",
}

var docsServerField = map[string]string{
	"host": "",
	"port": "",
}

var docsHTTPAuthField = map[string]string{
	"bearer_token": "",
	"basic": "",
	"api_key": "",
}

var docsHTTPAuthBasicField = map[string]string{
	"username": "",
	"password": "",
}

var docsHTTPAuthAPIKeyField = map[string]string{
	"key": "",
	"token": "",
}

var docsSMTP = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"sender": "",
	"server": "",
	"authentication": "",
}

var docsSMTPAuthField = map[string]string{
	"username": "",
	"password": "",
	"method": "",
}

var docsSumoLogic = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"http_source_url": "The URL associated with an HTTP Hosted collector",
	"audit_enabled": "",
	"audit_filters": "",
	"troubleshoot_log_enabled": "",
}

var docsTelesign = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"customer_id": "The unique Telesign account Customer ID",
	"api_key": "The unique Telesign API key",
}

var docsTraceable = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"secret_key": "The Traceable secret key.",
	"eu_region": "EU(Europe) Region deployment of Traceable platform.",
}

var docsTwilioCore = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"account_sid": "",
	"senders": "",
	"authentication": "",
}

var docsTwilioCoreSendersField = map[string]string{
	"sms": "",
	"voice": "",
}

var docsTwilioCoreSendersSMSField = map[string]string{
	"phone_number": "",
	"messaging_service_sid": "",
}

var docsTwilioCoreSendersVoiceField = map[string]string{
	"phone_number": "",
}

var docsTwilioAuthField = map[string]string{
	"auth_token": "",
	"api_key": "",
	"api_secret": "",
}

var docsTwilioVerify = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"account_sid": "",
	"service_sid": "",
	"sender": "",
	"authentication": "",
}

var docsVeriff = map[string]string{
	"name": "A custom name for your connector.",
	"description": "A description of what your connector is used for.",
	"api_key": "The Veriff Public API Key, you can find under Veriff Station - Integrations.",
	"secret_key": "The Veriff Shared secret key, you can find under Veriff Station - Integrations.",
	"base_url": "The base URL of the Veriff API, default value is https://stationapi.veriff.com.",
}

var docsFlow = map[string]string{
	"data": "",
}

var docsStyles = map[string]string{
	"data": "",
}

var docsJWTTemplate = map[string]string{
	"name": "",
	"description": "",
	"auth_schema": "",
	"conformance_issuer": "",
	"template": "",
}

var docsJWTTemplates = map[string]string{
	"user_templates": "",
	"access_key_templates": "",
}

var docsSettings = map[string]string{
	"cookie_policy": "",
	"domain": "",
	"enable_inactivity": "",
	"inactivity_time": "",
	"refresh_token_expiration": "",
	"user_jwt_template": "",
	"access_key_jwt_template": "",
}

var docsEmailService = map[string]string{
	"connector": "",
	"templates": "",
}

var docsEmailTemplate = map[string]string{
	"active": "",
	"name": "",
	"subject": "",
	"html_body": "",
	"plain_text_body": "",
	"use_plain_text_body": "",
}

var docsTextService = map[string]string{
	"connector": "",
	"templates": "",
}

var docsTextTemplate = map[string]string{
	"active": "",
	"name": "",
	"body": "",
}

var docsVoiceService = map[string]string{
	"connector": "",
	"templates": "",
}

var docsVoiceTemplate = map[string]string{
	"active": "",
	"name": "",
	"body": "",
}

# Stage 2: standalone resource configuration for the objects the legacy resource owned.
# Apply this with the provider pinned to 2.0.0, together with imports.tf.
# Server ids are written as literals so the adoption plan is evaluable before anything has been applied.
# Replacing a literal id with a reference to the resource that owns it is a no-op refactor to do afterwards.
# This file belongs in the module that declares descope_project.main.

resource "descope_project" "main" {
  name        = "acme-prod"
  environment = "production"
  tags        = ["team-auth"]
}

resource "descope_access_key_attribute" "scope" {
  project_id     = "P2abcdefghijklmnop"
  id             = "scope"
  name           = "Scope"
  select_options = []
  type           = "string"
}

resource "descope_admin_portal" "admin_portal" {
  project_id = "P2abcdefghijklmnop"
  enabled    = true
  style_id   = ""
  widgets = [
    {
      type      = "user-management"
      widget_id = "user-management"
    },
  ]
}

resource "descope_app_permission" "read" {
  project_id  = "P2abcdefghijklmnop"
  app_id      = "APP-oidc"
  name        = "read"
  description = "read"
}

resource "descope_app_role" "viewer" {
  project_id     = "P2abcdefghijklmnop"
  app_id         = "APP-oidc"
  name           = "viewer"
  description    = ""
  permission_ids = ["AP-read"]
  role_mappings  = []
}

resource "descope_email_template" "magiclink_email" {
  project_id          = "P2abcdefghijklmnop"
  method              = "magiclink"
  name                = "magiclink email"
  html_body           = "<p>hi</p>"
  plain_text_body     = ""
  subject             = "Sign in"
  use_plain_text_body = false
}

resource "descope_email_template" "otp_email" {
  project_id          = "P2abcdefghijklmnop"
  method              = "otp"
  name                = "otp email"
  html_body           = "<p>hi</p>"
  plain_text_body     = ""
  subject             = "Your code"
  use_plain_text_body = false
}

resource "descope_embeddedlink_settings" "embeddedlink_settings" {
  project_id      = "P2abcdefghijklmnop"
  expiration_time = "1 hour"
}

resource "descope_fga_schema" "fga_schema" {
  project_id = "P2abcdefghijklmnop"
  schema     = "model AuthZ 1.0\ntype user"
}

locals {
  descope_flow_flows = {
    sign-up-or-in = "payloads/descope_flow.flows.sign_up_or_in_7c0256edde1e.data.json"
    step-up       = "payloads/descope_flow.flows.step_up_72e408983887.data.json"
  }
}

resource "descope_flow" "flows" {
  for_each   = local.descope_flow_flows
  project_id = "P2abcdefghijklmnop"
  flow_id    = each.key
  data       = file("${path.module}/${each.value}")
}

resource "descope_invite_settings" "invite_settings" {
  project_id          = "P2abcdefghijklmnop"
  add_magiclink_token = false
  email_service = {
    connector_id = "C-smtp"
    templates = [
      {
        active              = true
        html_body           = "<p>hi</p>"
        name                = "invite"
        plain_text_body     = ""
        subject             = "You are invited"
        use_plain_text_body = false
      },
    ]
  }
  expire_invited_users = false
  invite_expiration    = "1 week"
  invite_url           = "https://app.acme.test/invite"
  require_invitation   = true
  send_email           = true
  send_text            = false
}

resource "descope_jwt_template" "user_claims" {
  project_id               = "P2abcdefghijklmnop"
  name                     = "user-claims"
  add_jti_claim            = false
  auth_schema              = "default"
  auto_tenant_claim        = false
  conformance_issuer       = false
  description              = ""
  empty_claim_policy       = "none"
  enforce_issuer           = false
  exclude_permission_claim = false
  override_subject_claim   = false
  template                 = file("${path.module}/payloads/descope_jwt_template.user_claims.template.json")
  type                     = "user"
}

resource "descope_list" "allowlist" {
  project_id  = "P2abcdefghijklmnop"
  name        = "allowlist"
  description = ""
  texts       = ["a@acme.test"]
}

resource "descope_list" "config" {
  project_id  = "P2abcdefghijklmnop"
  name        = "config"
  description = ""
  json        = file("${path.module}/payloads/descope_list.config.json.json")
}

resource "descope_list" "office_ips" {
  project_id  = "P2abcdefghijklmnop"
  name        = "office-ips"
  description = ""
  ips         = ["10.0.0.0/8"]
}

resource "descope_magiclink_settings" "magiclink_settings" {
  project_id = "P2abcdefghijklmnop"
  disabled   = false
  email_service = {
    connector_id = "C-smtp"
  }
  email_template_id = "T-ml-email"
  expiration_time   = "3 minutes"
  redirect_url      = "https://app.acme.test/ml"
}

resource "descope_oauth_provider" "oauth_provider" {
  project_id             = "P2abcdefghijklmnop"
  id                     = "google"
  allowed_grant_types    = []
  callback_domain        = ""
  client_id              = "google-client"
  client_secret          = var.descope_oauth_provider_oauth_provider_client_secret_eb34aef53128
  disable_jit_updates    = false
  disabled               = false
  manage_provider_tokens = false
  merge_user_accounts    = true
  native_client_id       = ""
  prompts                = []
  redirect_url           = ""
  scopes                 = ["email"]
}

resource "descope_oauth_settings" "oauth_settings" {
  project_id = "P2abcdefghijklmnop"
  disabled   = false
}

resource "descope_oidc_app" "acme_portal" {
  project_id                          = "P2abcdefghijklmnop"
  name                                = "Acme Portal"
  approved_redirect_urls              = ["https://app.acme.test/cb"]
  authorization_code_disabled         = false
  backchannel_logout_url              = ""
  claims                              = ["email"]
  client_credentials_disabled         = false
  client_id                           = "portal-client"
  client_secret                       = var.descope_oidc_app_acme_portal_client_secret_a46e64d65d8d
  client_type                         = "confidential"
  custom_idp_initiated_login_page_url = ""
  default_audience                    = ""
  description                         = ""
  device_code_disabled                = false
  disabled                            = false
  force_authentication                = false
  force_pkce                          = false
  jwt_bearer_disabled                 = false
  login_page_url                      = "https://app.acme.test/login"
  logo                                = ""
  refresh_token_disabled              = false
  trusted_apps_audience               = ""
}

resource "descope_otp_settings" "otp_settings" {
  project_id = "P2abcdefghijklmnop"
  disabled   = false
  domain     = ""
  email_service = {
    connector_id = "C-smtp"
  }
  email_template_id = "T-otp-email"
  expiration_time   = "3 minutes"
  text_service = {
    connector_id = "C-twilio"
  }
  text_template_id = "T-otp-text"
  voice_service = {
    connector_id = "C-twilio"
  }
  voice_template_id = "T-otp-voice"
}

resource "descope_passkey_settings" "passkey_settings" {
  project_id       = "P2abcdefghijklmnop"
  disabled         = false
  top_level_domain = "acme.test"
}

resource "descope_password_settings" "password_settings" {
  project_id       = "P2abcdefghijklmnop"
  disabled         = false
  expiration       = false
  expiration_weeks = 0
  lock             = false
  lock_attempts    = 5
  lowercase        = true
  min_length       = 8
  non_alphanumeric = true
  number           = true
  reuse            = false
  reuse_amount     = 3
  uppercase        = true
}

resource "descope_project_settings" "project_settings" {
  project_id                          = "P2abcdefghijklmnop"
  allow_auth_hosting_iframe_embedding = false
  app_url                             = "https://app.acme.test"
  approved_domains                    = ["acme.test"]
  custom_domain                       = ""
  default_no_sso_apps                 = false
  tenant_user_isolation               = false
  test_users_loginid_regexp           = ""
  test_users_static_otp               = ""
  test_users_verifier_regexp          = ""
}

resource "descope_role" "admin" {
  project_id  = "P2abcdefghijklmnop"
  name        = "Admin"
  default     = false
  description = "admins"
  permissions = ["Manage Users"]
  private     = false
}

resource "descope_sendgrid_connector" "acme_sendgrid" {
  project_id   = "P2abcdefghijklmnop"
  name         = "acme-sendgrid"
  api_key      = var.descope_sendgrid_connector_acme_sendgrid_api_key_24cf138f3fd8
  description  = ""
  sender_email = "mail@acme.test"
  sender_name  = "Acme Mail"
}

resource "descope_session_migration" "session_migration" {
  project_id                 = "P2abcdefghijklmnop"
  api_token                  = var.descope_session_migration_session_migration_api_token_2d7f85e28b43
  audience                   = "aud"
  client_id                  = "client"
  domain                     = "https://tenant.example"
  issuer                     = "https://tenant.example/"
  loginid_matched_attributes = ["email"]
  user_mapping = [
    {
      descope_key  = "loginId"
      external_key = "sub"
    },
  ]
  user_sync_type = "jit"
  vendor         = "auth0"
}

resource "descope_session_settings" "session_settings" {
  project_id                          = "P2abcdefghijklmnop"
  access_key_session_token_expiration = "10 minutes"
  enable_inactivity                   = false
  inactivity_time                     = "12 minutes"
  refresh_token_cookie_domain         = ""
  refresh_token_cookie_policy         = "none"
  refresh_token_expiration            = "4 weeks"
  refresh_token_response_method       = "response_body"
  refresh_token_rotation              = true
  session_token_cookie_domain         = ""
  session_token_cookie_policy         = "none"
  session_token_expiration            = "10 minutes"
  session_token_response_method       = "response_body"
  step_up_token_expiration            = "10 minutes"
  trusted_device_token_expiration     = "365 days"
  user_jwt_template                   = "JT-user"
}

resource "descope_smtp_connector" "acme_smtp" {
  project_id     = "P2abcdefghijklmnop"
  name           = "acme-smtp"
  auth_method    = "plain"
  description    = ""
  host           = "smtp.acme.test"
  password       = var.descope_smtp_connector_acme_smtp_password_812244e2f589
  port           = 587
  sender_email   = "no-reply@acme.test"
  sender_name    = "Acme"
  use_static_ips = false
  username       = "acme"
}

resource "descope_styles" "styles" {
  project_id = "P2abcdefghijklmnop"
  data       = file("${path.module}/payloads/descope_styles.styles.data.json")
}

resource "descope_tenant_attribute" "tier" {
  project_id     = "P2abcdefghijklmnop"
  id             = "tier"
  name           = "Tier"
  select_options = []
  type           = "string"
}

resource "descope_text_template" "otp_text" {
  project_id = "P2abcdefghijklmnop"
  method     = "otp"
  name       = "otp text"
  body       = "code {{code}}"
}

resource "descope_totp_settings" "totp_settings" {
  project_id = "P2abcdefghijklmnop"
  disabled   = false
}

resource "descope_twilio_core_connector" "acme_twilio" {
  project_id            = "P2abcdefghijklmnop"
  name                  = "acme-twilio"
  account_sid           = "AC123"
  auth_token            = var.descope_twilio_core_connector_acme_twilio_auth_token_5d40f97eaa3d
  description           = ""
  from_phone            = "+15551230000"
  from_phone_voice      = "+15551230001"
  messaging_service_sid = ""
}

resource "descope_user_attribute" "nickname" {
  project_id     = "P2abcdefghijklmnop"
  id             = "nickname"
  name           = "Nickname"
  select_options = []
  type           = "string"
  widget_authorization = {
    edit_permissions = []
    view_permissions = []
  }
}

resource "descope_voice_template" "otp_voice" {
  project_id = "P2abcdefghijklmnop"
  method     = "otp"
  name       = "otp voice"
  body       = "your code is {{code}}"
}

locals {
  descope_widget_widgets = {
    user-management = "payloads/descope_widget.widgets.user_management_f18711d7c220.data.json"
  }
}

resource "descope_widget" "widgets" {
  for_each   = local.descope_widget_widgets
  project_id = "P2abcdefghijklmnop"
  widget_id  = each.key
  data       = file("${path.module}/${each.value}")
}

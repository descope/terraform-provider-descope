resource "descope_magiclink_settings" "main" {
  project_id         = descope_project.test.id
  expiration_time    = "6 minutes"
  email_template_id  = descope_email_template.sink.id
  email_connector_id = descope_smtp_connector.sink.id
}

resource "descope_otp_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "6 minutes"
}

resource "descope_password_settings" "main" {
  project_id = descope_project.test.id
  min_length = 9
}

resource "descope_session_settings" "main" {
  project_id               = descope_project.test.id
  session_token_expiration = "25 minutes"
}

resource "descope_invite_settings" "main" {
  project_id        = descope_project.test.id
  invite_expiration = "4 weeks"
}

resource "descope_project_settings" "main" {
  project_id       = descope_project.test.id
  approved_domains = ["sink.example.com"]
}

resource "descope_totp_settings" "main" {
  project_id    = descope_project.test.id
  service_label = "Kitchen Sink"
}

resource "descope_enchantedlink_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "8 minutes"
}

resource "descope_smtp_connector" "sink" {
  project_id   = descope_project.test.id
  name         = "sink-smtp"
  host         = "smtp.example.com"
  port         = 587
  username     = "mailer"
  password     = "not-a-real-password"
  sender_email = "auth@example.com"
}

resource "descope_email_template" "sink" {
  project_id = descope_project.test.id
  method     = "magiclink"
  name       = "sink-template"
  subject    = "Sign in to {{.projectName}}"
  html_body  = "Follow the link to sign in"
}

resource "descope_permission" "sink" {
  project_id = descope_project.test.id
  name       = "sink.use"
}

resource "descope_role" "sink" {
  project_id  = descope_project.test.id
  name        = "sink-role"
  permissions = [descope_permission.sink.name]
}

resource "descope_oidc_app" "sink" {
  project_id          = descope_project.test.id
  name                = "sink-portal"
  deletion_protection = false
}

resource "descope_app_permission" "sink" {
  project_id = descope_project.test.id
  app_id     = descope_oidc_app.sink.id
  name       = "sink.app.use"
}

resource "descope_app_role" "sink" {
  project_id     = descope_project.test.id
  app_id         = descope_oidc_app.sink.id
  name           = "sink-app-role"
  permission_ids = [descope_app_permission.sink.id]
  role_mappings  = [descope_role.sink.id]
}

resource "descope_jwt_template" "sink" {
  project_id = descope_project.test.id
  type       = "user"
  name       = "sink-jwt"
  template   = jsonencode({ sink = true })
}

resource "descope_user_attribute" "sink" {
  project_id = descope_project.test.id
  id         = "sinkAttr"
  name       = "Sink Attribute"
  type       = "string"
}

resource "descope_list" "sink" {
  project_id = descope_project.test.id
  name       = "sink-list"
  texts      = ["a", "b"]
}

resource "descope_inbound_app" "sink" {
  project_id             = descope_project.test.id
  name                   = "sink-partner"
  approved_callback_urls = ["https://partner.example.com/cb"]
  deletion_protection    = false
}

resource "descope_oauth_provider" "sink" {
  project_id             = descope_project.test.id
  id                     = "sink-custom"
  client_id              = "sink-client"
  client_secret          = "not-a-real-secret"
  allowed_grant_types    = ["authorization_code"]
  authorization_endpoint = "https://idp.example.com/authorize"
  token_endpoint         = "https://idp.example.com/token"
  user_info_endpoint     = "https://idp.example.com/userinfo"
}

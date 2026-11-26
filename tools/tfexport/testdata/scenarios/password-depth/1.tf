resource "descope_smtp_connector" "reset" {
  project_id   = descope_project.test.id
  name         = "reset-smtp"
  host         = "smtp.example.com"
  port         = 587
  username     = "mailer"
  password     = "not-a-real-password"
  sender_email = "reset@example.com"
}

resource "descope_email_template" "reset" {
  project_id = descope_project.test.id
  method     = "password"
  name       = "reset-template"
  subject    = "Reset your password"
  html_body  = "Follow the link to reset your password"
}

resource "descope_password_settings" "main" {
  project_id         = descope_project.test.id
  min_length         = 12
  uppercase          = true
  lowercase          = true
  number             = true
  non_alphanumeric   = true
  expiration         = true
  expiration_weeks   = 13
  lock               = true
  lock_attempts      = 7
  reuse              = true
  reuse_amount       = 6
  email_template_id  = descope_email_template.reset.id
  email_connector_id = descope_smtp_connector.reset.id
}

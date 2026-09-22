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
  subject    = "Reset your password (updated)"
  html_body  = "Follow the new link to reset your password"
}

resource "descope_password_settings" "main" {
  project_id         = descope_project.test.id
  min_length         = 14
  uppercase          = true
  expiration         = false
  lock               = true
  lock_attempts      = 5
  reuse              = false
  email_template_id  = descope_email_template.reset.id
  email_connector_id = descope_smtp_connector.reset.id
}

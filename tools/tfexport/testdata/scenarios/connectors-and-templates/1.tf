resource "descope_smtp_connector" "email" {
  project_id   = descope_project.test.id
  name         = "roundtrip-smtp"
  host         = "smtp.example.com"
  port         = 587
  username     = "mailer"
  password     = "not-a-real-password"
  sender_email = "auth@example.com"
  sender_name  = "Example Auth"
}

resource "descope_email_template" "magiclink" {
  project_id = descope_project.test.id
  method     = "magiclink"
  name       = "roundtrip-template"
  subject    = "Sign in to {{.projectName}}"
  html_body  = "Follow the link in this email to sign in"
}

resource "descope_magiclink_settings" "main" {
  project_id         = descope_project.test.id
  email_template_id  = descope_email_template.magiclink.id
  email_connector_id = descope_smtp_connector.email.id
}

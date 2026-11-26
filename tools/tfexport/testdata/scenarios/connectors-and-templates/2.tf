resource "descope_smtp_connector" "email" {
  project_id   = descope_project.test.id
  name         = "roundtrip-smtp-renamed"
  host         = "smtp2.example.com"
  port         = 465
  username     = "mailer2"
  password     = "not-a-real-password"
  sender_email = "auth@example.com"
}

resource "descope_email_template" "magiclink" {
  project_id = descope_project.test.id
  method     = "magiclink"
  name       = "roundtrip-template"
  subject    = "Sign in to {{.projectName}} (updated)"
  html_body  = "Follow the updated link in this email to sign in"
}

resource "descope_magiclink_settings" "main" {
  project_id         = descope_project.test.id
  email_template_id  = descope_email_template.magiclink.id
  email_connector_id = descope_smtp_connector.email.id
}

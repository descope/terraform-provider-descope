resource "descope_smtp_connector" "mail" {
  project_id   = descope_project.test.id
  name         = "mail-beta"
  host         = "smtp.example.com"
  port         = 587
  username     = "mailer"
  password     = "not-a-real-password"
  sender_email = "auth@example.com"
}

resource "descope_magiclink_settings" "main" {
  project_id         = descope_project.test.id
  email_connector_id = descope_smtp_connector.mail.id
}

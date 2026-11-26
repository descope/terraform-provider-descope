resource "descope_email_template" "enchanted" {
  project_id = descope_project.test.id
  method     = "enchantedlink"
  name       = "matrix-enchanted"
  subject    = "Enchanted sign in"
  html_body  = "Follow the enchanted link"
}

resource "descope_email_template" "magic" {
  project_id = descope_project.test.id
  method     = "magiclink"
  name       = "matrix-magic"
  subject    = "Magic sign in"
  html_body  = "Follow the magic link"
}

resource "descope_email_template" "otp" {
  project_id = descope_project.test.id
  method     = "otp"
  name       = "matrix-otp"
  subject    = "Your code"
  html_body  = "Use the code to sign in"
}

resource "descope_email_template" "password" {
  project_id = descope_project.test.id
  method     = "password"
  name       = "matrix-password"
  subject    = "Reset password"
  html_body  = "Follow the link to reset"
}

resource "descope_email_template" "sso" {
  project_id = descope_project.test.id
  method     = "sso"
  name       = "matrix-sso"
  subject    = "SSO verification"
  html_body  = "Follow the link to verify"
}

resource "descope_text_template" "magic" {
  project_id = descope_project.test.id
  method     = "magiclink"
  name       = "matrix-text-magic"
  body       = "Tap to sign in"
}

resource "descope_text_template" "otp" {
  project_id = descope_project.test.id
  method     = "otp"
  name       = "matrix-text-otp"
  body       = "Code: {{.code}}"
}

resource "descope_voice_template" "otp" {
  project_id = descope_project.test.id
  method     = "otp"
  name       = "matrix-voice-otp"
  body       = "Your code is {{.code}}"
}

resource "descope_magiclink_settings" "main" {
  project_id        = descope_project.test.id
  email_template_id = descope_email_template.magic.id
  text_template_id  = descope_text_template.magic.id
}

resource "descope_otp_settings" "main" {
  project_id        = descope_project.test.id
  email_template_id = descope_email_template.otp.id
  text_template_id  = descope_text_template.otp.id
  voice_template_id = descope_voice_template.otp.id
}

resource "descope_enchantedlink_settings" "main" {
  project_id        = descope_project.test.id
  email_template_id = descope_email_template.enchanted.id
}

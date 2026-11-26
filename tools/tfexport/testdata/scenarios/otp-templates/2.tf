resource "descope_text_template" "otp" {
  project_id = descope_project.test.id
  method     = "otp"
  name       = "roundtrip-text"
  body       = "Use code {{.code}} to sign in"
}

resource "descope_voice_template" "otp" {
  project_id = descope_project.test.id
  method     = "otp"
  name       = "roundtrip-voice-renamed"
  body       = "Your code is {{.code}}"
}

resource "descope_otp_settings" "main" {
  project_id        = descope_project.test.id
  expiration_time   = "10 minutes"
  text_template_id  = descope_text_template.otp.id
  voice_template_id = descope_voice_template.otp.id
}

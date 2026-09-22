resource "descope_embeddedlink_settings" "main" {
  project_id      = descope_project.test.id
  disabled        = false
  expiration_time = "9 minutes"
}

resource "descope_oauth_settings" "main" {
  project_id = descope_project.test.id
  disabled   = false
}

resource "descope_otp_settings" "main" {
  project_id = descope_project.test.id
  disabled   = false
}

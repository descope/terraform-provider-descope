resource "descope_embeddedlink_settings" "main" {
  project_id      = descope_project.test.id
  disabled        = true
  expiration_time = "4 minutes"
}

resource "descope_oauth_settings" "main" {
  project_id = descope_project.test.id
  disabled   = true
}

resource "descope_otp_settings" "main" {
  project_id = descope_project.test.id
  disabled   = true
}

resource "descope_magiclink_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "10 minutes"
}

resource "descope_password_settings" "main" {
  project_id = descope_project.test.id
  min_length = 12
  lowercase  = true
}

resource "descope_session_settings" "main" {
  project_id               = descope_project.test.id
  session_token_expiration = "20 minutes"
}

resource "descope_otp_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "5 minutes"
}

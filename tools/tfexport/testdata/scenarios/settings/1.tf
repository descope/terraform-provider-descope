resource "descope_magiclink_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "5 minutes"
  redirect_url    = "https://example.com/magiclink"
}

resource "descope_password_settings" "main" {
  project_id       = descope_project.test.id
  min_length       = 10
  uppercase        = true
  non_alphanumeric = true
}

resource "descope_session_settings" "main" {
  project_id               = descope_project.test.id
  session_token_expiration = "15 minutes"
}

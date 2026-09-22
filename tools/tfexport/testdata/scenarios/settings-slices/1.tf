resource "descope_invite_settings" "main" {
  project_id         = descope_project.test.id
  invite_expiration  = "2 weeks"
  invite_url         = "https://example.com/invite"
  require_invitation = true
}

resource "descope_project_settings" "main" {
  project_id       = descope_project.test.id
  app_url          = "https://app.example.com"
  approved_domains = ["example.com", "example.org"]
}

resource "descope_session_settings" "main" {
  project_id               = descope_project.test.id
  session_token_expiration = "12 minutes"
  refresh_token_expiration = "3 weeks"
  refresh_token_rotation   = true
}

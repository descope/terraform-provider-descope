resource "descope_invite_settings" "main" {
  project_id        = descope_project.test.id
  invite_expiration = "3 weeks"
}

resource "descope_session_migration" "main" {
  project_id                 = descope_project.test.id
  vendor                     = "okta"
  client_id                  = "0oa1b2c3d4e5f6g7h8i9"
  issuer                     = "https://dev-123456.okta.com/oauth2/default"
  api_token                  = "not-a-real-token"
  loginid_matched_attributes = ["email"]
}

resource "descope_oauth_provider" "custom" {
  project_id             = descope_project.test.id
  id                     = "roundtrip-custom"
  description            = "Custom OIDC provider"
  client_id              = "roundtrip-client"
  client_secret          = "not-a-real-secret"
  allowed_grant_types    = ["authorization_code"]
  authorization_endpoint = "https://example.com/idp/authorize"
  token_endpoint         = "https://example.com/idp/token"
  user_info_endpoint     = "https://example.com/idp/userinfo"
  scopes                 = ["openid", "email"]
}

resource "descope_oauth_provider" "google" {
  project_id          = descope_project.test.id
  id                  = "google"
  client_id           = "roundtrip-google-client"
  client_secret       = "not-a-real-secret"
  allowed_grant_types = ["authorization_code"]
}

resource "descope_oauth_provider" "custom" {
  project_id             = descope_project.test.id
  id                     = "roundtrip-custom"
  description            = "Custom OIDC provider (updated)"
  client_id              = "roundtrip-client-v2"
  client_secret          = "not-a-real-secret"
  allowed_grant_types    = ["authorization_code", "implicit"]
  authorization_endpoint = "https://example.com/idp2/authorize"
  token_endpoint         = "https://example.com/idp2/token"
  user_info_endpoint     = "https://example.com/idp2/userinfo"
  scopes                 = ["openid", "email", "profile"]
}

resource "descope_oauth_provider" "google" {
  project_id          = descope_project.test.id
  id                  = "google"
  client_id           = "roundtrip-google-client-v2"
  client_secret       = "not-a-real-secret"
  allowed_grant_types = ["authorization_code", "implicit"]
  prompts             = ["select_account"]
}

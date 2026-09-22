resource "descope_oauth_provider" "custom" {
  project_id             = descope_project.test.id
  id                     = "roundtrip-custom"
  description            = "Custom OIDC provider (updated)"
  client_id              = "roundtrip-client-v2"
  client_secret          = "not-a-real-secret"
  allowed_grant_types    = ["authorization_code", "implicit"]
  authorization_endpoint = "https://idp2.example.com/authorize"
  token_endpoint         = "https://idp2.example.com/token"
  user_info_endpoint     = "https://idp2.example.com/userinfo"
  scopes                 = ["openid", "email", "profile"]
}

resource "descope_oauth_provider" "custom" {
  project_id             = descope_project.test.id
  id                     = "roundtrip-custom"
  description            = "Custom OIDC provider"
  client_id              = "roundtrip-client"
  client_secret          = "not-a-real-secret"
  allowed_grant_types    = ["authorization_code"]
  authorization_endpoint = "https://idp.example.com/authorize"
  token_endpoint         = "https://idp.example.com/token"
  user_info_endpoint     = "https://idp.example.com/userinfo"
  scopes                 = ["openid", "email"]
}

resource "descope_sso_settings" "main" {
  project_id              = descope_project.test.id
  allow_override_roles    = true
  allow_duplicate_domains = true
}

resource "descope_enchantedlink_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "15 minutes"
}

resource "descope_totp_settings" "main" {
  project_id    = descope_project.test.id
  service_label = "Roundtrip Service v2"
}

resource "descope_invite_settings" "main" {
  project_id        = descope_project.test.id
  invite_expiration = "10 days"
}

resource "descope_embeddedlink_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "5 minutes"
}

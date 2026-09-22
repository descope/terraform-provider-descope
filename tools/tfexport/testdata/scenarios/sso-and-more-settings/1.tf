resource "descope_sso_settings" "main" {
  project_id           = descope_project.test.id
  allow_override_roles = true
  groups_priority      = true
}

resource "descope_enchantedlink_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "10 minutes"
  redirect_url    = "https://example.com/enchanted"
}

resource "descope_totp_settings" "main" {
  project_id    = descope_project.test.id
  service_label = "Roundtrip Service"
}

resource "descope_invite_settings" "main" {
  project_id        = descope_project.test.id
  invite_expiration = "2 weeks"
}

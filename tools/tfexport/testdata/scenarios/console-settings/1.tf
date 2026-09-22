resource "descope_magiclink_settings" "main" {
  project_id      = descope_project.test.id
  expiration_time = "5 minutes"
}

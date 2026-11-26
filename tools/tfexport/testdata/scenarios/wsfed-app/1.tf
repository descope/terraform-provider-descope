resource "descope_wsfed_app" "legacy" {
  project_id          = descope_project.test.id
  name                = "wsfed-legacy"
  description         = "WS-Federation app"
  deletion_protection = false
}

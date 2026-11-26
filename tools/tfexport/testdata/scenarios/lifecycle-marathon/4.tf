resource "descope_role" "base" {
  project_id = descope_project.test.id
  name       = "marathon-base-renamed"
}

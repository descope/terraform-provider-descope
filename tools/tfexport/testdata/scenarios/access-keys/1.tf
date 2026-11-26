resource "descope_role" "automation" {
  project_id = descope_project.test.id
  name       = "roundtrip-automation"
}

resource "descope_access_key" "ci" {
  project_id = descope_project.test.id
  name       = "roundtrip-ci"
  roles      = [descope_role.automation.name]
}

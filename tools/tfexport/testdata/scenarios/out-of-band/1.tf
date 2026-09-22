resource "descope_role" "seeded" {
  project_id  = descope_project.test.id
  name        = "roundtrip-seeded"
  description = "Created by terraform"
}

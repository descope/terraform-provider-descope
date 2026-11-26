resource "descope_role" "keeper" {
  project_id = descope_project.test.id
  name       = "roundtrip-keeper"
}

resource "descope_jwt_template" "keeper" {
  project_id = descope_project.test.id
  type       = "user"
  name       = "roundtrip-keeper-jwt"
  template   = jsonencode({ kept = true })
}

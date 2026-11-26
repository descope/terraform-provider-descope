resource "descope_role" "keeper" {
  project_id = descope_project.test.id
  name       = "roundtrip-keeper"
}

resource "descope_role" "doomed" {
  project_id = descope_project.test.id
  name       = "roundtrip-doomed"
}

resource "descope_permission" "temp" {
  project_id = descope_project.test.id
  name       = "roundtrip.temp"
}

resource "descope_list" "temp" {
  project_id = descope_project.test.id
  name       = "roundtrip-temp"
  texts      = ["one", "two"]
}

resource "descope_jwt_template" "keeper" {
  project_id = descope_project.test.id
  type       = "user"
  name       = "roundtrip-keeper-jwt"
  template   = jsonencode({ kept = true })
}

resource "descope_permission" "base" {
  project_id = descope_project.test.id
  name       = "marathon.base"
}

resource "descope_role" "base" {
  project_id  = descope_project.test.id
  name        = "marathon-base-renamed"
  permissions = [descope_permission.base.name]
}

resource "descope_list" "base" {
  project_id = descope_project.test.id
  name       = "marathon-list"
  texts      = ["one", "two"]
}

resource "descope_jwt_template" "added" {
  project_id = descope_project.test.id
  type       = "user"
  name       = "marathon-jwt"
  template   = jsonencode({ marathon = 1 })
}

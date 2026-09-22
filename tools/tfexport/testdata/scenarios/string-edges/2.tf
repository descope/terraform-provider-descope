resource "descope_role" "unicode" {
  project_id  = descope_project.test.id
  name        = "rôle-בדיקה-役割"
  description = "changed — still ünïcödé 🎯"
}

resource "descope_role" "interp" {
  project_id  = descope_project.test.id
  name        = "edge case"
  description = "now with %%{ other } and $${vars} reversed"
}

resource "descope_role" "collide" {
  project_id = descope_project.test.id
  name       = "edge-case"
}

resource "descope_permission" "weird" {
  project_id  = descope_project.test.id
  name        = "edge.perm"
  description = "line1\nline2 changed"
}

resource "descope_list" "unicode" {
  project_id = descope_project.test.id
  name       = "edge-json"
  json = jsonencode({
    "ключ" = "новое значение"
    emoji  = "🚀"
    num    = 2.25
  })
}

resource "descope_email_template" "edge" {
  project_id = descope_project.test.id
  method     = "magiclink"
  name       = "edge-template"
  subject    = "Hi {{.projectName}} — now $${10} 🎯"
  html_body  = "<a href=\"https://example.com?a=1&b=2\">link $${literal}</a>\nline two changed"
}

resource "descope_jwt_template" "edge" {
  project_id = descope_project.test.id
  type       = "user"
  name       = "edge-jwt"
  template = jsonencode({
    "weird key" = "new va\"lue"
    dollar      = "$${yep}"
  })
}

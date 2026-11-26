resource "descope_role" "unicode" {
  project_id  = descope_project.test.id
  name        = "rôle-בדיקה-役割"
  description = "Ünïcödé — emoji 🎯 and \"double\" and 'single' and \\backslash\\"
}

resource "descope_role" "interp" {
  project_id  = descope_project.test.id
  name        = "edge case"
  description = "literal $${not_a_var} and %%{ not_a_directive } and trailing dollar $"
}

resource "descope_role" "collide" {
  project_id  = descope_project.test.id
  name        = "edge-case"
  description = "sanitizes to the same label as 'edge case'"
}

resource "descope_permission" "weird" {
  project_id  = descope_project.test.id
  name        = "edge.perm"
  description = "line1\nline2\ttabbed <html> &amp; entities"
}

resource "descope_list" "unicode" {
  project_id = descope_project.test.id
  name       = "edge-json"
  json = jsonencode({
    "ключ"  = "значение"
    emoji   = "🎯🚀"
    num     = 1.5
    neg     = -3
    big     = 9007199254740991
    escaped = "va\"lue with $${literal} and \\slash"
    nested  = { arr = ["a", "\"b\"", "c\nd"] }
  })
}

resource "descope_email_template" "edge" {
  project_id = descope_project.test.id
  method     = "magiclink"
  name       = "edge-template"
  subject    = "Hi {{.projectName}} — cost is $${5} & <b>bold</b> 🎯"
  html_body  = "<a href=\"https://example.com?a=1&b=2\">link $${literal}</a>\nline two — ünïcödé"
}

resource "descope_jwt_template" "edge" {
  project_id = descope_project.test.id
  type       = "user"
  name       = "edge-jwt"
  template = jsonencode({
    "weird key" = "va\"lue"
    path        = "C:\\dir\\file"
    tpl         = "{{user.email}}"
    dollar      = "$${nope}"
    unicode     = "héllo 🎯"
  })
}

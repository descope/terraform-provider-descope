resource "descope_permission" "old" {
  project_id = descope_project.test.id
  name       = "roundtrip.old"
}

resource "descope_role" "renamer" {
  project_id  = descope_project.test.id
  name        = "roundtrip-renamer"
  permissions = [descope_permission.old.name]
}

resource "descope_email_template" "dangle" {
  project_id = descope_project.test.id
  method     = "magiclink"
  name       = "dangle-template"
  subject    = "Sign in"
  html_body  = "Follow the link to sign in"
}

resource "descope_magiclink_settings" "main" {
  project_id        = descope_project.test.id
  email_template_id = descope_email_template.dangle.id
}

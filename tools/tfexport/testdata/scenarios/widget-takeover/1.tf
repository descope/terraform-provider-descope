resource "descope_widget" "profile" {
  project_id = descope_project.test.id
  widget_id  = "user-profile"
  data       = file("${path.module}/roundtrip-widget.json")
}

resource "descope_admin_portal" "main" {
  project_id = descope_project.test.id
  enabled    = true
  widgets = [
    { widget_id = descope_widget.profile.widget_id, type = "user-profile" },
  ]
}

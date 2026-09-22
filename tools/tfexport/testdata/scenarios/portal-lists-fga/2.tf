resource "descope_admin_portal" "main" {
  project_id = descope_project.test.id
  enabled    = true
  widgets = [
    { widget_id = "user-management", type = "user-management" },
    { widget_id = "audit-management", type = "audit-management" },
  ]
}

resource "descope_list" "blocked" {
  project_id = descope_project.test.id
  name       = "roundtrip-blocked-ips"
  ips        = ["10.0.0.1", "192.168.0.0/24", "172.16.0.0/12"]
}

resource "descope_list" "config" {
  project_id = descope_project.test.id
  name       = "roundtrip-config"
  json       = jsonencode({ feature = true, tier = 3 })
}

resource "descope_fga_schema" "main" {
  project_id = descope_project.test.id
  schema     = "model AuthZ 1.0\n\ntype user\n\ntype doc\n  relation owner: user\n  relation editor: user\n  permission view: owner\n  permission edit: editor\n"
}

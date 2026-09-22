resource "descope_permission" "read" {
  project_id  = descope_project.test.id
  name        = "roundtrip.read"
  description = "Read access"
}

resource "descope_permission" "write" {
  project_id  = descope_project.test.id
  name        = "roundtrip.write"
  description = "Write access (updated)"
}

resource "descope_role" "viewer" {
  project_id  = descope_project.test.id
  name        = "roundtrip-viewer"
  description = "Read only"
  permissions = [descope_permission.read.name]
}

resource "descope_role" "editor" {
  project_id  = descope_project.test.id
  name        = "roundtrip-editor"
  permissions = [descope_permission.read.name, descope_permission.write.name]
}

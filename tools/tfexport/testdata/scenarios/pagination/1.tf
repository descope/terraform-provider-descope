resource "descope_permission" "bulk" {
  count      = 25
  project_id = descope_project.test.id
  name       = "bulk.perm.${count.index}"
}

resource "descope_role" "bulk" {
  count       = 35
  project_id  = descope_project.test.id
  name        = "bulk-role-${count.index}"
  permissions = [descope_permission.bulk[count.index % 25].name]
}

resource "descope_list" "bulk" {
  count      = 10
  project_id = descope_project.test.id
  name       = "bulk-list-${count.index}"
  texts      = ["item-${count.index}"]
}

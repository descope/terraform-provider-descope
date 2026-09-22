resource "descope_permission" "view" {
  project_id = descope_project.test.id
  name       = "attr.view"
}

resource "descope_user_attribute" "text" {
  project_id = descope_project.test.id
  id         = "fuzzText"
  name       = "Text Attribute (renamed)"
  type       = "string"
  widget_authorization = {
    view_permissions = [descope_permission.view.name]
    edit_permissions = [descope_permission.view.name]
  }
}

resource "descope_user_attribute" "select" {
  depends_on     = [descope_user_attribute.text]
  project_id     = descope_project.test.id
  id             = "fuzzSelect"
  name           = "Select Attribute"
  type           = "singleselect"
  select_options = ["red", "green", "blue", "yellow"]
}

resource "descope_user_attribute" "multi" {
  depends_on     = [descope_user_attribute.select]
  project_id     = descope_project.test.id
  id             = "fuzzMulti"
  name           = "Multi Attribute"
  type           = "multiselect"
  select_options = ["a", "b"]
}

resource "descope_user_attribute" "flag" {
  depends_on = [descope_user_attribute.multi]
  project_id = descope_project.test.id
  id         = "fuzzFlag"
  name       = "Flag Attribute"
  type       = "boolean"
}

resource "descope_user_attribute" "when" {
  depends_on = [descope_user_attribute.flag]
  project_id = descope_project.test.id
  id         = "fuzzWhen"
  name       = "Date Attribute"
  type       = "date"
}

resource "descope_tenant_attribute" "tier" {
  depends_on = [descope_user_attribute.when]
  project_id = descope_project.test.id
  id         = "fuzzTier"
  name       = "Tier Attribute"
  type       = "number"
}

resource "descope_access_key_attribute" "owner" {
  depends_on = [descope_tenant_attribute.tier]
  project_id = descope_project.test.id
  id         = "fuzzOwner"
  name       = "Owner Attribute"
  type       = "string"
}

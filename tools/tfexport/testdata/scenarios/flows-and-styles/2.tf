resource "descope_flow" "passkeys" {
  project_id = descope_project.test.id
  flow_id    = "add-passkeys"
  data       = file("${path.module}/roundtrip-flow.json")
}

resource "descope_styles" "main" {
  project_id = descope_project.test.id
  data       = file("${path.module}/roundtrip-styles-v2.json")
}

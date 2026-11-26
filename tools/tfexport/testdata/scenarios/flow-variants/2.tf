resource "descope_flow" "signin" {
  project_id = descope_project.test.id
  flow_id    = "sign-in"
  data       = file("${path.module}/roundtrip-flow-signin.json")
}

resource "descope_flow" "copy" {
  project_id = descope_project.test.id
  flow_id    = "passkeys-copy"
  data       = file("${path.module}/roundtrip-flow.json")
}

resource "descope_user_attribute" "plan" {
  project_id = descope_project.test.id
  id         = "roundtripPlan"
  name       = "Roundtrip Plan"
  type       = "string"
}

resource "descope_jwt_template" "access" {
  project_id  = descope_project.test.id
  type        = "user"
  name        = "roundtrip-jwt"
  description = "Roundtrip template"
  template = jsonencode({
    plan = "{{user.customAttributes.roundtripPlan}}"
  })
}

resource "descope_user_attribute" "plan" {
  project_id = descope_project.test.id
  id         = "roundtripPlan"
  name       = "Roundtrip Plan Level"
  type       = "string"
}

resource "descope_tenant_attribute" "tier" {
  project_id = descope_project.test.id
  id         = "roundtripTier"
  name       = "Roundtrip Tier"
  type       = "number"
}

resource "descope_jwt_template" "access" {
  project_id  = descope_project.test.id
  type        = "user"
  name        = "roundtrip-jwt"
  description = "Roundtrip template (updated)"
  template = jsonencode({
    plan = "{{user.customAttributes.roundtripPlan}}"
    v    = 2
  })
}

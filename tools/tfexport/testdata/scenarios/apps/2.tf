resource "descope_role" "member" {
  project_id = descope_project.test.id
  name       = "roundtrip-member"
}

resource "descope_oidc_app" "portal" {
  project_id          = descope_project.test.id
  name                = "roundtrip-portal-v2"
  description         = "Customer portal (renamed)"
  login_page_url      = "https://example.com/login"
  deletion_protection = false
}

resource "descope_saml_app" "legacy" {
  project_id          = descope_project.test.id
  name                = "roundtrip-legacy"
  login_page_url      = "https://example.com/saml"
  deletion_protection = false
  manual_configuration = {
    acs_url   = "https://legacy2.example.com/acs"
    entity_id = "legacy-entity"
  }
}

resource "descope_app_permission" "reports" {
  project_id  = descope_project.test.id
  app_id      = descope_oidc_app.portal.id
  name        = "reports.view"
  description = "View reports"
}

resource "descope_app_permission" "exports" {
  project_id = descope_project.test.id
  app_id     = descope_oidc_app.portal.id
  name       = "reports.export"
}

resource "descope_app_role" "analyst" {
  project_id     = descope_project.test.id
  app_id         = descope_oidc_app.portal.id
  name           = "analyst"
  permission_ids = [descope_app_permission.reports.id, descope_app_permission.exports.id]
  role_mappings  = [descope_role.member.id]
}

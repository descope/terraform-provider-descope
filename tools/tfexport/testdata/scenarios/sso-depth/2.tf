resource "descope_email_template" "sso" {
  project_id = descope_project.test.id
  method     = "sso"
  name       = "sso-template"
  subject    = "Verify your SSO sign in"
  html_body  = "Follow the link to verify"
}

resource "descope_sso_settings" "main" {
  project_id                    = descope_project.test.id
  merge_users                   = true
  require_sso_domains           = true
  require_groups_attribute_name = true
  email_template_id             = descope_email_template.sso.id
  mandatory_user_attributes = [
    { id = "email" },
  ]
}

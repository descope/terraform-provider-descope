resource "descope_inbound_app_attribute" "managing_org" {
  project_id = descope_project.test.id
  id         = "managingOrg"
  name       = "Managing Organization"
  type       = "string"
}

resource "descope_inbound_app" "partner" {
  project_id             = descope_project.test.id
  name                   = "roundtrip-partner"
  client_type            = "confidential"
  description            = "Partner integration (updated)"
  approved_callback_urls = ["https://partner.example.com/callback", "https://partner.example.com/callback2"]
  force_pkce             = true
  custom_attributes      = jsonencode({ (descope_inbound_app_attribute.managing_org.id) = "org-2" })
  deletion_protection    = false
}

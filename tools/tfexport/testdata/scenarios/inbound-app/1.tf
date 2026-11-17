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
  description            = "Partner integration"
  approved_callback_urls = ["https://partner.example.com/callback"]
  custom_attributes      = jsonencode({ (descope_inbound_app_attribute.managing_org.id) = "org-1" })
  deletion_protection    = false
}

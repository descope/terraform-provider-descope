resource "descope_inbound_app" "partner" {
  project_id             = descope_project.test.id
  name                   = "roundtrip-partner"
  client_type            = "confidential"
  description            = "Partner integration"
  approved_callback_urls = ["https://partner.example.com/callback"]
  deletion_protection    = false
}

resource "descope_inbound_app" "partner" {
  project_id             = descope_project.test.id
  name                   = "roundtrip-partner"
  description            = "Partner integration (updated)"
  approved_callback_urls = ["https://partner.example.com/callback", "https://partner.example.com/callback2"]
  force_pkce             = true
  deletion_protection    = false
}

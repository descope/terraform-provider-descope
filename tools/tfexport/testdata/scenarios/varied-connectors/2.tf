resource "descope_audit_webhook_connector" "audit" {
  project_id  = descope_project.test.id
  name        = "roundtrip-audit"
  description = "Sends audit events"
  base_url    = "https://audit2.example.com"
  authentication = {
    bearer_token = "not-a-real-token"
  }
  headers = {
    "X-Source" = "descope"
    "X-Env"    = "test"
  }
  insecure = true
}

resource "descope_amplitude_connector" "analytics" {
  project_id = descope_project.test.id
  name       = "roundtrip-amplitude"
  api_key    = "not-a-real-key"
  server_url = "https://api.eu.amplitude.com/2/httpapi"
}

resource "descope_audit_webhook_connector" "audit" {
  project_id = descope_project.test.id
  name       = "roundtrip-audit"
  base_url   = "https://audit.example.com"
  authentication = {
    bearer_token = "not-a-real-token"
  }
  headers = {
    "X-Source" = "descope"
  }
}

resource "descope_amplitude_connector" "analytics" {
  project_id = descope_project.test.id
  name       = "roundtrip-amplitude"
  api_key    = "not-a-real-key"
  server_url = "https://api.eu.amplitude.com/2/httpapi"
}

resource "descope_http_connector" "webhook" {
  project_id  = descope_project.test.id
  name        = "roundtrip-http"
  base_url    = "https://hooks.example.com"
  hmac_secret = "not-a-real-secret"
  authentication = {
    bearer_token = "not-a-real-token"
  }
  secret_headers = {
    "X-Api-Key" = "not-a-real-key"
  }
}

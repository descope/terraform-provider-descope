# Stage 2: native import blocks adopting the existing Descope objects.
# Import blocks are evaluated in the root module, so this file belongs in the root module even when the
# resources it adopts are declared in a child module.
# The legacy project resource is imported back here too, after stage 1 detached it from state.

import {
  to = descope_project.main
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_access_key_attribute.scope
  id = "P2abcdefghijklmnop/scope"
}

import {
  to = descope_admin_portal.admin_portal
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_app_permission.read
  id = "P2abcdefghijklmnop/APP-oidc/AP-read"
}

import {
  to = descope_app_role.viewer
  id = "P2abcdefghijklmnop/APP-oidc/AR-viewer"
}

import {
  to = descope_email_template.magiclink_email
  id = "P2abcdefghijklmnop/magiclink/T-ml-email"
}

import {
  to = descope_email_template.otp_email
  id = "P2abcdefghijklmnop/otp/T-otp-email"
}

import {
  to = descope_embeddedlink_settings.embeddedlink_settings
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_fga_schema.fga_schema
  id = "P2abcdefghijklmnop"
}

import {
  for_each = {
    sign-up-or-in = "P2abcdefghijklmnop/sign-up-or-in"
    step-up       = "P2abcdefghijklmnop/step-up"
  }
  to = descope_flow.flows[each.key]
  id = each.value
}

import {
  to = descope_invite_settings.invite_settings
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_jwt_template.user_claims
  id = "P2abcdefghijklmnop/JT-user"
}

import {
  to = descope_list.allowlist
  id = "P2abcdefghijklmnop/L-allow"
}

import {
  to = descope_list.config
  id = "P2abcdefghijklmnop/L-json"
}

import {
  to = descope_list.office_ips
  id = "P2abcdefghijklmnop/L-ips"
}

import {
  to = descope_magiclink_settings.magiclink_settings
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_oauth_provider.oauth_provider
  id = "P2abcdefghijklmnop/google"
}

import {
  to = descope_oauth_settings.oauth_settings
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_oidc_app.acme_portal
  id = "P2abcdefghijklmnop/APP-oidc"
}

import {
  to = descope_otp_settings.otp_settings
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_passkey_settings.passkey_settings
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_password_settings.password_settings
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_project_settings.project_settings
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_role.admin
  id = "P2abcdefghijklmnop/R-admin"
}

import {
  to = descope_sendgrid_connector.acme_sendgrid
  id = "P2abcdefghijklmnop/C-sg"
}

import {
  to = descope_session_migration.session_migration
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_session_settings.session_settings
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_smtp_connector.acme_smtp
  id = "P2abcdefghijklmnop/C-smtp"
}

import {
  to = descope_styles.styles
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_tenant_attribute.tier
  id = "P2abcdefghijklmnop/tier"
}

import {
  to = descope_text_template.otp_text
  id = "P2abcdefghijklmnop/otp/T-otp-text"
}

import {
  to = descope_totp_settings.totp_settings
  id = "P2abcdefghijklmnop"
}

import {
  to = descope_twilio_core_connector.acme_twilio
  id = "P2abcdefghijklmnop/C-twilio"
}

import {
  to = descope_user_attribute.nickname
  id = "P2abcdefghijklmnop/nickname"
}

import {
  to = descope_voice_template.otp_voice
  id = "P2abcdefghijklmnop/otp/T-otp-voice"
}

import {
  for_each = {
    user-management = "P2abcdefghijklmnop/user-management"
  }
  to = descope_widget.widgets[each.key]
  id = each.value
}

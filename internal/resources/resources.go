package resources

import (
	"github.com/descope/terraform-provider-descope/internal/models/accesskey"
	"github.com/descope/terraform-provider-descope/internal/models/descoper"
	"github.com/descope/terraform-provider-descope/internal/models/engine"
	"github.com/descope/terraform-provider-descope/internal/models/inboundapp"
	"github.com/descope/terraform-provider-descope/internal/models/managementkey"
	"github.com/descope/terraform-provider-descope/internal/models/settings"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewAccessKeyResource() resource.Resource {
	return newInfraResource[accesskey.AccessKeyModel]("access_key", accesskey.Schema)
}

func NewDescoperResource() resource.Resource {
	return newInfraResource[descoper.DescoperModel]("descoper", descoper.Schema)
}

func NewManagementKeyResource() resource.Resource {
	return newInfraResource[managementkey.ManagementKeyModel]("management_key", managementkey.Schema)
}

func NewInboundAppResource() resource.Resource {
	return newInfraResource[inboundapp.InboundAppModel]("inbound_app", inboundapp.Schema)
}

func NewEngineResource() resource.Resource {
	return newInfraResource[engine.EngineModel]("engine", engine.Schema)
}

func NewOAuthSettingsResource() resource.Resource {
	return newSettingsResource[settings.OAuthSettingsModel]("oauth_settings", settings.OAuthSettingsSchema, "/v1/mgmt/oauth/settings")
}

func NewPasskeySettingsResource() resource.Resource {
	return newSettingsResource[settings.PasskeySettingsModel]("passkey_settings", settings.PasskeySettingsSchema, "/v1/mgmt/passkey/settings")
}

func NewTOTPSettingsResource() resource.Resource {
	return newSettingsResource[settings.TOTPSettingsModel]("totp_settings", settings.TOTPSettingsSchema, "/v1/mgmt/totp/settings")
}

func NewAdminPortalResource() resource.Resource {
	return newSettingsResource[settings.AdminPortalModel]("admin_portal", settings.AdminPortalSchema, "/v1/mgmt/adminportal/settings")
}

func NewOTPSettingsResource() resource.Resource {
	return newSettingsResource[settings.OTPSettingsModel]("otp_settings", settings.OTPSettingsSchema, "/v1/mgmt/otp/settings")
}

func NewMagicLinkSettingsResource() resource.Resource {
	return newSettingsResource[settings.MagicLinkSettingsModel]("magiclink_settings", settings.MagicLinkSettingsSchema, "/v1/mgmt/magiclink/settings")
}

func NewInviteSettingsResource() resource.Resource {
	return newSettingsResource[settings.InviteSettingsModel]("invite_settings", settings.InviteSettingsSchema, "/v1/mgmt/project/settings")
}

func NewProjectSettingsResource() resource.Resource {
	return newSettingsResource[settings.ProjectSettingsModel]("project_settings", settings.ProjectSettingsSchema, "/v1/mgmt/project/settings")
}

func NewSessionSettingsResource() resource.Resource {
	return newSettingsResource[settings.SessionSettingsModel]("session_settings", settings.SessionSettingsSchema, "/v1/mgmt/project/settings")
}

func NewSessionMigrationResource() resource.Resource {
	return newSettingsResource[settings.SessionMigrationSettingsModel]("session_migration", settings.SessionMigrationSettingsSchema, "/v1/mgmt/project/settings")
}

func NewEnchantedLinkSettingsResource() resource.Resource {
	return newSettingsResource[settings.EnchantedLinkSettingsModel]("enchantedlink_settings", settings.EnchantedLinkSettingsSchema, "/v1/mgmt/enchantedlink/settings")
}

func NewEmbeddedLinkSettingsResource() resource.Resource {
	return newSettingsResource[settings.EmbeddedLinkSettingsModel]("embeddedlink_settings", settings.EmbeddedLinkSettingsSchema, "/v1/mgmt/embeddedlink/settings")
}

func NewPasswordSettingsResource() resource.Resource {
	return newSettingsResource[settings.PasswordSettingsModel]("password_settings", settings.PasswordSettingsSchema, "/v2/mgmt/password/settings")
}

func NewSSOSettingsResource() resource.Resource {
	return newSettingsResource[settings.SSOSettingsModel]("sso_settings", settings.SSOSettingsSchema, "/v3/mgmt/sso/settings")
}

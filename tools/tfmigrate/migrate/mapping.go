package migrate

// The mapping between the legacy descope_project nested configuration and the standalone resources it was split into.
//
// Attribute translation is by name against the destination's live schema, because the split lifted the nested models
// into standalone models without renaming their attributes. Only the deliberate exceptions are restated here, and a
// populated legacy attribute that matches neither the destination schema nor an exception fails the migration closed.

// family describes one legacy collection or object that became a standalone resource.
type family struct {
	// path is the legacy attribute path inside the descope_project instance, dot separated.
	path string
	// resourceType is the destination Terraform type.
	resourceType string
	// shape is how the legacy value holds its objects.
	shape shape
	// idAttr is the legacy attribute holding the backend id. Empty means the collection key is the id.
	idAttr string
	// nameAttr is the legacy attribute holding the human name, used for labels and reference resolution.
	nameAttr string
	// keyAttr, when set, receives the collection key of a map shaped family (for example flow_id).
	keyAttr string
	// fixed are destination attributes with a value the legacy shape implied rather than stored.
	fixed map[string]string
	// rename maps a legacy attribute to a differently named destination attribute.
	rename map[string]string
	// flatten maps a dotted legacy path inside a nested object to a flat destination attribute.
	flatten map[string]string
	// consumed are legacy attributes handled by dedicated logic elsewhere, and therefore not translated by name.
	consumed []string
	// payload are destination attributes written to an external file and referenced with file().
	payload []string
	// serviceRefs names legacy delivery service objects whose connector name must become a connector id.
	serviceRefs []string
	// jwtRefs names destination attributes holding a JWT template id that legacy stored as a template name.
	jwtRefs []string
	// forceSecrets are schema-optional write-only attributes that legacy values make conditionally required.
	forceSecrets []string
	// group, when set, generates the whole family as one for_each resource with this label.
	group string
	// labelPrefix is used when the natural name cannot produce a valid HCL identifier.
	labelPrefix string
}

type shape int

const (
	// shapeObject is a single nested object that became a singleton resource.
	shapeObject shape = iota
	// shapeList is a list of objects, each of which became its own resource.
	shapeList
	// shapeMap is a map of objects keyed by the entity id.
	shapeMap
)

// retainedAttributes stay on descope_project after the split.
var retainedAttributes = []Retained{
	{Attribute: "id", Detail: "The project id remains the descope_project resource's own id."},
	{Attribute: "name", Detail: "Still a descope_project attribute."},
	{Attribute: "environment", Detail: "Still a descope_project attribute."},
	{Attribute: "tags", Detail: "Still a descope_project attribute."},
	{Attribute: "deletion_protection", Detail: "Still a descope_project attribute when present in source state."},
}

// authMethods are the legacy authentication.* objects that became settings singletons. Each carries the delivery
// service objects whose nested message templates became standalone method scoped template resources.
var authMethods = []authMethod{
	{legacy: "otp", resourceType: "descope_otp_settings", method: "otp", email: true, text: true, voice: true},
	{legacy: "magic_link", resourceType: "descope_magiclink_settings", method: "magiclink", email: true, text: true},
	{legacy: "enchanted_link", resourceType: "descope_enchantedlink_settings", method: "enchantedlink", email: true},
	{legacy: "embedded_link", resourceType: "descope_embeddedlink_settings", method: "embeddedlink"},
	{legacy: "password", resourceType: "descope_password_settings", method: "password", email: true},
	{legacy: "sso", resourceType: "descope_sso_settings", method: "sso", email: true},
	{legacy: "totp", resourceType: "descope_totp_settings", method: "totp"},
	{legacy: "passkeys", resourceType: "descope_passkey_settings", method: "passkey"},
}

// authMethod is one legacy authentication method object and the delivery services it could configure.
type authMethod struct {
	legacy       string
	resourceType string
	// method is the value the standalone template resources use in their method attribute and import id.
	method string
	email  bool
	text   bool
	voice  bool
}

// serviceKind is one of the three delivery services a legacy auth method could configure.
type serviceKind struct {
	// legacy is the legacy attribute name of the service object.
	legacy string
	// resourceType is the standalone template resource the nested templates became.
	resourceType string
	// templateIDAttr is the settings attribute that selects the active template by id.
	templateIDAttr string
}

var (
	emailService = serviceKind{legacy: "email_service", resourceType: "descope_email_template", templateIDAttr: "email_template_id"}
	textService  = serviceKind{legacy: "text_service", resourceType: "descope_text_template", templateIDAttr: "text_template_id"}
	voiceService = serviceKind{legacy: "voice_service", resourceType: "descope_voice_template", templateIDAttr: "voice_template_id"}
)

// projectSettingsSplit lists the singletons the legacy project_settings object was split into. Attributes are routed
// by the destination schema: every legacy key must land in exactly one of these, or in sessionMigrationFamily.
var projectSettingsSplit = []string{"descope_project_settings", "descope_session_settings"}

// sessionMigrationFamily is the nested project_settings.session_migration object.
var sessionMigrationFamily = family{
	path:         "project_settings.session_migration",
	resourceType: "descope_session_migration",
	shape:        shapeObject,
}

// jwtTemplateKinds map the two legacy jwt_templates lists onto the type attribute of the standalone resource.
var jwtTemplateKinds = []family{
	{
		path:         "jwt_templates.user_templates",
		resourceType: "descope_jwt_template",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		fixed:        map[string]string{"type": "user"},
		payload:      []string{"template"},
		labelPrefix:  "user_jwt_template",
	},
	{
		path:         "jwt_templates.access_key_templates",
		resourceType: "descope_jwt_template",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		fixed:        map[string]string{"type": "key"},
		payload:      []string{"template"},
		labelPrefix:  "key_jwt_template",
	},
}

// applicationFamilies are the three SSO application lists. Their nested permissions and roles become app scoped
// resources, so they are consumed rather than translated onto the application itself.
var applicationFamilies = []family{
	{
		path:         "applications.oidc_applications",
		resourceType: "descope_oidc_app",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		consumed:     []string{"permissions", "roles"},
		labelPrefix:  "oidc_app",
	},
	{
		path:         "applications.saml_applications",
		resourceType: "descope_saml_app",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		consumed:     []string{"permissions", "roles"},
		labelPrefix:  "saml_app",
	},
	{
		path:         "applications.wsfed_applications",
		resourceType: "descope_wsfed_app",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		consumed:     []string{"permissions", "roles"},
		labelPrefix:  "wsfed_app",
	},
}

// appPermissionFamily and appRoleFamily are the app scoped children of an SSO application.
var (
	appPermissionFamily = family{
		resourceType: "descope_app_permission",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		labelPrefix:  "app_permission",
	}
	appRoleFamily = family{
		resourceType: "descope_app_role",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		rename:       map[string]string{"permissions": "permission_ids"},
		labelPrefix:  "app_role",
	}
)

// simpleFamilies are the legacy collections whose attributes translate by name onto their destination.
var simpleFamilies = []family{
	{
		path:         "authorization.roles",
		resourceType: "descope_role",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		consumed:     []string{"key"},
		labelPrefix:  "role",
	},
	{
		path:         "authorization.permissions",
		resourceType: "descope_permission",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		labelPrefix:  "permission",
	},
	{
		path:         "attributes.user",
		resourceType: "descope_user_attribute",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		labelPrefix:  "user_attribute",
	},
	{
		path:         "attributes.tenant",
		resourceType: "descope_tenant_attribute",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		labelPrefix:  "tenant_attribute",
	},
	{
		path:         "attributes.access_key",
		resourceType: "descope_access_key_attribute",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		labelPrefix:  "access_key_attribute",
	},
	{
		path:         "lists",
		resourceType: "descope_list",
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		consumed:     []string{"type", "data"},
		payload:      []string{"json"},
		labelPrefix:  "list",
	},
	// flows and widgets were keyed maps whose key is also the backend id, so they stay a keyed collection: one
	// for_each resource and one for_each import block reproduce the shape the configuration already had.
	{
		path:         "flows",
		resourceType: "descope_flow",
		shape:        shapeMap,
		keyAttr:      "flow_id",
		payload:      []string{"data"},
		group:        "flows",
		labelPrefix:  "flow",
	},
	{
		path:         "widgets",
		resourceType: "descope_widget",
		shape:        shapeMap,
		keyAttr:      "widget_id",
		payload:      []string{"data"},
		group:        "widgets",
		labelPrefix:  "widget",
	},
	{
		path:         "styles",
		resourceType: "descope_styles",
		shape:        shapeObject,
		payload:      []string{"data"},
	},
	{
		path:         "admin_portal",
		resourceType: "descope_admin_portal",
		shape:        shapeObject,
	},
	{
		path:         "invite_settings",
		resourceType: "descope_invite_settings",
		shape:        shapeObject,
		// the invitation settings kept their nested message templates, and only the connector became an id
		serviceRefs: []string{"email_service"},
	},
}

// fgaFamily is the legacy authorization.fga string, which became the schema attribute of its own resource.
var fgaFamily = family{
	path:         "authorization.fga",
	resourceType: "descope_fga_schema",
	shape:        shapeObject,
	rename:       map[string]string{"fga": "schema"},
}

// oauthSettingsFamily is the legacy authentication.oauth object minus its providers.
var oauthSettingsFamily = family{
	path:         "authentication.oauth",
	resourceType: "descope_oauth_settings",
	shape:        shapeObject,
	consumed:     []string{"system", "custom"},
}

// oauthProviderFamily is one entry of authentication.oauth.system or authentication.oauth.custom. The provider name
// is its backend id, so the id is the key rather than a stored attribute.
var oauthProviderFamily = family{
	resourceType: "descope_oauth_provider",
	shape:        shapeMap,
	keyAttr:      "id",
	labelPrefix:  "oauth_provider",
}

// oauthSystemProviders are the built-in providers the legacy schema exposed as fixed attributes.
var oauthSystemProviders = []string{"apple", "discord", "facebook", "github", "gitlab", "google", "linkedin", "microsoft", "slack"}

// connectorFamily builds the family for one connector, using the legacy attribute name and its destination type.
// Legacy connectors nested some configuration that the standalone connectors flattened, so the exceptions below
// carry those values across. A legacy attribute with no destination fails closed like any other.
func connectorFamily(legacyName, resourceType string) family {
	f := family{
		path:         "connectors." + legacyName,
		resourceType: resourceType,
		shape:        shapeList,
		idAttr:       "id",
		nameAttr:     "name",
		labelPrefix:  legacyName + "_connector",
	}
	if rename, ok := connectorRenames[legacyName]; ok {
		f.rename = rename
	}
	if flatten, ok := connectorFlattening[legacyName]; ok {
		f.flatten = flatten
	}
	return f
}

// connectorRenames are the connector attributes the split renamed. `project_id` became `config_project_id` wherever a
// connector's own configuration had a project id, because the standalone resource needs `project_id` for its scope.
var connectorRenames = map[string]map[string]string{
	"google_cloud_translation": {"project_id": "config_project_id"},
	"mixpanel":                 {"project_id": "config_project_id"},
	"recaptcha_enterprise":     {"project_id": "config_project_id"},
	"lokalise":                 {"project_id": "config_project_id"},
	"ses":                      {"secret": "secret_access_key"},
	// sns also carried a deprecated alias of origination_number.
	"sns": {"secret": "secret_access_key", "organization_number": "origination_number"},
}

// connectorFlattening lifts a legacy nested connector object into the flat attributes the standalone connector uses.
// Each entry maps a legacy nested path to its destination attribute.
var connectorFlattening = map[string]map[string]string{
	"smtp": {
		"sender.email":            "sender_email",
		"sender.name":             "sender_name",
		"server.host":             "host",
		"server.port":             "port",
		"authentication.username": "username",
		"authentication.password": "password",
		"authentication.method":   "auth_method",
	},
	"sendgrid": {
		"sender.email":           "sender_email",
		"sender.name":            "sender_name",
		"authentication.api_key": "api_key",
	},
	"ses": {
		"sender.email": "sender_email",
		"sender.name":  "sender_name",
	},
	"twilio_core": {
		"senders.sms.phone_number":          "from_phone",
		"senders.sms.messaging_service_sid": "messaging_service_sid",
		"senders.voice.phone_number":        "from_phone_voice",
		"authentication.auth_token":         "auth_token",
		"authentication.api_key":            "api_key",
		"authentication.api_secret":         "api_secret",
	},
	"twilio_verify": {
		"authentication.auth_token": "auth_token",
		"authentication.api_key":    "api_key",
		"authentication.api_secret": "api_secret",
	},
}

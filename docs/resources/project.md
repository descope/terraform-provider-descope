---
page_title: "descope_project Resource - descope"
subcategory: ""
description: |-
  Manages a Descope project's full configuration: authentication methods, flows, roles, permissions, connectors, applications, JWT templates, and more.
---

# descope_project (Resource)

Manages the configuration of a Descope project. A project is the core entity in Descope—it contains all authentication settings, user flows, roles, connectors, and other configuration for your application.

This resource manages _project configuration_, not users or tenants. For user management, use the [Descope Management API](https://docs.descope.com/api/openapi) or [SDKs](https://docs.descope.com).

For a full reference of all supported connectors, see the [connectors reference](https://docs.descope.com/connectors) in the Descope documentation.

## Example Usage

### Basic Project

```hcl
resource "descope_project" "example" {
  name        = "my-app"
  environment = "production"
  tags        = ["prod", "v2"]
}
```

### Authentication Methods

Enable and configure the authentication methods your users will use:

```hcl
resource "descope_project" "example" {
  name = "my-app"

  authentication = {
    # Email magic link (passwordless)
    magic_link = {
      expiration_time = "1 hour"
    }

    # Password-based login with lockout
    password = {
      lock          = true
      lock_attempts = 5
      min_length    = 12
    }

    # OTP via email or SMS
    otp = {
      expiration_time = "5 minutes"
    }

    # Passkeys (WebAuthn)
    passkeys = {
      disabled = false
    }
  }
}
```

### Roles and Permissions (RBAC)

Define roles and permissions for your users:

```hcl
resource "descope_project" "example" {
  name = "my-app"

  authorization = {
    permissions = [
      { name = "read:data",   description = "Read access to application data" },
      { name = "write:data",  description = "Write access to application data" },
      { name = "admin:panel", description = "Access to the admin panel" },
    ]

    roles = [
      {
        name        = "viewer"
        description = "Can read data"
        permissions = ["read:data"]
      },
      {
        name        = "editor"
        description = "Can read and write data"
        permissions = ["read:data", "write:data"]
      },
      {
        name        = "admin"
        description = "Full access"
        permissions = ["read:data", "write:data", "admin:panel"]
      },
    ]
  }
}
```

### Connectors

Integrate with third-party services to enrich flows and send notifications:

```hcl
resource "descope_project" "example" {
  name = "my-app"

  connectors = {
    # Generic HTTP webhook with bearer token auth
    http = [{
      name        = "User Eligibility Check"
      description = "Checks if a new user is allowed to register"
      base_url    = "https://api.example.com"
      authentication = {
        bearer_token = var.webhook_secret
      }
    }]

    # SendGrid for email delivery
    sendgrid = [{
      name = "Transactional Email"
      sender = {
        email = "noreply@example.com"
        name  = "My App"
      }
      authentication = {
        api_key = var.sendgrid_api_key
      }
    }]

    # Twilio for SMS OTP
    twilio_core = [{
      name        = "SMS OTP"
      account_sid = var.twilio_account_sid
      senders = {
        sms = {
          phone_number = "+15551234567"
        }
      }
      authentication = {
        auth_token = var.twilio_auth_token
      }
    }]
  }
}
```

### Flows and Styles

Load a custom authentication flow from a JSON file (exported from the Descope console):

```hcl
resource "descope_project" "example" {
  name = "my-app"

  flows = {
    "sign-up-or-in" = {
      data = file("${path.module}/flows/sign-up-or-in.json")
    }
    "forgot-password" = {
      data = file("${path.module}/flows/forgot-password.json")
    }
  }

  styles = {
    data = file("${path.module}/flows/styles.json")
  }
}
```

### Session Settings

Configure token lifetimes and session behavior:

```hcl
resource "descope_project" "example" {
  name = "my-app"

  project_settings = {
    refresh_token_expiration     = "3 weeks"
    session_token_expiration     = "15 minutes"
    refresh_token_rotation       = true
    enable_inactivity            = true
    inactivity_time              = "30 minutes"
    custom_domain                = "auth.example.com"
    approved_domains             = ["example.com", "app.example.com"]
  }
}
```

### OIDC Applications

Register an OIDC application for SSO:

```hcl
resource "descope_project" "example" {
  name = "my-app"

  applications = {
    oidc_applications = [
      {
        name          = "My Web App"
        description   = "Primary web application"
        login_page_url = "https://app.example.com/login"
      }
    ]
  }
}
```

### JWT Templates

Customize the JWT claims added to session tokens:

```hcl
resource "descope_project" "example" {
  name = "my-app"

  jwt_templates = {
    user_templates = [
      {
        name     = "app-claims"
        template = jsonencode({
          tier   = "@user.customAttributes.subscriptionTier"
          org_id = "@user.tenants[0].tenantId"
        })
      }
    ]
  }

  project_settings = {
    user_jwt_template = "app-claims"
  }
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the Descope project.

### Optional

- `admin_portal` (Attributes) Admin portal configuration - A hosted page for end users to access and use Descope Widgets (see [below for nested schema](#nestedatt--admin_portal))
- `applications` (Attributes) Applications that are registered with the project. (see [below for nested schema](#nestedatt--applications))
- `attributes` (Attributes) Custom attributes that can be attached to users and tenants. (see [below for nested schema](#nestedatt--attributes))
- `authentication` (Attributes) Settings for each authentication method. (see [below for nested schema](#nestedatt--authentication))
- `authorization` (Attributes) Define Role-Based Access Control (RBAC) for your users by creating roles and permissions. (see [below for nested schema](#nestedatt--authorization))
- `connectors` (Attributes) Enrich your flows by interacting with third party services. (see [below for nested schema](#nestedatt--connectors))
- `environment` (String) This can be set to `production` to mark production projects, otherwise this should be left unset for development or staging projects.
- `flows` (Attributes Map) Custom authentication flows to use in this project. (see [below for nested schema](#nestedatt--flows))
- `invite_settings` (Attributes) User invitation settings and behavior. (see [below for nested schema](#nestedatt--invite_settings))
- `jwt_templates` (Attributes) Defines templates for JSON Web Tokens (JWT) used for authentication. (see [below for nested schema](#nestedatt--jwt_templates))
- `lists` (Attributes List) Lists that can be used for various purposes in the project, such as IP allowlists, text lists, or custom JSON data. (see [below for nested schema](#nestedatt--lists))
- `project_settings` (Attributes) General settings for the Descope project. (see [below for nested schema](#nestedatt--project_settings))
- `styles` (Attributes) Custom styles that can be applied to the project's authentication flows. (see [below for nested schema](#nestedatt--styles))
- `tags` (Set of String) Descriptive tags for your Descope project. Each tag must be no more than 50 characters long.
- `widgets` (Attributes Map) Embeddable components designed to facilitate the delegation of operations to tenant admins and end users. (see [below for nested schema](#nestedatt--widgets))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--admin_portal"></a>
### Nested Schema for `admin_portal`

Optional:

- `enabled` (Boolean) Whether the Admin Portal is enabled
- `style_id` (String) The style id to use
- `widgets` (Attributes List) The widgets to show in the Admin Portal (see [below for nested schema](#nestedatt--admin_portal--widgets))

<a id="nestedatt--admin_portal--widgets"></a>
### Nested Schema for `admin_portal.widgets`

Required:

- `type` (String) The type of the Widget
- `widget_id` (String) The unique identifier of the Widget



<a id="nestedatt--applications"></a>
### Nested Schema for `applications`

Optional:

- `oidc_applications` (Attributes List) Applications using OpenID Connect (OIDC) for authentication. (see [below for nested schema](#nestedatt--applications--oidc_applications))
- `saml_applications` (Attributes List) Applications using SAML for authentication. (see [below for nested schema](#nestedatt--applications--saml_applications))

<a id="nestedatt--applications--oidc_applications"></a>
### Nested Schema for `applications.oidc_applications`

Required:

- `name` (String) A name for the OIDC application.

Optional:

- `claims` (List of String) A list of supported claims. e.g. `sub`, `email`, `exp`.
- `description` (String) A description for the OIDC application.
- `disabled` (Boolean) Whether the application should be enabled or disabled.
- `force_authentication` (Boolean) This configuration overrides the default behavior of the SSO application and forces the user to authenticate via the Descope flow, regardless of the SP's request.
- `id` (String) An optional identifier for the OIDC application.
- `login_page_url` (String) The Flow Hosting URL. Read more about using this parameter with custom domain [here](https://docs.descope.com/sso-integrations/applications/saml-apps).
- `logo` (String) A logo for the OIDC application. Should be a hosted image URL.


<a id="nestedatt--applications--saml_applications"></a>
### Nested Schema for `applications.saml_applications`

Required:

- `name` (String) A name for the SAML application.

Optional:

- `acs_allowed_callback_urls` (Set of String) A list of allowed ACS callback URLS. This configuration is used when the default ACS URL value is unreachable. Supports wildcards.
- `attribute_mapping` (Attributes List) The `AttributeMapping` object. Read the description below. (see [below for nested schema](#nestedatt--applications--saml_applications--attribute_mapping))
- `default_relay_state` (String) The default relay state. When using IdP-initiated authentication, this value may be used as a URL to a resource in the Service Provider.
- `description` (String) A description for the SAML application.
- `disabled` (Boolean) Whether the application should be enabled or disabled.
- `dynamic_configuration` (Attributes) The `DynamicConfiguration` object. Read the description below. (see [below for nested schema](#nestedatt--applications--saml_applications--dynamic_configuration))
- `force_authentication` (Boolean) This configuration overrides the default behavior of the SSO application and forces the user to authenticate via the Descope flow, regardless of the SP's request.
- `id` (String) An optional identifier for the SAML application.
- `login_page_url` (String) The Flow Hosting URL. Read more about using this parameter with custom domain [here](https://docs.descope.com/sso-integrations/applications/saml-apps).
- `logo` (String) A logo for the SAML application. Should be a hosted image URL.
- `manual_configuration` (Attributes) The `ManualConfiguration` object. Read the description below. (see [below for nested schema](#nestedatt--applications--saml_applications--manual_configuration))
- `subject_name_id_format` (String) The subject name id format. Choose one of "", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress", "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent", "urn:oasis:names:tc:SAML:2.0:nameid-format:transient". Read more about this configuration [here](https://docs.descope.com/sso-integrations/applications/saml-apps).
- `subject_name_id_type` (String) The subject name id type. Choose one of "", "email", "phone". Read more about this configuration [here](https://docs.descope.com/sso-integrations/applications/saml-apps).

<a id="nestedatt--applications--saml_applications--attribute_mapping"></a>
### Nested Schema for `applications.saml_applications.attribute_mapping`

Required:

- `name` (String) The name of the attribute.
- `value` (String) The value of the attribute.


<a id="nestedatt--applications--saml_applications--dynamic_configuration"></a>
### Nested Schema for `applications.saml_applications.dynamic_configuration`

Required:

- `metadata_url` (String) The metadata URL when retrieving the connection details dynamically.


<a id="nestedatt--applications--saml_applications--manual_configuration"></a>
### Nested Schema for `applications.saml_applications.manual_configuration`

Required:

- `acs_url` (String) Enter the `ACS URL` from the SP.
- `entity_id` (String) Enter the `Entity Id` from the SP.

Optional:

- `certificate` (String) Enter the `Certificate` from the SP.




<a id="nestedatt--attributes"></a>
### Nested Schema for `attributes`

Optional:

- `access_key` (Attributes List) A list of custom attributes for storing additional details about each access key in the project. (see [below for nested schema](#nestedatt--attributes--access_key))
- `tenant` (Attributes List) A list of custom attributes for storing additional details about each tenant in the project. (see [below for nested schema](#nestedatt--attributes--tenant))
- `user` (Attributes List) A list of custom attributes for storing additional details about each user in the project. (see [below for nested schema](#nestedatt--attributes--user))

<a id="nestedatt--attributes--access_key"></a>
### Nested Schema for `attributes.access_key`

Required:

- `name` (String) The name of the attribute. This value is called `Display Name` in the Descope console.
- `type` (String) The type of the attribute. Choose one of "string", "number", "boolean", "singleselect", "multiselect", "date".

Optional:

- `id` (String) An optional identifier for the attribute. This value is called `Machine Name` in the Descope console. If a value is not provided then an appropriate one will be created from the value of `name`.
- `select_options` (Set of String) When the attribute type is "multiselect". A list of options to choose from.
- `widget_authorization` (Attributes) Determines the permissions access key are required to have to access this attribute in the access key management widget. (see [below for nested schema](#nestedatt--attributes--access_key--widget_authorization))

<a id="nestedatt--attributes--access_key--widget_authorization"></a>
### Nested Schema for `attributes.access_key.widget_authorization`

Optional:

- `edit_permissions` (Set of String) The permissions users are required to have to edit this attribute in the access key management widget.
- `view_permissions` (Set of String) The permissions users are required to have to view this attribute in the access key management widget.



<a id="nestedatt--attributes--tenant"></a>
### Nested Schema for `attributes.tenant`

Required:

- `name` (String) The name of the attribute. This value is called `Display Name` in the Descope console.
- `type` (String) The type of the attribute. Choose one of "string", "number", "boolean", "singleselect", "multiselect", "date".

Optional:

- `authorization` (Attributes) Determines the required permissions for this tenant. (see [below for nested schema](#nestedatt--attributes--tenant--authorization))
- `id` (String) An optional identifier for the attribute. This value is called `Machine Name` in the Descope console. If a value is not provided then an appropriate one will be created from the value of `name`.
- `select_options` (Set of String) When the attribute type is "multiselect". A list of options to choose from.

<a id="nestedatt--attributes--tenant--authorization"></a>
### Nested Schema for `attributes.tenant.authorization`

Optional:

- `view_permissions` (Set of String) Determines the required permissions for this tenant.



<a id="nestedatt--attributes--user"></a>
### Nested Schema for `attributes.user`

Required:

- `name` (String) The name of the attribute. This value is called `Display Name` in the Descope console.
- `type` (String) The type of the attribute. Choose one of "string", "number", "boolean", "singleselect", "multiselect", "date".

Optional:

- `id` (String) An optional identifier for the attribute. This value is called `Machine Name` in the Descope console. If a value is not provided then an appropriate one will be created from the value of `name`.
- `select_options` (Set of String) When the attribute type is "multiselect". A list of options to choose from.
- `widget_authorization` (Attributes) Determines the permissions users are required to have to access this attribute in the user management widget. (see [below for nested schema](#nestedatt--attributes--user--widget_authorization))

<a id="nestedatt--attributes--user--widget_authorization"></a>
### Nested Schema for `attributes.user.widget_authorization`

Optional:

- `edit_permissions` (Set of String) The permissions users are required to have to edit this attribute in the user management widget.
- `view_permissions` (Set of String) The permissions users are required to have to view this attribute in the user management widget.




<a id="nestedatt--authentication"></a>
### Nested Schema for `authentication`

Optional:

- `embedded_link` (Attributes) Make the authentication experience smoother for the user by generating their initial token in a way that does not require the end user to initiate the process, requiring only verification. (see [below for nested schema](#nestedatt--authentication--embedded_link))
- `enchanted_link` (Attributes) An enhanced and more secure version of Magic Link, enabling users to start the authentication process on one device and execute the verification on another. (see [below for nested schema](#nestedatt--authentication--enchanted_link))
- `magic_link` (Attributes) An authentication method where a user receives a unique link via email to log in. (see [below for nested schema](#nestedatt--authentication--magic_link))
- `oauth` (Attributes) Authentication using Open Authorization, which allows users to authenticate with various external services. (see [below for nested schema](#nestedatt--authentication--oauth))
- `otp` (Attributes) A dynamically generated set of numbers, granting the user one-time access. (see [below for nested schema](#nestedatt--authentication--otp))
- `passkeys` (Attributes) Device-based passwordless authentication, using fingerprint, face scan, and more. (see [below for nested schema](#nestedatt--authentication--passkeys))
- `password` (Attributes) The classic username and password combination used for authentication. (see [below for nested schema](#nestedatt--authentication--password))
- `sso` (Attributes) Single Sign-On (SSO) authentication method that enables users to access multiple applications with a single set of credentials. (see [below for nested schema](#nestedatt--authentication--sso))
- `totp` (Attributes) A one-time code generated for the user using a shared secret and time. (see [below for nested schema](#nestedatt--authentication--totp))

<a id="nestedatt--authentication--embedded_link"></a>
### Nested Schema for `authentication.embedded_link`

Optional:

- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `expiration_time` (String) How long the embedded link remains valid before it expires.


<a id="nestedatt--authentication--enchanted_link"></a>
### Nested Schema for `authentication.enchanted_link`

Optional:

- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `email_service` (Attributes) Settings related to sending emails as part of the enchanted link authentication. (see [below for nested schema](#nestedatt--authentication--enchanted_link--email_service))
- `expiration_time` (String) How long the enchanted link remains valid before it expires.
- `redirect_url` (String) The URL to redirect users to after they log in using the enchanted link.

<a id="nestedatt--authentication--enchanted_link--email_service"></a>
### Nested Schema for `authentication.enchanted_link.email_service`

Required:

- `connector` (String) The name of the email connector to use for sending emails.

Optional:

- `templates` (Attributes List) A list of email templates for different authentication flows. (see [below for nested schema](#nestedatt--authentication--enchanted_link--email_service--templates))

<a id="nestedatt--authentication--enchanted_link--email_service--templates"></a>
### Nested Schema for `authentication.enchanted_link.email_service.templates`

Required:

- `name` (String) Unique name for this email template.
- `subject` (String) Subject line of the email message.

Optional:

- `active` (Boolean) Whether this email template is currently active and in use.
- `html_body` (String) HTML content of the email message body, required if `use_plain_text_body` isn't set.
- `plain_text_body` (String) Plain text version of the email message body, required if `use_plain_text_body` is set to `true`.
- `use_plain_text_body` (Boolean) Whether to use the plain text body instead of HTML for the email.

Read-Only:

- `id` (String)




<a id="nestedatt--authentication--magic_link"></a>
### Nested Schema for `authentication.magic_link`

Optional:

- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `email_service` (Attributes) Settings related to sending emails as part of the magic link authentication. (see [below for nested schema](#nestedatt--authentication--magic_link--email_service))
- `expiration_time` (String) How long the magic link remains valid before it expires.
- `redirect_url` (String) The URL to redirect users to after they log in using the magic link.
- `text_service` (Attributes) Settings related to sending SMS messages as part of the magic link authentication. (see [below for nested schema](#nestedatt--authentication--magic_link--text_service))

<a id="nestedatt--authentication--magic_link--email_service"></a>
### Nested Schema for `authentication.magic_link.email_service`

Required:

- `connector` (String) The name of the email connector to use for sending emails.

Optional:

- `templates` (Attributes List) A list of email templates for different authentication flows. (see [below for nested schema](#nestedatt--authentication--magic_link--email_service--templates))

<a id="nestedatt--authentication--magic_link--email_service--templates"></a>
### Nested Schema for `authentication.magic_link.email_service.templates`

Required:

- `name` (String) Unique name for this email template.
- `subject` (String) Subject line of the email message.

Optional:

- `active` (Boolean) Whether this email template is currently active and in use.
- `html_body` (String) HTML content of the email message body, required if `use_plain_text_body` isn't set.
- `plain_text_body` (String) Plain text version of the email message body, required if `use_plain_text_body` is set to `true`.
- `use_plain_text_body` (Boolean) Whether to use the plain text body instead of HTML for the email.

Read-Only:

- `id` (String)



<a id="nestedatt--authentication--magic_link--text_service"></a>
### Nested Schema for `authentication.magic_link.text_service`

Required:

- `connector` (String) The name of the SMS/text connector to use for sending text messages.

Optional:

- `templates` (Attributes List) A list of text message templates for different authentication flows. (see [below for nested schema](#nestedatt--authentication--magic_link--text_service--templates))

<a id="nestedatt--authentication--magic_link--text_service--templates"></a>
### Nested Schema for `authentication.magic_link.text_service.templates`

Required:

- `body` (String) The content of the text message.
- `name` (String) Unique name for this text template.

Optional:

- `active` (Boolean) Whether this text template is currently active and in use.

Read-Only:

- `id` (String)




<a id="nestedatt--authentication--oauth"></a>
### Nested Schema for `authentication.oauth`

Optional:

- `custom` (Attributes Map) Custom OAuth providers configured for this project. (see [below for nested schema](#nestedatt--authentication--oauth--custom))
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `system` (Attributes) Custom configurations for builtin OAuth providers such as Apple, Google, GitHub, Facebook, etc. (see [below for nested schema](#nestedatt--authentication--oauth--system))

<a id="nestedatt--authentication--oauth--custom"></a>
### Nested Schema for `authentication.oauth.custom`

Optional:

- `allowed_grant_types` (List of String) The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens from the OAuth provider.
- `apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic apple client secret for applications. (see [below for nested schema](#nestedatt--authentication--oauth--custom--apple_key_generator))
- `authorization_endpoint` (String) The URL that users are redirected to for authorization with the OAuth provider.
- `callback_domain` (String) Use a custom domain in your OAuth verification screen.
- `claim_mapping` (Map of String) Maps OAuth provider claims to Descope user attributes.
- `client_id` (String) The client ID for the OAuth provider, used to identify the application to the provider.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider, used to authenticate the application with the provider.
- `description` (String) A brief description of the OAuth provider.
- `disable_jit_updates` (Boolean) By default the user attribute mapping configuration is used to update the user's attributes automatically during sign in. Disable this if you want this to happen only during user creation.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `issuer` (String) The issuer identifier for the OAuth provider.
- `jwks_endpoint` (String) The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.
- `logo` (String) The URL of the logo associated with the OAuth provider.
- `manage_provider_tokens` (Boolean) Whether to enable provider token management for this OAuth provider.
- `merge_user_accounts` (Boolean) Whether to merge existing user accounts with new ones created through OAuth authentication.
- `native_apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic native apple client secret for mobile apps. (see [below for nested schema](#nestedatt--authentication--oauth--custom--native_apple_key_generator))
- `native_client_id` (String) The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.
- `native_client_secret` (String, Sensitive) The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.
- `prompts` (List of String) Custom prompts or consent screens that users may see during OAuth authentication.
- `provider_token_management` (Attributes) This attribute is deprecated, use the `manage_provider_tokens`, `callback_domain`, and `redirect_url` fields instead. (see [below for nested schema](#nestedatt--authentication--oauth--custom--provider_token_management))
- `redirect_url` (String) Users will be directed to this URL after authentication. If redirect URL is specified in the SDK/API call, it will override this value.
- `scopes` (List of String) Scopes of access that the application requests from the user's account on the OAuth provider.
- `token_endpoint` (String) The URL where the application requests an access token from the OAuth provider.
- `use_client_assertion` (Boolean) Use private key JWT (client assertion) instead of client secret.
- `user_info_endpoint` (String) The URL where the application retrieves user information from the OAuth provider.

<a id="nestedatt--authentication--oauth--custom--apple_key_generator"></a>
### Nested Schema for `authentication.oauth.custom.apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--custom--native_apple_key_generator"></a>
### Nested Schema for `authentication.oauth.custom.native_apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--custom--provider_token_management"></a>
### Nested Schema for `authentication.oauth.custom.provider_token_management`



<a id="nestedatt--authentication--oauth--system"></a>
### Nested Schema for `authentication.oauth.system`

Optional:

- `apple` (Attributes) Apple's OAuth provider, allowing users to authenticate with their Apple Account. (see [below for nested schema](#nestedatt--authentication--oauth--system--apple))
- `discord` (Attributes) Discord's OAuth provider, allowing users to authenticate with their Discord account. (see [below for nested schema](#nestedatt--authentication--oauth--system--discord))
- `facebook` (Attributes) Facebook's OAuth provider, allowing users to authenticate with their Facebook account. (see [below for nested schema](#nestedatt--authentication--oauth--system--facebook))
- `github` (Attributes) GitHub's OAuth provider, allowing users to authenticate with their GitHub account. (see [below for nested schema](#nestedatt--authentication--oauth--system--github))
- `gitlab` (Attributes) GitLab's OAuth provider, allowing users to authenticate with their GitLab account. (see [below for nested schema](#nestedatt--authentication--oauth--system--gitlab))
- `google` (Attributes) Google's OAuth provider, allowing users to authenticate with their Google account. (see [below for nested schema](#nestedatt--authentication--oauth--system--google))
- `linkedin` (Attributes) LinkedIn's OAuth provider, allowing users to authenticate with their LinkedIn account. (see [below for nested schema](#nestedatt--authentication--oauth--system--linkedin))
- `microsoft` (Attributes) Microsoft's OAuth provider, allowing users to authenticate with their Microsoft account. (see [below for nested schema](#nestedatt--authentication--oauth--system--microsoft))
- `slack` (Attributes) Slack's OAuth provider, allowing users to authenticate with their Slack account. (see [below for nested schema](#nestedatt--authentication--oauth--system--slack))

<a id="nestedatt--authentication--oauth--system--apple"></a>
### Nested Schema for `authentication.oauth.system.apple`

Optional:

- `allowed_grant_types` (List of String) The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens from the OAuth provider.
- `apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic apple client secret for applications. (see [below for nested schema](#nestedatt--authentication--oauth--system--apple--apple_key_generator))
- `authorization_endpoint` (String) The URL that users are redirected to for authorization with the OAuth provider.
- `callback_domain` (String) Use a custom domain in your OAuth verification screen.
- `claim_mapping` (Map of String) Maps OAuth provider claims to Descope user attributes.
- `client_id` (String) The client ID for the OAuth provider, used to identify the application to the provider.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider, used to authenticate the application with the provider.
- `description` (String) A brief description of the OAuth provider.
- `disable_jit_updates` (Boolean) By default the user attribute mapping configuration is used to update the user's attributes automatically during sign in. Disable this if you want this to happen only during user creation.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `issuer` (String) The issuer identifier for the OAuth provider.
- `jwks_endpoint` (String) The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.
- `logo` (String) The URL of the logo associated with the OAuth provider.
- `manage_provider_tokens` (Boolean) Whether to enable provider token management for this OAuth provider.
- `merge_user_accounts` (Boolean) Whether to merge existing user accounts with new ones created through OAuth authentication.
- `native_apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic native apple client secret for mobile apps. (see [below for nested schema](#nestedatt--authentication--oauth--system--apple--native_apple_key_generator))
- `native_client_id` (String) The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.
- `native_client_secret` (String, Sensitive) The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.
- `prompts` (List of String) Custom prompts or consent screens that users may see during OAuth authentication.
- `provider_token_management` (Attributes) This attribute is deprecated, use the `manage_provider_tokens`, `callback_domain`, and `redirect_url` fields instead. (see [below for nested schema](#nestedatt--authentication--oauth--system--apple--provider_token_management))
- `redirect_url` (String) Users will be directed to this URL after authentication. If redirect URL is specified in the SDK/API call, it will override this value.
- `scopes` (List of String) Scopes of access that the application requests from the user's account on the OAuth provider.
- `token_endpoint` (String) The URL where the application requests an access token from the OAuth provider.
- `use_client_assertion` (Boolean) Use private key JWT (client assertion) instead of client secret.
- `user_info_endpoint` (String) The URL where the application retrieves user information from the OAuth provider.

<a id="nestedatt--authentication--oauth--system--apple--apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.apple.apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--apple--native_apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.apple.native_apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--apple--provider_token_management"></a>
### Nested Schema for `authentication.oauth.system.apple.provider_token_management`



<a id="nestedatt--authentication--oauth--system--discord"></a>
### Nested Schema for `authentication.oauth.system.discord`

Optional:

- `allowed_grant_types` (List of String) The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens from the OAuth provider.
- `apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic apple client secret for applications. (see [below for nested schema](#nestedatt--authentication--oauth--system--discord--apple_key_generator))
- `authorization_endpoint` (String) The URL that users are redirected to for authorization with the OAuth provider.
- `callback_domain` (String) Use a custom domain in your OAuth verification screen.
- `claim_mapping` (Map of String) Maps OAuth provider claims to Descope user attributes.
- `client_id` (String) The client ID for the OAuth provider, used to identify the application to the provider.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider, used to authenticate the application with the provider.
- `description` (String) A brief description of the OAuth provider.
- `disable_jit_updates` (Boolean) By default the user attribute mapping configuration is used to update the user's attributes automatically during sign in. Disable this if you want this to happen only during user creation.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `issuer` (String) The issuer identifier for the OAuth provider.
- `jwks_endpoint` (String) The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.
- `logo` (String) The URL of the logo associated with the OAuth provider.
- `manage_provider_tokens` (Boolean) Whether to enable provider token management for this OAuth provider.
- `merge_user_accounts` (Boolean) Whether to merge existing user accounts with new ones created through OAuth authentication.
- `native_apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic native apple client secret for mobile apps. (see [below for nested schema](#nestedatt--authentication--oauth--system--discord--native_apple_key_generator))
- `native_client_id` (String) The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.
- `native_client_secret` (String, Sensitive) The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.
- `prompts` (List of String) Custom prompts or consent screens that users may see during OAuth authentication.
- `provider_token_management` (Attributes) This attribute is deprecated, use the `manage_provider_tokens`, `callback_domain`, and `redirect_url` fields instead. (see [below for nested schema](#nestedatt--authentication--oauth--system--discord--provider_token_management))
- `redirect_url` (String) Users will be directed to this URL after authentication. If redirect URL is specified in the SDK/API call, it will override this value.
- `scopes` (List of String) Scopes of access that the application requests from the user's account on the OAuth provider.
- `token_endpoint` (String) The URL where the application requests an access token from the OAuth provider.
- `use_client_assertion` (Boolean) Use private key JWT (client assertion) instead of client secret.
- `user_info_endpoint` (String) The URL where the application retrieves user information from the OAuth provider.

<a id="nestedatt--authentication--oauth--system--discord--apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.discord.apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--discord--native_apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.discord.native_apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--discord--provider_token_management"></a>
### Nested Schema for `authentication.oauth.system.discord.provider_token_management`



<a id="nestedatt--authentication--oauth--system--facebook"></a>
### Nested Schema for `authentication.oauth.system.facebook`

Optional:

- `allowed_grant_types` (List of String) The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens from the OAuth provider.
- `apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic apple client secret for applications. (see [below for nested schema](#nestedatt--authentication--oauth--system--facebook--apple_key_generator))
- `authorization_endpoint` (String) The URL that users are redirected to for authorization with the OAuth provider.
- `callback_domain` (String) Use a custom domain in your OAuth verification screen.
- `claim_mapping` (Map of String) Maps OAuth provider claims to Descope user attributes.
- `client_id` (String) The client ID for the OAuth provider, used to identify the application to the provider.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider, used to authenticate the application with the provider.
- `description` (String) A brief description of the OAuth provider.
- `disable_jit_updates` (Boolean) By default the user attribute mapping configuration is used to update the user's attributes automatically during sign in. Disable this if you want this to happen only during user creation.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `issuer` (String) The issuer identifier for the OAuth provider.
- `jwks_endpoint` (String) The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.
- `logo` (String) The URL of the logo associated with the OAuth provider.
- `manage_provider_tokens` (Boolean) Whether to enable provider token management for this OAuth provider.
- `merge_user_accounts` (Boolean) Whether to merge existing user accounts with new ones created through OAuth authentication.
- `native_apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic native apple client secret for mobile apps. (see [below for nested schema](#nestedatt--authentication--oauth--system--facebook--native_apple_key_generator))
- `native_client_id` (String) The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.
- `native_client_secret` (String, Sensitive) The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.
- `prompts` (List of String) Custom prompts or consent screens that users may see during OAuth authentication.
- `provider_token_management` (Attributes) This attribute is deprecated, use the `manage_provider_tokens`, `callback_domain`, and `redirect_url` fields instead. (see [below for nested schema](#nestedatt--authentication--oauth--system--facebook--provider_token_management))
- `redirect_url` (String) Users will be directed to this URL after authentication. If redirect URL is specified in the SDK/API call, it will override this value.
- `scopes` (List of String) Scopes of access that the application requests from the user's account on the OAuth provider.
- `token_endpoint` (String) The URL where the application requests an access token from the OAuth provider.
- `use_client_assertion` (Boolean) Use private key JWT (client assertion) instead of client secret.
- `user_info_endpoint` (String) The URL where the application retrieves user information from the OAuth provider.

<a id="nestedatt--authentication--oauth--system--facebook--apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.facebook.apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--facebook--native_apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.facebook.native_apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--facebook--provider_token_management"></a>
### Nested Schema for `authentication.oauth.system.facebook.provider_token_management`



<a id="nestedatt--authentication--oauth--system--github"></a>
### Nested Schema for `authentication.oauth.system.github`

Optional:

- `allowed_grant_types` (List of String) The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens from the OAuth provider.
- `apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic apple client secret for applications. (see [below for nested schema](#nestedatt--authentication--oauth--system--github--apple_key_generator))
- `authorization_endpoint` (String) The URL that users are redirected to for authorization with the OAuth provider.
- `callback_domain` (String) Use a custom domain in your OAuth verification screen.
- `claim_mapping` (Map of String) Maps OAuth provider claims to Descope user attributes.
- `client_id` (String) The client ID for the OAuth provider, used to identify the application to the provider.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider, used to authenticate the application with the provider.
- `description` (String) A brief description of the OAuth provider.
- `disable_jit_updates` (Boolean) By default the user attribute mapping configuration is used to update the user's attributes automatically during sign in. Disable this if you want this to happen only during user creation.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `issuer` (String) The issuer identifier for the OAuth provider.
- `jwks_endpoint` (String) The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.
- `logo` (String) The URL of the logo associated with the OAuth provider.
- `manage_provider_tokens` (Boolean) Whether to enable provider token management for this OAuth provider.
- `merge_user_accounts` (Boolean) Whether to merge existing user accounts with new ones created through OAuth authentication.
- `native_apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic native apple client secret for mobile apps. (see [below for nested schema](#nestedatt--authentication--oauth--system--github--native_apple_key_generator))
- `native_client_id` (String) The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.
- `native_client_secret` (String, Sensitive) The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.
- `prompts` (List of String) Custom prompts or consent screens that users may see during OAuth authentication.
- `provider_token_management` (Attributes) This attribute is deprecated, use the `manage_provider_tokens`, `callback_domain`, and `redirect_url` fields instead. (see [below for nested schema](#nestedatt--authentication--oauth--system--github--provider_token_management))
- `redirect_url` (String) Users will be directed to this URL after authentication. If redirect URL is specified in the SDK/API call, it will override this value.
- `scopes` (List of String) Scopes of access that the application requests from the user's account on the OAuth provider.
- `token_endpoint` (String) The URL where the application requests an access token from the OAuth provider.
- `use_client_assertion` (Boolean) Use private key JWT (client assertion) instead of client secret.
- `user_info_endpoint` (String) The URL where the application retrieves user information from the OAuth provider.

<a id="nestedatt--authentication--oauth--system--github--apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.github.apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--github--native_apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.github.native_apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--github--provider_token_management"></a>
### Nested Schema for `authentication.oauth.system.github.provider_token_management`



<a id="nestedatt--authentication--oauth--system--gitlab"></a>
### Nested Schema for `authentication.oauth.system.gitlab`

Optional:

- `allowed_grant_types` (List of String) The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens from the OAuth provider.
- `apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic apple client secret for applications. (see [below for nested schema](#nestedatt--authentication--oauth--system--gitlab--apple_key_generator))
- `authorization_endpoint` (String) The URL that users are redirected to for authorization with the OAuth provider.
- `callback_domain` (String) Use a custom domain in your OAuth verification screen.
- `claim_mapping` (Map of String) Maps OAuth provider claims to Descope user attributes.
- `client_id` (String) The client ID for the OAuth provider, used to identify the application to the provider.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider, used to authenticate the application with the provider.
- `description` (String) A brief description of the OAuth provider.
- `disable_jit_updates` (Boolean) By default the user attribute mapping configuration is used to update the user's attributes automatically during sign in. Disable this if you want this to happen only during user creation.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `issuer` (String) The issuer identifier for the OAuth provider.
- `jwks_endpoint` (String) The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.
- `logo` (String) The URL of the logo associated with the OAuth provider.
- `manage_provider_tokens` (Boolean) Whether to enable provider token management for this OAuth provider.
- `merge_user_accounts` (Boolean) Whether to merge existing user accounts with new ones created through OAuth authentication.
- `native_apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic native apple client secret for mobile apps. (see [below for nested schema](#nestedatt--authentication--oauth--system--gitlab--native_apple_key_generator))
- `native_client_id` (String) The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.
- `native_client_secret` (String, Sensitive) The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.
- `prompts` (List of String) Custom prompts or consent screens that users may see during OAuth authentication.
- `provider_token_management` (Attributes) This attribute is deprecated, use the `manage_provider_tokens`, `callback_domain`, and `redirect_url` fields instead. (see [below for nested schema](#nestedatt--authentication--oauth--system--gitlab--provider_token_management))
- `redirect_url` (String) Users will be directed to this URL after authentication. If redirect URL is specified in the SDK/API call, it will override this value.
- `scopes` (List of String) Scopes of access that the application requests from the user's account on the OAuth provider.
- `token_endpoint` (String) The URL where the application requests an access token from the OAuth provider.
- `use_client_assertion` (Boolean) Use private key JWT (client assertion) instead of client secret.
- `user_info_endpoint` (String) The URL where the application retrieves user information from the OAuth provider.

<a id="nestedatt--authentication--oauth--system--gitlab--apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.gitlab.apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--gitlab--native_apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.gitlab.native_apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--gitlab--provider_token_management"></a>
### Nested Schema for `authentication.oauth.system.gitlab.provider_token_management`



<a id="nestedatt--authentication--oauth--system--google"></a>
### Nested Schema for `authentication.oauth.system.google`

Optional:

- `allowed_grant_types` (List of String) The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens from the OAuth provider.
- `apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic apple client secret for applications. (see [below for nested schema](#nestedatt--authentication--oauth--system--google--apple_key_generator))
- `authorization_endpoint` (String) The URL that users are redirected to for authorization with the OAuth provider.
- `callback_domain` (String) Use a custom domain in your OAuth verification screen.
- `claim_mapping` (Map of String) Maps OAuth provider claims to Descope user attributes.
- `client_id` (String) The client ID for the OAuth provider, used to identify the application to the provider.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider, used to authenticate the application with the provider.
- `description` (String) A brief description of the OAuth provider.
- `disable_jit_updates` (Boolean) By default the user attribute mapping configuration is used to update the user's attributes automatically during sign in. Disable this if you want this to happen only during user creation.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `issuer` (String) The issuer identifier for the OAuth provider.
- `jwks_endpoint` (String) The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.
- `logo` (String) The URL of the logo associated with the OAuth provider.
- `manage_provider_tokens` (Boolean) Whether to enable provider token management for this OAuth provider.
- `merge_user_accounts` (Boolean) Whether to merge existing user accounts with new ones created through OAuth authentication.
- `native_apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic native apple client secret for mobile apps. (see [below for nested schema](#nestedatt--authentication--oauth--system--google--native_apple_key_generator))
- `native_client_id` (String) The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.
- `native_client_secret` (String, Sensitive) The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.
- `prompts` (List of String) Custom prompts or consent screens that users may see during OAuth authentication.
- `provider_token_management` (Attributes) This attribute is deprecated, use the `manage_provider_tokens`, `callback_domain`, and `redirect_url` fields instead. (see [below for nested schema](#nestedatt--authentication--oauth--system--google--provider_token_management))
- `redirect_url` (String) Users will be directed to this URL after authentication. If redirect URL is specified in the SDK/API call, it will override this value.
- `scopes` (List of String) Scopes of access that the application requests from the user's account on the OAuth provider.
- `token_endpoint` (String) The URL where the application requests an access token from the OAuth provider.
- `use_client_assertion` (Boolean) Use private key JWT (client assertion) instead of client secret.
- `user_info_endpoint` (String) The URL where the application retrieves user information from the OAuth provider.

<a id="nestedatt--authentication--oauth--system--google--apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.google.apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--google--native_apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.google.native_apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--google--provider_token_management"></a>
### Nested Schema for `authentication.oauth.system.google.provider_token_management`



<a id="nestedatt--authentication--oauth--system--linkedin"></a>
### Nested Schema for `authentication.oauth.system.linkedin`

Optional:

- `allowed_grant_types` (List of String) The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens from the OAuth provider.
- `apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic apple client secret for applications. (see [below for nested schema](#nestedatt--authentication--oauth--system--linkedin--apple_key_generator))
- `authorization_endpoint` (String) The URL that users are redirected to for authorization with the OAuth provider.
- `callback_domain` (String) Use a custom domain in your OAuth verification screen.
- `claim_mapping` (Map of String) Maps OAuth provider claims to Descope user attributes.
- `client_id` (String) The client ID for the OAuth provider, used to identify the application to the provider.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider, used to authenticate the application with the provider.
- `description` (String) A brief description of the OAuth provider.
- `disable_jit_updates` (Boolean) By default the user attribute mapping configuration is used to update the user's attributes automatically during sign in. Disable this if you want this to happen only during user creation.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `issuer` (String) The issuer identifier for the OAuth provider.
- `jwks_endpoint` (String) The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.
- `logo` (String) The URL of the logo associated with the OAuth provider.
- `manage_provider_tokens` (Boolean) Whether to enable provider token management for this OAuth provider.
- `merge_user_accounts` (Boolean) Whether to merge existing user accounts with new ones created through OAuth authentication.
- `native_apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic native apple client secret for mobile apps. (see [below for nested schema](#nestedatt--authentication--oauth--system--linkedin--native_apple_key_generator))
- `native_client_id` (String) The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.
- `native_client_secret` (String, Sensitive) The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.
- `prompts` (List of String) Custom prompts or consent screens that users may see during OAuth authentication.
- `provider_token_management` (Attributes) This attribute is deprecated, use the `manage_provider_tokens`, `callback_domain`, and `redirect_url` fields instead. (see [below for nested schema](#nestedatt--authentication--oauth--system--linkedin--provider_token_management))
- `redirect_url` (String) Users will be directed to this URL after authentication. If redirect URL is specified in the SDK/API call, it will override this value.
- `scopes` (List of String) Scopes of access that the application requests from the user's account on the OAuth provider.
- `token_endpoint` (String) The URL where the application requests an access token from the OAuth provider.
- `use_client_assertion` (Boolean) Use private key JWT (client assertion) instead of client secret.
- `user_info_endpoint` (String) The URL where the application retrieves user information from the OAuth provider.

<a id="nestedatt--authentication--oauth--system--linkedin--apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.linkedin.apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--linkedin--native_apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.linkedin.native_apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--linkedin--provider_token_management"></a>
### Nested Schema for `authentication.oauth.system.linkedin.provider_token_management`



<a id="nestedatt--authentication--oauth--system--microsoft"></a>
### Nested Schema for `authentication.oauth.system.microsoft`

Optional:

- `allowed_grant_types` (List of String) The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens from the OAuth provider.
- `apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic apple client secret for applications. (see [below for nested schema](#nestedatt--authentication--oauth--system--microsoft--apple_key_generator))
- `authorization_endpoint` (String) The URL that users are redirected to for authorization with the OAuth provider.
- `callback_domain` (String) Use a custom domain in your OAuth verification screen.
- `claim_mapping` (Map of String) Maps OAuth provider claims to Descope user attributes.
- `client_id` (String) The client ID for the OAuth provider, used to identify the application to the provider.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider, used to authenticate the application with the provider.
- `description` (String) A brief description of the OAuth provider.
- `disable_jit_updates` (Boolean) By default the user attribute mapping configuration is used to update the user's attributes automatically during sign in. Disable this if you want this to happen only during user creation.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `issuer` (String) The issuer identifier for the OAuth provider.
- `jwks_endpoint` (String) The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.
- `logo` (String) The URL of the logo associated with the OAuth provider.
- `manage_provider_tokens` (Boolean) Whether to enable provider token management for this OAuth provider.
- `merge_user_accounts` (Boolean) Whether to merge existing user accounts with new ones created through OAuth authentication.
- `native_apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic native apple client secret for mobile apps. (see [below for nested schema](#nestedatt--authentication--oauth--system--microsoft--native_apple_key_generator))
- `native_client_id` (String) The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.
- `native_client_secret` (String, Sensitive) The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.
- `prompts` (List of String) Custom prompts or consent screens that users may see during OAuth authentication.
- `provider_token_management` (Attributes) This attribute is deprecated, use the `manage_provider_tokens`, `callback_domain`, and `redirect_url` fields instead. (see [below for nested schema](#nestedatt--authentication--oauth--system--microsoft--provider_token_management))
- `redirect_url` (String) Users will be directed to this URL after authentication. If redirect URL is specified in the SDK/API call, it will override this value.
- `scopes` (List of String) Scopes of access that the application requests from the user's account on the OAuth provider.
- `token_endpoint` (String) The URL where the application requests an access token from the OAuth provider.
- `use_client_assertion` (Boolean) Use private key JWT (client assertion) instead of client secret.
- `user_info_endpoint` (String) The URL where the application retrieves user information from the OAuth provider.

<a id="nestedatt--authentication--oauth--system--microsoft--apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.microsoft.apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--microsoft--native_apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.microsoft.native_apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--microsoft--provider_token_management"></a>
### Nested Schema for `authentication.oauth.system.microsoft.provider_token_management`



<a id="nestedatt--authentication--oauth--system--slack"></a>
### Nested Schema for `authentication.oauth.system.slack`

Optional:

- `allowed_grant_types` (List of String) The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens from the OAuth provider.
- `apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic apple client secret for applications. (see [below for nested schema](#nestedatt--authentication--oauth--system--slack--apple_key_generator))
- `authorization_endpoint` (String) The URL that users are redirected to for authorization with the OAuth provider.
- `callback_domain` (String) Use a custom domain in your OAuth verification screen.
- `claim_mapping` (Map of String) Maps OAuth provider claims to Descope user attributes.
- `client_id` (String) The client ID for the OAuth provider, used to identify the application to the provider.
- `client_secret` (String, Sensitive) The client secret for the OAuth provider, used to authenticate the application with the provider.
- `description` (String) A brief description of the OAuth provider.
- `disable_jit_updates` (Boolean) By default the user attribute mapping configuration is used to update the user's attributes automatically during sign in. Disable this if you want this to happen only during user creation.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `issuer` (String) The issuer identifier for the OAuth provider.
- `jwks_endpoint` (String) The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.
- `logo` (String) The URL of the logo associated with the OAuth provider.
- `manage_provider_tokens` (Boolean) Whether to enable provider token management for this OAuth provider.
- `merge_user_accounts` (Boolean) Whether to merge existing user accounts with new ones created through OAuth authentication.
- `native_apple_key_generator` (Attributes) The apple key generator object describing how to create a dynamic native apple client secret for mobile apps. (see [below for nested schema](#nestedatt--authentication--oauth--system--slack--native_apple_key_generator))
- `native_client_id` (String) The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.
- `native_client_secret` (String, Sensitive) The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.
- `prompts` (List of String) Custom prompts or consent screens that users may see during OAuth authentication.
- `provider_token_management` (Attributes) This attribute is deprecated, use the `manage_provider_tokens`, `callback_domain`, and `redirect_url` fields instead. (see [below for nested schema](#nestedatt--authentication--oauth--system--slack--provider_token_management))
- `redirect_url` (String) Users will be directed to this URL after authentication. If redirect URL is specified in the SDK/API call, it will override this value.
- `scopes` (List of String) Scopes of access that the application requests from the user's account on the OAuth provider.
- `token_endpoint` (String) The URL where the application requests an access token from the OAuth provider.
- `use_client_assertion` (Boolean) Use private key JWT (client assertion) instead of client secret.
- `user_info_endpoint` (String) The URL where the application retrieves user information from the OAuth provider.

<a id="nestedatt--authentication--oauth--system--slack--apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.slack.apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--slack--native_apple_key_generator"></a>
### Nested Schema for `authentication.oauth.system.slack.native_apple_key_generator`

Required:

- `key_id` (String) The apple generator key id produced by Apple.
- `private_key` (String, Sensitive) The apple generator private key produced by Apple.
- `team_id` (String) The apple generator team id assigned to the key by Apple.


<a id="nestedatt--authentication--oauth--system--slack--provider_token_management"></a>
### Nested Schema for `authentication.oauth.system.slack.provider_token_management`





<a id="nestedatt--authentication--otp"></a>
### Nested Schema for `authentication.otp`

Optional:

- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `domain` (String) The domain to embed in OTP messages.
- `email_service` (Attributes) Settings related to sending emails with OTP codes. (see [below for nested schema](#nestedatt--authentication--otp--email_service))
- `expiration_time` (String) The amount of time that an OTP code will be valid for.
- `text_service` (Attributes) Settings related to sending SMS messages with OTP codes. (see [below for nested schema](#nestedatt--authentication--otp--text_service))
- `voice_service` (Attributes) Settings related to voice calls with OTP codes. (see [below for nested schema](#nestedatt--authentication--otp--voice_service))

<a id="nestedatt--authentication--otp--email_service"></a>
### Nested Schema for `authentication.otp.email_service`

Required:

- `connector` (String) The name of the email connector to use for sending emails.

Optional:

- `templates` (Attributes List) A list of email templates for different authentication flows. (see [below for nested schema](#nestedatt--authentication--otp--email_service--templates))

<a id="nestedatt--authentication--otp--email_service--templates"></a>
### Nested Schema for `authentication.otp.email_service.templates`

Required:

- `name` (String) Unique name for this email template.
- `subject` (String) Subject line of the email message.

Optional:

- `active` (Boolean) Whether this email template is currently active and in use.
- `html_body` (String) HTML content of the email message body, required if `use_plain_text_body` isn't set.
- `plain_text_body` (String) Plain text version of the email message body, required if `use_plain_text_body` is set to `true`.
- `use_plain_text_body` (Boolean) Whether to use the plain text body instead of HTML for the email.

Read-Only:

- `id` (String)



<a id="nestedatt--authentication--otp--text_service"></a>
### Nested Schema for `authentication.otp.text_service`

Required:

- `connector` (String) The name of the SMS/text connector to use for sending text messages.

Optional:

- `templates` (Attributes List) A list of text message templates for different authentication flows. (see [below for nested schema](#nestedatt--authentication--otp--text_service--templates))

<a id="nestedatt--authentication--otp--text_service--templates"></a>
### Nested Schema for `authentication.otp.text_service.templates`

Required:

- `body` (String) The content of the text message.
- `name` (String) Unique name for this text template.

Optional:

- `active` (Boolean) Whether this text template is currently active and in use.

Read-Only:

- `id` (String)



<a id="nestedatt--authentication--otp--voice_service"></a>
### Nested Schema for `authentication.otp.voice_service`

Required:

- `connector` (String) The name of the voice connector to use for making voice calls.

Optional:

- `templates` (Attributes List) A list of voice message templates for different purposes. (see [below for nested schema](#nestedatt--authentication--otp--voice_service--templates))

<a id="nestedatt--authentication--otp--voice_service--templates"></a>
### Nested Schema for `authentication.otp.voice_service.templates`

Required:

- `body` (String) The content of the voice message that will be spoken.
- `name` (String) Unique name for this voice template.

Optional:

- `active` (Boolean) Whether this voice template is currently active and in use.

Read-Only:

- `id` (String)




<a id="nestedatt--authentication--passkeys"></a>
### Nested Schema for `authentication.passkeys`

Optional:

- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `top_level_domain` (String) Passkeys will be usable in the following domain and all its subdomains.


<a id="nestedatt--authentication--password"></a>
### Nested Schema for `authentication.password`

Optional:

- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `email_service` (Attributes) Settings related to sending password reset emails as part of the password feature. (see [below for nested schema](#nestedatt--authentication--password--email_service))
- `expiration` (Boolean) Whether users are required to change their password periodically.
- `expiration_weeks` (Number) The number of weeks after which a user's password expires and they need to replace it.
- `lock` (Boolean) Whether the user account should be locked after a specified number of failed login attempts.
- `lock_attempts` (Number) The number of failed login attempts allowed before an account is locked.
- `lowercase` (Boolean) Whether passwords must contain at least one lowercase letter.
- `mask_errors` (Boolean) Prevents information about user accounts from being revealed in error messages, e.g., whether a user already exists.
- `min_length` (Number) The minimum length of the password that users are required to use. The maximum length is always `64`.
- `non_alphanumeric` (Boolean) Whether passwords must contain at least one non-alphanumeric character (e.g. `!`, `@`, `#`).
- `number` (Boolean) Whether passwords must contain at least one number.
- `reuse` (Boolean) Whether to forbid password reuse when users change their password.
- `reuse_amount` (Number) The number of previous passwords whose hashes are kept to prevent users from reusing old passwords.
- `temporary_lock` (Boolean) Whether the user account should be temporarily locked after a specified number of failed login attempts.
- `temporary_lock_attempts` (Number) The number of failed login attempts allowed before an account is temporarily locked.
- `temporary_lock_duration` (String) The amount of time before the user can sign in again after the account is temporarily locked.
- `uppercase` (Boolean) Whether passwords must contain at least one uppercase letter.

<a id="nestedatt--authentication--password--email_service"></a>
### Nested Schema for `authentication.password.email_service`

Required:

- `connector` (String) The name of the email connector to use for sending emails.

Optional:

- `templates` (Attributes List) A list of email templates for different authentication flows. (see [below for nested schema](#nestedatt--authentication--password--email_service--templates))

<a id="nestedatt--authentication--password--email_service--templates"></a>
### Nested Schema for `authentication.password.email_service.templates`

Required:

- `name` (String) Unique name for this email template.
- `subject` (String) Subject line of the email message.

Optional:

- `active` (Boolean) Whether this email template is currently active and in use.
- `html_body` (String) HTML content of the email message body, required if `use_plain_text_body` isn't set.
- `plain_text_body` (String) Plain text version of the email message body, required if `use_plain_text_body` is set to `true`.
- `use_plain_text_body` (Boolean) Whether to use the plain text body instead of HTML for the email.

Read-Only:

- `id` (String)




<a id="nestedatt--authentication--sso"></a>
### Nested Schema for `authentication.sso`

Optional:

- `allow_duplicate_domains` (Boolean) Whether to allow duplicate SSO domains across tenants.
- `allow_override_roles` (Boolean) Whether to allow overriding user's roles with SSO related roles.
- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `groups_priority` (Boolean) Whether to enable groups priority.
- `merge_users` (Boolean) Whether to merge existing user accounts with new ones created through SSO authentication.
- `redirect_url` (String) The URL the end user is redirected to after a successful authentication. If one is specified in tenant level settings or SDK/API call, they will override this value.
- `sso_suite_settings` (Attributes) Configuration block for the SSO Suite. (see [below for nested schema](#nestedatt--authentication--sso--sso_suite_settings))

<a id="nestedatt--authentication--sso--sso_suite_settings"></a>
### Nested Schema for `authentication.sso.sso_suite_settings`

Optional:

- `force_domain_verification` (Boolean) Setting this to `true` will allow only verified domains to be used.
- `hide_domains` (Boolean) Setting this to `true` will hide the domains configuration section in the SSO Suite interface.
- `hide_groups_mapping` (Boolean) Setting this to `true` will hide the groups mapping configuration section in the SSO Suite interface.
- `hide_oidc` (Boolean) Setting this to `true` will hide the OIDC configuration option.
- `hide_saml` (Boolean) Setting this to `true` will hide the SAML configuration option.
- `hide_scim` (Boolean) Setting this to `true` will hide the SCIM configuration in the SSO Suite interface.
- `style_id` (String) Specifies the style ID to apply in the SSO Suite. Ensure a style with this ID exists in the console for it to be used.



<a id="nestedatt--authentication--totp"></a>
### Nested Schema for `authentication.totp`

Optional:

- `disabled` (Boolean) Setting this to `true` will disallow using this authentication method directly via API and SDK calls. Note that this does not affect authentication flows that are configured to use this authentication method.
- `service_label` (String) The template for the service issuer label (issuer) shown in the authenticator app.



<a id="nestedatt--authorization"></a>
### Nested Schema for `authorization`

Optional:

- `permissions` (Attributes List) A list of `Permission` objects. (see [below for nested schema](#nestedatt--authorization--permissions))
- `roles` (Attributes List) A list of `Role` objects. (see [below for nested schema](#nestedatt--authorization--roles))

<a id="nestedatt--authorization--permissions"></a>
### Nested Schema for `authorization.permissions`

Required:

- `name` (String) A name for the permission.

Optional:

- `description` (String) A description for the permission.

Read-Only:

- `id` (String)


<a id="nestedatt--authorization--roles"></a>
### Nested Schema for `authorization.roles`

Required:

- `name` (String) A name for the role.

Optional:

- `default` (Boolean) Whether this role should automatically be assigned to users that are created without any roles.
- `description` (String) A description for the role.
- `key` (String) A persistent value that identifies a role uniquely across plan changes and configuration updates. It is used exclusively by the Terraform provider during planning, to ensure that user roles are maintained consistently even when role names or other details are changed. Once the `key` is set it should never be changed, otherwise the role will be removed and a new one will be created instead.
- `permissions` (Set of String) A list of permissions by name to be included in the role.
- `private` (Boolean) Whether this role should not be displayed to tenant admins.

Read-Only:

- `id` (String)



<a id="nestedatt--connectors"></a>
### Nested Schema for `connectors`

Optional:

- `abuseipdb` (Attributes List) Utilize IP threat intelligence to block malicious login attempts with the AbuseIPDB connector. (see [below for nested schema](#nestedatt--connectors--abuseipdb))
- `amplitude` (Attributes List) Track user activity and traits at any point in your user journey with the Amplitude connector. (see [below for nested schema](#nestedatt--connectors--amplitude))
- `arkose` (Attributes List) Use the Arkose connector to integrate with Arkose's bot and fraud detection. (see [below for nested schema](#nestedatt--connectors--arkose))
- `audit_webhook` (Attributes List) Send audit events to a custom webhook. (see [below for nested schema](#nestedatt--connectors--audit_webhook))
- `aws_s3` (Attributes List) Stream authentication audit logs with the Amazon S3 connector. (see [below for nested schema](#nestedatt--connectors--aws_s3))
- `aws_translate` (Attributes List) Localize the language of your login and user journey screens with the Amazon Translate connector. (see [below for nested schema](#nestedatt--connectors--aws_translate))
- `bitsight` (Attributes List) Utilize threat intelligence to block malicious login attempts or check leaks with the Bitsight Threat Intelligence connector. (see [below for nested schema](#nestedatt--connectors--bitsight))
- `coralogix` (Attributes List) Send audit events and troubleshooting logs to Coralogix. (see [below for nested schema](#nestedatt--connectors--coralogix))
- `darwinium` (Attributes List) Connect to Darwinium API for fraud detection and device intelligence. (see [below for nested schema](#nestedatt--connectors--darwinium))
- `datadog` (Attributes List) Stream authentication audit logs with the Datadog connector. (see [below for nested schema](#nestedatt--connectors--datadog))
- `devrev_grow` (Attributes List) DevRev Grow is a Growth CRM that brings salespeople, product marketers, and PMs onto an AI-native platform to follow the journey of a visitor to a lead, to a contact, and then to a user - to create a champion, not a churned user. (see [below for nested schema](#nestedatt--connectors--devrev_grow))
- `docebo` (Attributes List) Get user information from Docebo in your Descope user journeys with the Docebo connector. (see [below for nested schema](#nestedatt--connectors--docebo))
- `eight_by_eight_viber` (Attributes List) Send Viber messages to the user. (see [below for nested schema](#nestedatt--connectors--eight_by_eight_viber))
- `eight_by_eight_whatsapp` (Attributes List) Send WhatsApp messages to the user. (see [below for nested schema](#nestedatt--connectors--eight_by_eight_whatsapp))
- `elephant` (Attributes List) Use this connector to obtain an identity trust score. (see [below for nested schema](#nestedatt--connectors--elephant))
- `external_token_http` (Attributes List) A generic HTTP token connector. (see [below for nested schema](#nestedatt--connectors--external_token_http))
- `fingerprint` (Attributes List) Prevent fraud by adding device intelligence with the Fingerprint connector. (see [below for nested schema](#nestedatt--connectors--fingerprint))
- `fingerprint_descope` (Attributes List) Descope Fingerprint capabilities for fraud detection and risk assessment. (see [below for nested schema](#nestedatt--connectors--fingerprint_descope))
- `firebase_admin` (Attributes List) Firebase connector enables you to utilize Firebase's APIs to generate a Firebase user token for a given Descope user. (see [below for nested schema](#nestedatt--connectors--firebase_admin))
- `forter` (Attributes List) Leverage ML-based risk scores for fraud prevention with the Forter connector. (see [below for nested schema](#nestedatt--connectors--forter))
- `generic_email_gateway` (Attributes List) Send emails using a generic Email gateway. (see [below for nested schema](#nestedatt--connectors--generic_email_gateway))
- `generic_sms_gateway` (Attributes List) Send messages using a generic SMS gateway. (see [below for nested schema](#nestedatt--connectors--generic_sms_gateway))
- `google_cloud_logging` (Attributes List) Stream logs and audit events with the Google Cloud Logging connector. (see [below for nested schema](#nestedatt--connectors--google_cloud_logging))
- `google_cloud_translation` (Attributes List) Localize the language of your login and user journey screens with the Google Cloud Translation connector. (see [below for nested schema](#nestedatt--connectors--google_cloud_translation))
- `google_maps_places` (Attributes List) Get address autocompletions from Place Autocomplete Data API. (see [below for nested schema](#nestedatt--connectors--google_maps_places))
- `hcaptcha` (Attributes List) hCaptcha can help protect your applications from bots, spam, and other forms of automated abuse. (see [below for nested schema](#nestedatt--connectors--hcaptcha))
- `hibp` (Attributes List) Check if passwords have been previously exposed in data breaches with the Have I Been Pwned connector. (see [below for nested schema](#nestedatt--connectors--hibp))
- `http` (Attributes List) A general purpose HTTP client (see [below for nested schema](#nestedatt--connectors--http))
- `hubspot` (Attributes List) Orchestrate customer identity information from your Descope user journey with the HubSpot connector. (see [below for nested schema](#nestedatt--connectors--hubspot))
- `incode` (Attributes List) Use the Incode connection to run identity verification processes like document checks or facial recognition. (see [below for nested schema](#nestedatt--connectors--incode))
- `intercom` (Attributes List) Orchestrate customer identity information from your Descope user journey with the Intercom connector. (see [below for nested schema](#nestedatt--connectors--intercom))
- `ldap` (Attributes List) Use this connector to authenticate users against an LDAP directory server with support for both password and mTLS authentication. (see [below for nested schema](#nestedatt--connectors--ldap))
- `lokalise` (Attributes List) Localize the language of your login and user journey screens with the Lokalise connector. (see [below for nested schema](#nestedatt--connectors--lokalise))
- `mixpanel` (Attributes List) Stream authentication audit logs and troubleshoot logs to Mixpanel. (see [below for nested schema](#nestedatt--connectors--mixpanel))
- `mparticle` (Attributes List) Track and send user event data (e.g. page views, purchases, etc.) across connected tools using the mParticle connector. (see [below for nested schema](#nestedatt--connectors--mparticle))
- `newrelic` (Attributes List) Stream authentication audit logs with the New Relic connector. (see [below for nested schema](#nestedatt--connectors--newrelic))
- `opentelemetry` (Attributes List) Send audit events and troubleshooting logs to an OpenTelemetry-compatible endpoint using OTLP over HTTP or gRPC. (see [below for nested schema](#nestedatt--connectors--opentelemetry))
- `ping_directory` (Attributes List) Authenticate against PingDirectory. (see [below for nested schema](#nestedatt--connectors--ping_directory))
- `postmark` (Attributes List) Send emails using Postmark (see [below for nested schema](#nestedatt--connectors--postmark))
- `radar` (Attributes List) Get address autocompletions from Radar Autocomplete API. (see [below for nested schema](#nestedatt--connectors--radar))
- `recaptcha` (Attributes List) Prevent bot attacks on your login pages with the reCAPTCHA v3 connector. (see [below for nested schema](#nestedatt--connectors--recaptcha))
- `recaptcha_enterprise` (Attributes List) Mitigate fraud using advanced risk analysis and add adaptive MFA with the reCAPTCHA Enterprise connector. (see [below for nested schema](#nestedatt--connectors--recaptcha_enterprise))
- `rekognition` (Attributes List) Add image recognition capabilities for identity verification and fraud prevention with the Amazon Rekognition connector. (see [below for nested schema](#nestedatt--connectors--rekognition))
- `salesforce` (Attributes List) Run SQL queries to retrieve user roles, profiles, account status, and more with the Salesforce connector. (see [below for nested schema](#nestedatt--connectors--salesforce))
- `salesforce_marketing_cloud` (Attributes List) Send transactional messages with the Salesforce Marketing Cloud connector. (see [below for nested schema](#nestedatt--connectors--salesforce_marketing_cloud))
- `sardine` (Attributes List) Evaluate customer risk using Sardine (see [below for nested schema](#nestedatt--connectors--sardine))
- `segment` (Attributes List) Orchestrate customer identity traits and signals from your Descope user journey with the Segment connector. (see [below for nested schema](#nestedatt--connectors--segment))
- `sendgrid` (Attributes List) SendGrid is a cloud-based SMTP provider that allows you to send emails without having to maintain email servers. (see [below for nested schema](#nestedatt--connectors--sendgrid))
- `ses` (Attributes List) Amazon Simple Email Service (SES) for sending emails through AWS infrastructure. (see [below for nested schema](#nestedatt--connectors--ses))
- `slack` (Attributes List) Send updates to your team on Slack. (see [below for nested schema](#nestedatt--connectors--slack))
- `smartling` (Attributes List) Localize the language of your login and user journey screens with the Smartling connector. (see [below for nested schema](#nestedatt--connectors--smartling))
- `smtp` (Attributes List) Simple Mail Transfer Protocol (SMTP) server for sending emails. (see [below for nested schema](#nestedatt--connectors--smtp))
- `sns` (Attributes List) Amazon Simple Notification Service (SNS) for sending SMS messages through AWS. (see [below for nested schema](#nestedatt--connectors--sns))
- `splunk` (Attributes List) Stream logs and audit events with the Splunk HTTP Event Collector (HEC). (see [below for nested schema](#nestedatt--connectors--splunk))
- `sql` (Attributes List) SQL connector for relational databases including PostgreSQL, MySQL, MariaDB, Microsoft SQL Server (MSSQL), Oracle, CockroachDB, and Amazon Redshift. (see [below for nested schema](#nestedatt--connectors--sql))
- `sumologic` (Attributes List) Stream logs and audit events with the Sumo Logic connector. (see [below for nested schema](#nestedatt--connectors--sumologic))
- `supabase` (Attributes List) Generate external tokens for user authentication in Supabase projects. (see [below for nested schema](#nestedatt--connectors--supabase))
- `telesign` (Attributes List) Verify phone numbers and leverage granular risk scores for adaptive MFA with the Telesign Intelligence connector. (see [below for nested schema](#nestedatt--connectors--telesign))
- `traceable` (Attributes List) Identify and respond to fraudulent login activity with the Traceable Digital Fraud Prevention connector. (see [below for nested schema](#nestedatt--connectors--traceable))
- `turnstile` (Attributes List) Prevent bot attacks on your login pages with the Turnstile connector. (see [below for nested schema](#nestedatt--connectors--turnstile))
- `twilio_core` (Attributes List) Twilio is a cloud-based communication provider of communication tools for making and receiving phone calls, sending and receiving text messages, and performing other communication functions. (see [below for nested schema](#nestedatt--connectors--twilio_core))
- `twilio_verify` (Attributes List) Twilio Verify is an OTP service that can be used via text messages, instant messaging platforms, voice and e-mail. Choose this connector only if you are a Twilio Verify customer. (see [below for nested schema](#nestedatt--connectors--twilio_verify))
- `unibeam` (Attributes List) SIM-based authentication and approval using Unibeam's OnSim technology for passwordless authentication and transaction approval. (see [below for nested schema](#nestedatt--connectors--unibeam))
- `zerobounce` (Attributes List) Email validation with ZeroBounce (see [below for nested schema](#nestedatt--connectors--zerobounce))

<a id="nestedatt--connectors--abuseipdb"></a>
### Nested Schema for `connectors.abuseipdb`

Required:

- `api_key` (String, Sensitive) The unique AbuseIPDB API key.
- `name` (String) A custom name for your connector.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--amplitude"></a>
### Nested Schema for `connectors.amplitude`

Required:

- `api_key` (String, Sensitive) The Amplitude API Key generated for the Descope service.
- `name` (String) A custom name for your connector.

Optional:

- `description` (String) A description of what your connector is used for.
- `server_url` (String) The server URL of the Amplitude API, when using different api or a custom domain in Amplitude.
- `server_zone` (String) `EU` or `US`. Sets the Amplitude server zone. Set this to `EU` for Amplitude projects created in `EU` data center. Default is `US`.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--arkose"></a>
### Nested Schema for `connectors.arkose`

Required:

- `name` (String) A custom name for your connector.
- `private_key` (String, Sensitive) The private key that can be copied from the Keys screen in the Arkose Labs portal.
- `public_key` (String) The public key that's shown in the Keys screen in the Arkose Labs portal.

Optional:

- `client_base_url` (String) A custom base URL to use when loading the Arkose Labs client script. If not provided, the default value of `https://client-api.arkoselabs.com/v2` will be used.
- `description` (String) A description of what your connector is used for.
- `verify_base_url` (String) A custom base URL to use when verifying the session token using the Arkose Labs Verify API. If not provided, the default value of `https://verify-api.arkoselabs.com/api/v4` will be used.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--audit_webhook"></a>
### Nested Schema for `connectors.audit_webhook`

Required:

- `base_url` (String) The base URL to fetch
- `name` (String) A custom name for your connector.

Optional:

- `audit_filters` (Attributes List) Specify which events will be sent to the external audit service (including tenant selection). (see [below for nested schema](#nestedatt--connectors--audit_webhook--audit_filters))
- `authentication` (Attributes) Authentication Information (see [below for nested schema](#nestedatt--connectors--audit_webhook--authentication))
- `description` (String) A description of what your connector is used for.
- `headers` (Map of String) The headers to send with the request
- `hmac_secret` (String, Sensitive) HMAC is a method for message signing with a symmetrical key. This secret will be used to sign the payload, and the resulting signature will be sent in the `x-descope-webhook-s256` header. The receiving service should use this secret to verify the integrity and authenticity of the payload by checking the provided signature
- `insecure` (Boolean) Will ignore certificate errors raised by the client

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--audit_webhook--audit_filters"></a>
### Nested Schema for `connectors.audit_webhook.audit_filters`

Required:

- `key` (String) The field name to filter on (either 'actions' or 'tenants').
- `operator` (String) The filter operation to apply ('includes' or 'excludes').
- `values` (List of String) The list of values to match against for the filter.


<a id="nestedatt--connectors--audit_webhook--authentication"></a>
### Nested Schema for `connectors.audit_webhook.authentication`

Optional:

- `api_key` (Attributes) API key authentication configuration. (see [below for nested schema](#nestedatt--connectors--audit_webhook--authentication--api_key))
- `basic` (Attributes) Basic authentication credentials (username and password). (see [below for nested schema](#nestedatt--connectors--audit_webhook--authentication--basic))
- `bearer_token` (String, Sensitive) Bearer token for HTTP authentication.

<a id="nestedatt--connectors--audit_webhook--authentication--api_key"></a>
### Nested Schema for `connectors.audit_webhook.authentication.api_key`

Required:

- `key` (String) The API key.
- `token` (String, Sensitive) The API secret.


<a id="nestedatt--connectors--audit_webhook--authentication--basic"></a>
### Nested Schema for `connectors.audit_webhook.authentication.basic`

Required:

- `password` (String, Sensitive) Password for basic HTTP authentication.
- `username` (String) Username for basic HTTP authentication.




<a id="nestedatt--connectors--aws_s3"></a>
### Nested Schema for `connectors.aws_s3`

Required:

- `bucket` (String) The AWS S3 bucket. This bucket should already exist for the connector to work.
- `name` (String) A custom name for your connector.
- `region` (String) The AWS S3 region, e.g. `us-east-1`.

Optional:

- `access_key_id` (String, Sensitive) The unique AWS access key ID.
- `audit_enabled` (Boolean) Whether to enable streaming of audit events.
- `audit_filters` (Attributes List) Specify which events will be sent to the external audit service (including tenant selection). (see [below for nested schema](#nestedatt--connectors--aws_s3--audit_filters))
- `auth_type` (String) The authentication type to use.
- `description` (String) A description of what your connector is used for.
- `external_id` (String) The external ID to use when assuming the role.
- `role_arn` (String) The Amazon Resource Name (ARN) of the role to assume.
- `secret_access_key` (String, Sensitive) The secret AWS access key.
- `troubleshoot_log_enabled` (Boolean) Whether to send troubleshooting events.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--aws_s3--audit_filters"></a>
### Nested Schema for `connectors.aws_s3.audit_filters`

Required:

- `key` (String) The field name to filter on (either 'actions' or 'tenants').
- `operator` (String) The filter operation to apply ('includes' or 'excludes').
- `values` (List of String) The list of values to match against for the filter.



<a id="nestedatt--connectors--aws_translate"></a>
### Nested Schema for `connectors.aws_translate`

Required:

- `access_key_id` (String) AWS access key ID.
- `name` (String) A custom name for your connector.
- `region` (String) The AWS region to which this client will send requests. (e.g. us-east-1.)
- `secret_access_key` (String, Sensitive) AWS secret access key.

Optional:

- `description` (String) A description of what your connector is used for.
- `session_token` (String, Sensitive) (Optional) A security or session token to use with these credentials. Usually present for temporary credentials.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--bitsight"></a>
### Nested Schema for `connectors.bitsight`

Required:

- `client_id` (String) API Client ID issued when you create the credentials in Bitsight Threat Intelligence.
- `client_secret` (String, Sensitive) Client secret issued when you create the credentials in Bitsight Threat Intelligence.
- `name` (String) A custom name for your connector.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--coralogix"></a>
### Nested Schema for `connectors.coralogix`

Required:

- `bearer_token` (String, Sensitive) Bearer token issued by Coralogix as Send-Your-Data API key
- `endpoint` (String) The ingress OpenTelemetry endpoint URL.
- `name` (String) A custom name for your connector.

Optional:

- `audit_enabled` (Boolean) Whether to enable streaming of audit events.
- `audit_filters` (Attributes List) Specify which events will be sent to the external audit service (including tenant selection). (see [below for nested schema](#nestedatt--connectors--coralogix--audit_filters))
- `description` (String) A description of what your connector is used for.
- `troubleshoot_log_enabled` (Boolean) Whether to send troubleshooting events.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--coralogix--audit_filters"></a>
### Nested Schema for `connectors.coralogix.audit_filters`

Required:

- `key` (String) The field name to filter on (either 'actions' or 'tenants').
- `operator` (String) The filter operation to apply ('includes' or 'excludes').
- `values` (List of String) The list of values to match against for the filter.



<a id="nestedatt--connectors--darwinium"></a>
### Nested Schema for `connectors.darwinium`

Required:

- `journey_name` (String) The name of the Darwinium journey to use for profiling.
- `name` (String) A custom name for your connector.
- `node_name` (String) The name of the Darwinium node.
- `pem_certificate` (String, Sensitive) The PEM certificate for client authentication.
- `private_key` (String, Sensitive) The private key for client authentication.
- `web_api_name` (String) The name of the Darwinium Web API to use.

Optional:

- `default_result` (String) The default result to return if no result is available.
- `description` (String) A description of what your connector is used for.
- `native_api_name` (String) The name of the Darwinium Native Mobile API to use.
- `native_blob_key_name` (String) The key name for the native profiling blob sent via the client parameter. If not provided, the default key of 'nativeProfilingBlob' will be used.
- `passphrase` (String, Sensitive) The passphrase for the PEM certificate, if applicable.
- `profiling_tags_script_url` (String) The custom URL where the Darwinium Tags script is hosted. If not provided, the default Darwinium script URL will be used.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--datadog"></a>
### Nested Schema for `connectors.datadog`

Required:

- `api_key` (String, Sensitive) The unique Datadog organization key.
- `name` (String) A custom name for your connector.

Optional:

- `audit_enabled` (Boolean) Whether to enable streaming of audit events.
- `audit_filters` (Attributes List) Specify which events will be sent to the external audit service (including tenant selection). (see [below for nested schema](#nestedatt--connectors--datadog--audit_filters))
- `description` (String) A description of what your connector is used for.
- `mask_pii` (Boolean) Whether to mask personally identifiable information in the logs.
- `site` (String) The Datadog site to send logs to. Default is `datadoghq.com`. European, free tier and other customers should set their site accordingly.
- `troubleshoot_log_enabled` (Boolean) Whether to send troubleshooting events.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--datadog--audit_filters"></a>
### Nested Schema for `connectors.datadog.audit_filters`

Required:

- `key` (String) The field name to filter on (either 'actions' or 'tenants').
- `operator` (String) The filter operation to apply ('includes' or 'excludes').
- `values` (List of String) The list of values to match against for the filter.



<a id="nestedatt--connectors--devrev_grow"></a>
### Nested Schema for `connectors.devrev_grow`

Required:

- `api_key` (String, Sensitive) Authentication to DevRev APIs requires a personal access token (PAT).
- `name` (String) A custom name for your connector.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--docebo"></a>
### Nested Schema for `connectors.docebo`

Required:

- `base_url` (String) The Docebo api base url.
- `client_id` (String) The Docebo OAuth 2.0 app client ID.
- `client_secret` (String, Sensitive) The Docebo OAuth 2.0 app client secret.
- `name` (String) A custom name for your connector.
- `password` (String, Sensitive) The Docebo user's password.
- `username` (String) The Docebo username.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--eight_by_eight_viber"></a>
### Nested Schema for `connectors.eight_by_eight_viber`

Required:

- `api_key` (String) The 8x8 API key for authentication.
- `name` (String) A custom name for your connector.
- `sub_account_id` (String) The 8x8 sub-account ID is required for the Messaging API.

Optional:

- `country` (String) The country code or region where your Viber messaging service is configured.
- `description` (String) A description of what your connector is used for.
- `use_static_ips` (Boolean) Whether the connector should send all requests from specific static IPs.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--eight_by_eight_whatsapp"></a>
### Nested Schema for `connectors.eight_by_eight_whatsapp`

Required:

- `api_key` (String) The 8x8 API key for authentication.
- `name` (String) A custom name for your connector.
- `sub_account_id` (String) The 8x8 sub-account ID is required for the Messaging API.
- `template_id` (String) The ID of a WhatsApp message template.

Optional:

- `country` (String) The country code or region where your Viber messaging service is configured.
- `description` (String) A description of what your connector is used for.
- `use_static_ips` (Boolean) Whether the connector should send all requests from specific static IPs.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--elephant"></a>
### Nested Schema for `connectors.elephant`

Required:

- `access_key` (String, Sensitive) The Elephant access key.
- `name` (String) A custom name for your connector.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--external_token_http"></a>
### Nested Schema for `connectors.external_token_http`

Required:

- `endpoint` (String) The endpoint to get the token from (Using POST method). Descope will send the user information in the body of the request, and should return a JSON response with a 'token' string field.
- `name` (String) A custom name for your connector.

Optional:

- `authentication` (Attributes) Authentication Information (see [below for nested schema](#nestedatt--connectors--external_token_http--authentication))
- `description` (String) A description of what your connector is used for.
- `headers` (Map of String) The headers to send with the request
- `hmac_secret` (String, Sensitive) HMAC is a method for message signing with a symmetrical key. This secret will be used to sign the base64 encoded payload, and the resulting signature will be sent in the `x-descope-webhook-s256` header. The receiving service should use this secret to verify the integrity and authenticity of the payload by checking the provided signature
- `insecure` (Boolean) Will ignore certificate errors raised by the client
- `use_static_ips` (Boolean) Whether the connector should send all requests from specific static IPs.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--external_token_http--authentication"></a>
### Nested Schema for `connectors.external_token_http.authentication`

Optional:

- `api_key` (Attributes) API key authentication configuration. (see [below for nested schema](#nestedatt--connectors--external_token_http--authentication--api_key))
- `basic` (Attributes) Basic authentication credentials (username and password). (see [below for nested schema](#nestedatt--connectors--external_token_http--authentication--basic))
- `bearer_token` (String, Sensitive) Bearer token for HTTP authentication.

<a id="nestedatt--connectors--external_token_http--authentication--api_key"></a>
### Nested Schema for `connectors.external_token_http.authentication.api_key`

Required:

- `key` (String) The API key.
- `token` (String, Sensitive) The API secret.


<a id="nestedatt--connectors--external_token_http--authentication--basic"></a>
### Nested Schema for `connectors.external_token_http.authentication.basic`

Required:

- `password` (String, Sensitive) Password for basic HTTP authentication.
- `username` (String) Username for basic HTTP authentication.




<a id="nestedatt--connectors--fingerprint"></a>
### Nested Schema for `connectors.fingerprint`

Required:

- `name` (String) A custom name for your connector.
- `public_api_key` (String) The Fingerprint public API key.
- `secret_api_key` (String, Sensitive) The Fingerprint secret API key.

Optional:

- `cloudflare_endpoint_url` (String) The Cloudflare integration Endpoint URL.
- `cloudflare_script_url` (String) The Cloudflare integration Script URL.
- `description` (String) A description of what your connector is used for.
- `use_cloudflare_integration` (Boolean) Enable to configure the relevant Cloudflare integration parameters if Cloudflare integration is set in your Fingerprint account.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--fingerprint_descope"></a>
### Nested Schema for `connectors.fingerprint_descope`

Required:

- `name` (String) A custom name for your connector.

Optional:

- `custom_domain` (String) The custom domain to fetch
- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--firebase_admin"></a>
### Nested Schema for `connectors.firebase_admin`

Required:

- `name` (String) A custom name for your connector.
- `service_account` (String, Sensitive) The Firebase service account JSON.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--forter"></a>
### Nested Schema for `connectors.forter`

Required:

- `name` (String) A custom name for your connector.
- `secret_key` (String, Sensitive) The Forter secret key.
- `site_id` (String) The Forter site ID.

Optional:

- `api_version` (String) The Forter API version.
- `description` (String) A description of what your connector is used for.
- `override_ip_address` (String) Override the user IP address.
- `override_user_email` (String) Override the user email.
- `overrides` (Boolean) Override the user's IP address or email so that Forter can provide a specific decision or recommendation. Contact the Forter team for further details. Note: Overriding the user IP address or email is intended for testing purpose and should not be utilized in production environments.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--generic_email_gateway"></a>
### Nested Schema for `connectors.generic_email_gateway`

Required:

- `name` (String) A custom name for your connector.
- `post_url` (String) The URL of the post email request

Optional:

- `authentication` (Attributes) Authentication Information (see [below for nested schema](#nestedatt--connectors--generic_email_gateway--authentication))
- `description` (String) A description of what your connector is used for.
- `headers` (Map of String) The headers to send with the request
- `hmac_secret` (String, Sensitive) HMAC is a method for message signing with a symmetrical key. This secret will be used to sign the base64 encoded payload, and the resulting signature will be sent in the `x-descope-webhook-s256` header. The receiving service should use this secret to verify the integrity and authenticity of the payload by checking the provided signature
- `insecure` (Boolean) Will ignore certificate errors raised by the client
- `sender` (String) The sender address
- `use_static_ips` (Boolean) Whether the connector should send all requests from specific static IPs.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--generic_email_gateway--authentication"></a>
### Nested Schema for `connectors.generic_email_gateway.authentication`

Optional:

- `api_key` (Attributes) API key authentication configuration. (see [below for nested schema](#nestedatt--connectors--generic_email_gateway--authentication--api_key))
- `basic` (Attributes) Basic authentication credentials (username and password). (see [below for nested schema](#nestedatt--connectors--generic_email_gateway--authentication--basic))
- `bearer_token` (String, Sensitive) Bearer token for HTTP authentication.

<a id="nestedatt--connectors--generic_email_gateway--authentication--api_key"></a>
### Nested Schema for `connectors.generic_email_gateway.authentication.api_key`

Required:

- `key` (String) The API key.
- `token` (String, Sensitive) The API secret.


<a id="nestedatt--connectors--generic_email_gateway--authentication--basic"></a>
### Nested Schema for `connectors.generic_email_gateway.authentication.basic`

Required:

- `password` (String, Sensitive) Password for basic HTTP authentication.
- `username` (String) Username for basic HTTP authentication.




<a id="nestedatt--connectors--generic_sms_gateway"></a>
### Nested Schema for `connectors.generic_sms_gateway`

Required:

- `name` (String) A custom name for your connector.
- `post_url` (String) The URL of the post message request

Optional:

- `authentication` (Attributes) Authentication Information (see [below for nested schema](#nestedatt--connectors--generic_sms_gateway--authentication))
- `description` (String) A description of what your connector is used for.
- `headers` (Map of String) The headers to send with the request
- `hmac_secret` (String, Sensitive) HMAC is a method for message signing with a symmetrical key. This secret will be used to sign the base64 encoded payload, and the resulting signature will be sent in the `x-descope-webhook-s256` header. The receiving service should use this secret to verify the integrity and authenticity of the payload by checking the provided signature
- `insecure` (Boolean) Will ignore certificate errors raised by the client
- `sender` (String) The sender number
- `use_static_ips` (Boolean) Whether the connector should send all requests from specific static IPs.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--generic_sms_gateway--authentication"></a>
### Nested Schema for `connectors.generic_sms_gateway.authentication`

Optional:

- `api_key` (Attributes) API key authentication configuration. (see [below for nested schema](#nestedatt--connectors--generic_sms_gateway--authentication--api_key))
- `basic` (Attributes) Basic authentication credentials (username and password). (see [below for nested schema](#nestedatt--connectors--generic_sms_gateway--authentication--basic))
- `bearer_token` (String, Sensitive) Bearer token for HTTP authentication.

<a id="nestedatt--connectors--generic_sms_gateway--authentication--api_key"></a>
### Nested Schema for `connectors.generic_sms_gateway.authentication.api_key`

Required:

- `key` (String) The API key.
- `token` (String, Sensitive) The API secret.


<a id="nestedatt--connectors--generic_sms_gateway--authentication--basic"></a>
### Nested Schema for `connectors.generic_sms_gateway.authentication.basic`

Required:

- `password` (String, Sensitive) Password for basic HTTP authentication.
- `username` (String) Username for basic HTTP authentication.




<a id="nestedatt--connectors--google_cloud_logging"></a>
### Nested Schema for `connectors.google_cloud_logging`

Required:

- `name` (String) A custom name for your connector.
- `service_account_key` (String, Sensitive) A Service Account Key JSON file created from a service account on your Google Cloud project. This file is used to authenticate and authorize the connector to access Google Cloud Logging. The service account this key belongs to must have the appropriate permissions to write logs.

Optional:

- `audit_enabled` (Boolean) Whether to enable streaming of audit events.
- `audit_filters` (Attributes List) Specify which events will be sent to the external audit service (including tenant selection). (see [below for nested schema](#nestedatt--connectors--google_cloud_logging--audit_filters))
- `description` (String) A description of what your connector is used for.
- `troubleshoot_log_enabled` (Boolean) Whether to send troubleshooting events.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--google_cloud_logging--audit_filters"></a>
### Nested Schema for `connectors.google_cloud_logging.audit_filters`

Required:

- `key` (String) The field name to filter on (either 'actions' or 'tenants').
- `operator` (String) The filter operation to apply ('includes' or 'excludes').
- `values` (List of String) The list of values to match against for the filter.



<a id="nestedatt--connectors--google_cloud_translation"></a>
### Nested Schema for `connectors.google_cloud_translation`

Required:

- `name` (String) A custom name for your connector.
- `project_id` (String) The Google Cloud project ID where the Google Cloud Translation is managed.
- `service_account_json` (String, Sensitive) Service Account JSON associated with the current project.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--google_maps_places"></a>
### Nested Schema for `connectors.google_maps_places`

Required:

- `name` (String) A custom name for your connector.
- `public_api_key` (String) The Google Maps Places public API key.

Optional:

- `address_types` (String) The address types to return.
- `description` (String) A description of what your connector is used for.
- `language` (String) The language in which to return results.
- `region` (String) The region code, specified as a CLDR two-character region code.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--hcaptcha"></a>
### Nested Schema for `connectors.hcaptcha`

Required:

- `name` (String) A custom name for your connector.
- `secret_key` (String, Sensitive) The secret key authorizes communication between Descope backend and the hCaptcha server to verify the user's response.
- `site_key` (String) The site key is used to invoke hCaptcha service on your site or mobile application.

Optional:

- `assessment_score` (Number) When configured, the hCaptcha action will return the score without assessing the request. The score ranges between 0 and 1, where 1 is a human interaction and 0 is a bot.
- `bot_threshold` (Number) The bot threshold is used to determine whether the request is a bot or a human. The score ranges between 0 and 1, where 1 is a human interaction and 0 is a bot. If the score is below this threshold, the request is considered a bot.
- `description` (String) A description of what your connector is used for.
- `override_assessment` (Boolean) Override the default assessment model. Note: Overriding assessment is intended for automated testing and should not be utilized in production environments.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--hibp"></a>
### Nested Schema for `connectors.hibp`

Required:

- `name` (String) A custom name for your connector.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--http"></a>
### Nested Schema for `connectors.http`

Required:

- `base_url` (String) The base URL to fetch
- `name` (String) A custom name for your connector.

Optional:

- `authentication` (Attributes) Authentication Information (see [below for nested schema](#nestedatt--connectors--http--authentication))
- `description` (String) A description of what your connector is used for.
- `headers` (Map of String) The headers to send with the request
- `hmac_secret` (String, Sensitive) HMAC is a method for message signing with a symmetrical key. This secret will be used to sign the base64 encoded payload, and the resulting signature will be sent in the `x-descope-webhook-s256` header. The receiving service should use this secret to verify the integrity and authenticity of the payload by checking the provided signature
- `include_headers_in_context` (Boolean) The connector response context will also include the headers. The context will have a "body" attribute and a "headers" attribute. See more details in the help guide
- `insecure` (Boolean) Will ignore certificate errors raised by the client
- `use_static_ips` (Boolean) Whether the connector should send all requests from specific static IPs.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--http--authentication"></a>
### Nested Schema for `connectors.http.authentication`

Optional:

- `api_key` (Attributes) API key authentication configuration. (see [below for nested schema](#nestedatt--connectors--http--authentication--api_key))
- `basic` (Attributes) Basic authentication credentials (username and password). (see [below for nested schema](#nestedatt--connectors--http--authentication--basic))
- `bearer_token` (String, Sensitive) Bearer token for HTTP authentication.

<a id="nestedatt--connectors--http--authentication--api_key"></a>
### Nested Schema for `connectors.http.authentication.api_key`

Required:

- `key` (String) The API key.
- `token` (String, Sensitive) The API secret.


<a id="nestedatt--connectors--http--authentication--basic"></a>
### Nested Schema for `connectors.http.authentication.basic`

Required:

- `password` (String, Sensitive) Password for basic HTTP authentication.
- `username` (String) Username for basic HTTP authentication.




<a id="nestedatt--connectors--hubspot"></a>
### Nested Schema for `connectors.hubspot`

Required:

- `access_token` (String, Sensitive) The HubSpot private API access token generated for the Descope service.
- `name` (String) A custom name for your connector.

Optional:

- `base_url` (String) The base URL of the HubSpot API, when using a custom domain in HubSpot, default value is https://api.hubapi.com .
- `description` (String) A description of what your connector is used for.
- `use_static_ips` (Boolean) Whether the connector should send all requests from specific static IPs.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--incode"></a>
### Nested Schema for `connectors.incode`

Required:

- `api_key` (String, Sensitive) Your InCode API key.
- `api_url` (String) The base URL of the Incode API
- `flow_id` (String) Your wanted InCode's flow ID.
- `name` (String) A custom name for your connector.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--intercom"></a>
### Nested Schema for `connectors.intercom`

Required:

- `name` (String) A custom name for your connector.
- `token` (String, Sensitive) The Intercom access token.

Optional:

- `description` (String) A description of what your connector is used for.
- `region` (String) Regional Hosting - US, EU, or AU. default: US

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--ldap"></a>
### Nested Schema for `connectors.ldap`

Required:

- `name` (String) A custom name for your connector.
- `server_url` (String) The LDAP server URL (e.g., ldap://localhost:389 or ldaps://localhost:636 for SSL/TLS).

Optional:

- `bind_dn` (String) The Distinguished Name to bind with for searching.
- `bind_password` (String, Sensitive) The password for the bind DN.
- `ca_certificate` (String, Sensitive) The Certificate Authority certificate in PEM format for validating the server certificate.
- `client_certificate` (String, Sensitive) The client certificate in PEM format for mTLS authentication.
- `client_key` (String, Sensitive) The client private key in PEM format for mTLS authentication.
- `description` (String) A description of what your connector is used for.
- `reject_unauthorized` (Boolean) Reject connections to LDAP servers with invalid certificates.
- `use_mtls` (Boolean) Enable mutual TLS authentication for LDAP connection.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--lokalise"></a>
### Nested Schema for `connectors.lokalise`

Required:

- `api_token` (String, Sensitive) Lokalise API token.
- `name` (String) A custom name for your connector.
- `project_id` (String) Lokalise project ID.

Optional:

- `card_id` (String) (Optional) The ID of the payment card to use for translation orders. If not provided, the team credit will be used.
- `description` (String) A description of what your connector is used for.
- `team_id` (String) Lokalise team ID. If not provided, the oldest available team will be used.
- `translation_provider` (String) The translation provider to use ('gengo', 'google', 'lokalise', 'deepl'), default is 'deepl'.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--mixpanel"></a>
### Nested Schema for `connectors.mixpanel`

Required:

- `name` (String) A custom name for your connector.
- `project_token` (String) The unique Mixpanel project token used to identify the project where data will be sent.

Optional:

- `api_secret` (String, Sensitive) The Mixpanel API secret key used for authenticating API requests.
- `audit_enabled` (Boolean) Whether to enable streaming of audit events.
- `audit_filters` (Attributes List) Specify which events will be sent to the external audit service (including tenant selection). (see [below for nested schema](#nestedatt--connectors--mixpanel--audit_filters))
- `description` (String) A description of what your connector is used for.
- `eu_residency` (Boolean) Indicates if your Mixpanel project data is stored in the EU region.
- `logs_prefix` (String) Specify a custom prefix for all log fields. The default prefix is `descope.`.
- `override_logs_prefix` (Boolean) Enable this option to use a custom prefix for log fields.
- `project_id` (String) The unique identifier for your Mixpanel project.
- `service_account_secret` (String, Sensitive) The Mixpanel service account secret used for integration.
- `service_account_username` (String) The Mixpanel service account username used for integration.
- `troubleshoot_log_enabled` (Boolean) Whether to send troubleshooting events.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--mixpanel--audit_filters"></a>
### Nested Schema for `connectors.mixpanel.audit_filters`

Required:

- `key` (String) The field name to filter on (either 'actions' or 'tenants').
- `operator` (String) The filter operation to apply ('includes' or 'excludes').
- `values` (List of String) The list of values to match against for the filter.



<a id="nestedatt--connectors--mparticle"></a>
### Nested Schema for `connectors.mparticle`

Required:

- `api_key` (String, Sensitive) The mParticle Server to Server Key generated for the Descope service.
- `api_secret` (String, Sensitive) The mParticle Server to Server Secret generated for the Descope service.
- `name` (String) A custom name for your connector.

Optional:

- `base_url` (String) The base URL of the mParticle API, when using a custom domain in mParticle. default value is https://s2s.mparticle.com/
- `default_environment` (String) The default environment of which connector send data to, either “production” or “development“. default value: production. This field can be overridden per event (see at flows).
- `description` (String) A description of what your connector is used for.
- `use_static_ips` (Boolean) Whether the connector should send all requests from specific static IPs.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--newrelic"></a>
### Nested Schema for `connectors.newrelic`

Required:

- `api_key` (String, Sensitive) Ingest License Key of the account you want to report data to.
- `name` (String) A custom name for your connector.

Optional:

- `audit_enabled` (Boolean) Whether to enable streaming of audit events.
- `audit_filters` (Attributes List) Specify which events will be sent to the external audit service (including tenant selection). (see [below for nested schema](#nestedatt--connectors--newrelic--audit_filters))
- `data_center` (String) The New Relic data center the account belongs to. Possible values are: `US`, `EU`, `FedRAMP`. Default is `US`.
- `description` (String) A description of what your connector is used for.
- `logs_prefix` (String) Specify a custom prefix for all log fields. The default prefix is `descope.`.
- `override_logs_prefix` (Boolean) Enable this option to use a custom prefix for log fields.
- `troubleshoot_log_enabled` (Boolean) Whether to send troubleshooting events.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--newrelic--audit_filters"></a>
### Nested Schema for `connectors.newrelic.audit_filters`

Required:

- `key` (String) The field name to filter on (either 'actions' or 'tenants').
- `operator` (String) The filter operation to apply ('includes' or 'excludes').
- `values` (List of String) The list of values to match against for the filter.



<a id="nestedatt--connectors--opentelemetry"></a>
### Nested Schema for `connectors.opentelemetry`

Required:

- `endpoint` (String) The OTLP endpoint URL.
- `name` (String) A custom name for your connector.

Optional:

- `audit_enabled` (Boolean) Whether to enable streaming of audit events.
- `audit_filters` (Attributes List) Specify which events will be sent to the external audit service (including tenant selection). (see [below for nested schema](#nestedatt--connectors--opentelemetry--audit_filters))
- `authentication` (Attributes) Authentication Information (see [below for nested schema](#nestedatt--connectors--opentelemetry--authentication))
- `description` (String) A description of what your connector is used for.
- `headers` (Map of String) The headers to send with the request
- `insecure` (Boolean) Will ignore certificate errors raised by the client
- `protocol` (String) Protocol to use for OTLP: http or grpc.
- `troubleshoot_log_enabled` (Boolean) Whether to send troubleshooting events.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--opentelemetry--audit_filters"></a>
### Nested Schema for `connectors.opentelemetry.audit_filters`

Required:

- `key` (String) The field name to filter on (either 'actions' or 'tenants').
- `operator` (String) The filter operation to apply ('includes' or 'excludes').
- `values` (List of String) The list of values to match against for the filter.


<a id="nestedatt--connectors--opentelemetry--authentication"></a>
### Nested Schema for `connectors.opentelemetry.authentication`

Optional:

- `api_key` (Attributes) API key authentication configuration. (see [below for nested schema](#nestedatt--connectors--opentelemetry--authentication--api_key))
- `basic` (Attributes) Basic authentication credentials (username and password). (see [below for nested schema](#nestedatt--connectors--opentelemetry--authentication--basic))
- `bearer_token` (String, Sensitive) Bearer token for HTTP authentication.

<a id="nestedatt--connectors--opentelemetry--authentication--api_key"></a>
### Nested Schema for `connectors.opentelemetry.authentication.api_key`

Required:

- `key` (String) The API key.
- `token` (String, Sensitive) The API secret.


<a id="nestedatt--connectors--opentelemetry--authentication--basic"></a>
### Nested Schema for `connectors.opentelemetry.authentication.basic`

Required:

- `password` (String, Sensitive) Password for basic HTTP authentication.
- `username` (String) Username for basic HTTP authentication.




<a id="nestedatt--connectors--ping_directory"></a>
### Nested Schema for `connectors.ping_directory`

Required:

- `host` (String) PingDirectory's REST API host.
- `name` (String) A custom name for your connector.
- `port` (Number) PingDirectory's REST API port.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--postmark"></a>
### Nested Schema for `connectors.postmark`

Required:

- `email_from` (String) The email address that will appear in the 'From' field of the sent email
- `message_stream_id` (String) The ID of the message stream to use for the email
- `name` (String) A custom name for your connector.
- `server_api_token` (String, Sensitive) The API token for authenticating with the Postmark server

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--radar"></a>
### Nested Schema for `connectors.radar`

Required:

- `name` (String) A custom name for your connector.
- `public_api_key` (String) The Radar publishable API key.

Optional:

- `address_types` (String) The address types to return.
- `description` (String) A description of what your connector is used for.
- `language` (String) The language in which to return results.
- `limit` (Number) The maximum number of results to return.
- `region` (String) The region code, specified as a two-letter ISO 3166 code.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--recaptcha"></a>
### Nested Schema for `connectors.recaptcha`

Required:

- `name` (String) A custom name for your connector.
- `secret_key` (String, Sensitive) The secret key authorizes communication between Descope backend and the reCAPTCHA server to verify the user's response.
- `site_key` (String) The site key is used to invoke reCAPTCHA service on your site or mobile application.

Optional:

- `assessment_score` (Number) When configured, the Recaptcha action will return the score without assessing the request. The score ranges between 0 and 1, where 1 is a human interaction and 0 is a bot.
- `description` (String) A description of what your connector is used for.
- `override_assessment` (Boolean) Override the default assessment model. Note: Overriding assessment is intended for automated testing and should not be utilized in production environments.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--recaptcha_enterprise"></a>
### Nested Schema for `connectors.recaptcha_enterprise`

Required:

- `api_key` (String, Sensitive) API key associated with the current project.
- `name` (String) A custom name for your connector.
- `project_id` (String) The Google Cloud project ID where the reCAPTCHA Enterprise is managed.
- `site_key` (String) The site key is used to invoke reCAPTCHA Enterprise service on your site or mobile application.

Optional:

- `assessment_score` (Number) When configured, the Recaptcha action will return the score without assessing the request. The score ranges between 0 and 1, where 1 is a human interaction and 0 is a bot.
- `base_url` (String) Apply a custom url to the reCAPTCHA Enterprise scripts. This is useful when attempting to use reCAPTCHA globally. Defaults to https://www.google.com
- `bot_threshold` (Number) The bot threshold is used to determine whether the request is a bot or a human. The score ranges between 0 and 1, where 1 is a human interaction and 0 is a bot. If the score is below this threshold, the request is considered a bot.
- `description` (String) A description of what your connector is used for.
- `override_assessment` (Boolean) Override the default assessment model. Note: Overriding assessment is intended for automated testing and should not be utilized in production environments.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--rekognition"></a>
### Nested Schema for `connectors.rekognition`

Required:

- `access_key_id` (String) The AWS access key ID
- `collection_id` (String) The collection to store registered users in. Should match `[a-zA-Z0-9_.-]+` pattern. Changing this will cause losing existing users.
- `name` (String) A custom name for your connector.
- `secret_access_key` (String, Sensitive) The AWS secret access key

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--salesforce"></a>
### Nested Schema for `connectors.salesforce`

Required:

- `base_url` (String) The Salesforce API base URL.
- `client_id` (String) The consumer key of the connected app.
- `client_secret` (String, Sensitive) The consumer secret of the connected app.
- `name` (String) A custom name for your connector.
- `version` (String) REST API Version.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--salesforce_marketing_cloud"></a>
### Nested Schema for `connectors.salesforce_marketing_cloud`

Required:

- `client_id` (String) Client ID issued when you create the API integration in Installed Packages.
- `client_secret` (String, Sensitive) Client secret issued when you create the API integration in Installed Packages.
- `name` (String) A custom name for your connector.
- `subdomain` (String) The Salesforce Marketing Cloud endpoint subdomain.

Optional:

- `account_id` (String) Account identifier, or MID, of the target business unit.
- `description` (String) A description of what your connector is used for.
- `scope` (String) Space-separated list of data-access permissions for your connector.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--sardine"></a>
### Nested Schema for `connectors.sardine`

Required:

- `base_url` (String) The base URL for the Sardine API, e.g.: https://api.sandbox.sardine.ai, https://api.sardine.ai, https://api.eu.sardine.ai.
- `client_id` (String) The Sardine Client ID.
- `client_secret` (String, Sensitive) The Sardine Client Secret.
- `name` (String) A custom name for your connector.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--segment"></a>
### Nested Schema for `connectors.segment`

Required:

- `name` (String) A custom name for your connector.
- `write_key` (String, Sensitive) The Segment Write Key generated for the Descope service.

Optional:

- `description` (String) A description of what your connector is used for.
- `host` (String) The base URL of the Segment API, when using a custom domain in Segment.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--sendgrid"></a>
### Nested Schema for `connectors.sendgrid`

Required:

- `authentication` (Attributes) SendGrid API authentication configuration. (see [below for nested schema](#nestedatt--connectors--sendgrid--authentication))
- `name` (String) A custom name for your connector.
- `sender` (Attributes) The sender details that should be displayed in the email message. (see [below for nested schema](#nestedatt--connectors--sendgrid--sender))

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--sendgrid--authentication"></a>
### Nested Schema for `connectors.sendgrid.authentication`

Required:

- `api_key` (String, Sensitive) SendGrid API key for authentication.


<a id="nestedatt--connectors--sendgrid--sender"></a>
### Nested Schema for `connectors.sendgrid.sender`

Required:

- `email` (String) The email address that will appear as the sender of the email.

Optional:

- `name` (String) The display name that will appear as the sender of the email.



<a id="nestedatt--connectors--ses"></a>
### Nested Schema for `connectors.ses`

Required:

- `name` (String) A custom name for your connector.
- `region` (String) AWS region to send requests to (e.g. `us-west-2`).
- `sender` (Attributes) The sender details that should be displayed in the email message. (see [below for nested schema](#nestedatt--connectors--ses--sender))

Optional:

- `access_key_id` (String, Sensitive) AWS Access key ID.
- `auth_type` (String) The authentication type to use.
- `description` (String) A description of what your connector is used for.
- `endpoint` (String) An optional endpoint URL (hostname only or fully qualified URI).
- `external_id` (String) The external ID to use when assuming the role.
- `role_arn` (String) The Amazon Resource Name (ARN) of the role to assume.
- `secret` (String, Sensitive) AWS Secret Access Key.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--ses--sender"></a>
### Nested Schema for `connectors.ses.sender`

Required:

- `email` (String) The email address that will appear as the sender of the email.

Optional:

- `name` (String) The display name that will appear as the sender of the email.



<a id="nestedatt--connectors--slack"></a>
### Nested Schema for `connectors.slack`

Required:

- `name` (String) A custom name for your connector.
- `token` (String, Sensitive) The OAuth token for Slack's Bot User, used to authenticate API requests.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--smartling"></a>
### Nested Schema for `connectors.smartling`

Required:

- `account_uid` (String) The account UID for the Smartling account.
- `name` (String) A custom name for your connector.
- `user_identifier` (String) The user identifier for the Smartling account.
- `user_secret` (String, Sensitive) The user secret for the Smartling account.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--smtp"></a>
### Nested Schema for `connectors.smtp`

Required:

- `authentication` (Attributes) SMTP server authentication credentials and method. (see [below for nested schema](#nestedatt--connectors--smtp--authentication))
- `name` (String) A custom name for your connector.
- `sender` (Attributes) The sender details that should be displayed in the email message. (see [below for nested schema](#nestedatt--connectors--smtp--sender))
- `server` (Attributes) SMTP server connection details including hostname and port. (see [below for nested schema](#nestedatt--connectors--smtp--server))

Optional:

- `description` (String) A description of what your connector is used for.
- `use_static_ips` (Boolean) Whether the connector should send all requests from specific static IPs.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--smtp--authentication"></a>
### Nested Schema for `connectors.smtp.authentication`

Required:

- `password` (String, Sensitive) Password for SMTP server authentication.
- `username` (String) Username for SMTP server authentication.

Optional:

- `method` (String) SMTP authentication method (`plain` or `login`).


<a id="nestedatt--connectors--smtp--sender"></a>
### Nested Schema for `connectors.smtp.sender`

Required:

- `email` (String) The email address that will appear as the sender of the email.

Optional:

- `name` (String) The display name that will appear as the sender of the email.


<a id="nestedatt--connectors--smtp--server"></a>
### Nested Schema for `connectors.smtp.server`

Required:

- `host` (String) The hostname or IP address of the SMTP server.

Optional:

- `port` (Number) The port number to connect to on the SMTP server.



<a id="nestedatt--connectors--sns"></a>
### Nested Schema for `connectors.sns`

Required:

- `access_key_id` (String, Sensitive) AWS Access key ID.
- `name` (String) A custom name for your connector.
- `region` (String) AWS region to send requests to (e.g. `us-west-2`).
- `secret` (String, Sensitive) AWS Secret Access Key.

Optional:

- `description` (String) A description of what your connector is used for.
- `endpoint` (String) An optional endpoint URL (hostname only or fully qualified URI).
- `entity_id` (String) The entity ID or principal entity (PE) ID for sending text messages to recipients in India.
- `organization_number` (String, Deprecated) Use the `origination_number` attribute instead.
- `origination_number` (String) An optional phone number from which the text messages are going to be sent. Make sure it is registered properly in your server.
- `sender_id` (String) The name of the sender from which the text message is going to be sent (see SNS documentation regarding acceptable IDs and supported regions/countries).
- `template_id` (String) The template for sending text messages to recipients in India. The template ID must be associated with the sender ID.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--splunk"></a>
### Nested Schema for `connectors.splunk`

Required:

- `hec_token` (String, Sensitive) An HTTP Event Collector token configured on your Splunk project.
- `hec_url` (String) The URL to be used accessing your Splunk system, including the appropriate port
- `name` (String) A custom name for your connector.

Optional:

- `audit_enabled` (Boolean) Whether to enable streaming of audit events.
- `audit_filters` (Attributes List) Specify which events will be sent to the external audit service (including tenant selection). (see [below for nested schema](#nestedatt--connectors--splunk--audit_filters))
- `description` (String) A description of what your connector is used for.
- `index` (String) An optional index to use for all sent events
- `troubleshoot_log_enabled` (Boolean) Whether to send troubleshooting events.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--splunk--audit_filters"></a>
### Nested Schema for `connectors.splunk.audit_filters`

Required:

- `key` (String) The field name to filter on (either 'actions' or 'tenants').
- `operator` (String) The filter operation to apply ('includes' or 'excludes').
- `values` (List of String) The list of values to match against for the filter.



<a id="nestedatt--connectors--sql"></a>
### Nested Schema for `connectors.sql`

Required:

- `engine_name` (String) The database engine type.
- `host` (String) The database host.
- `name` (String) A custom name for your connector.
- `password` (String, Sensitive) The database password.
- `username` (String) The database username.

Optional:

- `database_name` (String) The database name.
- `description` (String) A description of what your connector is used for.
- `port` (Number) The database port. If not specified, the default port for the selected engine will be used.
- `service_name` (String) The Oracle service name (required for Oracle only).

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--sumologic"></a>
### Nested Schema for `connectors.sumologic`

Required:

- `http_source_url` (String, Sensitive) The URL associated with an HTTP Hosted collector
- `name` (String) A custom name for your connector.

Optional:

- `audit_enabled` (Boolean) Whether to enable streaming of audit events.
- `audit_filters` (Attributes List) Specify which events will be sent to the external audit service (including tenant selection). (see [below for nested schema](#nestedatt--connectors--sumologic--audit_filters))
- `description` (String) A description of what your connector is used for.
- `troubleshoot_log_enabled` (Boolean) Whether to send troubleshooting events.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--sumologic--audit_filters"></a>
### Nested Schema for `connectors.sumologic.audit_filters`

Required:

- `key` (String) The field name to filter on (either 'actions' or 'tenants').
- `operator` (String) The filter operation to apply ('includes' or 'excludes').
- `values` (List of String) The list of values to match against for the filter.



<a id="nestedatt--connectors--supabase"></a>
### Nested Schema for `connectors.supabase`

Required:

- `name` (String) A custom name for your connector.

Optional:

- `auth_type` (String) The authentication type to use.
- `create_users` (Boolean) Enable to automatically create users in Supabase when generating tokens. Will only create a new user if one does not already exist. When disabled, only JWT tokens will be generated, WITHOUT user creation.
- `custom_claims_mapping` (Map of String) A mapping of Descope user fields or JWT claims to Supabase custom claims
- `description` (String) A description of what your connector is used for.
- `expiration_time` (Number) The duration in minutes for which the token is valid.
- `private_key` (String, Sensitive) The private key in JWK format used to sign the JWT. You can generate a key using tools like `npx supabase gen signing-key --algorithm ES256`. Make sure to use the ES256 algorithm.
- `project_base_url` (String) Your Supabase Project's API base URL, e.g.: https://<your-project-id>.supabase.co.
- `service_role_api_key` (String, Sensitive) The service role API key for your Supabase project, required to create users.
- `signing_secret` (String, Sensitive) The signing secret for your Supabase project.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--telesign"></a>
### Nested Schema for `connectors.telesign`

Required:

- `api_key` (String, Sensitive) The unique Telesign API key
- `customer_id` (String) The unique Telesign account Customer ID
- `name` (String) A custom name for your connector.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--traceable"></a>
### Nested Schema for `connectors.traceable`

Required:

- `name` (String) A custom name for your connector.
- `secret_key` (String, Sensitive) The Traceable secret key.

Optional:

- `description` (String) A description of what your connector is used for.
- `eu_region` (Boolean) EU(Europe) Region deployment of Traceable platform.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--turnstile"></a>
### Nested Schema for `connectors.turnstile`

Required:

- `name` (String) A custom name for your connector.
- `secret_key` (String, Sensitive) The secret key authorizes communication between Descope backend and the Turnstile server to verify the user's response.
- `site_key` (String) The site key is used to invoke Turnstile service on your site or mobile application.

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--twilio_core"></a>
### Nested Schema for `connectors.twilio_core`

Required:

- `account_sid` (String) Twilio Account SID from your Twilio Console.
- `authentication` (Attributes) Twilio authentication credentials (either auth token or API key/secret). (see [below for nested schema](#nestedatt--connectors--twilio_core--authentication))
- `name` (String) A custom name for your connector.
- `senders` (Attributes) Configuration for SMS and voice message senders. (see [below for nested schema](#nestedatt--connectors--twilio_core--senders))

Optional:

- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--twilio_core--authentication"></a>
### Nested Schema for `connectors.twilio_core.authentication`

Optional:

- `api_key` (String, Sensitive) Twilio API Key for authentication (used with API Secret).
- `api_secret` (String, Sensitive) Twilio API Secret for authentication (used with API Key).
- `auth_token` (String, Sensitive) Twilio Auth Token for authentication.


<a id="nestedatt--connectors--twilio_core--senders"></a>
### Nested Schema for `connectors.twilio_core.senders`

Required:

- `sms` (Attributes) SMS sender configuration using either a phone number or messaging service. (see [below for nested schema](#nestedatt--connectors--twilio_core--senders--sms))

Optional:

- `voice` (Attributes) Voice call sender configuration. (see [below for nested schema](#nestedatt--connectors--twilio_core--senders--voice))

<a id="nestedatt--connectors--twilio_core--senders--sms"></a>
### Nested Schema for `connectors.twilio_core.senders.sms`

Optional:

- `messaging_service_sid` (String) Twilio Messaging Service SID for sending SMS messages.
- `phone_number` (String) Twilio phone number for sending SMS messages.


<a id="nestedatt--connectors--twilio_core--senders--voice"></a>
### Nested Schema for `connectors.twilio_core.senders.voice`

Required:

- `phone_number` (String) Twilio phone number for making voice calls.




<a id="nestedatt--connectors--twilio_verify"></a>
### Nested Schema for `connectors.twilio_verify`

Required:

- `account_sid` (String) Twilio Account SID from your Twilio Console.
- `authentication` (Attributes) Twilio authentication credentials (either auth token or API key/secret). (see [below for nested schema](#nestedatt--connectors--twilio_verify--authentication))
- `name` (String) A custom name for your connector.
- `service_sid` (String) Twilio Verify Service SID for verification services.

Optional:

- `description` (String) A description of what your connector is used for.
- `sender` (String) Optional sender identifier for verification messages.

Read-Only:

- `id` (String)

<a id="nestedatt--connectors--twilio_verify--authentication"></a>
### Nested Schema for `connectors.twilio_verify.authentication`

Optional:

- `api_key` (String, Sensitive) Twilio API Key for authentication (used with API Secret).
- `api_secret` (String, Sensitive) Twilio API Secret for authentication (used with API Key).
- `auth_token` (String, Sensitive) Twilio Auth Token for authentication.



<a id="nestedatt--connectors--unibeam"></a>
### Nested Schema for `connectors.unibeam`

Required:

- `base_url` (String) Unibeam API base URL.
- `client_id` (String) OAuth2 client ID for authentication.
- `client_secret` (String, Sensitive) OAuth2 client secret for authentication.
- `customer_id` (String) Your Unibeam customer ID.
- `hmac_secret` (String, Sensitive) HMAC secret supplied by Unibeam for securing communications.
- `name` (String) A custom name for your connector.

Optional:

- `default_message` (String) Default message to display when no message is provided in the command.
- `description` (String) A description of what your connector is used for.

Read-Only:

- `id` (String)


<a id="nestedatt--connectors--zerobounce"></a>
### Nested Schema for `connectors.zerobounce`

Required:

- `api_key` (String, Sensitive) The ZeroBounce API key.
- `name` (String) A custom name for your connector.

Optional:

- `description` (String) A description of what your connector is used for.
- `region` (String) ZeroBounce platform region.

Read-Only:

- `id` (String)



<a id="nestedatt--flows"></a>
### Nested Schema for `flows`

Required:

- `data` (String) The JSON data defining the authentication flow configuration, including metadata, screens, contents, and references.


<a id="nestedatt--invite_settings"></a>
### Nested Schema for `invite_settings`

Optional:

- `add_magiclink_token` (Boolean) Whether to include a magic link token in invitation messages.
- `email_service` (Attributes) Settings related to sending invitation emails. (see [below for nested schema](#nestedatt--invite_settings--email_service))
- `expire_invited_users` (Boolean) Expire the user account if the invitation is not accepted within the expiration time.
- `invite_expiration` (String) The expiry time for the invitation, meant to be used together with `expire_invited_users` and/or `add_magiclink_token`. Use values such as "2 weeks", "4 days", etc. The minimum value is "1 hour".
- `invite_url` (String) Custom URL to include in the message sent to invited users.
- `require_invitation` (Boolean) Whether users must be invited before they can sign up to the project.
- `send_email` (Boolean) Whether to send invitation emails to users.
- `send_text` (Boolean) Whether to send invitation SMS messages to users.

<a id="nestedatt--invite_settings--email_service"></a>
### Nested Schema for `invite_settings.email_service`

Required:

- `connector` (String) The name of the email connector to use for sending emails.

Optional:

- `templates` (Attributes List) A list of email templates for different authentication flows. (see [below for nested schema](#nestedatt--invite_settings--email_service--templates))

<a id="nestedatt--invite_settings--email_service--templates"></a>
### Nested Schema for `invite_settings.email_service.templates`

Required:

- `name` (String) Unique name for this email template.
- `subject` (String) Subject line of the email message.

Optional:

- `active` (Boolean) Whether this email template is currently active and in use.
- `html_body` (String) HTML content of the email message body, required if `use_plain_text_body` isn't set.
- `plain_text_body` (String) Plain text version of the email message body, required if `use_plain_text_body` is set to `true`.
- `use_plain_text_body` (Boolean) Whether to use the plain text body instead of HTML for the email.

Read-Only:

- `id` (String)




<a id="nestedatt--jwt_templates"></a>
### Nested Schema for `jwt_templates`

Optional:

- `access_key_templates` (Attributes List) A list of `Access Key` type JWT Templates. (see [below for nested schema](#nestedatt--jwt_templates--access_key_templates))
- `user_templates` (Attributes List) A list of `User` type JWT Templates. (see [below for nested schema](#nestedatt--jwt_templates--user_templates))

<a id="nestedatt--jwt_templates--access_key_templates"></a>
### Nested Schema for `jwt_templates.access_key_templates`

Required:

- `name` (String) Name of the JWT Template.
- `template` (String) The JSON template defining the structure and claims of the JWT token. This is expected to be a valid JSON object given as a `string` value.

Optional:

- `auth_schema` (String) The authorization claims format - `default`, `tenantOnly` or `none`. Read more about schema types [here](https://docs.descope.com/project-settings/jwt-templates).
- `auto_tenant_claim` (Boolean) When a user is associated with a single tenant, the tenant will be set as the user's active tenant, using the `dct` (Descope Current Tenant) claim in their JWT.
- `conformance_issuer` (Boolean) Whether to use OIDC conformance for the JWT issuer field.
- `description` (String) Description of the JWT Template.
- `empty_claim_policy` (String) Policy for empty claims - `none`, `nil` or `delete`.
- `enforce_issuer` (Boolean) Whether to enforce that the JWT issuer matches the project configuration.

Read-Only:

- `id` (String)


<a id="nestedatt--jwt_templates--user_templates"></a>
### Nested Schema for `jwt_templates.user_templates`

Required:

- `name` (String) Name of the JWT Template.
- `template` (String) The JSON template defining the structure and claims of the JWT token. This is expected to be a valid JSON object given as a `string` value.

Optional:

- `auth_schema` (String) The authorization claims format - `default`, `tenantOnly` or `none`. Read more about schema types [here](https://docs.descope.com/project-settings/jwt-templates).
- `auto_tenant_claim` (Boolean) When a user is associated with a single tenant, the tenant will be set as the user's active tenant, using the `dct` (Descope Current Tenant) claim in their JWT.
- `conformance_issuer` (Boolean) Whether to use OIDC conformance for the JWT issuer field.
- `description` (String) Description of the JWT Template.
- `empty_claim_policy` (String) Policy for empty claims - `none`, `nil` or `delete`.
- `enforce_issuer` (Boolean) Whether to enforce that the JWT issuer matches the project configuration.

Read-Only:

- `id` (String)



<a id="nestedatt--lists"></a>
### Nested Schema for `lists`

Required:

- `data` (String) The JSON data for the list. The format depends on the `type`: - For `"texts"` and `"ips"` types: Must be a JSON array of strings (e.g., `["item1", "item2"]`) - For `"ips"` type: Each string must be a valid IP address or CIDR range - For `"json"` type: Must be a JSON object (e.g., `{"key": "value"}`)
- `name` (String) The name of the list. Maximum length is 100 characters.
- `type` (String) The type of list. Must be one of: - `"texts"` - A list of text strings - `"ips"` - A list of IP addresses or CIDR ranges - `"json"` - A JSON object

Optional:

- `description` (String) An optional description for the list. Defaults to an empty string if not provided.

Read-Only:

- `id` (String)


<a id="nestedatt--project_settings"></a>
### Nested Schema for `project_settings`

Optional:

- `access_key_jwt_template` (String) Name of the access key JWT Template.
- `access_key_session_token_expiration` (String) The expiry time for access key session tokens. Use values such as "10 minutes", "4 hours", etc. The value needs to be at least 3 minutes and can't be longer than 4 weeks.
- `app_url` (String) The URL which your application resides on.
- `approved_domains` (Set of String) The list of approved domains that are allowed for redirect and verification URLs for different authentication methods.
- `custom_domain` (String) A custom CNAME that's configured to point to `cname.descope.com`. Read more about custom domains and cookie policy [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).
- `default_no_sso_apps` (Boolean) Define whether a user created with no federated apps, will have access to all apps, or will not have access to any app.
- `enable_inactivity` (Boolean) Use `True` to enable session inactivity. To read more about session inactivity click [here](https://docs.descope.com/project-settings#session-inactivity).
- `inactivity_time` (String) The session inactivity time. Use values such as "15 minutes", "1 hour", etc. The minimum value is "10 minutes".
- `refresh_token_cookie_domain` (String) The domain name for refresh token cookies. To read more about custom domain and cookie policy click [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).
- `refresh_token_cookie_policy` (String) Use `strict`, `lax` or `none`. Read more about custom domains and cookie policy [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).
- `refresh_token_expiration` (String) The expiry time for the refresh token, after which the user must log in again. Use values such as "4 weeks", "14 days", etc. The minimum value is "3 minutes".
- `refresh_token_response_method` (String) Configure how refresh tokens are managed by the Descope SDKs. Must be either `response_body` or `cookies`. The default value is `response_body`.
- `refresh_token_rotation` (Boolean) Every time the user refreshes their session token via their refresh token, the refresh token itself is also updated to a new one.
- `session_migration` (Attributes) Configure seamless migration of existing user sessions from another vendor to Descope. (see [below for nested schema](#nestedatt--project_settings--session_migration))
- `session_token_cookie_domain` (String) The domain name for session token cookies. To read more about custom domain and cookie policy click [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).
- `session_token_cookie_policy` (String) Use `strict`, `lax` or `none`. Read more about custom domains and cookie policy [here](https://docs.descope.com/how-to-deploy-to-production/custom-domain).
- `session_token_expiration` (String) The expiry time of the session token, used for accessing the application's resources. The value needs to be at least 3 minutes and can't be longer than the refresh token expiration.
- `session_token_response_method` (String) Configure how sessions tokens are managed by the Descope SDKs. Must be either `response_body` or `cookies`. The default value is `response_body`.
- `step_up_token_expiration` (String) The expiry time for the step up token, after which it will not be valid and the user will automatically go back to the session token.
- `test_users_loginid_regexp` (String) Define a regular expression so that whenever a user is created with a matching login ID it will automatically be marked as a test user.
- `test_users_static_otp` (String) A 6 digit static OTP code for use with test users.
- `test_users_verifier_regexp` (String) The pattern of the verifiers that will be used for testing.
- `trusted_device_token_expiration` (String) The expiry time for the trusted device token. The minimum value is "3 minutes".
- `user_jwt_template` (String) Name of the user JWT Template.

<a id="nestedatt--project_settings--session_migration"></a>
### Nested Schema for `project_settings.session_migration`

Optional:

- `audience` (String) The audience value if needed by the vendor.
- `client_id` (String) The unique client ID for the vendor.
- `domain` (String) The domain value if needed by the vendor.
- `issuer` (String) An issuer URL if needed by the vendor.
- `loginid_matched_attributes` (Set of String) A set of attributes from the vendor's user that should be used to match with the Descope user's login ID.
- `vendor` (String) The name of the vendor the sessions are migrated from, in all lowercase.



<a id="nestedatt--styles"></a>
### Nested Schema for `styles`

Required:

- `data` (String) The JSON data defining the visual styling and theme configuration used for authentication, widgets, etc.


<a id="nestedatt--widgets"></a>
### Nested Schema for `widgets`

Required:

- `data` (String) The JSON data defining the widget. This will usually be exported as a `.json` file from the Descope console, and set in the `.tf` file using the `data = file("...")` syntax.

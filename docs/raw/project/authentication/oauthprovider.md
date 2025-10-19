
OAuthProviderCustom
===================



disabled
--------

- Type: `bool` 

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



client_id
---------

- Type: `string` 

The client ID for the OAuth provider, used to identify the application to the provider.



client_secret
-------------

- Type: `secret` 

The client secret for the OAuth provider, used to authenticate the application with the provider.



manage_provider_tokens
----------------------

- Type: `bool` 

Whether to enable provider token management for this OAuth provider.



callback_domain
---------------

- Type: `string` 

Use a custom domain in your OAuth verification screen.



redirect_url
------------

- Type: `string` 

Users will be directed to this URL after authentication. If redirect URL is specified in
the SDK/API call, it will override this value.



prompts
-------

- Type: `list` of `string` 

Custom prompts or consent screens that users may see during OAuth authentication.



allowed_grant_types
-------------------

- Type: `list` of `string` 

The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens
from the OAuth provider.



scopes
------

- Type: `list` of `string` 

Scopes of access that the application requests from the user's account on the OAuth provider.



merge_user_accounts
-------------------

- Type: `bool` 
- Default: `true`

Whether to merge existing user accounts with new ones created through OAuth authentication.



description
-----------

- Type: `string` 

A brief description of the OAuth provider.



logo
----

- Type: `string` 

The URL of the logo associated with the OAuth provider.



issuer
------

- Type: `string` 

The issuer identifier for the OAuth provider.



authorization_endpoint
----------------------

- Type: `string` 

The URL that users are redirected to for authorization with the OAuth provider.



token_endpoint
--------------

- Type: `string` 

The URL where the application requests an access token from the OAuth provider.



user_info_endpoint
------------------

- Type: `string` 

The URL where the application retrieves user information from the OAuth provider.



jwks_endpoint
-------------

- Type: `string` 

The URL where the application can retrieve JSON Web Key Sets (JWKS) for the OAuth provider.



use_client_assertion
--------------------

- Type: `bool` 

Use private key JWT (client assertion) instead of client secret.



claim_mapping
-------------

- Type: `map` of `string` 

Maps OAuth provider claims to Descope user attributes.





OAuthProviderSystem
===================



disabled
--------

- Type: `bool` 

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



client_id
---------

- Type: `string` 

The client ID for the OAuth provider, used to identify the application to the provider.



client_secret
-------------

- Type: `secret` 

The client secret for the OAuth provider, used to authenticate the application with the provider.



manage_provider_tokens
----------------------

- Type: `bool` 

Whether to enable provider token management for this OAuth provider.



callback_domain
---------------

- Type: `string` 

Use a custom domain in your OAuth verification screen.



redirect_url
------------

- Type: `string` 

Users will be directed to this URL after authentication. If redirect URL is specified in
the SDK/API call, it will override this value.



prompts
-------

- Type: `list` of `string` 

Custom prompts or consent screens that users may see during OAuth authentication.



allowed_grant_types
-------------------

- Type: `list` of `string` 

The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens
from the OAuth provider.



scopes
------

- Type: `list` of `string` 

Scopes of access that the application requests from the user's account on the OAuth provider.



merge_user_accounts
-------------------

- Type: `bool` 
- Default: `true`

Whether to merge existing user accounts with new ones created through OAuth authentication.





OAuthProviderApple
==================



disabled
--------

- Type: `bool` 

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



client_id
---------

- Type: `string` 

The client ID for the OAuth provider, used to identify the application to the provider.



client_secret
-------------

- Type: `secret` 

The client secret for the OAuth provider, used to authenticate the application with the provider.



manage_provider_tokens
----------------------

- Type: `bool` 

Whether to enable provider token management for this OAuth provider.



callback_domain
---------------

- Type: `string` 

Use a custom domain in your OAuth verification screen.



redirect_url
------------

- Type: `string` 

Users will be directed to this URL after authentication. If redirect URL is specified in
the SDK/API call, it will override this value.



prompts
-------

- Type: `list` of `string` 

Custom prompts or consent screens that users may see during OAuth authentication.



allowed_grant_types
-------------------

- Type: `list` of `string` 

The type of grants (`authorization_code` or `implicit`) to allow when requesting access tokens
from the OAuth provider.



scopes
------

- Type: `list` of `string` 

Scopes of access that the application requests from the user's account on the OAuth provider.



merge_user_accounts
-------------------

- Type: `bool` 
- Default: `true`

Whether to merge existing user accounts with new ones created through OAuth authentication.



native_client_id
----------------

- Type: `string` 

The client ID for the OAuth provider, used for Sign in with Apple in mobile apps.



native_client_secret
--------------------

- Type: `secret` 

The client secret for the OAuth provider, used for Sign in with Apple in mobile apps.



apple_key_generator
-------------------

- Type: `object` of `authentication.OAuthProviderAppleKeyGenerator` 

The apple key generator object describing how to create a dynamic apple client secret for applications.



native_apple_key_generator
--------------------------

- Type: `object` of `authentication.OAuthProviderAppleKeyGenerator` 

The apple key generator object describing how to create a dynamic naive apple client secret for mobile apps.





OAuthProviderAppleKeyGenerator
==============================



key_id
------

- Type: `string` (required)

The apple generator key id produced by Apple.



team_id
-------

- Type: `string` (required)

The apple generator team id assigned to the key by Apple.



private_key
-----------

- Type: `secret` (required)

The apple generator private key produced by Apple.

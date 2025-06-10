
OAuth
=====



disabled
--------

- Type: `bool` 

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



system
------

- Type: `object` of `authentication.OAuthSystemProvider` 

Custom configurations for builtin OAuth providers such as Apple, Google, GitHub, Facebook, etc.



custom
------

- Type: `map` of `authentication.OAuthProvider` 

Custom OAuth providers configured for this project.





OAuthSystemProvider
===================



apple
-----

- Type: `object` of `authentication.OAuthProvider` 

Apple's OAuth provider, allowing users to authenticate with their Apple Account.



discord
-------

- Type: `object` of `authentication.OAuthProvider` 

Discord's OAuth provider, allowing users to authenticate with their Discord account.



facebook
--------

- Type: `object` of `authentication.OAuthProvider` 

Facebook's OAuth provider, allowing users to authenticate with their Facebook account.



github
------

- Type: `object` of `authentication.OAuthProvider` 

GitHub's OAuth provider, allowing users to authenticate with their GitHub account.



gitlab
------

- Type: `object` of `authentication.OAuthProvider` 

GitLab's OAuth provider, allowing users to authenticate with their GitLab account.



google
------

- Type: `object` of `authentication.OAuthProvider` 

Google's OAuth provider, allowing users to authenticate with their Google account.



linkedin
--------

- Type: `object` of `authentication.OAuthProvider` 

LinkedIn's OAuth provider, allowing users to authenticate with their LinkedIn account.



microsoft
---------

- Type: `object` of `authentication.OAuthProvider` 

Microsoft's OAuth provider, allowing users to authenticate with their Microsoft account.



slack
-----

- Type: `object` of `authentication.OAuthProvider` 

Slack's OAuth provider, allowing users to authenticate with their Slack account.





OAuthProvider
=============



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



provider_token_management
-------------------------

- Type: `object` of `authentication.OAuthProviderTokenManagementAttribute` 

Settings related to token management for the OAuth provider.



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



claim_mapping
-------------

- Type: `map` of `string` 

Maps OAuth provider claims to Descope user attributes.

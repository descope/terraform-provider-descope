
OAuth
=====



disabled
--------

- Type: `bool` 

// description for disabled



system
------

- Type: `object` of `authentication.OAuthSystemProvider` 

// description for system



custom
------

- Type: `map` of `authentication.OAuthProvider` 

// description for custom





OAuthSystemProvider
===================



apple
-----

- Type: `object` of `authentication.OAuthProvider` 

// description for apple



discord
-------

- Type: `object` of `authentication.OAuthProvider` 

// description for discord



facebook
--------

- Type: `object` of `authentication.OAuthProvider` 

// description for facebook



github
------

- Type: `object` of `authentication.OAuthProvider` 

// description for github



gitlab
------

- Type: `object` of `authentication.OAuthProvider` 

// description for gitlab



google
------

- Type: `object` of `authentication.OAuthProvider` 

// description for google



linkedin
--------

- Type: `object` of `authentication.OAuthProvider` 

// description for linkedin



microsoft
---------

- Type: `object` of `authentication.OAuthProvider` 

// description for microsoft



slack
-----

- Type: `object` of `authentication.OAuthProvider` 

// description for slack





OAuthProvider
=============



disabled
--------

- Type: `bool` 

// description for disabled



client_id
---------

- Type: `string` 

// description for client_id



client_secret
-------------

- Type: `secret` 

// description for client_secret



provider_token_management
-------------------------

- Type: `object` of `authentication.OAuthProviderTokenManagementAttribute` 

// description for provider_token_management



prompts
-------

- Type: `list` of `string` 

// description for prompts



scopes
------

- Type: `list` of `string` 

// description for scopes



merge_user_accounts
-------------------

- Type: `bool` 
- Default: `true`

// description for merge_user_accounts



description
-----------

- Type: `string` 

// description for description



logo
----

- Type: `string` 

// description for logo



grant_type
----------

- Type: `string` 

// description for grant_type



issuer
------

- Type: `string` 

// description for issuer



authorization_endpoint
----------------------

- Type: `string` 

// description for authorization_endpoint



token_endpoint
--------------

- Type: `string` 

// description for token_endpoint



user_info_endpoint
------------------

- Type: `string` 

// description for user_info_endpoint



jwks_endpoint
-------------

- Type: `string` 

// description for jwks_endpoint



claim_mapping
-------------

- Type: `map` of `string` 

// description for claim_mapping


SSO
====



disabled
--------

- Type: `bool` 

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



merge_users
-----------

- Type: `bool` 

Whether to merge existing user accounts with new ones created through SSO authentication.



redirect_url
------------

- Type: `string` 

The URL the end user is redirected to after a successful authentication. If one is specified
in tenant level settings or SDK/API call, they will override this value.

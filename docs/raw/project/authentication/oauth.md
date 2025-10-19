
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

- Type: `map` of `authentication.OAuthProviderCustom` 

Custom OAuth providers configured for this project.


Passkeys
========



disabled
--------

- Type: `bool`

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



display_name
------------

- Type: `string`

The human-friendly name shown to users when they create or use a passkey. Some password
managers display this name, while others display the top level domain instead. When left
empty, the project name is used.



top_level_domain
----------------

- Type: `string`

Passkeys will be usable in the following domain and all its subdomains.



android_fingerprints
--------------------

- Type: `set` of `string`

A list of SHA-256 APK key hash fingerprints (colon-separated hex, e.g. `AB:CD:EF:...`) that are
allowed as passkey origins for Android apps. When set, only Android apps with a matching fingerprint
will be permitted to use passkey authentication.

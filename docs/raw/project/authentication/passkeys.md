
Passkeys
========



disabled
--------

- Type: `bool`

Setting this to `true` will disallow using this authentication method directly via
API and SDK calls. Note that this does not affect authentication flows that are
configured to use this authentication method.



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

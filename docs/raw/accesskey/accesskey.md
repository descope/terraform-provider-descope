
AccessKey
=========



project_id
----------

- Type: `string` (required)

The ID of the Descope project this access key belongs to. Changing this value will require the
resource to be deleted and recreated.



name
----

- Type: `string` (required)

A name for the access key.



description
-----------

- Type: `string`

A description for the access key.



status
------

- Type: `string`
- Default: `"active"`

The status of the access key. Must be either `active` or `inactive`. A new access key cannot be
created with an `inactive` status.



expire_time
-----------

- Type: `int`

The expiration time of the access key as a Unix timestamp. If not set, the key will not expire.
Changing this value after creation will require the access key to be replaced.



bound_user_id
-------------

- Type: `string`

The ID of a user to bind this access key to. When the key is exchanged for a session JWT, the
session acts on behalf of the bound user. Changing this value after creation will require the
access key to be replaced.



roles
-----

- Type: `list` of `string`

A list of project-level roles to grant to the access key. Cannot be used together with `tenants`.



tenants
-------

- Type: `list` of `accesskey.AccessKeyTenant`

A list of tenants to associate with the access key, each with its own set of roles. Cannot be
used together with `roles`.



custom_claims
-------------

- Type: `string`
- Default: `"{}"`

A JSON-encoded object of custom claims to add to the JWT created when the access key is exchanged.



custom_attributes
-----------------

- Type: `string`
- Default: `"{}"`

A JSON-encoded object of custom attribute values for the access key. The attributes must be defined
in the project's access key custom attribute schema.



permitted_ips
-------------

- Type: `list` of `string`

A list of IP addresses or CIDR ranges that are allowed to use this access key. If not set, the key
can be used from any IP address.



created_time
------------

- Type: `int`

The time the access key was created, as a Unix timestamp. This value is set by the server and is
read-only.



created_by
----------

- Type: `string`

The ID of the user or management key that created the access key. This value is set by the server
and is read-only.



cleartext
---------

- Type: `secret`

The plaintext value of the access key. This is only available after the key is created and cannot
be retrieved later. Store this value securely as it is required to exchange the key for a JWT.

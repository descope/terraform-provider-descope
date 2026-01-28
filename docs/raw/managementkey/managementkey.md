
ManagementKey
=============



name
----

- Type: `string` (required)

A name for the management key.



description
-----------

- Type: `string` 

A description for the management key.



status
------

- Type: `string` 
- Default: `"active"`

The status of the management key. Must be either `active` or `inactive`.



expire_time
-----------

- Type: `int` 

The expiration time of the management key as a Unix timestamp. If not set,
the key will not expire. Changing this value after creation will require
the management key to be replaced.



permitted_ips
-------------

- Type: `list` of `string` 

A list of IP addresses or CIDR ranges that are allowed to use this management key.
If not set, the key can be used from any IP address.



rebac
-----

- Type: `object` of `managementkey.ReBac` (required)

Access control settings for the management key. This defines the permissions granted
to the management key, either at the company level or for specific projects or for
project tags. Changing this value after creation will require the management key
to be replaced.



cleartext
---------

- Type: `secret` 

The plaintext value of the management key. This is only available after the key is
created and cannot be retrieved later. Store this value securely as it is required
to authenticate API requests.

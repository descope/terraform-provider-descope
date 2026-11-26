
List
====



project_id
----------

- Type: `string` (required)

The ID of the project that this list belongs to. Changing this value will require the resource
to be deleted and recreated.



name
----

- Type: `string` (required)

The name of the list. Maximum length is 100 characters.



description
-----------

- Type: `string`

An optional description for the list. Defaults to an empty string if not provided.



texts
-----

- Type: `set` of `string`

A set of text values. Exactly one of `texts`, `ips` or `json` must be set, and which one is set
determines the kind of list. Duplicate values are collapsed by the set, matching the backend which
stores each value once.



ips
----

- Type: `set` of `string`

A set of IP addresses and CIDR ranges. Each value must be a valid IP address or CIDR range.
Exactly one of `texts`, `ips` or `json` must be set.



json
----

- Type: `string`

An arbitrary JSON object, given as a JSON string. Exactly one of `texts`, `ips` or `json` must be
set. The document is compared by content, so reindenting it or reordering its keys is not a change.

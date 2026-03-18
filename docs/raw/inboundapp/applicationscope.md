
ApplicationScope
================



name
----

- Type: `string` (required)

A name for the scope.



description
-----------

- Type: `string` (required)

A description for the scope.



optional
--------

- Type: `bool`

Whether this scope is optional. When `false`, the scope is mandatory and must be granted during
authorization. When `true`, the user may choose to withhold it.



values
------

- Type: `list` of `string`

The identifiers of the relevant permission, attribute or connection scopes.

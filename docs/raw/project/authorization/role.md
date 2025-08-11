
Role
====



name
----

- Type: `string` (required)

A name for the role.



description
-----------

- Type: `string` 

A description for the role.



permissions
-----------

- Type: `set` of `string` 

A list of permissions by name to be included in the role.



default
-------

- Type: `bool` 

Whether this role should automatically be assigned to users that are created without any roles.



private
-------

- Type: `bool` 

Whether this role should not be displayed to tenant admins.

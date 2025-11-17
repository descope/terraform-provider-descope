
Role
====



key
----

- Type: `string` 

A persistent value that identifies a role uniquely across plan changes and configuration updates. It
is used exclusively by the Terraform provider during planning, to ensure that user roles are maintained
consistently even when role names or other details are changed. Once the `key` is set it should never be
changed, otherwise the role will be removed and a new one will be created instead.



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

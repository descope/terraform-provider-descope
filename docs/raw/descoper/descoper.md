
Descoper
========



email
-----

- Type: `string` (required)

The email address of the Descope console user.



phone
-----

- Type: `string` 

The phone number of the Descope console user.



name
----

- Type: `string` 

The display name of the Descope console user.



rbac
----

- Type: `object` of `descoper.RBac` (required)

Access control settings for the Descope console user. This defines the permissions
granted to the user, either as a company admin or for specific projects or project tags.


Project
=======



name
----

- Type: `string` (required)

The name of the Descope project.



environment
-----------

- Type: `string`

This can be set to `production` to mark production projects, otherwise this should be
left unset for development or staging projects.



deletion_protection
-------------------

- Type: `bool`

Protects the project from being accidentally destroyed. When this attribute isn't set,
deletion protection is enabled automatically for projects that have their `environment`
attribute set to `production`. To destroy a protected project, set this attribute to
`false` and apply the change first. Note that this only guards operations performed
through this provider, so removing the resource from the Terraform state is not prevented.



tags
----

- Type: `set` of `string`

Descriptive tags for your Descope project. Each tag must be no more than 50 characters long.

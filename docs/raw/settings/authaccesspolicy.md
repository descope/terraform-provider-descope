
AuthAccessPolicy
================



project_id
----------

- Type: `string` (required)

The ID of the project that this policy belongs to. Changing this value will require the
resource to be deleted and recreated.



default_action
--------------

- Type: `string`
- Default: `"allow"`

Whether an authentication that matches no rule is `allow`ed or `block`ed.



rule_ids
--------

- Type: `list` of `string`

The evaluation order of the project's rules, as a list of `descope_auth_access_rule` ids.
The list must include every rule in the project exactly once. Reference the `id` attribute
of each rule resource (e.g., `descope_auth_access_rule.example.id`) so that Terraform orders
the operations correctly, and avoid `create_before_destroy` on rule resources, which makes
Terraform update the order before a removed rule is deleted.

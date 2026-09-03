Resource Target
===============

type
----

- Type: `string` (required)

The resource type matched by this target. Must be `api`, `mcp`, `outbound_app`, or `any`.

all_of_type
-----------

- Type: `bool`
- Default: `false`

Whether this target matches every resource of the selected type. This capability may require a project feature flag.

ids
---

- Type: `list`
- Element: `string`

The specific resource IDs matched by this target when `all_of_type` is false.

Grant
=====

scopes
------

- Type: `list`
- Element: `string`

The scopes granted when this rule permits a request.

allowed_audiences
-----------------

- Type: `list`
- Element: `string`

The audiences allowed by this grant.

all_scopes
----------

- Type: `bool`
- Default: `false`

Whether this grant includes all scopes.

Condition
=========

key
---

- Type: `string` (required)

The policy context key evaluated by this condition.

operator
--------

- Type: `string` (required)

The comparison operator. Must be `equal`, `notEqual`, `contains`, `notContains`, `in`, or `notIn`.

value_json
----------

- Type: `string` (required)

A JSON-encoded condition value. The `in` and `notIn` operators require a JSON array; the other operators require a scalar JSON value.

Policy Rule
===========

project_id
----------

- Type: `string` (required)

The ID of the Descope project this policy rule belongs to. Changing this value requires replacement.

version
-------

- Type: `int`

The server-managed policy rule version used for optimistic concurrency. This value is read-only.

name
----

- Type: `string` (required)

A name for the policy rule.

description
-----------

- Type: `string`

A description for the policy rule.

enabled
-------

- Type: `bool` (required)

Whether the policy rule participates in authorization decisions.

rule_family
-----------

- Type: `string` (required)

The policy rule family. Must be `resource_access`, `outbound_access`, or `token_exchange`.

action_kinds
------------

- Type: `list` (required)
- Element: `string`

The action kinds matched by the rule. Supported values are `user_access`, `client_access`, `exchange_token`, and `fetch_outbound_token`.

effect
------

- Type: `string` (required)

The rule effect. Must be `permit` or `forbid`; `forbid` may require a project feature flag.

principal_type
--------------

- Type: `string` (required)

The principal type matched by the rule. Must be `any`, `user`, or `client`; `user` may require a project feature flag.

principal_selector
------------------

- Type: `list`
- Element: `string`

Optional principal IDs matched by the rule. An empty list matches any principal of the selected type.

resource_targets
----------------

- Type: `list`
- Element: `ResourceTargetModel`

The resource targets matched by the rule.

grants
------

- Type: `list`
- Element: `GrantModel`

The scopes and audiences granted by a permit rule. Token-exchange permit rules require at least one grant.

conditions
----------

- Type: `list`
- Element: `ConditionModel`

Structured conditions that are combined with logical AND.

cedar_text
----------

- Type: `string`

The Cedar policy source generated and validated by Descope. This value is read-only.

created_time
------------

- Type: `int`

The rule creation time as a Unix timestamp. This value is read-only.

modified_time
-------------

- Type: `int`

The rule modification time as a Unix timestamp. This value is read-only.

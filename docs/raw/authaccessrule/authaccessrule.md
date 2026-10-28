
AuthAccessRule
==============



project_id
----------

- Type: `string` (required)

The ID of the project that this rule belongs to. Changing this value will require the
resource to be deleted and recreated.



name
----

- Type: `string`

A name for the rule, shown in the console.



description
-----------

- Type: `string`

An optional description of what the rule is for.



enabled
-------

- Type: `bool`
- Default: `true`

Whether the rule is evaluated. A disabled rule keeps its place in the evaluation order
but never matches.



effect
------

- Type: `string` (required)

Whether an authentication that matches the rule is `allow`ed or `block`ed.



conditions
----------

- Type: `list` of `authaccessrule.AuthAccessCondition`

The conditions that must all match for the rule to apply. A rule with no conditions
matches every authentication.



expire_time
-----------

- Type: `int`

The time after which the rule is no longer evaluated, in seconds since the Unix epoch.
Set to `0` for a rule that never expires.





AuthAccessCondition
===================



key
----

- Type: `string` (required)

The attribute the condition checks. One of `request.ip`, `request.country`, `user.email`,
`user.emailDomain`, `user.loginId` or `user.phone`.



operator
--------

- Type: `string` (required)

How the attribute is compared with the values. One of `equal`, `notEqual`, `in`, `notIn`,
`inCidr` or `notInCidr`. The CIDR operators are only supported with the `request.ip` key.



values
------

- Type: `list` of `string` (required)

The values to compare with. The `equal`, `notEqual`, `inCidr` and `notInCidr` operators
take exactly one value. With `equal` and `notEqual` the value may instead be a reference of
the form `{{jwtClaims.<name>}}`, which compares the attribute with that custom claim of the
token the user presented. An absent claim never matches. Values must be in canonical form: lowercase emails and domains, uppercase
country codes, phone numbers in E.164 form and IP addresses in their shortest form.

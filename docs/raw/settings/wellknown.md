
WellKnownSettings
=================



project_id
----------

- Type: `string` (required)

The ID of the project that these settings belong to. Changing this value will require the resource
to be deleted and recreated.



enabled
-------

- Type: `bool`

Setting this to `true` applies the customization to the project's discovery documents. While `false` the
customization is kept but the default documents are served.



hidden_fields
-------------

- Type: `set` of `string`

Fields to omit from the discovery documents. Hiding a field only stops advertising it, the endpoint or
capability itself stays enabled. Supported values are `backchannel_authentication_endpoint`,
`pushed_authorization_request_endpoint`, `request_parameter_supported`, `registration_endpoint`,
`end_session_endpoint`, `introspection_endpoint`, `device_authorization_endpoint`, `revocation_endpoint`,
`userinfo_endpoint` and `backchannel_logout_supported`.



removed_values
--------------

- Type: `set` of `settings.WellKnownRemovedValues`

Values to remove from list fields in the discovery documents, such as a grant type from
`grant_types_supported`. A list is never emptied, removing every value keeps the full list.



scopes_supported
----------------

- Type: `object` of `settings.WellKnownValueDelta`

Scopes to add to or exclude from `scopes_supported`.



claims_supported
----------------

- Type: `object` of `settings.WellKnownValueDelta`

Claims to add to or exclude from `claims_supported`.





WellKnownRemovedValues
======================



field
-----

- Type: `string` (required)

The name of a list field in the discovery documents, such as `grant_types_supported`,
`token_endpoint_auth_methods_supported` or `code_challenge_methods_supported`.



values
------

- Type: `set` of `string` (required)

The values to remove from the list.





WellKnownValueDelta
===================



extra
-----

- Type: `set` of `string`

Values to add to the list.



excluded
--------

- Type: `set` of `string`

Values to remove from the list.

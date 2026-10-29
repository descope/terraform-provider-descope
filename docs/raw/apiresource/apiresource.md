
APIResource
===========



project_id
----------

- Type: `string` (required)

The ID of the project that this API resource belongs to. Changing this value will require the resource
to be deleted and recreated.



name
----

- Type: `string` (required)

A name for the API resource.



description
-----------

- Type: `string`

An optional description for the API resource.



uri
----

- Type: `string` (required)

The URI that identifies the API or MCP server, used as the token audience.



user_access
-----------

- Type: `string`

Whether users may be granted access to the resource: `allow`, `deny` or `allowViaPolicy`. Empty uses the
project default.



client_access
-------------

- Type: `string`

Whether clients may be granted access to the resource: `allow`, `deny` or `allowViaPolicy`. Empty uses the
project default.



permissions_scopes
------------------

- Type: `list` of `inboundapp.ApplicationScope`

Permission scopes that inbound applications can request for this resource.



attributes_scopes
-----------------

- Type: `list` of `inboundapp.ApplicationScope`

User attribute scopes that inbound applications can request for this resource.



connections_scopes
------------------

- Type: `list` of `inboundapp.ApplicationScope`

Connection scopes that inbound applications can request for this resource.



id_jag_target
-------------

- Type: `object` of `apiresource.IDJAGTarget`

Cross App Access settings for an external server protected by another authorization server. Inbound
applications that exchange a token for this resource receive an ID-JAG assertion for that authorization server.

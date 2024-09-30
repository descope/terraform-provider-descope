
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



project_settings
----------------

- Type: `object` of `settings.Settings` 

General settings for the Descope project.



authentication
--------------

- Type: `object` of `authentication.Authentication` 

Settings for each authentication method.



authorization
-------------

- Type: `object` of `authorization.Authorization` 

Define Role-Based Access Control (RBAC) for your users by creating roles and permissions.



attributes
----------

- Type: `object` of `attributes.Attributes` 

Custom attributes that can be attached to users and tenants.



connectors
----------

- Type: `object` of `connectors.Connectors` 

Enrich your flows by interacting with third party services.



applications
------------

- Type: `object` of `applications.Application` 

Applications that are registered with the project.



jwt_templates
-------------

- Type: `object` of `jwttemplates.JWTTemplates` 

Defines templates for JSON Web Tokens (JWT) used for authentication.



styles
------

- Type: `object` of `flows.Styles` 

Custom styles that can be applied to the project's authentication flows.



flows
-----

- Type: `map` of `flows.Flow` 

Custom authentication flows to use in this project.

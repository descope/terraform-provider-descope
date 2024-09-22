
Project
=======



name
----

- Type: `string` (required)

The name of the Descope project.



environment
-----------

- Type: `string` 

This can be set to `production` to mark production projects, otherwise this should be left unset for development or staging projects.



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

// Specify the `Authorization` object.



attributes
----------

- Type: `object` of `attributes.Attributes` 

// Specify the `Attributes` object.



connectors
----------

- Type: `object` of `connectors.Connectors` 

// Specify the `Connectors` object.



applications
------------

- Type: `object` of `applications.Application` 

// Specify the `Application` object.



jwt_templates
-------------

- Type: `object` of `jwttemplates.JWTTemplates` 

// Specify the `JWTTemplates` object.



styles
------

- Type: `object` of `flows.Styles` 

// Specify the `Styles` object.



flows
-----

- Type: `map` of `flows.Flow` 

// Specify the `Flow` object.

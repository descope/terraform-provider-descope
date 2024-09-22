
SAML
====



id
----

- Type: `string` 

// Specify an identifier for the SAML application.



name
----

- Type: `string` (required)

// Specify a name for the SAML application.



description
-----------

- Type: `string` 

// Specify a description for the SAML application.



logo
----

- Type: `string` 

// Specify a logo for the SAML application. Should be a hosted image URL.



disabled
--------

- Type: `bool` 

// Specify whether the application should be enabled or disabled.



login_page_url
--------------

- Type: `string` 

// Specify the Flow Hosting URL. Read more about using this parameter with custom domain [here](https://docs.descope.com/sso-integrations/applications/saml-apps).



dynamic_configuration
---------------------

- Type: `object` of `applications.DynamicConfiguration` 

// Specify the `DynamicConfiguration` object. Read the description below.



manual_configuration
--------------------

- Type: `object` of `applications.ManualConfiguration` 

// Specify the `ManualConfiguration` object. Read the description below.



acs_allowed_callback_urls
-------------------------

- Type: `list` of `string` 

// Specify a list of allowed ACS callback URLS. This configuration is used when the default ACS URL value is unreachable. Supports wildcards.



subject_name_id_type
--------------------

- Type: `string` 

// Specify the subject name id type. Choose one of "", "email", "phone". Read more about this configuration [here](https://docs.descope.com/sso-integrations/applications/saml-apps).



subject_name_id_format
----------------------

- Type: `string` 

// Specify the subject name id format. Choose one of "", "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified", "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress", "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent", "urn:oasis:names:tc:SAML:2.0:nameid-format:transient". Read more about this configuration [here](https://docs.descope.com/sso-integrations/applications/saml-apps).



default_relay_state
-------------------

- Type: `string` 

// Specify the default relay state. When using IdP-initiated authentication, this value may be used as a URL to a resource in the Service Provider.



attribute_mapping
-----------------

- Type: `list` of `applications.AttributeMapping` 

// Specify the `AttributeMapping` object. Read the description below.





AttributeMapping
================



name
----

- Type: `string` (required)

// Specify the name of the attribute.



value
-----

- Type: `string` (required)

// Specify the value of the attribute.





DynamicConfiguration
====================



metadata_url
------------

- Type: `string` (required)

// Specify the metadata URL when retrieving the connection details dynamically.





ManualConfiguration
===================



acs_url
-------

- Type: `string` (required)

// Enter the `ACS URL` from the SP.



entity_id
---------

- Type: `string` (required)

// Enter the `Entity Id` from the SP.



certificate
-----------

- Type: `string` (required)

// Enter the `Certificate` from the SP.

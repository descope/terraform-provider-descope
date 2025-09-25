
SNS
====



access_key_id
-------------

- Type: `secret` (required)

AWS Access key ID.



secret
------

- Type: `secret` (required)

AWS Secret Access Key.



region
------

- Type: `string` (required)

AWS region to send requests to (e.g. `us-west-2`).



endpoint
--------

- Type: `string` 

An optional endpoint URL (hostname only or fully qualified URI).



origination_number
------------------

- Type: `string` 

An optional phone number from which the text messages are going to be sent. Make sure it
is registered properly in your server.



sender_id
---------

- Type: `string` 

The name of the sender from which the text message is going to be sent (see SNS documentation
regarding acceptable IDs and supported regions/countries).



entity_id
---------

- Type: `string` 

The entity ID or principal entity (PE) ID for sending text messages to recipients in India.



template_id
-----------

- Type: `string` 

The template for sending text messages to recipients in India. The template ID must be
associated with the sender ID.



organization_number
-------------------

- Type: `string` 

Deprecated, use the `origination_number` attribute instead.


HttpStaticIP
============



base_url
--------

- Type: `string` (required)

The base URL to fetch



authentication
--------------

- Type: `object` of `connectors.HTTPAuthField` 

Authentication Information



headers
-------

- Type: `map` of `string` 

The headers to send with the request



hmac_secret
-----------

- Type: `secret` 

HMAC is a method for message signing with a symmetrical key. This secret will be
used to sign the base64 encoded payload, and the resulting signature will be
sent in the `x-descope-webhook-s256` header. The receiving service should use
this secret to verify the integrity and authenticity of the payload by checking
the provided signature



insecure
--------

- Type: `bool` 

Will ignore certificate errors raised by the client



include_headers_in_context
--------------------------

- Type: `bool` 

The connector response context will also include the headers. The context will
have a "body" attribute and a "headers" attribute. See more details in the help
guide


List
====



name
----

- Type: `string` (required)

The name of the list. Maximum length is 100 characters.



description
-----------

- Type: `string` 

An optional description for the list. Defaults to an empty string if not provided.



type
----

- Type: `string` (required)

The type of list. Must be one of:
- `"texts"` - A list of text strings
- `"ips"` - A list of IP addresses or CIDR ranges
- `"json"` - A JSON object



data
----

- Type: `string` (required)

The JSON data for the list. The format depends on the `type`:
- For `"texts"` and `"ips"` types: Must be a JSON array of strings (e.g., `["item1", "item2"]`)
- For `"ips"` type: Each string must be a valid IP address or CIDR range
- For `"json"` type: Must be a JSON object (e.g., `{"key": "value"}`)

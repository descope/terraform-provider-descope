
## types

This package provides generic object and collection attribute types to augment the basic types
in [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework/tree/main/types).

The `objtype` package is meant to be used instead of `types.Object` to represent model objects,
whereas collections of model objects can be represented using the `listtype`, `settype`, and
`maptype` packages.

The `valuelisttype`, `valuesettype`, and `valuemaptype` packages provide types for simple
collections of plain values, in particular for `types.String` values.

For troubleshooting and references, you can consult the Terraform documentation and
similar implementations:

- https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes
- https://github.com/hashicorp/terraform-provider-aws/tree/main/internal/framework/types
- https://github.com/cloudflare/terraform-provider-cloudflare/tree/main/internal/customfield

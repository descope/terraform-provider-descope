package helpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const importKey = "descopeImport"

// Sets a private key in the response to indicate that the resource is being imported.
func MarkImportState(ctx context.Context, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Setting import key for resource")
	d := resp.Private.SetKey(ctx, importKey, []byte(`{}`))
	resp.Diagnostics.Append(d...)
}

// Used in read operations to check if the import key is set in the state and mark the context accordingly.
func ContextWithImportState(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) context.Context {
	if value, _ := req.Private.GetKey(ctx, importKey); value != nil {
		tflog.Info(ctx, "Unsetting import key for reading resource")
		ctx = context.WithValue(ctx, importKey, true)
		d := resp.Private.SetKey(ctx, importKey, nil)
		resp.Diagnostics.Append(d...)
	}
	return ctx
}

// Checks if we're currently reading a source as part of an import operation.
func isImportState(ctx context.Context) bool {
	value, _ := ctx.Value(importKey).(bool)
	return value
}

// If an attribute value is already set to Null we do not overwrite it, so we don't not cause
// provider errors due to inconsistent values before and after the plan is applied. We also
// don't want to read/write state the plan writer isn't interested in, and left the attribute
// as its default value of Null.
//
// When importing resources though the attribute values are Null because there's no state yet,
// unlike during other operations where they are Unknown if not set and the attribute doesn't
// have a default value. In this case we want to update the state and set the object.
func ShouldSetAttributeValue(ctx context.Context, v attr.Value) bool {
	if !v.IsNull() {
		return true
	}
	if isImportState(ctx) {
		return true
	}
	return false
}

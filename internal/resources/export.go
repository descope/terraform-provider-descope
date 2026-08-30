package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExportableResource lets the tfexport tool drive the resource's import read path without a provider server; like DestroyChecker
// this is tooling support, not part of the resource's Terraform-facing behavior.
type ExportableResource interface {
	ExportName() string
	ExportSchema() schema.Schema
	ExportSingleton() bool
	ExportRead(ctx context.Context, client *infra.Client, projectID, scope, id string) (model any, found bool, diags diag.Diagnostics)
	ExportWireType(ctx context.Context) string
	ExportValidate(ctx context.Context, object types.Object) diag.Diagnostics
}

func (r *baseResource[T, M]) ExportName() string {
	return r.name
}

func (r *baseResource[T, M]) ExportSchema() schema.Schema {
	return r.schema
}

func (r *baseResource[T, M]) ExportSingleton() bool {
	return r.singleton
}

func (r *baseResource[T, M]) ExportWireType(ctx context.Context) (wireType string) {
	defer func() {
		_ = recover() // some models can't produce Values from a zero value (e.g. flows)
	}()
	var diags diag.Diagnostics
	model := M(new(T))
	data := model.Values(helpers.NewHandler(ctx, &diags))
	wireType, _ = data["type"].(string)
	return
}

func (r *baseResource[T, M]) ExportValidate(ctx context.Context, object types.Object) diag.Diagnostics {
	var diags diag.Diagnostics
	model := M(new(T))
	validatable, ok := any(model).(validatableModel)
	if !ok {
		return diags
	}
	value, err := object.ToTerraformValue(ctx)
	if err != nil {
		diags.AddError("Error converting "+r.name+" value", err.Error())
		return diags
	}
	state := tfsdk.State{Schema: r.schema, Raw: value}
	diags.Append(state.Get(ctx, model)...)
	if diags.HasError() {
		return diags
	}
	validatable.Validate(helpers.NewHandler(ctx, &diags))
	return diags
}

func (r *baseResource[T, M]) ExportRead(ctx context.Context, client *infra.Client, projectID, scope, id string) (any, bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	ctx = helpers.MarkImportContext(ctx)

	if r.singleton {
		id = projectID
	}

	var data map[string]any
	var err error
	if r.ops.ScopedRead != nil {
		data, err = r.ops.ScopedRead(ctx, client, projectID, scope, id)
	} else {
		data, err = r.ops.Read(ctx, client, projectID, id)
	}
	if infra.AsNotFoundError(err) {
		return nil, false, diags
	}
	if err != nil {
		diags.AddError("Error reading "+r.name, err.Error())
		return nil, false, diags
	}

	model := M(new(T))
	model.SetID(types.StringValue(id))
	model.SetValues(helpers.NewHandler(ctx, &diags), data)
	return model, true, diags
}

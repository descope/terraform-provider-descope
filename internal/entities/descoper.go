package entities

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/descoper"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var DescoperSchema = schema.Schema{
	Attributes: descoper.DescoperAttributes,
}

type DescoperEntity struct {
	Model       *descoper.DescoperModel
	Diagnostics *diag.Diagnostics
}

// Creates a new descoper entity by loading data from the source Terraform plan or state.
func NewDescoperEntity(ctx context.Context, source entitySource, diagnostics *diag.Diagnostics) *DescoperEntity {
	e := &DescoperEntity{Model: &descoper.DescoperModel{}, Diagnostics: diagnostics}
	load(ctx, source, e.Model, e.Diagnostics)
	return e
}

// Saves the descoper entity data to the target Terraform state.
func (e *DescoperEntity) Save(ctx context.Context, target entityTarget) {
	save(ctx, target, e.Model, e.Diagnostics)
}

// Returns a representation of the descoper entity data for sending in an infra API request.
func (e *DescoperEntity) Values(ctx context.Context) map[string]any {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	return e.Model.Values(handler)
}

// Updates the descoper entity with the data received in an infra API response.
func (e *DescoperEntity) SetValues(ctx context.Context, data map[string]any) {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	e.Model.SetValues(handler, data)
}

// Returns the ID value from the model.
func (e *DescoperEntity) ID(_ context.Context) string {
	return e.Model.ID.ValueString()
}

// Sets the ID value in the model, for use after creation.
func (e *DescoperEntity) SetID(_ context.Context, id string) {
	e.Model.ID = types.StringValue(id)
}

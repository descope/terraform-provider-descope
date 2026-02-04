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

func NewDescoperEntity(ctx context.Context, source entitySource, diagnostics *diag.Diagnostics) *DescoperEntity {
	e := &DescoperEntity{Model: &descoper.DescoperModel{}, Diagnostics: diagnostics}
	load(ctx, source, e.Model, e.Diagnostics)
	return e
}

func (e *DescoperEntity) Save(ctx context.Context, target entityTarget) {
	save(ctx, target, e.Model, e.Diagnostics)
}

func (e *DescoperEntity) Values(ctx context.Context) map[string]any {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	return e.Model.Values(handler)
}

func (e *DescoperEntity) SetValues(ctx context.Context, data map[string]any) {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	e.Model.SetValues(handler, data)
}

func (e *DescoperEntity) ID(_ context.Context) string {
	return e.Model.ID.ValueString()
}

func (e *DescoperEntity) SetID(_ context.Context, id string) {
	e.Model.ID = types.StringValue(id)
}

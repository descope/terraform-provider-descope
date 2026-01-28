package entities

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/managementkey"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ManagementKeySchema = schema.Schema{
	Attributes: managementkey.ManagementKeyAttributes,
}

type ManagementKeyEntity struct {
	Model       *managementkey.ManagementKeyModel
	Diagnostics *diag.Diagnostics
}

func NewManagementKeyEntity(ctx context.Context, source entitySource, diagnostics *diag.Diagnostics) *ManagementKeyEntity {
	e := &ManagementKeyEntity{Model: &managementkey.ManagementKeyModel{}, Diagnostics: diagnostics}
	load(ctx, source, e.Model, e.Diagnostics)
	return e
}

func (e *ManagementKeyEntity) Save(ctx context.Context, target entityTarget) {
	save(ctx, target, e.Model, e.Diagnostics)
}

func (e *ManagementKeyEntity) Values(ctx context.Context) map[string]any {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	return e.Model.Values(handler)
}

func (e *ManagementKeyEntity) SetValues(ctx context.Context, data map[string]any) {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	e.Model.SetValues(handler, data)
}

func (e *ManagementKeyEntity) ID(_ context.Context) string {
	return e.Model.ID.ValueString()
}

func (e *ManagementKeyEntity) SetID(_ context.Context, id string) {
	e.Model.ID = types.StringValue(id)
}

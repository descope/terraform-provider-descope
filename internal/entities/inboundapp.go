package entities

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/inboundapp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var InboundAppSchema = schema.Schema{
	Attributes: inboundapp.InboundAppAttributes,
}

type InboundAppEntity struct {
	Model       *inboundapp.InboundAppModel
	Diagnostics *diag.Diagnostics
}

func NewInboundAppEntity(ctx context.Context, source entitySource, diagnostics *diag.Diagnostics) *InboundAppEntity {
	e := &InboundAppEntity{Model: &inboundapp.InboundAppModel{}, Diagnostics: diagnostics}
	load(ctx, source, e.Model, e.Diagnostics)
	return e
}

func (e *InboundAppEntity) Save(ctx context.Context, target entityTarget) {
	save(ctx, target, e.Model, e.Diagnostics)
}

func (e *InboundAppEntity) Values(ctx context.Context) map[string]any {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	return e.Model.Values(handler)
}

func (e *InboundAppEntity) SetValues(ctx context.Context, data map[string]any) {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	e.Model.SetValues(handler, data)
}

func (e *InboundAppEntity) ID(_ context.Context) string {
	return e.Model.ID.ValueString()
}

func (e *InboundAppEntity) SetID(_ context.Context, id string) {
	e.Model.ID = types.StringValue(id)
}

func (e *InboundAppEntity) ProjectID(_ context.Context) string {
	return e.Model.ProjectID.ValueString()
}

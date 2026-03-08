package entities

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/accesskey"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AccessKeySchema = schema.Schema{
	Attributes: accesskey.AccessKeyAttributes,
}

type AccessKeyEntity struct {
	Model       *accesskey.AccessKeyModel
	Diagnostics *diag.Diagnostics
}

func NewAccessKeyEntity(ctx context.Context, source entitySource, diagnostics *diag.Diagnostics) *AccessKeyEntity {
	e := &AccessKeyEntity{Model: &accesskey.AccessKeyModel{}, Diagnostics: diagnostics}
	load(ctx, source, e.Model, e.Diagnostics)
	return e
}

func (e *AccessKeyEntity) Save(ctx context.Context, target entityTarget) {
	save(ctx, target, e.Model, e.Diagnostics)
}

func (e *AccessKeyEntity) Values(ctx context.Context) map[string]any {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	return e.Model.Values(handler)
}

func (e *AccessKeyEntity) SetValues(ctx context.Context, data map[string]any) {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	e.Model.SetValues(handler, data)
}

func (e *AccessKeyEntity) ID(_ context.Context) string {
	return e.Model.ID.ValueString()
}

func (e *AccessKeyEntity) SetID(_ context.Context, id string) {
	e.Model.ID = types.StringValue(id)
}

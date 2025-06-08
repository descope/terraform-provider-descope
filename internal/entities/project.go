package entities

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/project"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ProjectSchema = schema.Schema{
	Attributes: project.ProjectAttributes,
}

type ProjectEntity struct {
	Model       *project.ProjectModel
	Diagnostics *diag.Diagnostics
}

// Creates a new project entity by loading data from the source Terraform plan or state.
func NewProjectEntity(ctx context.Context, source entitySource, diagnostics *diag.Diagnostics) *ProjectEntity {
	e := &ProjectEntity{Model: &project.ProjectModel{}, Diagnostics: diagnostics}
	load(ctx, source, e.Model, e.Diagnostics)
	return e
}

// Saves the project entity data to the target Tarraform state.
func (e *ProjectEntity) Save(ctx context.Context, target entityTarget) {
	save(ctx, target, e.Model, e.Diagnostics)
}

// Returns a representation of the project entity data for sending in an infra API request.
func (e *ProjectEntity) Values(ctx context.Context) map[string]any {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	// collect all existing references from the plan
	e.Model.CollectReferences(handler)
	// convert the model to a backend request format
	values := e.Model.Values(handler)
	return values
}

// Updates the project entity with the data received in an infra API response.
func (e *ProjectEntity) SetValues(ctx context.Context, data map[string]any) {
	handler := helpers.NewHandler(ctx, e.Diagnostics)
	// collect all existing references from the plan or state
	e.Model.CollectReferences(handler)
	// update the model with the new values from the backend response
	e.Model.SetValues(handler, data)
	// collect the references again after the model has been updated
	e.Model.CollectReferences(handler)
	// apply the references to replace any server IDs
	e.Model.UpdateReferences(handler)
}

// Returns the projectID value from the model.
func (e *ProjectEntity) ProjectID(_ context.Context) string {
	return e.Model.ID.ValueString()
}

// Sets the projectID value in the model, for use after project creation.
func (e *ProjectEntity) SetProjectID(_ context.Context, id string) {
	e.Model.ID = types.StringValue(id)
}

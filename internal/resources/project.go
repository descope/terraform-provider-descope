package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/project"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// The project container resource on its dedicated /v1/mgmt/project route. The resource id is the project
// id: reads, updates and deletes scope the bearer token with it, while create runs with a bare key and
// gets the new id from the response.
func NewProjectResource() resource.Resource {
	const path = "/v1/mgmt/project"
	return newResource[project.ProjectModel]("project", project.Schema, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			body, err := c.PostData(ctx, projectID, path, data)
			if err != nil {
				return "", nil, err
			}
			id, _ := body["id"].(string)
			return id, body, nil
		},
		Read: func(ctx context.Context, c *infra.Client, projectID, _ string) (map[string]any, error) {
			return c.Get(ctx, projectID, path, nil)
		},
		Update: func(ctx context.Context, c *infra.Client, projectID, _ string, data map[string]any) (map[string]any, error) {
			return c.PutData(ctx, projectID, path, data)
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, _ string) error {
			return c.Del(ctx, projectID, path, nil)
		},
	})
}

package resources

import (
	"context"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// The list endpoints nest the entity under a "list" key in responses, read with a path parameter, and delete with the id in a POST body.
func NewListResource() resource.Resource {
	return newResource[list.ListModel]("list", list.Schema, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			body, err := c.PostData(ctx, projectID, "/v1/mgmt/list", data)
			if err != nil {
				return "", nil, err
			}
			entity, err := unwrapEntity(body, "list")
			if err != nil {
				return "", nil, err
			}
			id, _ := entity["id"].(string)
			return id, entity, nil
		},
		Read: func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
			body, err := c.Get(ctx, projectID, "/v1/mgmt/list/"+id, nil)
			if err != nil {
				return nil, err
			}
			return unwrapEntity(body, "list")
		},
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			payload := maps.Clone(data)
			payload["id"] = id
			body, err := c.PostData(ctx, projectID, "/v1/mgmt/list/update", payload)
			if err != nil {
				return nil, err
			}
			return unwrapEntity(body, "list")
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Post(ctx, projectID, "/v1/mgmt/list/delete", map[string]any{"id": id})
		},
	})
}

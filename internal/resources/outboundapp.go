package resources

import (
	"context"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/outboundapp"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Outbound apps live on dedicated CRUD endpoints with a server-assigned id. Every response wraps the entity
// in an `app` object, the load endpoint takes the id in the path rather than a query, and update expects the
// whole entity nested under `app` alongside its id.
func NewOutboundAppResource() resource.Resource {
	const path = "/v1/mgmt/outbound/app"

	read := func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
		body, err := c.Get(ctx, projectID, path+"/"+id, nil)
		if err != nil {
			return nil, err
		}
		return unwrapEntity(body, "app")
	}

	return newResource[outboundapp.OutboundAppModel]("outbound_app", outboundapp.OutboundAppSchema, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			body, err := c.PostData(ctx, projectID, path+"/create", data)
			if err != nil {
				return "", nil, err
			}
			entity, err := unwrapEntity(body, "app")
			if err != nil {
				return "", nil, err
			}
			id, _ := entity["id"].(string)
			return id, entity, nil
		},
		Read: read,
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			app := maps.Clone(data)
			app["id"] = id
			body, err := c.PostData(ctx, projectID, path+"/update", map[string]any{"app": app})
			if err != nil {
				return nil, err
			}
			return unwrapEntity(body, "app")
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Post(ctx, projectID, path+"/delete", map[string]any{"id": id})
		},
	})
}

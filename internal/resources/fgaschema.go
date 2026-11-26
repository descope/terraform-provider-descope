package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/fgaschema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// The FGA schema is a singleton DSL document; clearing it goes through the authz schema delete endpoint,
// both on destroy and when the dsl attribute is set to an empty string.
func NewFGASchemaResource() resource.Resource {
	read := func(ctx context.Context, c *infra.Client, projectID, _ string) (map[string]any, error) {
		return c.Get(ctx, projectID, "/v1/mgmt/fga/schema", nil)
	}
	clearSchema := func(ctx context.Context, c *infra.Client, projectID string) error {
		return c.Post(ctx, projectID, "/v1/mgmt/authz/schema/delete", map[string]any{})
	}
	return newSingletonResource[fgaschema.FGASchemaModel]("fga_schema", fgaschema.Schema, operations{
		Read: read,
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			if dsl, _ := data["dsl"].(string); dsl == "" {
				if err := clearSchema(ctx, c, projectID); err != nil {
					return nil, err
				}
			} else if err := c.Post(ctx, projectID, "/v1/mgmt/fga/schema", data); err != nil {
				return nil, err
			}
			return read(ctx, c, projectID, id)
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, _ string) error {
			return clearSchema(ctx, c, projectID)
		},
	})
}

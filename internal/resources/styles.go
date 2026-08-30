package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/styles"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Styles are a singleton stored behind the theme export/import endpoints, with the entity nested under a
// "theme" key in both request and response bodies.
func NewStylesResource() resource.Resource {
	read := func(ctx context.Context, c *infra.Client, projectID, _ string) (map[string]any, error) {
		body, err := c.PostData(ctx, projectID, "/v2/mgmt/theme/export", map[string]any{})
		if err != nil {
			return nil, err
		}
		return unwrapEntity(body, "theme")
	}
	return newSingletonResource[styles.StylesModel]("styles", styles.Schema, operations{
		Read: read,
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			if err := c.Post(ctx, projectID, "/v2/mgmt/theme/import", map[string]any{"theme": data}); err != nil {
				return nil, err
			}
			return read(ctx, c, projectID, id)
		},
		Delete: noDelete,
	})
}

package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/widget"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Widgets are owned by a user-chosen id the model sends as widgetId in the body, but reads and deletes use
// an `id` query param. Writes return no body, so they re-read to pick up server-populated fields.
func NewWidgetResource() resource.Resource {
	const path = "/v1/mgmt/widget"
	read := func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
		return c.Get(ctx, projectID, path, map[string]string{"id": id})
	}
	write := func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
		if err := c.Post(ctx, projectID, path, data); err != nil {
			return nil, err
		}
		return read(ctx, c, projectID, id)
	}
	return newResource[widget.WidgetModel]("widget", widget.Schema, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			id, _ := data["widgetId"].(string)
			entity, err := write(ctx, c, projectID, id, data)
			return id, entity, err
		},
		Read:   read,
		Update: write,
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Del(ctx, projectID, path, map[string]string{"id": id})
		},
	})
}

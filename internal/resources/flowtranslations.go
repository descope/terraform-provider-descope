package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/flowtranslations"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewFlowTranslationsResource() resource.Resource {
	const path = "/v1/mgmt/localization/flow"
	read := func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
		return c.PostData(ctx, projectID, path+"/export", map[string]any{"flowId": id})
	}
	write := func(ctx context.Context, c *infra.Client, projectID, _ string, data map[string]any) (map[string]any, error) {
		return c.PostData(ctx, projectID, path+"/import", data)
	}
	return newResource[flowtranslations.FlowTranslationsModel]("flow_translations", flowtranslations.Schema, operations{
		CreateOverwrites: true,
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			id, _ := data["flowId"].(string)
			entity, err := write(ctx, c, projectID, id, data)
			return id, entity, err
		},
		Read:   read,
		Update: write,
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			empty := map[string]any{"flowId": id, "overrides": map[string]any{"flowId": id, "languages": map[string]any{}}}
			if _, err := write(ctx, c, projectID, id, empty); err != nil {
				if _, readErr := read(ctx, c, projectID, id); infra.AsNotFoundError(readErr) {
					return nil
				}
				return err
			}
			return nil
		},
	})
}

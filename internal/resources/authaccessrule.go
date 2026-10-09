package resources

import (
	"context"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/authaccessrule"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewAuthAccessRuleResource() resource.Resource {
	const path = "/v1/mgmt/authaccess/rule"
	return newResource[authaccessrule.AuthAccessRuleModel]("auth_access_rule", authaccessrule.Schema, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			body, err := c.PostData(ctx, projectID, path, data)
			if err != nil {
				return "", nil, err
			}
			id, _ := body["id"].(string)
			return id, body, nil
		},
		Read: func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
			return c.Get(ctx, projectID, path, map[string]string{"id": id})
		},
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			body := maps.Clone(data)
			body["id"] = id
			return c.PutData(ctx, projectID, path, body)
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Del(ctx, projectID, path, map[string]string{"id": id})
		},
	})
}

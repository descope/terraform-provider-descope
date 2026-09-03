package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/oauthprovider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// OAuth providers are owned by a user-chosen id the model sends in the body, and the write endpoint echoes
// the stored entity so no follow-up read is needed.
func NewOAuthProviderResource() resource.Resource {
	const path = "/v1/mgmt/oauth/provider"
	write := func(ctx context.Context, c *infra.Client, projectID, _ string, data map[string]any) (map[string]any, error) {
		return c.PostData(ctx, projectID, path, data)
	}
	return newResource[oauthprovider.OAuthProviderModel]("oauth_provider", oauthprovider.Schema, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			id, _ := data["id"].(string)
			entity, err := write(ctx, c, projectID, id, data)
			return id, entity, err
		},
		Read: func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
			return c.Get(ctx, projectID, path, map[string]string{"id": id})
		},
		Update: write,
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Del(ctx, projectID, path, map[string]string{"id": id})
		},
	})
}

package resources

import (
	"context"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/jwttemplate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// JWT templates use POST-per-verb endpoints that nest the entity under a "template" key in both request and response bodies.
func NewJWTTemplateResource() resource.Resource {
	return newResource[jwttemplate.JWTTemplateModel]("jwt_template", jwttemplate.Schema, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			body, err := c.PostData(ctx, projectID, "/v1/mgmt/jwt/templates/create", map[string]any{"template": data})
			if err != nil {
				return "", nil, err
			}
			entity, err := unwrapEntity(body, "template")
			if err != nil {
				return "", nil, err
			}
			id, _ := entity["id"].(string)
			return id, entity, nil
		},
		Read: func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
			body, err := c.PostData(ctx, projectID, "/v1/mgmt/jwt/templates/load", map[string]any{"id": id})
			if err != nil {
				return nil, err
			}
			return unwrapEntity(body, "template")
		},
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			template := maps.Clone(data)
			template["id"] = id
			body, err := c.PostData(ctx, projectID, "/v1/mgmt/jwt/templates/update", map[string]any{"template": template})
			if err != nil {
				return nil, err
			}
			return unwrapEntity(body, "template")
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Post(ctx, projectID, "/v1/mgmt/jwt/templates/delete", map[string]any{"id": id})
		},
	})
}

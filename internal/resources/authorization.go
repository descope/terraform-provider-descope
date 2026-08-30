package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/permission"
	"github.com/descope/terraform-provider-descope/internal/models/role"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewRoleResource() resource.Resource {
	// roles have no single-entity load endpoint, so reads go through the search endpoint and match by id
	read := func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
		body, err := c.PostData(ctx, projectID, "/v1/mgmt/role/search", map[string]any{"roleIds": []string{id}})
		if err != nil {
			return nil, err
		}
		roles, _ := body["roles"].([]any)
		for _, r := range roles {
			if entity, ok := r.(map[string]any); ok && entity["id"] == id {
				return entity, nil
			}
		}
		return nil, fmt.Errorf("role with id %s: %w", id, infra.ErrNotFound)
	}

	return newResource[role.RoleModel]("role", role.Schema, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			entity, err := c.PostData(ctx, projectID, "/v1/mgmt/role/create", data)
			if err != nil {
				return "", nil, err
			}
			id, _ := entity["id"].(string)
			return id, entity, nil
		},
		Read: read,
		// the update endpoint identifies the role by exactly one of id or name and expects the planned name in newName
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			body := maps.Clone(data)
			body["id"] = id
			body["newName"] = body["name"]
			delete(body, "name")
			return c.PostData(ctx, projectID, "/v1/mgmt/role/update", body)
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Post(ctx, projectID, "/v1/mgmt/role/delete", map[string]any{"id": id})
		},
	})
}

func NewPermissionResource() resource.Resource {
	read := func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
		return findPermission(ctx, c, projectID, "id", id)
	}

	return newResource[permission.PermissionModel]("permission", permission.Schema, operations{
		// the create endpoint returns no body, so the new permission is located by its name to discover the server-assigned id
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			if err := c.Post(ctx, projectID, "/v1/mgmt/permission/create", data); err != nil {
				return "", nil, err
			}
			name, _ := data["name"].(string)
			entity, err := findPermission(ctx, c, projectID, "name", name)
			if err != nil {
				return "", nil, err
			}
			id, _ := entity["id"].(string)
			return id, entity, nil
		},
		Read: read,
		// the update endpoint identifies the permission by exactly one of id or name and expects the planned name in newName
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			body := maps.Clone(data)
			body["id"] = id
			body["newName"] = body["name"]
			delete(body, "name")
			if err := c.Post(ctx, projectID, "/v1/mgmt/permission/update", body); err != nil {
				return nil, err
			}
			return read(ctx, c, projectID, id)
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Post(ctx, projectID, "/v1/mgmt/permission/delete", map[string]any{"id": id})
		},
	})
}

// findPermission locates a permission in the full list because there's no single-entity load endpoint.
func findPermission(ctx context.Context, c *infra.Client, projectID, key, value string) (map[string]any, error) {
	body, err := c.Get(ctx, projectID, "/v1/mgmt/permission/all", nil)
	if err != nil {
		return nil, err
	}
	permissions, _ := body["permissions"].([]any)
	for _, p := range permissions {
		if entity, ok := p.(map[string]any); ok && entity[key] == value {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("permission with %s %s: %w", key, value, infra.ErrNotFound)
}

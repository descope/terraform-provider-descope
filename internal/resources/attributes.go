package resources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/attribute"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewUserAttributeResource() resource.Resource {
	return newResource[attribute.UserAttributeModel]("user_attribute", attribute.UserAttributeSchema, attributeOps("/v1/mgmt/user/attribute"))
}

func NewTenantAttributeResource() resource.Resource {
	return newResource[attribute.TenantAttributeModel]("tenant_attribute", attribute.TenantAttributeSchema, attributeOps("/v1/mgmt/tenant/attribute"))
}

func NewAccessKeyAttributeResource() resource.Resource {
	return newResource[attribute.AccessKeyAttributeModel]("access_key_attribute", attribute.AccessKeyAttributeSchema, attributeOps("/v1/mgmt/accesskey/attribute"))
}

// attributeOps is for custom attribute resources owning one entry in a backend collection, addressed by a user-chosen name that
// the model sends in the body and that reads and deletes pass as the `name` query param. Writes return no body, so they re-read.
func attributeOps(path string) operations {
	read := func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
		return c.Get(ctx, projectID, path, map[string]string{"name": id})
	}
	write := func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
		if err := c.Post(ctx, projectID, path, data); err != nil {
			return nil, err
		}
		return read(ctx, c, projectID, id)
	}
	return operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			id, _ := data["name"].(string)
			entity, err := write(ctx, c, projectID, id, data)
			return id, entity, err
		},
		Read:   read,
		Update: write,
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Del(ctx, projectID, path, map[string]string{"name": id})
		},
	}
}

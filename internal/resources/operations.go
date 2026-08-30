package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Create returns the new entity's id along with its data; Read and Update only return the entity data
// because the caller already knows the id.
type createFunc func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error)
type readFunc func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error)
type updateFunc func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error)
type deleteFunc func(ctx context.Context, c *infra.Client, projectID, id string) error
type scopedReadFunc func(ctx context.Context, c *infra.Client, projectID, scope, id string) (map[string]any, error)
type scopedDeleteFunc func(ctx context.Context, c *infra.Client, projectID, scope, id string) error

// Defines how a resource's CRUD requests reach the backend. The shared resource families below have their
// own constructors; every other resource writes its operations longhand in a dedicated file, as plain
// closures over the infra.Client, so its endpoint quirks stay visible in one place.
//
// Singleton resources leave Create nil - baseResource creates them by applying Update with the project id
// as the entity id. ScopedRead and ScopedDelete are set instead of Read/Delete for entities addressed by
// both a scope (the model's app_id or method attribute) and their own id; Create and Update need no scoped
// variants because the model emits the scope in the request body.
type operations struct {
	Create       createFunc
	Read         readFunc
	Update       updateFunc
	Delete       deleteFunc
	ScopedRead   scopedReadFunc
	ScopedDelete scopedDeleteFunc
}

// noDelete is for resources with no delete endpoint: destroy only removes them from state.
// unwrapEntity treats an empty response envelope as a missing entity, so a read removes the resource from state instead of reporting a clean refresh.
func unwrapEntity(body map[string]any, key string) (map[string]any, error) {
	entity, ok := body[key].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: the response has no %q object", infra.ErrNotFound, key)
	}
	return entity, nil
}

func noDelete(_ context.Context, _ *infra.Client, _, _ string) error {
	return nil
}

// Creates a resource driven through the legacy generic `/v1/mgmt/infra` entity endpoint (frozen): only for
// resources already shipped on it. A schema with a `project_id` attribute makes it a project-level resource,
// otherwise it's assumed to be a company-level resource.
func newInfraResource[T any, M helpers.ResourceModel[T]](name string, sc schema.Schema) resource.Resource {
	return newResource[T, M](name, sc, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			res, err := c.Create(ctx, projectID, name, data)
			if err != nil {
				return "", nil, err
			}
			return res.ID, res.Data, nil
		},
		Read: func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
			res, err := c.Read(ctx, projectID, name, id)
			if err != nil {
				return nil, err
			}
			return res.Data, nil
		},
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			res, err := c.Update(ctx, projectID, name, id, data)
			if err != nil {
				return nil, err
			}
			return res.Data, nil
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Delete(ctx, projectID, name, id)
		},
	})
}

// Creates a singleton resource for settings on one GET/POST endpoint (e.g. /v1/mgmt/oauth/settings); there
// is no reset endpoint, so delete is a no-op.
func newSettingsResource[T any, M helpers.ResourceModel[T]](name string, sc schema.Schema, path string) resource.Resource {
	read := func(ctx context.Context, c *infra.Client, projectID, _ string) (map[string]any, error) {
		return c.Get(ctx, projectID, path, nil)
	}
	return newSingletonResource[T, M](name, sc, operations{
		Read: read,
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			if err := c.Post(ctx, projectID, path, data); err != nil {
				return nil, err
			}
			return read(ctx, c, projectID, id)
		},
		Delete: noDelete,
	})
}

// Returns a constructor for a standalone connector resource, used by the generated connectors.go registration
// file. Connectors live on dedicated CRUD endpoints and are identified by a server-assigned id.
func newConnectorResource[T any, M helpers.ResourceModel[T]](name string, sc schema.Schema) func() resource.Resource {
	const path = "/v1/mgmt/connector"
	return func() resource.Resource {
		return newResource[T, M](name, sc, operations{
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
				return c.PostData(ctx, projectID, path+"/update", body)
			},
			Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
				return c.Del(ctx, projectID, path, map[string]string{"id": id})
			},
		})
	}
}

package resources

import (
	"context"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/apppermission"
	"github.com/descope/terraform-provider-descope/internal/models/approle"
	"github.com/descope/terraform-provider-descope/internal/models/apps"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewOIDCAppResource() resource.Resource {
	read := ssoAppRead(true)
	return newResource[apps.OIDCAppModel]("oidc_app", apps.OIDCAppSchema, operations{
		Create: ssoAppCreate("oidc", read),
		Read:   read,
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			body := maps.Clone(data)
			body["id"] = id
			delete(body, "clientSecret")
			if err := c.Post(ctx, projectID, "/v1/mgmt/sso/idp/app/oidc/update", body); err != nil {
				return nil, err
			}
			return read(ctx, c, projectID, id)
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Post(ctx, projectID, "/v1/mgmt/sso/idp/app/delete", map[string]any{"id": id})
		},
	})
}

func NewSAMLAppResource() resource.Resource {
	read := ssoAppRead(false)
	return newResource[apps.SAMLAppModel]("saml_app", apps.SAMLAppSchema, operations{
		Create: ssoAppCreate("saml", read),
		Read:   read,
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			body := maps.Clone(data)
			body["id"] = id
			if err := c.Post(ctx, projectID, "/v1/mgmt/sso/idp/app/saml/update", body); err != nil {
				return nil, err
			}
			return read(ctx, c, projectID, id)
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Post(ctx, projectID, "/v1/mgmt/sso/idp/app/delete", map[string]any{"id": id})
		},
	})
}

func NewWSFedAppResource() resource.Resource {
	read := ssoAppRead(false)
	return newResource[apps.WSFedAppModel]("wsfed_app", apps.WSFedAppSchema, operations{
		Create: ssoAppCreate("wsfed", read),
		Read:   read,
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			body := maps.Clone(data)
			body["id"] = id
			if err := c.Post(ctx, projectID, "/v1/mgmt/sso/idp/app/wsfed/update", body); err != nil {
				return nil, err
			}
			return read(ctx, c, projectID, id)
		},
		Delete: func(ctx context.Context, c *infra.Client, projectID, id string) error {
			return c.Post(ctx, projectID, "/v1/mgmt/sso/idp/app/delete", map[string]any{"id": id})
		},
	})
}

// ssoAppCreate creates the app and re-reads it because the create response only returns the new id.
func ssoAppCreate(kind string, read readFunc) createFunc {
	return func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
		body, err := c.PostData(ctx, projectID, "/v1/mgmt/sso/idp/app/"+kind+"/create", data)
		if err != nil {
			return "", nil, err
		}
		id, _ := body["id"].(string)
		entity, err := read(ctx, c, projectID, id)
		return id, entity, err
	}
}

// ssoAppRead also fetches the OIDC client secret cleartext on every read (including import) because the load endpoint always returns
// it empty; the model only applies it when the state has no secret already.
func ssoAppRead(withSecret bool) readFunc {
	return func(ctx context.Context, c *infra.Client, projectID, id string) (map[string]any, error) {
		data, err := c.Get(ctx, projectID, "/v1/mgmt/sso/idp/app/load", map[string]string{"id": id})
		if err != nil {
			return nil, err
		}
		if withSecret {
			secret, err := c.Get(ctx, projectID, "/v1/mgmt/sso/idp/app/secret", map[string]string{"id": id})
			if err != nil {
				return nil, err
			}
			if settings, ok := data["oidcSettings"].(map[string]any); ok {
				if cleartext, ok := secret["cleartext"].(string); ok && cleartext != "" {
					settings["clientSecret"] = cleartext
				}
			}
		}
		return data, nil
	}
}

func NewAppPermissionResource() resource.Resource {
	const path = "/v1/mgmt/sso/idp/app/permission"
	return newResource[apppermission.AppPermissionModel]("app_permission", apppermission.Schema, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			entity, err := c.PostData(ctx, projectID, path, data)
			if err != nil {
				return "", nil, err
			}
			id, _ := entity["id"].(string)
			return id, entity, nil
		},
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			body := maps.Clone(data)
			body["id"] = id
			return c.PostData(ctx, projectID, path+"/update", body)
		},
		ScopedRead: func(ctx context.Context, c *infra.Client, projectID, appID, id string) (map[string]any, error) {
			return c.Get(ctx, projectID, path, map[string]string{"appId": appID, "id": id})
		},
		ScopedDelete: func(ctx context.Context, c *infra.Client, projectID, appID, id string) error {
			return c.Del(ctx, projectID, path, map[string]string{"appId": appID, "id": id})
		},
	})
}

func NewAppRoleResource() resource.Resource {
	const path = "/v1/mgmt/sso/idp/app/role"
	return newResource[approle.AppRoleModel]("app_role", approle.Schema, operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			entity, err := c.PostData(ctx, projectID, path, data)
			if err != nil {
				return "", nil, err
			}
			id, _ := entity["id"].(string)
			return id, entity, nil
		},
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			body := maps.Clone(data)
			body["id"] = id
			return c.PostData(ctx, projectID, path+"/update", body)
		},
		ScopedRead: func(ctx context.Context, c *infra.Client, projectID, appID, id string) (map[string]any, error) {
			return c.Get(ctx, projectID, path, map[string]string{"appId": appID, "id": id})
		},
		ScopedDelete: func(ctx context.Context, c *infra.Client, projectID, appID, id string) error {
			return c.Del(ctx, projectID, path, map[string]string{"appId": appID, "id": id})
		},
	})
}

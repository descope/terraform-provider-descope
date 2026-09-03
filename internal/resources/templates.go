package resources

import (
	"context"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/emailtemplate"
	"github.com/descope/terraform-provider-descope/internal/models/texttemplate"
	"github.com/descope/terraform-provider-descope/internal/models/voicetemplate"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func NewEmailTemplateResource() resource.Resource {
	return newResource[emailtemplate.EmailTemplateModel]("email_template", emailtemplate.Schema, templateOps("email"))
}

func NewTextTemplateResource() resource.Resource {
	return newResource[texttemplate.TextTemplateModel]("text_template", texttemplate.Schema, templateOps("text"))
}

func NewVoiceTemplateResource() resource.Resource {
	return newResource[voicetemplate.VoiceTemplateModel]("voice_template", voicetemplate.Schema, templateOps("voice"))
}

// templateOps is for messaging templates on the shared template endpoints, addressed by fixed type, auth method (from the model's
// method attribute) and a server-assigned id. Create and Update carry the method in the body, so only reads and deletes are scoped.
func templateOps(typ string) operations {
	const path = "/v1/mgmt/template"
	return operations{
		Create: func(ctx context.Context, c *infra.Client, projectID string, data map[string]any) (string, map[string]any, error) {
			body := maps.Clone(data)
			body["type"] = typ
			entity, err := c.PostData(ctx, projectID, path, body)
			if err != nil {
				return "", nil, err
			}
			id, _ := entity["id"].(string)
			return id, entity, nil
		},
		Update: func(ctx context.Context, c *infra.Client, projectID, id string, data map[string]any) (map[string]any, error) {
			body := maps.Clone(data)
			body["id"] = id
			body["type"] = typ
			return c.PostData(ctx, projectID, path+"/update", body)
		},
		ScopedRead: func(ctx context.Context, c *infra.Client, projectID, method, id string) (map[string]any, error) {
			return c.Get(ctx, projectID, path, map[string]string{"type": typ, "method": method, "id": id})
		},
		ScopedDelete: func(ctx context.Context, c *infra.Client, projectID, method, id string) error {
			return c.Del(ctx, projectID, path, map[string]string{"type": typ, "method": method, "id": id})
		},
	}
}

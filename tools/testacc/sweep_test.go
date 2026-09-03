package testacc

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

const sweepMinAge = 2 * time.Hour

func init() {
	resource.AddTestSweepers("descope_entities", &resource.Sweeper{
		Name: "descope_entities",
		F:    sweepEntities,
	})
}

type sweepTarget struct {
	name   string
	list   func(ctx context.Context, c *infra.Client, projectID string) ([]map[string]any, error)
	delete func(ctx context.Context, c *infra.Client, projectID string, entity map[string]any) error
}

// Apps go first so their scoped roles and permissions cascade away before the project-level roles they might reference are deleted.
var sweepTargets = []sweepTarget{
	{name: "sso_app", list: listWithGet("/v1/mgmt/sso/idp/apps/load", "apps"), delete: deleteByField("/v1/mgmt/sso/idp/app/delete", "id")},
	{name: "role", list: listWithGet("/v1/mgmt/role/all", "roles"), delete: deleteByField("/v1/mgmt/role/delete", "name")},
	{name: "permission", list: listWithGet("/v1/mgmt/permission/all", "permissions"), delete: deleteByField("/v1/mgmt/permission/delete", "name")},
	{name: "jwt_template", list: listWithPost("/v1/mgmt/jwt/templates/list", "templates"), delete: deleteByField("/v1/mgmt/jwt/templates/delete", "id")},
	{name: "access_key", list: listWithPost("/v1/mgmt/accesskey/search", "keys"), delete: deleteInfraEntity("access_key")},
}

// Removes leftover testacc- entities from the shared test project; leaked dynamic projects are removed by the testcleanup make target.
func sweepEntities(_ string) error {
	managementKey := os.Getenv("DESCOPE_MANAGEMENT_KEY")
	baseURL := os.Getenv("DESCOPE_BASE_URL")
	projectID := os.Getenv("DESCOPE_TESTACC_PROJECT_ID")
	if managementKey == "" || baseURL == "" || projectID == "" {
		return fmt.Errorf("the DESCOPE_MANAGEMENT_KEY, DESCOPE_BASE_URL and DESCOPE_TESTACC_PROJECT_ID environment variables must be set for sweeping")
	}

	ctx := context.Background()
	client := infra.NewClient("testacc", managementKey, baseURL)
	for _, target := range sweepTargets {
		entities, err := target.list(ctx, client, projectID)
		if err != nil {
			return fmt.Errorf("listing %s entities: %w", target.name, err)
		}
		for _, entity := range entities {
			name, _ := entity["name"].(string)
			if !strings.HasPrefix(name, aliasPrefix) {
				continue
			}
			if created, ok := aliasCreatedAt(name, time.Now()); ok && time.Since(created) < sweepMinAge {
				continue
			}
			if err := target.delete(ctx, client, projectID, entity); err != nil {
				return fmt.Errorf("deleting %s entity %q: %w", target.name, name, err)
			}
			fmt.Printf("Swept %s entity %q\n", target.name, name)
		}
	}
	return nil
}

func aliasCreatedAt(name string, now time.Time) (time.Time, bool) {
	match := aliasStampPattern.FindStringSubmatch(name)
	if match == nil {
		return time.Time{}, false
	}
	stamp, err := time.ParseInLocation("01021504", match[1], time.Local)
	if err != nil {
		return time.Time{}, false
	}
	created := stamp.AddDate(now.Year(), 0, 0)
	if created.After(now.Add(time.Hour)) {
		created = created.AddDate(-1, 0, 0)
	}
	return created, true
}

func listWithGet(path, key string) func(ctx context.Context, c *infra.Client, projectID string) ([]map[string]any, error) {
	return func(ctx context.Context, c *infra.Client, projectID string) ([]map[string]any, error) {
		body, err := c.Get(ctx, projectID, path, nil)
		if err != nil {
			return nil, err
		}
		return extractEntities(body, key), nil
	}
}

func listWithPost(path, key string) func(ctx context.Context, c *infra.Client, projectID string) ([]map[string]any, error) {
	return func(ctx context.Context, c *infra.Client, projectID string) ([]map[string]any, error) {
		body, err := c.PostData(ctx, projectID, path, map[string]any{})
		if err != nil {
			return nil, err
		}
		return extractEntities(body, key), nil
	}
}

func extractEntities(body map[string]any, key string) []map[string]any {
	var entities []map[string]any
	list, _ := body[key].([]any)
	for _, item := range list {
		if entity, ok := item.(map[string]any); ok {
			entities = append(entities, entity)
		}
	}
	return entities
}

func deleteByField(path, field string) func(ctx context.Context, c *infra.Client, projectID string, entity map[string]any) error {
	return func(ctx context.Context, c *infra.Client, projectID string, entity map[string]any) error {
		return c.Post(ctx, projectID, path, map[string]any{field: entity[field]})
	}
}

func deleteInfraEntity(entity string) func(ctx context.Context, c *infra.Client, projectID string, e map[string]any) error {
	return func(ctx context.Context, c *infra.Client, projectID string, e map[string]any) error {
		id, _ := e["id"].(string)
		return c.Delete(ctx, projectID, entity, id)
	}
}

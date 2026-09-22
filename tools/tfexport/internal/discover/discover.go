package discover

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/warn"
)

func Snapshot(ctx context.Context, client *infra.Client, projectID string) (map[string]any, error) {
	data, err := client.PostData(ctx, projectID, "/v1/mgmt/project/snapshot/export", map[string]any{"format": "plain"})
	if err != nil {
		return nil, err
	}
	files, ok := data["files"].(map[string]any)
	if !ok {
		return nil, errors.New("unexpected snapshot response: no files object")
	}
	return files, nil
}

func InboundApps(ctx context.Context, client *infra.Client, projectID string) ([]Instance, error) {
	data, err := client.Get(ctx, projectID, "/v1/mgmt/thirdparty/apps/load", nil)
	if err != nil {
		return nil, err
	}
	var instances []Instance
	for _, app := range anyObjectList(data["apps"]) {
		id, _ := app["id"].(string)
		name, _ := app["name"].(string)
		if id != "" {
			instances = append(instances, Instance{Resource: "descope_inbound_app", ID: id, Name: name})
		}
	}
	return instances, nil
}

func ResolveAuthorizationIDs(ctx context.Context, client *infra.Client, projectID string, instances []Instance) (resolved []Instance, warnings []warn.Warning) {
	live := map[string]map[string]string{}
	for resource, load := range map[string]func() (map[string]string, error){
		"descope_role":       func() (map[string]string, error) { return loadRoles(ctx, client, projectID) },
		"descope_permission": func() (map[string]string, error) { return loadPermissions(ctx, client, projectID) },
	} {
		if !slices.ContainsFunc(instances, func(i Instance) bool { return i.Resource == resource }) {
			continue
		}
		ids, err := load()
		if err != nil {
			warnings = append(warnings, warn.Lost("Failed to resolve %s ids: %s", resource, err.Error()))
			continue
		}
		live[resource] = ids
	}

	for _, instance := range instances {
		if ids, ok := live[instance.Resource]; ok {
			id, ok := ids[instance.Name]
			if !ok {
				warnings = append(warnings, warn.Lost("Skipped %s %q: no entity with that name exists", instance.Resource, instance.Name))
				continue
			}
			instance.ID = id
		}
		resolved = append(resolved, instance)
	}
	return resolved, warnings
}

func loadRoles(ctx context.Context, client *infra.Client, projectID string) (map[string]string, error) {
	data, err := client.PostData(ctx, projectID, "/v1/mgmt/role/search", map[string]any{})
	if err != nil {
		return nil, err
	}
	return idsByName(data["roles"]), nil
}

func loadPermissions(ctx context.Context, client *infra.Client, projectID string) (map[string]string, error) {
	data, err := client.Get(ctx, projectID, "/v1/mgmt/permission/all", nil)
	if err != nil {
		return nil, err
	}
	return idsByName(data["permissions"]), nil
}

func idsByName(value any) map[string]string {
	ids := map[string]string{}
	for _, entry := range anyObjectList(value) {
		id, _ := entry["id"].(string)
		name, _ := entry["name"].(string)
		if id != "" && name != "" {
			ids[name] = id
		}
	}
	return ids
}

func DumpIndex(w io.Writer, files map[string]any) {
	for _, path := range sortedKeys(files) {
		_, _ = fmt.Fprintf(w, "%s\n", path)
		dumpValue(w, files[path], 1)
	}
}

func dumpValue(w io.Writer, value any, depth int) {
	indent := strings.Repeat("    ", depth)
	switch v := value.(type) {
	case map[string]any:
		if depth > 3 {
			_, _ = fmt.Fprintf(w, "%s{%d keys}\n", indent, len(v))
			return
		}
		for _, key := range sortedKeys(v) {
			_, _ = fmt.Fprintf(w, "%s%s: %s\n", indent, key, describe(v[key]))
			switch nested := v[key].(type) {
			case map[string]any:
				dumpValue(w, nested, depth+1)
			case []any:
				if len(nested) > 0 {
					dumpValue(w, nested[0], depth+1)
				}
			}
		}
	case []any:
		if len(v) > 0 {
			_, _ = fmt.Fprintf(w, "%sfirst of %d: %s\n", indent, len(v), describe(v[0]))
			if nested, ok := v[0].(map[string]any); ok {
				dumpValue(w, nested, depth+1)
			}
		}
	default:
		_, _ = fmt.Fprintf(w, "%s%s\n", indent, describe(v))
	}
}

func describe(value any) string {
	switch v := value.(type) {
	case string:
		if len(v) > 48 {
			return fmt.Sprintf("string(%d)", len(v))
		}
		return fmt.Sprintf("%q", v)
	case map[string]any:
		return fmt.Sprintf("object{%s}", strings.Join(sortedKeys(v), ", "))
	case []any:
		return fmt.Sprintf("array(%d)", len(v))
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

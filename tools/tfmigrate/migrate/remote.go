package migrate

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/descope/terraform-provider-descope/internal/infra"
)

// Environment variables the provider itself reads, reused so a migration needs no separate credentials.
const (
	managementKeyEnv = "DESCOPE_MANAGEMENT_KEY"
	baseURLEnv       = "DESCOPE_BASE_URL"
	unsafeLogsEnv    = "TF_UNSAFE_LOGS"
	// projectEntity is the entity name the generic management endpoint addresses a whole project with.
	projectEntity = "project"
)

// ProjectReader reads a project's current configuration. The only implementation issues a single HTTP GET.
type ProjectReader interface {
	ReadProject(ctx context.Context) (map[string]any, error)
}

// getOnlyReader reads the project through infra.Client.Read, which is the provider's GET of the generic management
// entity endpoint. It is the only remote call the migration is allowed to make: no create, update or delete path is
// reachable from here, and the response is used solely to recover ids the legacy state is missing.
type getOnlyReader struct {
	client    *infra.Client
	projectID string
}

// NewProjectReader builds the GET-only reader from the same credentials the provider uses.
func NewProjectReader(projectID, baseURL string) (ProjectReader, error) {
	if os.Getenv(unsafeLogsEnv) != "" {
		return nil, fmt.Errorf("unset %s before using -remote: the migration never permits API response payloads to be written to logs", unsafeLogsEnv)
	}
	if baseURL == "" {
		baseURL = os.Getenv(baseURLEnv)
	}
	if err := validateRemoteBaseURL(baseURL); err != nil {
		return nil, err
	}
	managementKey := os.Getenv(managementKeyEnv)
	if managementKey == "" {
		return nil, fmt.Errorf("%s is not set, and a GET-only backfill needs the same management key the provider uses", managementKeyEnv)
	}
	return &getOnlyReader{client: infra.NewClient(ToolName, managementKey, baseURL), projectID: projectID}, nil
}

func validateRemoteBaseURL(baseURL string) error {
	if baseURL == "" {
		return nil // the SDK default is https://api.descope.com with certificate verification
	}
	parsed, err := url.ParseRequestURI(baseURL)
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return fmt.Errorf("remote base URL must be an HTTPS URL without credentials, query parameters or fragments")
	}
	host := parsed.Hostname()
	if net.ParseIP(host) != nil || !strings.Contains(host, ".") || parsed.Port() != "" {
		return fmt.Errorf("remote base URL must use a DNS hostname and the default HTTPS port so the provider SDK verifies its certificate")
	}
	return nil
}

func (r *getOnlyReader) ReadProject(ctx context.Context) (map[string]any, error) {
	response, err := r.client.Read(ctx, r.projectID, projectEntity, r.projectID)
	if err != nil {
		return nil, fmt.Errorf("reading project %s: %w", r.projectID, err)
	}
	if response == nil || len(response.Data) == 0 {
		return nil, fmt.Errorf("project %s returned no configuration", r.projectID)
	}
	return response.Data, nil
}

// Backfill fills in the entity ids the legacy state is missing, using the project configuration the backend
// returns. Only ids are taken from the response: every configured value still comes from the Terraform state, so a
// value edited in the console is never silently adopted as if Terraform had written it.
func Backfill(ctx context.Context, legacy *LegacyResource, reader ProjectReader) ([]string, error) {
	project, err := reader.ReadProject(ctx)
	if err != nil {
		return nil, err
	}
	filled := []string{}
	targets := append([]backfillTarget(nil), backfillTargets...)
	registry, registryErr := NewRegistry(ctx)
	if registryErr != nil {
		return nil, fmt.Errorf("discovering connector resources for GET backfill: %w", registryErr)
	}
	for legacyName, resourceType := range registry.ConnectorFamilies() {
		destination, ok := registry.Lookup(resourceType)
		if !ok || destination.WireType == "" {
			continue
		}
		targets = append(targets, backfillTarget{
			statePath: "connectors." + legacyName,
			stateName: "name",
			wirePath:  "connectors." + destination.WireType,
			wireName:  "name",
		})
	}
	sort.SliceStable(targets, func(i, j int) bool { return targets[i].statePath < targets[j].statePath })
	for _, target := range targets {
		items := asObjects(valueAtPath(legacy.Attributes, target.statePath))
		if len(items) == 0 {
			continue
		}
		remote := asObjects(valueAtPath(project, target.wirePath))
		if len(remote) == 0 {
			continue
		}
		byName := map[string]string{}
		for _, entity := range remote {
			name, _ := entity[target.wireName].(string)
			id, _ := entity["id"].(string)
			if name == "" || id == "" {
				continue
			}
			if previous, duplicate := byName[name]; duplicate && previous != id {
				return nil, fmt.Errorf("GET backfill for %s found more than one backend id for name %q", target.statePath, name)
			}
			byName[name] = id
		}
		for _, item := range items {
			if id, _ := item["id"].(string); id != "" {
				continue
			}
			name, _ := item[target.stateName].(string)
			id, found := byName[name]
			if !found {
				continue
			}
			item["id"] = id
			filled = append(filled, fmt.Sprintf("%s %q", target.statePath, name))
		}
	}
	appFilled, appErr := backfillAppChildren(legacy.Attributes, project)
	if appErr != nil {
		return nil, appErr
	}
	filled = append(filled, appFilled...)
	sort.Strings(filled)
	return filled, nil
}

func backfillAppChildren(state, project map[string]any) ([]string, error) {
	filled := []string{}
	families := []struct {
		statePath string
		wirePath  string
	}{
		{statePath: "applications.oidc_applications", wirePath: "applications.oidc"},
		{statePath: "applications.saml_applications", wirePath: "applications.saml"},
		{statePath: "applications.wsfed_applications", wirePath: "applications.wsfed"},
	}
	for _, family := range families {
		stateApps := asObjects(valueAtPath(state, family.statePath))
		remoteApps := asObjects(valueAtPath(project, family.wirePath))
		remoteByName := map[string]map[string]any{}
		for _, app := range remoteApps {
			name, _ := app["name"].(string)
			if name == "" {
				continue
			}
			if _, duplicate := remoteByName[name]; duplicate {
				return nil, fmt.Errorf("GET backfill for %s found more than one application named %q", family.statePath, name)
			}
			remoteByName[name] = app
		}
		for _, app := range stateApps {
			name, _ := app["name"].(string)
			remoteApp := remoteByName[name]
			for _, childType := range []string{"permissions", "roles"} {
				children := asObjects(app[childType])
				remoteChildren := asObjects(remoteApp[childType])
				ids, err := idsByName(remoteChildren, family.statePath+"."+childType)
				if err != nil {
					return nil, err
				}
				for _, child := range children {
					if id, _ := child["id"].(string); id != "" {
						continue
					}
					childName, _ := child["name"].(string)
					id, found := ids[childName]
					if !found {
						continue
					}
					child["id"] = id
					filled = append(filled, fmt.Sprintf("%s %q %s %q", family.statePath, name, childType, childName))
				}
			}
		}
	}
	return filled, nil
}

func idsByName(items []map[string]any, source string) (map[string]string, error) {
	ids := map[string]string{}
	for _, item := range items {
		name, _ := item["name"].(string)
		id, _ := item["id"].(string)
		if name == "" || id == "" {
			continue
		}
		if previous, duplicate := ids[name]; duplicate && previous != id {
			return nil, fmt.Errorf("GET backfill for %s found more than one backend id for name %q", source, name)
		}
		ids[name] = id
	}
	return ids, nil
}

// backfillTarget maps a legacy state collection onto the same collection in the project configuration the backend
// returns, so a missing id can be matched by the name both sides agree on.
type backfillTarget struct {
	statePath string
	stateName string
	wirePath  string
	wireName  string
}

// backfillTargets covers the collections whose legacy id attribute was optional, and therefore the collections a
// state file can legitimately be missing an id for.
var backfillTargets = []backfillTarget{
	{statePath: "authorization.roles", stateName: "name", wirePath: "authorization.roles", wireName: "name"},
	{statePath: "authorization.permissions", stateName: "name", wirePath: "authorization.permissions", wireName: "name"},
	{statePath: "attributes.user", stateName: "name", wirePath: "attributes.user", wireName: "displayName"},
	{statePath: "attributes.tenant", stateName: "name", wirePath: "attributes.tenant", wireName: "displayName"},
	{statePath: "attributes.access_key", stateName: "name", wirePath: "attributes.accessKey", wireName: "displayName"},
	{statePath: "lists", stateName: "name", wirePath: "lists", wireName: "name"},
	{statePath: "applications.oidc_applications", stateName: "name", wirePath: "applications.oidc", wireName: "name"},
	{statePath: "applications.saml_applications", stateName: "name", wirePath: "applications.saml", wireName: "name"},
	{statePath: "applications.wsfed_applications", stateName: "name", wirePath: "applications.wsfed", wireName: "name"},
	{statePath: "jwt_templates.user_templates", stateName: "name", wirePath: "jwtTemplates.userTemplates", wireName: "name"},
	{statePath: "jwt_templates.access_key_templates", stateName: "name", wirePath: "jwtTemplates.keyTemplates", wireName: "name"},
	{statePath: "authentication.otp.email_service.templates", stateName: "name", wirePath: "authentication.otp.emailTemplates", wireName: "name"},
	{statePath: "authentication.otp.text_service.templates", stateName: "name", wirePath: "authentication.otp.textTemplates", wireName: "name"},
	{statePath: "authentication.otp.voice_service.templates", stateName: "name", wirePath: "authentication.otp.voiceTemplates", wireName: "name"},
	{statePath: "authentication.magic_link.email_service.templates", stateName: "name", wirePath: "authentication.magiclink.emailTemplates", wireName: "name"},
	{statePath: "authentication.magic_link.text_service.templates", stateName: "name", wirePath: "authentication.magiclink.textTemplates", wireName: "name"},
	{statePath: "authentication.enchanted_link.email_service.templates", stateName: "name", wirePath: "authentication.enchantedlink.emailTemplates", wireName: "name"},
	{statePath: "authentication.password.email_service.templates", stateName: "name", wirePath: "authentication.password.emailTemplates", wireName: "name"},
	{statePath: "authentication.sso.email_service.templates", stateName: "name", wirePath: "authentication.sso.emailTemplates", wireName: "name"},
	{statePath: "invite_settings.email_service.templates", stateName: "name", wirePath: "settings.inviteEmailTemplates", wireName: "name"},
}

func valueAtPath(root map[string]any, path string) any {
	var current any = root
	for _, segment := range strings.Split(path, ".") {
		object, ok := asObject(current)
		if !ok {
			return nil
		}
		current = object[segment]
	}
	return current
}

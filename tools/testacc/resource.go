package testacc

import (
	"fmt"
	"maps"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

func ProjectID(t *testing.T) string {
	projectID := os.Getenv("DESCOPE_TESTACC_PROJECT_ID")
	if projectID == "" && os.Getenv("TF_ACC") != "" {
		t.Fatal("Expected the DESCOPE_TESTACC_PROJECT_ID environment variable to be set to the ID of an existing project for running acceptance tests")
	}
	return projectID
}

func Project(t *testing.T) *Resource {
	return newResource(t, "project")
}

func Descoper(t *testing.T) *Resource {
	return newResource(t, "descoper")
}

func ManagementKey(t *testing.T) *Resource {
	return newResource(t, "management_key")
}

func AccessKey(t *testing.T) *Resource {
	return newResource(t, "access_key")
}

func InboundApp(t *testing.T) *Resource {
	return newResource(t, "inbound_app")
}

func Engine(t *testing.T) *Resource {
	return newResource(t, "engine")
}

func OAuthSettings(t *testing.T) *Resource {
	return newResource(t, "oauth_settings")
}

func OAuthProvider(t *testing.T) *Resource {
	return newResource(t, "oauth_provider")
}

func UserAttribute(t *testing.T) *Resource {
	return newResource(t, "user_attribute")
}

func TenantAttribute(t *testing.T) *Resource {
	return newResource(t, "tenant_attribute")
}

func AccessKeyAttribute(t *testing.T) *Resource {
	return newResource(t, "access_key_attribute")
}

func PasskeySettings(t *testing.T) *Resource {
	return newResource(t, "passkey_settings")
}

func TOTPSettings(t *testing.T) *Resource {
	return newResource(t, "totp_settings")
}

func AdminPortal(t *testing.T) *Resource {
	return newResource(t, "admin_portal")
}

func MagicLinkSettings(t *testing.T) *Resource {
	return newResource(t, "magiclink_settings")
}

func InviteSettings(t *testing.T) *Resource {
	return newResource(t, "invite_settings")
}

func ProjectSettings(t *testing.T) *Resource {
	return newResource(t, "project_settings")
}

func SessionSettings(t *testing.T) *Resource {
	return newResource(t, "session_settings")
}

func SessionMigration(t *testing.T) *Resource {
	return newResource(t, "session_migration")
}

func OTPSettings(t *testing.T) *Resource {
	return newResource(t, "otp_settings")
}

func EnchantedLinkSettings(t *testing.T) *Resource {
	return newResource(t, "enchantedlink_settings")
}

func EmbeddedLinkSettings(t *testing.T) *Resource {
	return newResource(t, "embeddedlink_settings")
}

func PasswordSettings(t *testing.T) *Resource {
	return newResource(t, "password_settings")
}

func SSOSettings(t *testing.T) *Resource {
	return newResource(t, "sso_settings")
}

func Flow(t *testing.T) *Resource {
	return newResource(t, "flow")
}

func Role(t *testing.T) *Resource {
	return newResource(t, "role")
}

func Permission(t *testing.T) *Resource {
	return newResource(t, "permission")
}

func OIDCApp(t *testing.T) *Resource {
	return newResource(t, "oidc_app")
}

func SAMLApp(t *testing.T) *Resource {
	return newResource(t, "saml_app")
}

func WSFedApp(t *testing.T) *Resource {
	return newResource(t, "wsfed_app")
}

func AppRole(t *testing.T) *Resource {
	return newResource(t, "app_role")
}

func AppPermission(t *testing.T) *Resource {
	return newResource(t, "app_permission")
}

func List(t *testing.T) *Resource {
	return newResource(t, "list")
}

func JWTTemplate(t *testing.T) *Resource {
	return newResource(t, "jwt_template")
}

func EmailTemplate(t *testing.T) *Resource {
	return newResource(t, "email_template")
}

func TextTemplate(t *testing.T) *Resource {
	return newResource(t, "text_template")
}

func VoiceTemplate(t *testing.T) *Resource {
	return newResource(t, "voice_template")
}

func Styles(t *testing.T) *Resource {
	return newResource(t, "styles")
}

func FGASchema(t *testing.T) *Resource {
	return newResource(t, "fga_schema")
}

func Widget(t *testing.T) *Resource {
	return newResource(t, "widget")
}

func NewResource(t *testing.T, typ string) *Resource {
	return newResource(t, typ)
}

func newResource(t *testing.T, typ string) *Resource {
	return &Resource{
		Type: typ,
		ID:   "test",
		Name: GenerateAlias(t),
	}
}

type Resource struct {
	Type string // the resource type without the 'descope_' prefix
	ID   string // the resource name in the Terraform config
	Name string // the value of the 'name' attribute
}

func (r *Resource) Path() string {
	return fmt.Sprintf(`descope_%s.%s`, r.Type, r.ID)
}

func (r *Resource) Variables(s ...string) string {
	return strings.Join(s, "\n")
}

func (r *Resource) Config(s ...string) string {
	n := fmt.Sprintf(`name = %q`, r.Name)
	s = append([]string{n}, s...)
	return fmt.Sprintf(resourceFormat, r.Type, r.ID, strings.Join(s, "\n	"))
}

// Block builds a resource block with only the provided fields, with no auto-injected name, for resources that have no name attribute.
func (r *Resource) Block(s ...string) string {
	return fmt.Sprintf(resourceFormat, r.Type, r.ID, strings.Join(s, "\n	"))
}

func (r *Resource) Check(checks map[string]any, extras ...resource.TestCheckFunc) resource.TestCheckFunc {
	path := r.Path()
	f := []resource.TestCheckFunc{}
	checks = flatten(checks, "")
	for k, v := range checks {
		if first := strings.TrimSuffix(k, ".=="); first != k {
			second, ok := v.(string)
			if !ok {
				panic(fmt.Sprintf("unexpected non-string argument of type %T in equality check: %v", v, v))
			}
			f = append(f, resource.TestCheckResourceAttrPair(path, first, path, second))
		} else if value, ok := v.(string); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k, value))
		} else if value, ok := v.(int); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k, strconv.Itoa(value)))
		} else if value, ok := v.(bool); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k, fmt.Sprintf("%t", value)))
		} else if value, ok := v.([]string); ok {
			f = append(f, resource.TestCheckResourceAttr(path, k+".#", strconv.Itoa(len(value))))
			for i := range value {
				f = append(f, resource.TestCheckTypeSetElemAttr(path, fmt.Sprintf("%s.*", k), value[i]))
			}
		} else if value, ok := v.(func(string) error); ok {
			f = append(f, resource.TestCheckResourceAttrWith(path, k, value))
		} else if v == AttributeIsSet {
			f = append(f, resource.TestCheckResourceAttrSet(path, k))
		} else if v == AttributeIsNotSet {
			f = append(f, resource.TestCheckNoResourceAttr(path, k))
		} else {
			panic(fmt.Sprintf("unexpected value of type %T in Check(): %v", v, v))
		}
	}
	f = append(f, extras...)
	return resource.ComposeAggregateTestCheckFunc(f...)
}

const resourceFormat = `
resource "descope_%s" "%s" {
	%s
}
`

const aliasPrefix = "testacc-"

var aliasStampPattern = regexp.MustCompile(`-(\d{8})-[0-9a-f]{8}$`)

func GenerateAlias(t *testing.T) string {
	prefix := strings.TrimSpace(os.Getenv("DESCOPE_TESTACC_PREFIX"))
	if prefix == "" {
		prefix = "testacc-local"
	}
	if !strings.HasPrefix(prefix, aliasPrefix) {
		prefix = aliasPrefix + prefix
	}
	test := strings.TrimPrefix(t.Name(), "Test")
	ts := time.Now().Format("01021504") // MMddHHmm
	rand, err := uuid.GenerateUUID()
	require.NoError(t, err)
	suffix := rand[len(rand)-8:]
	return fmt.Sprintf("%s-%s-%s-%s", prefix, test, ts, suffix)
}

func GenerateImportStateID(path string, attrs ...string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		resources, ok := state.RootModule().Resources[path]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", path) // nolint:forbidigo
		}
		var parts []string
		for _, attr := range attrs {
			v, ok := resources.Primary.Attributes[attr]
			if !ok || v == "" {
				return "", fmt.Errorf("attribute %q not found in %s", attr, path) // nolint:forbidigo
			}
			parts = append(parts, v)
		}
		return strings.Join(parts, "/"), nil
	}
}

func flatten(checks map[string]any, keypath string) map[string]any {
	result := map[string]any{}
	for k, v := range checks {
		if keypath != "" {
			k = keypath + "." + k
		}
		if m, ok := v.(map[string]any); ok {
			maps.Copy(result, flatten(m, k))
		} else {
			result[k] = v
		}
	}
	return result
}

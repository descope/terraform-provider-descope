package testacc

import (
	"os"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/require"
)

var protoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"descope": providerserver.NewProtocol6WithError(provider.NewDescopeProvider("test")()),
}

func preCheck(t *testing.T) {
	env := []string{"DESCOPE_PROJECT_ID", "DESCOPE_MANAGEMENT_KEY", "DESCOPE_BASE_URL"}
	for _, e := range env {
		require.NotEmpty(t, os.Getenv(e), "The following environment variables must be set for acceptance tests: "+strings.Join(env, ", "))
	}
}

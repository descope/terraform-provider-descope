package testacc

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/descope/terraform-provider-descope/internal/resources"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

func Run(t *testing.T, steps ...resource.TestStep) {
	if parallel, _ := strconv.ParseBool(os.Getenv("DESCOPE_TESTACC_PARALLEL")); parallel {
		t.Parallel()
	}
	resource.Test(t, TestCase(t, steps...))
}

func RunWithDestroyCheck(t *testing.T, resourceType string, steps ...resource.TestStep) {
	checker := lookupDestroyChecker(t, resourceType)
	if parallel, _ := strconv.ParseBool(os.Getenv("DESCOPE_TESTACC_PARALLEL")); parallel {
		t.Parallel()
	}
	testCase := TestCase(t, steps...)
	testCase.CheckDestroy = checkDestroyed(resourceType, checker)
	resource.Test(t, testCase)
}

func TestCase(t *testing.T, steps ...resource.TestStep) resource.TestCase {
	for i := range steps {
		steps[i] = applyStepThrottling(steps[i])
	}
	return resource.TestCase{
		PreCheck:                 func() { preCheck(t) },
		ProtoV6ProviderFactories: protoV6ProviderFactories,
		Steps:                    steps,
	}
}

// Helpers

func applyStepThrottling(step resource.TestStep) resource.TestStep {
	if seconds, _ := strconv.ParseInt(os.Getenv("DESCOPE_TESTACC_THROTTLE"), 10, 64); seconds != 0 {
		curr := step.PreConfig
		step.PreConfig = func() {
			time.Sleep(time.Duration(seconds) * time.Second)
			if curr != nil {
				curr()
			}
		}
	}
	return step
}

func lookupDestroyChecker(t *testing.T, resourceType string) resources.DestroyChecker {
	ctx := context.Background()
	for _, newResource := range provider.NewDescopeProvider("test")().Resources(ctx) {
		res := newResource()
		metadata := fwresource.MetadataResponse{}
		res.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "descope"}, &metadata)
		if metadata.TypeName != resourceType {
			continue
		}
		checker, ok := res.(resources.DestroyChecker)
		require.True(t, ok, "resource type %s does not support destroy checks", resourceType)
		return checker
	}
	require.Fail(t, "unknown resource type", "no registered resource named %s", resourceType)
	return nil
}

func checkDestroyed(resourceType string, checker resources.DestroyChecker) func(*terraform.State) error {
	return func(s *terraform.State) error {
		client := infra.NewClient("testacc", os.Getenv("DESCOPE_MANAGEMENT_KEY"), os.Getenv("DESCOPE_BASE_URL"))
		checked := 0
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			checked++
			projectID := rs.Primary.Attributes["project_id"]
			appID := rs.Primary.Attributes["app_id"]
			if err := checker.CheckDestroyed(context.Background(), client, projectID, appID, rs.Primary.ID); err != nil {
				return err
			}
		}
		// a mistyped resource type must not turn the destroy check into a silent pass
		if checked == 0 {
			return fmt.Errorf("no resources of type %s were found in the state for the destroy check", resourceType)
		}
		return nil
	}
}

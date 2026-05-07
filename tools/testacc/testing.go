package testacc

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Run(t *testing.T, steps ...resource.TestStep) {
	if parallel, _ := strconv.ParseBool(os.Getenv("DESCOPE_TESTACC_PARALLEL")); parallel {
		t.Parallel()
	}
	resource.Test(t, TestCase(t, steps...))
}

func RunIsolated(t *testing.T, steps ...resource.TestStep) {
	resource.Test(t, TestCase(t, steps...))
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

package infra

import (
	"context"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// A management key's authorization for a project comes from ReBAC tuples the backend writes asynchronously
// after the project is created, so calls in the same apply can be rejected until those tuples land.
const projectNotReadyErrorCode = "E111009"

var projectNotReadyDelays = []time.Duration{500 * time.Millisecond, time.Second, 2 * time.Second}

// Retries only the error above, which the backend returns before the request has any effect, so replaying
// it cannot duplicate a write. Any other failure, other forbidden responses included, is returned as is.
func retrying(ctx context.Context, call func() (*api.HTTPResponse, error)) (*api.HTTPResponse, error) {
	res, err := call()
	for _, delay := range projectNotReadyDelays {
		if descope.AsError(err, projectNotReadyErrorCode) == nil {
			break
		}
		tflog.Info(ctx, "Waiting for project authorization to propagate", map[string]any{"delay": delay.String()})
		select {
		case <-ctx.Done():
			return res, err
		case <-time.After(delay):
		}
		res, err = call()
	}
	return res, err
}

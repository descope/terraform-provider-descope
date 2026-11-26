package infra

import (
	"context"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Transient races with the backend's async post-creation work: E111009 = management key ReBAC tuples not yet
// written, E111604 = theme version conflict with the creation handler republishing it. Both clear on the next apply.
var retryableErrorCodes = []string{"E111009", "E111604"}

var retryDelays = []time.Duration{500 * time.Millisecond, time.Second, 2 * time.Second, 4 * time.Second}

// Safe to replay: the backend returns these codes before the request has any effect.
func retrying(ctx context.Context, call func() (*api.HTTPResponse, error)) (*api.HTTPResponse, error) {
	res, err := call()
	for _, delay := range retryDelays {
		de := descope.AsError(err, retryableErrorCodes...)
		if de == nil {
			break
		}
		tflog.Info(ctx, "Retrying after a transient backend error", map[string]any{"code": de.Code, "delay": delay.String()})
		select {
		case <-ctx.Done():
			return res, err
		case <-time.After(delay):
		}
		res, err = call()
	}
	return res, err
}

package infra

import (
	"context"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Transient races with the backend's async post-creation work: E111009 = management key ReBAC tuples not yet
// written, E111604 = theme version conflict with the creation handler republishing it. Both clear on the next apply.
var retryableErrorCodes = []string{"E111009", "E111604"}

var retryDelays = []time.Duration{500 * time.Millisecond, time.Second, 2 * time.Second, 4 * time.Second}

var rateLimitDelays = []time.Duration{5 * time.Second, 10 * time.Second, 20 * time.Second, 40 * time.Second, 60 * time.Second, 60 * time.Second}

var rateLimitMaxRetryAfter = 2 * time.Minute

var rateLimitJitter = func(delay time.Duration) time.Duration {
	return rand.N(delay/4 + 1)
}

// Safe to replay: the backend returns these codes before the request has any effect.
func retrying(ctx context.Context, call func() (*api.HTTPResponse, error)) (*api.HTTPResponse, error) {
	res, err := call()
	transientAttempts, rateLimitAttempts := 0, 0
	for {
		var delay time.Duration
		if de := descope.AsError(err, retryableErrorCodes...); de != nil && transientAttempts < len(retryDelays) {
			delay = retryDelays[transientAttempts]
			transientAttempts++
			tflog.Info(ctx, "Retrying after a transient backend error", map[string]any{"code": de.Code, "delay": delay.String()})
		} else if de := rateLimitError(err); de != nil {
			var ok bool
			if delay, ok = rateLimitDelay(de, rateLimitAttempts); !ok {
				break
			}
			rateLimitAttempts++
			tflog.Info(ctx, "Retrying after the backend rate limited the request", map[string]any{"code": de.Code, "delay": delay.String()})
		} else {
			break
		}
		select {
		case <-ctx.Done():
			return res, err
		case <-time.After(delay):
		}
		res, err = call()
	}
	return res, err
}

func rateLimitError(err error) *descope.Error {
	de := descope.AsError(err)
	if de == nil {
		return nil
	}
	if de.Code == descope.ErrRateLimitExceeded.Code || de.Info[descope.ErrorInfoKeys.HTTPResponseStatusCode] == http.StatusTooManyRequests {
		return de
	}
	return nil
}

func rateLimitDelay(de *descope.Error, attempt int) (time.Duration, bool) {
	if attempt >= len(rateLimitDelays) {
		return 0, false
	}
	delay := rateLimitDelays[attempt]
	if seconds, ok := de.Info[descope.ErrorInfoKeys.RateLimitExceededRetryAfter].(int); ok {
		retryAfter := time.Duration(seconds) * time.Second
		if retryAfter > rateLimitMaxRetryAfter {
			return 0, false
		}
		delay = max(delay, retryAfter)
	}
	return delay + rateLimitJitter(delay), true
}

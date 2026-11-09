package infra

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fastDelays(t *testing.T, count int) {
	previous := retryDelays
	retryDelays = make([]time.Duration, count)
	t.Cleanup(func() { retryDelays = previous })
}

func retryable(code string) error {
	return &descope.Error{Code: code, Description: "transient"}
}

func TestRetryingSucceedsOnceTheConditionClears(t *testing.T) {
	for _, code := range retryableErrorCodes {
		fastDelays(t, 4)
		calls := 0
		res, err := retrying(t.Context(), func() (*api.HTTPResponse, error) {
			calls++
			if calls < 3 {
				return nil, retryable(code)
			}
			return &api.HTTPResponse{BodyStr: "{}"}, nil
		})
		require.NoError(t, err, code)
		require.NotNil(t, res, code)
		assert.Equal(t, 3, calls, code)
	}
}

func TestRetryingGivesUpAndReturnsTheError(t *testing.T) {
	fastDelays(t, 4)
	calls := 0
	_, err := retrying(t.Context(), func() (*api.HTTPResponse, error) {
		calls++
		return nil, retryable("E111604")
	})
	require.Error(t, err)
	assert.Equal(t, "E111604", descope.AsError(err).Code)
	assert.Equal(t, 5, calls)
}

func TestRetryingLeavesOtherErrorsAlone(t *testing.T) {
	fastDelays(t, 4)
	for _, err := range []error{
		&descope.Error{Code: "E111008", Description: "some other forbidden"},
		&descope.Error{Code: "E111603", Description: "a neighbouring theme error"},
		errors.New("network is unreachable"),
	} {
		calls := 0
		_, got := retrying(t.Context(), func() (*api.HTTPResponse, error) {
			calls++
			return nil, err
		})
		require.ErrorIs(t, got, err)
		assert.Equal(t, 1, calls)
	}
}

func TestRetryingStopsOnCancelledContext(t *testing.T) {
	previous := retryDelays
	retryDelays = []time.Duration{time.Hour}
	t.Cleanup(func() { retryDelays = previous })

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	calls := 0
	_, err := retrying(ctx, func() (*api.HTTPResponse, error) {
		calls++
		return nil, retryable("E111009")
	})
	require.Error(t, err)
	assert.Equal(t, 1, calls)
}

func rateLimited(retryAfter int) error {
	err := &descope.Error{Code: descope.ErrRateLimitExceeded.Code, Description: "Rate limit exceeded"}
	if retryAfter > 0 {
		return err.WithInfo(descope.ErrorInfoKeys.RateLimitExceededRetryAfter, retryAfter)
	}
	return err
}

func fastRateLimitDelays(t *testing.T, count int) {
	previousDelays, previousJitter := rateLimitDelays, rateLimitJitter
	rateLimitDelays = make([]time.Duration, count)
	rateLimitJitter = func(time.Duration) time.Duration { return 0 }
	t.Cleanup(func() { rateLimitDelays, rateLimitJitter = previousDelays, previousJitter })
}

func TestRetryingSucceedsOnceTheRateLimitClears(t *testing.T) {
	fastRateLimitDelays(t, 6)
	calls := 0
	res, err := retrying(t.Context(), func() (*api.HTTPResponse, error) {
		calls++
		if calls < 5 {
			return nil, rateLimited(0)
		}
		return &api.HTTPResponse{BodyStr: "{}"}, nil
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, 5, calls)
}

func TestRetryingGivesUpOnAPersistentRateLimit(t *testing.T) {
	fastRateLimitDelays(t, 6)
	calls := 0
	_, err := retrying(t.Context(), func() (*api.HTTPResponse, error) {
		calls++
		return nil, rateLimited(0)
	})
	require.ErrorIs(t, err, descope.ErrRateLimitExceeded)
	assert.Equal(t, 7, calls)
}

func TestRetryingCountsRateLimitsAndTransientErrorsSeparately(t *testing.T) {
	fastDelays(t, 1)
	fastRateLimitDelays(t, 1)
	calls := 0
	res, err := retrying(t.Context(), func() (*api.HTTPResponse, error) {
		calls++
		switch calls {
		case 1:
			return nil, retryable("E111009")
		case 2:
			return nil, rateLimited(0)
		}
		return &api.HTTPResponse{BodyStr: "{}"}, nil
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, 3, calls)
}

func TestRetryingStopsWaitingForARateLimitOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	calls := 0
	_, err := retrying(ctx, func() (*api.HTTPResponse, error) {
		calls++
		return nil, rateLimited(60)
	})
	require.ErrorIs(t, err, descope.ErrRateLimitExceeded)
	assert.Equal(t, 1, calls)
}

func TestRateLimitDelay(t *testing.T) {
	previousDelays, previousJitter := rateLimitDelays, rateLimitJitter
	rateLimitDelays = []time.Duration{5 * time.Second, 10 * time.Second, 20 * time.Second, 40 * time.Second}
	rateLimitJitter = func(time.Duration) time.Duration { return 0 }
	t.Cleanup(func() { rateLimitDelays, rateLimitJitter = previousDelays, previousJitter })

	for _, tc := range []struct {
		name       string
		attempt    int
		retryAfter int
		want       time.Duration
	}{
		{name: "backs off exponentially without Retry-After", attempt: 0, want: 5 * time.Second},
		{name: "backs off exponentially without Retry-After on a later attempt", attempt: 3, want: 40 * time.Second},
		{name: "waits as long as Retry-After asks", attempt: 0, retryAfter: 60, want: 60 * time.Second},
		{name: "keeps backing off once the backoff passes Retry-After", attempt: 3, retryAfter: 30, want: 40 * time.Second},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := rateLimitDelay(descope.AsError(rateLimited(tc.retryAfter)), tc.attempt)
			require.True(t, ok)
			assert.Equal(t, tc.want, got)
		})
	}

	_, ok := rateLimitDelay(descope.AsError(rateLimited(60)), len(rateLimitDelays))
	assert.False(t, ok)
}

func TestRateLimitJitterSpreadsConcurrentRetries(t *testing.T) {
	for range 100 {
		jitter := rateLimitJitter(60 * time.Second)
		assert.GreaterOrEqual(t, jitter, time.Duration(0))
		assert.Less(t, jitter, 15*time.Second)
	}
}

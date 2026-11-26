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

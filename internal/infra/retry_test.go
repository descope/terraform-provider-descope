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
	previous := projectNotReadyDelays
	projectNotReadyDelays = make([]time.Duration, count)
	t.Cleanup(func() { projectNotReadyDelays = previous })
}

func notReady() error {
	return &descope.Error{Code: projectNotReadyErrorCode, Description: "Management key is forbidden"}
}

func TestRetryingSucceedsAfterProjectBecomesReady(t *testing.T) {
	fastDelays(t, 3)
	calls := 0
	res, err := retrying(t.Context(), func() (*api.HTTPResponse, error) {
		calls++
		if calls < 3 {
			return nil, notReady()
		}
		return &api.HTTPResponse{BodyStr: "{}"}, nil
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, 3, calls)
}

func TestRetryingGivesUpAndReturnsTheError(t *testing.T) {
	fastDelays(t, 3)
	calls := 0
	_, err := retrying(t.Context(), func() (*api.HTTPResponse, error) {
		calls++
		return nil, notReady()
	})
	require.Error(t, err)
	assert.Equal(t, projectNotReadyErrorCode, descope.AsError(err).Code)
	assert.Equal(t, 4, calls)
}

func TestRetryingLeavesOtherErrorsAlone(t *testing.T) {
	fastDelays(t, 3)
	for _, err := range []error{
		&descope.Error{Code: "E111008", Description: "some other forbidden"},
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
	previous := projectNotReadyDelays
	projectNotReadyDelays = []time.Duration{time.Hour}
	t.Cleanup(func() { projectNotReadyDelays = previous })

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	calls := 0
	_, err := retrying(ctx, func() (*api.HTTPResponse, error) {
		calls++
		return nil, notReady()
	})
	require.Error(t, err)
	assert.Equal(t, 1, calls)
}

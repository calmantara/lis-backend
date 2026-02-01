package utils

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetClientID(t *testing.T) {
	t.Run("returns empty string when ClientID is not set", func(t *testing.T) {
		ctx := t.Context()
		clientID := GetClientID(ctx)
		assert.Equal(t, "", clientID)
	})

	t.Run("returns empty string when ClientID is not a string", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), ClientID{}, 12345)
		clientID := GetClientID(ctx)
		assert.Equal(t, "", clientID)
	})

	t.Run("returns ClientID when it is set", func(t *testing.T) {
		expectedClientID := "test-client-id"
		ctx := SetClientID(t.Context(), expectedClientID)
		clientID := GetClientID(ctx)
		assert.Equal(t, expectedClientID, clientID)
	})
}

func TestSetClientID(t *testing.T) {
	t.Run("sets ClientID in context", func(t *testing.T) {
		expectedClientID := "test-client-id"
		ctx := SetClientID(t.Context(), expectedClientID)
		clientID := GetClientID(ctx)
		assert.Equal(t, expectedClientID, clientID)
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("returns uuid.Nil when UserID is not set", func(t *testing.T) {
		ctx := t.Context()
		userID := GetUserID(ctx)
		assert.Equal(t, uuid.Nil, userID)
	})

	t.Run("returns uuid.Nil when UserID is not a uuid.UUID", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), UserID{}, "not-a-uuid")
		userID := GetUserID(ctx)
		assert.Equal(t, uuid.Nil, userID)
	})

	t.Run("returns UserID when it is set", func(t *testing.T) {
		expectedUserID := uuid.New()
		ctx := SetUserID(t.Context(), expectedUserID)
		userID := GetUserID(ctx)
		assert.Equal(t, expectedUserID, userID)
	})
}

func TestSetUserID(t *testing.T) {
	t.Run("sets UserID in context", func(t *testing.T) {
		expectedUserID := uuid.New()
		ctx := SetUserID(t.Context(), expectedUserID)
		userID := GetUserID(ctx)
		assert.Equal(t, expectedUserID, userID)
	})
}

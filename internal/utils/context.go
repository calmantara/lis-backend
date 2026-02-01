package utils

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ClientID struct{}

func GetClientID(ctx context.Context) string {
	val := ctx.Value(ClientID{})
	if val == nil {
		return ""
	}
	clientID, ok := val.(string)
	if !ok {
		return ""
	}

	return clientID
}

func SetClientID(ctx context.Context, clientID string) context.Context {
	return context.WithValue(ctx, ClientID{}, clientID)
}

type UserID struct{}

func GetUserID(ctx context.Context) uuid.UUID {
	val := ctx.Value(UserID{})
	if val == nil {
		return uuid.Nil
	}
	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return userID
}

func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserID{}, userID)
}

type RequestID struct{}

func GetRequestID(ctx context.Context) string {
	val := ctx.Value(RequestID{})
	if val == nil {
		return ""
	}
	requestID, ok := val.(string)
	if !ok {
		return ""
	}

	return requestID
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestID{}, requestID)
}

type (
	EchoContextKey string
)

const (
	EchoClientID      EchoContextKey = "client_id"
	EchoUserID        EchoContextKey = "user_id"
	EchoApplicationID EchoContextKey = "application_id"
	EchoRequestID     EchoContextKey = "request_id"
)

func GetEchoContext(ctx echo.Context) (c context.Context) {
	clientID := ctx.Get(string(EchoClientID)).(string)
	c = SetClientID(ctx.Request().Context(), clientID)

	userID := uuid.Nil
	id, ok := ctx.Get(string(EchoUserID)).(string)
	if ok {
		userID, _ = uuid.Parse(id)
	}
	c = SetUserID(c, userID)

	requestID := ctx.Get(string(EchoRequestID)).(string)
	c = SetRequestID(c, requestID)

	return c
}

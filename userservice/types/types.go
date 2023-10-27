package types

import (
	"context"

	"go.opentelemetry.io/otel"
)

var UserServiceTracer = otel.Tracer("EchoWave/userservice")

type Service interface {
	// GetUserByUsername returns a user by username
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

package usecase

import "context"

type Auth interface {
	Register(ctx context.Context, in RegisterInput) (AuthOutput, error)
	Login(ctx context.Context, in LoginInput) (AuthOutput, error)
}

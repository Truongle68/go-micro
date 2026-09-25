package v1

import (
	"context"
	"user-service/internal/usecase"
)

type UserGRPCUsecase interface {
	GetProfile(ctx context.Context, id string) (*usecase.UserProfileDTO, error)
}

var _ UserGRPCUsecase = (*usecase.UserUC)(nil)

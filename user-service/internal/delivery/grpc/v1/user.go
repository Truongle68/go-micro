package v1

import (
	"context"
	"errors"
	"strings"
	"user-service/internal/domain"

	userv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/user/v1"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userv1.UnimplementedUserServiceServer
	uc UserGRPCUsecase
	l  logger.Interface
}

func NewUserServer(uc UserGRPCUsecase, l logger.Interface) *UserServer {
	return &UserServer{
		uc: uc,
		l:  l,
	}
}

func (s *UserServer) GetProfile(ctx context.Context, req *userv1.GetProfileRequest) (*userv1.GetProfileResponse, error) {
	if strings.TrimSpace(req.GetUserId()) == "" {
		return nil, status.Error(codes.InvalidArgument, domain.ErrEmptyUserID.Error())
	}

	dto, err := s.uc.GetProfile(ctx, req.GetUserId())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return nil, status.Errorf(codes.NotFound, "%v", err)
		default:
			s.l.Error("grpc.GetProfile: %v", err)
			return nil, status.Error(codes.Internal, "failed to get user profile")
		}
	}

	return &userv1.GetProfileResponse{
		UserId:          dto.ID,
		Email:           dto.Email,
		Phone:           dto.Phone,
		FullName:        dto.FullName,
		Status:          string(dto.Status),
		IsEmailVerified: dto.IsEmailVerified,
	}, nil
}

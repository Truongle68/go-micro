package grpc

import (
	"context"
	"fmt"
	"order-service/internal/client"
	"order-service/internal/domain"

	userv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type userGRPCClient struct {
	client userv1.UserServiceClient
}

func NewUserGRPCClient(target string) (client.UserClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user gRPC service at %s: %w", target, err)
	}

	return &userGRPCClient{
		client: userv1.NewUserServiceClient(conn),
	}, nil
}

func (c *userGRPCClient) GetProfile(ctx context.Context, id string) (*client.UserProfileDTO, error) {
	res, err := c.client.GetProfile(ctx, &userv1.GetProfileRequest{
		UserId: id,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, domain.ErrUserNotFound
		} else {
			return nil, err
		}
	}

	return &client.UserProfileDTO{
		UserID:          res.GetUserId(),
		Email:           res.GetEmail(),
		Phone:           res.GetPhone(),
		FullName:        res.GetFullName(),
		Status:          res.GetStatus(),
		IsEmailVerified: res.GetIsEmailVerified(),
	}, nil
}

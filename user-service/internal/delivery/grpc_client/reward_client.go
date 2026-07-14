package grpc_client

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"reward-service/pkg/pb"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// type RewardGRPCClient struct {
// 	client pb.RewardServiceClient
// 	conn   *grpc.ClientConn
// }

// func NewRewardGRPCClient(targetAddr string) (*RewardGRPCClient, error) {
// 	log.Printf("[UserService] Dialing Reward Service at: %s", targetAddr)
// 	conn, err := grpc.NewClient(targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		return nil, fmt.Errorf("could not connect to reward service: %w", err)
// 	}
// 	client := pb.NewRewardServiceClient(conn)
// 	return &RewardGRPCClient{
// 		client: client,
// 		conn:   conn,
// 	}, nil
// }

// func (c *RewardGRPCClient) Close() error {
// 	if c.conn != nil {
// 		return c.conn.Close()
// 	}
// 	return nil
// }

// // createRewardWallet implements the usecase.RewardGateway contract
// func (c *RewardGRPCClient) CreateRewardWallet(ctx context.Context, userID string, points int) error {
// 	log.Printf("[UserService] Calling Reward Service InitWallet for UserID: %s with Points: %d", userID, points)

// 	resp, err := c.client.InitWallet(ctx, &pb.InitWalletRequest{
// 		UserId:        userID,
// 		InitialPoints: int32(points),
// 	})
// 	if err != nil {
// 		return fmt.Errorf("gRPC InitWallet failed: %w", err)
// 	}

// 	log.Printf("[UserService] Reward Service InitWallet response: WalletID=%s, Status=%s", resp.GetWalletId(), resp.GetStatus())
// 	return nil
// }
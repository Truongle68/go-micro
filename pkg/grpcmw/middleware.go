package grpcmw

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	log.Printf("Incomming gRPC call: %s", info.FullMethod)

	res, err := handler(ctx, req)

	duration := time.Since(start)

	st, _ := status.FromError(err)
	statusCode := st.Code()

	if err != nil {
		log.Printf("FAIL: %s | Status: %s (%d) | Message: %s | Duration: %v", info.FullMethod, statusCode.String(), statusCode, st.Message(), duration)
	} else {
		log.Printf("SUCCESS: %s | Status: %s (%d) | Duration: %v", info.FullMethod, statusCode.String(), statusCode, duration)
	}

	return res, err
}

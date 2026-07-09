package grpcserver

import (
	"net"

	"google.golang.org/grpc"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

func ServerOptions(opts ...grpc.ServerOption) Option {
	return func(s *Server) {
		s.serverOpts = append(s.serverOpts, opts...)
	}
}

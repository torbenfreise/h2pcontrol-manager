package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/torbenfreise/h2pcontrol-manager/internal/registry"

	managergrpc "buf.build/gen/go/beyer-labs/h2pcontrol/grpc/go/h2pcontrol/manager/v1/managerv1grpc"
	managerpb "buf.build/gen/go/beyer-labs/h2pcontrol/protocolbuffers/go/h2pcontrol/manager/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
)

type server struct {
	managergrpc.UnimplementedManagerServiceServer
	sync.RWMutex
	registry   *registry.Registry
	grpcServer *grpc.Server
}

func (s *server) Register(ctx context.Context, in *managerpb.RegisterRequest) (*managerpb.RegisterResponse, error) {
	peerInfo, _ := peer.FromContext(ctx)
	addr := peerInfo.Addr.String()
	return s.registry.RegisterService(ctx, in, addr)
}

func (s *server) List(ctx context.Context, in *managerpb.ListRequest) (*managerpb.ListResponse, error) {
	return s.registry.List(ctx, in)
}

func (s *server) Heartbeat(stream grpc.BidiStreamingServer[managerpb.HeartbeatRequest, managerpb.HeartbeatResponse]) error {
	peerInfo, _ := peer.FromContext(stream.Context())
	addr := peerInfo.Addr.String()

	// Start a goroutine to send heartbeats to the client
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				pong := &managerpb.HeartbeatResponse{
					Healthy:   true,
					Timestamp: time.Now().Unix(),
				}
				if err := stream.Send(pong); err != nil {
					log.Printf("Error sending heartbeat pong: %v", err)
					close(done)
					return
				}
			case <-done:
				return
			}
		}
	}()

	for {
		ping, err := stream.Recv()
		if err != nil {
			log.Printf("Heartbeat stream closed from %v: %v", addr, err)
			s.registry.RemoveService(addr)
			close(done)
			return nil
		}

		s.registry.UpdateHeartbeat(addr)
		log.Printf("Received heartbeat from %v at %v", addr, time.Unix(ping.Timestamp, 0))
	}
}

func RunServer(port int) {
	lis, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := &server{
		registry: registry.NewRegistry(),
	}

	serverOpts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	s := grpc.NewServer(serverOpts...)
	srv.grpcServer = s
	managergrpc.RegisterManagerServiceServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Fatalf("Failed to serve: %v", err)
	}
}

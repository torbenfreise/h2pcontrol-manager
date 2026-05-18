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
	return s.registry.RegisterService(in, peerInfo.Addr.String())
}

func (s *server) List(_ context.Context, _ *managerpb.ListRequest) (*managerpb.ListResponse, error) {
	return s.registry.List()
}

func (s *server) Watch(_ *managerpb.WatchRequest, stream grpc.ServerStreamingServer[managerpb.WatchResponse]) error {
	id, ch := s.registry.Subscribe()
	defer s.registry.Unsubscribe(id)

	send := func() error {
		resp, err := s.registry.List()
		if err != nil {
			return err
		}
		return stream.Send(&managerpb.WatchResponse{Services: resp.GetServices()})
	}

	if err := send(); err != nil {
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ch:
			if err := send(); err != nil {
				return err
			}
		}
	}
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
				pong := &managerpb.HeartbeatResponse{}
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

		s.registry.UpdateHeartbeat(addr, ping.GetHealthy())
		log.Printf("Received heartbeat from %v at %v", addr, time.Now())
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

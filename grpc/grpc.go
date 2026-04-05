package grpc

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"h2pcontrol.manager/internal"

	managerv1 "buf.build/gen/go/beyer-labs/h2pcontrol/protocolbuffers/go/h2pcontrol/manager/v1"
	managerv1grpc "buf.build/gen/go/beyer-labs/h2pcontrol/grpc/go/h2pcontrol/manager/v1/managerv1grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	managerv1grpc.UnimplementedManagerServiceServer
	sync.RWMutex
	registry     *internal.ServerRegistry
	stub_service *internal.StubService
}

func (s *server) GetStub(ctx context.Context, in *managerv1.GetStubRequest) (*managerv1.GetStubResponse, error) {
	return s.stub_service.GetStub(ctx, in)
}

func (s *server) RegisterServer(ctx context.Context, in *managerv1.RegisterServerRequest) (*managerv1.RegisterServerResponse, error) {
	peerInfo, _ := peer.FromContext(ctx)
	addr := peerInfo.Addr.String()
	return s.registry.RegisterServer(ctx, in, addr)
}

func (s *server) FetchServers(ctx context.Context, in *managerv1.FetchServersRequest) (*managerv1.FetchServersResponse, error) {
	return s.registry.FetchServers(ctx, in)
}

func (s *server) FetchSpecificServer(ctx context.Context, in *managerv1.FetchSpecificServerRequest) (*managerv1.FetchSpecificServerResponse, error) {
	return s.registry.FetchSpecificServer(ctx, in)
}

func (s *server) Heartbeat(stream managerv1grpc.ManagerService_HeartbeatServer) error {
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
				pong := &managerv1.HeartbeatResponse{
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
			s.registry.RemoveServer(addr)
			close(done)
			return nil
		}

		s.registry.UpdateHeartbeat(addr)
		log.Printf("Received heartbeat from %v at %v", addr, time.Unix(ping.Timestamp, 0))
	}
}

func RunServer() {
	flag.Parse()

	lis, err := net.Listen("tcp4", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := &server{
		registry:     internal.NewServerRegistry(),
		stub_service: internal.NewStubService(),
	}

	serverOpts := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	}

	s := grpc.NewServer(serverOpts...)
	managerv1grpc.RegisterManagerServiceServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

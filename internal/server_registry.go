package internal

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	managerv1 "buf.build/gen/go/beyer-labs/h2pcontrol/protocolbuffers/go/h2pcontrol/manager/v1"
)

type ServerEntry struct {
	LastSeen  time.Time
	Metadata  *managerv1.ServerDefinition
	Heartbeat chan struct{}
}

type ServerRegistry struct {
	mu      sync.RWMutex
	servers map[string]*ServerEntry
}

func NewServerRegistry() *ServerRegistry {
	return &ServerRegistry{
		servers: make(map[string]*ServerEntry),
	}
}

func (r *ServerRegistry) RegisterServer(ctx context.Context, in *managerv1.RegisterServerRequest, addr string) (*managerv1.RegisterServerResponse, error) {

	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address format: %v", err)
	}

	// Create new address with peer's IP and original port
	newAddr := net.JoinHostPort(ip, port)

	entry := &ServerEntry{
		LastSeen:  time.Now(),
		Metadata:  in.Server,
		Heartbeat: make(chan struct{}),
	}

	log.Printf("Server wants to connect")

	r.mu.Lock()
	r.servers[newAddr] = entry
	r.mu.Unlock()

	log.Printf("Server connected: '%v' running '%v'", newAddr, in.Server.GetName())

	return &managerv1.RegisterServerResponse{
		Result: "Server registered successfully",
	}, nil
}

func (r *ServerRegistry) FetchServers(ctx context.Context, req *managerv1.FetchServersRequest) (*managerv1.FetchServersResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var serverList []*managerv1.FetchServerDefinition
	for addr, svc := range r.servers {
		// The port in addr is the port over which the manager is connected, not the service's open port.
		ip, _, _ := net.SplitHostPort(addr)
		server := managerv1.FetchServerDefinition{
			Name:        svc.Metadata.GetName(),
			Description: svc.Metadata.GetDescription(),
			Addr:        ip + ":" + svc.Metadata.GetPort(),
		}
		serverList = append(serverList, &server)
	}

	return &managerv1.FetchServersResponse{
		Servers: serverList,
	}, nil
}

func (r *ServerRegistry) RemoveServer(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.servers, addr)

}

func (r *ServerRegistry) UpdateHeartbeat(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entry, ok := r.servers[addr]; ok {
		entry.LastSeen = time.Now()
	}
}

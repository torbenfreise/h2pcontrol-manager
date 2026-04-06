package registry

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	manager "buf.build/gen/go/beyer-labs/h2pcontrol/protocolbuffers/go/h2pcontrol/manager/v1"
)

type Entry struct {
	LastSeen  time.Time
	Metadata  *manager.ServiceDefinition
	Heartbeat chan struct{}
}

type Registry struct {
	mu       sync.RWMutex
	services map[string]*Entry
}

func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]*Entry),
	}
}

func (r *Registry) RegisterService(ctx context.Context, in *manager.RegisterRequest, addr string) (*manager.RegisterResponse, error) {

	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address format: %v", err)
	}

	// Create new address with peer's IP and original port
	newAddr := net.JoinHostPort(ip, port)

	entry := &Entry{
		LastSeen:  time.Now(),
		Metadata:  in.Service,
		Heartbeat: make(chan struct{}),
	}

	log.Printf("Service wants to connect")

	r.mu.Lock()
	r.services[newAddr] = entry
	r.mu.Unlock()

	log.Printf("Service connected: '%v' running '%v'", newAddr, in.Service.GetName())

	return &manager.RegisterResponse{}, nil
}

func (r *Registry) List(ctx context.Context, req *manager.ListRequest) (*manager.ListResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var serverList []*manager.ServiceInfo
	for addr, svc := range r.services {
		// The port in addr is the port over which the manager is connected, not the service's open port.
		ip, _, _ := net.SplitHostPort(addr)
		server := manager.ServiceInfo{
			Name:        svc.Metadata.GetName(),
			Description: svc.Metadata.GetDescription(),
			Addr:        ip + ":" + svc.Metadata.GetPort(),
		}
		serverList = append(serverList, &server)
	}

	return &manager.ListResponse{
		Services: serverList,
	}, nil
}

func (r *Registry) RemoveService(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.services, addr)

}

func (r *Registry) UpdateHeartbeat(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entry, ok := r.services[addr]; ok {
		entry.LastSeen = time.Now()
	}
}

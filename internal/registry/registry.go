package registry

import (
	"log"
	"sync"
	"time"

	manager "buf.build/gen/go/beyer-labs/h2pcontrol/protocolbuffers/go/h2pcontrol/manager/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Entry struct {
	Definition *manager.ServiceDefinition
	LastSeen   time.Time
	Healthy    bool
	Heartbeat  chan struct{}
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

func (r *Registry) RegisterService(in *manager.RegisterRequest, addr string) (*manager.RegisterResponse, error) {
	entry := &Entry{
		LastSeen:   time.Now(),
		Definition: in.GetService(),
		Heartbeat:  make(chan struct{}),
	}

	r.mu.Lock()
	r.services[addr] = entry
	r.mu.Unlock()

	log.Printf("Service connected: '%v' ", in.GetService().GetName())

	return &manager.RegisterResponse{}, nil
}

func (r *Registry) List() (*manager.ListResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var serverList []*manager.ServiceInfo
	for _, svc := range r.services {
		server := manager.ServiceInfo{
			Definition: svc.Definition,
			Healthy:    svc.Healthy,
			LastSeen:   timestamppb.New(svc.LastSeen),
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

func (r *Registry) UpdateHeartbeat(addr string, healthy bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entry, ok := r.services[addr]; ok {
		entry.LastSeen = time.Now()
		entry.Healthy = healthy
	}
}

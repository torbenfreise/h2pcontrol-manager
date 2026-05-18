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
}

type Registry struct {
	mu          sync.RWMutex
	services    map[string]*Entry
	subscribers map[int]chan struct{}
	nextSubID   int
}

func NewRegistry() *Registry {
	return &Registry{
		services:    make(map[string]*Entry),
		subscribers: make(map[int]chan struct{}),
	}
}

// Subscribe returns an ID and a channel that receives a signal whenever the registry changes.
// The channel is buffered so a slow consumer misses at most one notification (not multiple).
func (r *Registry) Subscribe() (int, <-chan struct{}) {
	ch := make(chan struct{}, 1)
	r.mu.Lock()
	id := r.nextSubID
	r.nextSubID++
	r.subscribers[id] = ch
	r.mu.Unlock()
	return id, ch
}

func (r *Registry) Unsubscribe(id int) {
	r.mu.Lock()
	delete(r.subscribers, id)
	r.mu.Unlock()
}

// notify signals all subscribers that the registry has changed.
// Must be called while holding r.mu (write lock).
func (r *Registry) notify() {
	for _, ch := range r.subscribers {
		select {
		case ch <- struct{}{}:
		default: // already has a pending notification queued
		}
	}
}

func (r *Registry) RegisterService(in *manager.RegisterRequest, addr string) (*manager.RegisterResponse, error) {
	entry := &Entry{
		LastSeen:   time.Now(),
		Definition: in.GetService(),
	}

	r.mu.Lock()
	r.services[addr] = entry
	r.notify()
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
	r.notify()
}

func (r *Registry) UpdateHeartbeat(addr string, healthy bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if entry, ok := r.services[addr]; ok {
		entry.LastSeen = time.Now()
		if entry.Healthy != healthy {
			r.notify()
		}
		entry.Healthy = healthy
	}
}

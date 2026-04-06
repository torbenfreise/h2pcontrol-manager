# Welcome to the manager of the H2PControl project

The Manager is like the phonebook of the H2PControl project, it keeps track of all the servers that are registered to it and their status. It also keeps track of the status of the servers with a heartbeat.

### What the manager does

The h2pControl Manager is a gRPC service that acts as the discovery and contract hub for the H2PControl ecosystem.

It is responsible for:

- Registering servers when they come online, including server metadata (name, description, and port).
- Maintaining an in-memory server registry for fast lookup of available services.
- Monitoring server liveness through a bidirectional heartbeat stream and removing disconnected servers.
- Exposing APIs to fetch all registered servers with discoverable endpoint addresses.

In practice, this means each service in H2PControl can announce itself once and then be discovered by other components without manual endpoint tracking.

### Runtime flow

1. A service calls `RegisterServer` with metadata.
2. The manager stores the server entry in memory.
3. The service keeps a heartbeat stream open; each ping updates `LastSeen`.
4. If the stream closes, the manager removes that server from the registry.
5. Clients can call `FetchServers` to discover endpoint addresses.

### Core modules

- `main.go`: process entry point; starts the gRPC server.
- `grpc/grpc.go`: gRPC server setup, endpoint wiring, keepalive policy, heartbeat stream handling.
- `internal/server_registry.go`: registration, registry state, server listing, and heartbeat updates.

First, this is the repository structure. In this documentation we go through the most important files and folders, and how they contribute to the manager's role in H2PControl.

```bash
├── Dockerfile
├── go.mod
├── go.sum
├── grpc
│   └── grpc.go
├── h2pcontrol-manager-docs
│   ├── docs
│   │   ├── api-docs.md
│   │   └── index.md
│   └── mkdocs.yml
├── internal
│   └── server_registry.go
├── main.go
└── README.md

```

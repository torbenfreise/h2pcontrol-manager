# Welcome to the manager of the H2PControl project

The Manager is like the phonebook of the H2PControl project: it keeps track of all registered services and their status. It also tracks service health with a heartbeat stream.

### What the manager does

The h2pControl Manager is a gRPC service that acts as the discovery hub for the H2PControl ecosystem, exposed through a CLI application.

It is responsible for:

- Registering services when they come online, including service metadata (name, description, and port).
- Maintaining an in-memory service registry for fast lookup of available services.
- Monitoring service liveness through a bidirectional heartbeat stream and removing disconnected services.
- Exposing APIs to list all registered services with discoverable endpoint addresses.

In practice, this means each service in H2PControl can announce itself once and then be discovered by other components without manual endpoint tracking.

### Runtime flow

1. A service calls `Register` with metadata.
2. The manager stores the service entry in memory.
3. The service keeps a heartbeat stream open; each ping updates `LastSeen`.
4. If the stream closes, the manager removes that service from the registry.
5. Clients can call `List` to discover endpoint addresses.

### Core modules

- `main.go`: process entry point; executes the Cobra CLI.
- `cmd/`: CLI commands (`start`, `list`) and shared client setup.
- `internal/server/server.go`: gRPC server setup, endpoint wiring, keepalive policy, heartbeat stream handling.
- `internal/registry/registry.go`: registration, registry state, service listing, and heartbeat updates.

First, this is the repository structure. In this documentation we go through the most important files and folders, and how they contribute to the manager's role in H2PControl.

```bash
.
├── .goreleaser.yaml
├── cmd
│   ├── list.go
│   ├── root.go
│   └── start.go
├── go.mod
├── go.sum
├── h2pcontrol-manager-docs
│   ├── docs
│   │   ├── api-docs.md
│   │   ├── go-callvis-image.png
│   │   └── index.md
│   └── mkdocs.yml
├── internal
│   ├── registry
│   │   └── registry.go
│   └── server
│       └── server.go
├── main.go
└── README.md

```

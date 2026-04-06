# Welcome to the manager of the H2PControl project

The Manager is like the phonebook of the H2PControl project, it keeps track of all the servers that are registered to it and their status. It also keeps track of the status of the servers with a heartbeat.

### What the manager does

The h2pControl Manager is a gRPC service that acts as the discovery and contract hub for the H2PControl ecosystem.

It is responsible for:

- Registering servers when they come online, including server metadata (name, version, IP, port, and proto files).
- Maintaining an in-memory server registry for fast lookup of available services.
- Monitoring server liveness through a bidirectional heartbeat stream and removing disconnected servers.
- Exposing APIs to fetch all servers or inspect a specific server definition.
- Persisting uploaded proto files in a temporary local path so definitions can be retrieved later.
- Generating and returning language stubs (currently Python via `grpc_tools`/`betterproto`) as a zip archive on demand.

In practice, this means each service in H2PControl can announce itself once, then be discovered and consumed by other components without manual endpoint tracking.

### Runtime flow

1. A service calls `RegisterServer` with metadata and proto content.
2. The manager stores the server entry and writes proto files to temporary storage.
3. The service keeps a heartbeat stream open; each ping updates `LastSeen`.
4. If the stream closes, the manager removes that server from the registry.
5. Clients can call `FetchServers`/`FetchSpecificServer` to discover endpoints and definitions.
6. Clients can call `GetStub` to receive generated client code as a zip file.

### Core modules

- `main.go`: process entry point; starts the gRPC server.
- `grpc/grpc.go`: gRPC server setup, endpoint wiring, keepalive policy, heartbeat stream handling.
- `internal/server_registry.go`: registration, registry state, lookup endpoints, heartbeat updates, proto persistence.
- `pb/`: generated protobuf/grpc interfaces used by both server and clients.

First of all this is the repository structure, in this documentation we will go trough the most important files and folders of the project and how do they contribute to the general purpose of the manager and its role on H2PControl project.

```bash
├── Dockerfile
├── go.mod
├── go.sum
├── grpc
│   └── grpc.go
├── h2pcontrol-manager-docs
│   ├── docs
│   │   └── index.md
│   └── mkdocs.yml
├── internal
│   ├── server_registry.go
├── main.go
├── pb
│   ├── h2pcontrol_grpc.pb.go
│   └── h2pcontrol.pb.go
└── README.md

```

# API Documentation & Overview

This page explains every folder and file in the repository, what each one does, and how it contributes to the manager's main responsibilities:

- Service registration and discovery.
- Server liveness tracking via heartbeat.

## Repository map

Let's take the repository map from `index.md` and use it as a reference to explain the purpose of each file and folder in the project, so we do not need to jump back and forth.

```bash
├── Dockerfile
├── go.mod
├── go.sum
├── grpc
│   └── grpc.go
├── h2pcontrol-manager-docs
│   ├── docs
│   │   ├── api-docs.md
│   │   └── index.md
│   └── mkdocs.yml
├── internal
│   └── server_registry.go
├── main.go
└── README.md
```

## Root files

Let's start with the files that are in the root of the repository, these are the files that are not inside any folder.

- `main.go` - Minimal entry point that calls `grpc.RunServer()`.
- `Dockerfile` - Defines the container image and runtime for serving the manager.
- `go.mod` - Declares module path and direct/indirect dependencies.
- `go.sum` - Stores checksums for dependency versions.

And then the obvious files like `README.md` and `.gitignore` (not shown) for documentation and git configuration.

## `grpc` folder

#### `grpc/grpc.go`

This is the API transport layer and manager server bootstrap, it is the manager's network-facing API surface. Every external interaction enters through this file.

- Declares the `server` type implementing `managerv1grpc.ManagerServiceServer`.
- Wires RPC methods to internal business logic:
	- `RegisterServer` -> `internal.ServerRegistry.RegisterServer`
	- `FetchServers` -> `internal.ServerRegistry.FetchServers`
- Implements bidirectional `Heartbeat` stream:
	- Receives heartbeat pings from servers.
	- Sends periodic heartbeat pongs.
	- Updates last-seen timestamp.
	- Removes server when stream closes.
- Configures keepalive policy and starts gRPC listener on configurable port (default `50051`). With the `RunServer` function which is called in the `main.go` file.

## `internal` folder

#### `internal/server_registry.go`

This is the service registry and discovery core, this file gives the manager its "phonebook" behavior mostly.

- `ServerRegistry` holds an in-memory map of active servers.
- `RegisterServer`:
	- Accepts metadata from a server.
	- Normalizes address using the caller IP and connection port.
	- Stores server entry and timestamp.
- `FetchServers`:
	- Returns all registered services with discoverable endpoint addresses.
- `UpdateHeartbeat` and `RemoveServer` maintain liveness state.

## Protocol definitions

The project currently imports generated protobuf/gRPC packages from Buf (`buf.build/gen/...`). These generated definitions provide the service and message contracts used by the manager implementation.

## End-to-end summary

1. `main.go` starts the app.
2. `grpc/grpc.go` exposes and serves the gRPC API.
3. `internal/server_registry.go` stores server metadata and serves discoverable endpoints.
4. Buf-generated contracts enforce API compatibility.
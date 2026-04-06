# API Documentation & Overview

This page explains every folder and file in the repository, what each one does, and how it contributes to the manager's main responsibilities:

- Service registration and discovery.
- Server liveness tracking via heartbeat.
- Proto definition storage and retrieval.
- On-demand stub generation.

## Repository map

Lets take the map of the repository seen in the `index.md` file and use it as a reference to explain the purpose of each file and folder in the project. Also to not go back and forth between the index and this file.

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
├── helper
│   └── helper.go
├── internal
│   ├── server_registry.go
│   └── stub_service.go
├── main.go
├── pb
│   ├── h2pcontrol_grpc.pb.go
│   └── h2pcontrol.pb.go
└── README.md
```

## Root files

Let's start with the files that are in the root of the repository, these are the files that are not inside any folder.

- `main.go` - Minimal entry point that calls `grpc.RunServer()`.
- `Dockerfile` - Defines the container image, installing Go, Python, and necessary libraries for serving.
- `go.mod` - Declares module path and direct/indirect dependencies.
- `go.sum` - Stores checksums for dependency versions.

And then the obvious files like `README.md` and `.gitignore` (not shown) for documentation and git configuration.

## `grpc` folder

#### `grpc/grpc.go`

This is the API transport layer and manager server bootstrap, it is the manager's network-facing API surface. Every external interaction enters through this file.

- Declares the `server` type implementing `pb.ManagerServer`.
- Wires RPC methods to internal business logic:
	- `RegisterServer` -> `internal.ServerRegistry.RegisterServer`
	- `FetchServers` -> `internal.ServerRegistry.FetchServers`
	- `FetchSpecificServer` -> `internal.ServerRegistry.FetchSpecificServer`
	- `GetStub` -> `internal.StubService.GetStub` [REMOVED since stub generation was removed]
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
	- Normalizes address using caller port + provided server IP.
	- Stores server entry and timestamp.
	- Saves incoming proto files via `SaveProtoFiles`.
- `FetchServers`:
	- Returns all registered services with discoverable endpoint addresses.
- `FetchSpecificServer`:
	- Returns one server plus concatenated proto text read from temp storage.
- `UpdateHeartbeat` and `RemoveServer` maintain liveness state.
- `SaveProtoFiles` writes proto files under temp directory (`h2pcontrol_proto/...`).

#### `internal/stub_service.go` [REMOVED]
This file was responsible for generating client stubs from proto definitions and packaging them as zip files. It used `grpc_tools` or `betterproto` to compile the proto files into Python code, then created a zip archive to return to clients requesting stubs.

## `pb` folder

These files are generated code from `protoc-gen-go` and `protoc-gen-go-grpc` based on the protobuf definitions (not shown) that define the gRPC service and message contracts. They are used by both the server implementation and any clients to ensure type-safe interactions.

## End-to-end summary

1. `main.go` starts the app.
2. `grpc/grpc.go` exposes and serves the gRPC API.
3. `internal/server_registry.go` stores and serves service metadata/proto definitions.
4. `internal/stub_service.go` generates downloadable client stubs. [REMOVED]
5. `pb/*` enforces API contracts.


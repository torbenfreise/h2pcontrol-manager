# API Documentation & Overview

This page explains the key folders and files in the repository, what each one does, and how it contributes to the manager's main responsibilities:

- Service registration and discovery.
- Server liveness tracking via heartbeat.

## Repository map

Let's take the repository map from `index.md` and use it as a reference to explain the purpose of each file and folder in the project, so we do not need to jump back and forth.

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

## Root files

Let's start with the files that are in the root of the repository, these are the files that are not inside any folder.

- `main.go` - Minimal entry point that calls `cmd.Execute()`.
- `go.mod` - Declares module path and direct/indirect dependencies.
- `go.sum` - Stores checksums for dependency versions.
- `.goreleaser.yaml` - Release packaging configuration.

And then the obvious files like `README.md` and `.gitignore` (not shown) for documentation and git configuration.

## `cmd` folder

#### `cmd/root.go`

This is the CLI root command setup and shared client wiring.

- Defines the `h2pcontrol` root command.
- Configures manager address via flag/env (`--manager` and `H2PCONTROL_MANAGER`, default `127.0.0.1:50051`).
- Creates and stores a gRPC client in command context for subcommands.
- Closes the shared gRPC connection after command execution.

#### `cmd/start.go`

Server startup command.

- Exposes `start` command.
- Accepts `--port` flag (default `50051`).
- Calls `internal/server.RunServer(port)`.

#### `cmd/list.go`

Service discovery command.

- Exposes `list` command.
- Reuses shared client setup from `root.go` pre/post hooks.
- Calls `List` RPC and prints registered services.

## `internal` folder

#### `internal/server/server.go`

This is the gRPC transport and manager server runtime.

- Declares the server type implementing `ManagerServiceServer`.
- Wires RPC methods to registry logic:
	- `Register` -> `registry.RegisterService`
	- `List` -> `registry.List`
- Implements bidirectional `Heartbeat` stream:
	- Receives heartbeat pings from services.
	- Sends periodic heartbeat pongs.
	- Updates last-seen timestamp.
	- Removes service when stream closes.
- Configures keepalive and starts gRPC listener.

#### `internal/registry/registry.go`

This is the service registry and discovery core.

- `Registry` holds an in-memory map of active services.
- `RegisterService`:
	- Accepts metadata from a service.
	- Normalizes address using the caller IP and connection port.
	- Stores service entry and timestamp.
- `List`:
	- Returns all registered services with discoverable endpoint addresses.
- `UpdateHeartbeat` and `RemoveService` maintain liveness state.

## Protocol definitions

The project currently imports generated protobuf/gRPC packages from Buf (`buf.build/gen/...`). These generated definitions provide the service and message contracts used by the manager implementation.

## End-to-end summary

1. `main.go` starts the app.
2. `cmd/start.go` starts the manager server.
3. `internal/server/server.go` exposes and serves the gRPC API.
4. `internal/registry/registry.go` stores service metadata and serves discoverable endpoints.
5. `cmd/list.go` consumes the manager API to list active services.
6. Buf-generated contracts enforce API compatibility.
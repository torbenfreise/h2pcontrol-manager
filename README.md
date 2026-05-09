# h2pControl Manager

The h2pControl Manager is the central service registry for H2PControl. This repository provides a CLI application (
`h2pmanager`) that can:

- Start the manager gRPC server.
- Query and list currently registered services from a running manager.

For detailed documentation, see the [h2pcontrol-manager-docs](./h2pcontrol-manager-docs) folder.

## Installing

[Go](https://go.dev/dl/) must be installed. Then run:

```bash
go install github.com/torbenfreise/h2pcontrol-manager/cmd/h2pmanager@latest
```

This places the `h2pmanager` binary in your Go bin directory (typically `~/go/bin`). Make sure it is on your `PATH`:

**bash** — add this line to your `~/.bashrc`:
```bash
export PATH="$HOME/go/bin:$PATH"
```

**zsh** — add this line to your `~/.zshrc`:
```bash
export PATH="$HOME/go/bin:$PATH"
```

Then reload your shell (`source ~/.bashrc` or `source ~/.zshrc`) and `h2pmanager` will be available globally.

## Quick usage

Start the manager listening on the supplied port:
```bash
h2pmanager start --port 50051
```

List services from a running manager:

```bash
h2pmanager list --address 127.0.0.1:50051
```

You can also configure the manager address with the `H2PMANAGER_ADDRESS` environment variable.

## Documentation

For more detailed documentation, navigate to the `h2pcontrol-manager-docs` folder. It is an MkDocs site that can be
viewed locally.
Install MkDocs with `pip`, `uv`, or your preferred package manager.

```bash
pip install mkdocs
# or
uv pip install mkdocs
```

To view the documentation, you can use the command below, then you can open the local link in your browser to view the
documentation.

```bash
cd h2pcontrol-manager-docs
mkdocs serve

# You will see something like this in the terminal:
INFO    -  Building documentation...
INFO    -  Cleaning site directory
INFO    -  Documentation built in 0.05 seconds
INFO    -  [17:04:28] Watching paths for changes: 'docs', 'mkdocs.yml'
INFO    -  [17:04:28] Serving on http://127.0.0.1:8000/

```

## Formatting and Linting

This project uses [golangci-lint](https://golangci-lint.run) for linting. The same configuration runs in CI via
the [lint workflow](.github/workflows/lint.yml).
The linter can be configured in [.golangci.yml](.golangci.yml)

To run the linter locally:

```bash
golangci-lint run
```

To format the code:

```bash
gofmt -w .
```
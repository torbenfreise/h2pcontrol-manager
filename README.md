# h2pControl Manager

The h2pControl Manager is the central registry service for H2PControl. It exposes a gRPC API that lets servers register themselves, keeps an up-to-date list of active servers through heartbeat streaming, stores uploaded proto definitions, and returns generated client stubs (currently Python) as zip files so other services can integrate quickly.

For more in detail documentation, see the [h2pcontrol-manager-docs](./h2pcontrol-manager-docs) folder.

## Documentation

For more detailed documentation please navigate to the `h2pcontrol-manager-docs` folder. It is an mkdocs documentation that can be easily be view locally on a well formatted website.
You are going to need to install mkdocs with pip or uv or your desired package manager.

```bash
pip install mkdocs
# or
uv pip install mkdocs
```

To view the documentation, you can use the command below, then you can open the local link in your browser to view the documentation.

```bash
mkdocs serve

# You will see something like this in the terminal:
INFO    -  Building documentation...
INFO    -  Cleaning site directory
INFO    -  Documentation built in 0.05 seconds
INFO    -  [17:04:28] Watching paths for changes: 'docs', 'mkdocs.yml'
INFO    -  [17:04:28] Serving on http://127.0.0.1:8000/

```

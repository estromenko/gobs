# Gobs

Tool for infrastructure monitoring.

## Getting started

### Building from source

```bash
go build .
```

### Running

Single host:

````bash
./gobs -exporter -server
``````

Multiple hosts:

```bash
# On worker hosts
./gobs -exporter -address=manager-host-addr:3333

# On manager host
./gobs -server
```

### Recieving metrics:

```bash
curl localhost:3333/metrics
```

### Get info about available commands and settings:

```bash
./gobs --help
```

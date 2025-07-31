# kayakDB

Light-weight, distributed key–value store powered by the Raft consensus algorithm.

![kayakdb logo](./docs/images/logo.png)

---

## Table of contents

1.  [Why kayakDB?](#why-kayakdb)
2.  [Quick start](#quick-start)
3.  [Configuration](#configuration)
4.  [The API server](#the-api-server)
5.  [Raft under the hood](#raft-under-the-hood)
6.  [Worker-pool](#worker-pool)
7.  [Command-line client `kayakctl`](#command-line-client-kayakctl)
8.  [License](#license)

---

## Why kayakDB?

*   **Simple** – a single binary for the server and another one for the client.
*   **Eventually consistent** – replication powered by the [Raft](https://raft.github.io/) algorithm converges all nodes to the same state over time.
*   **Embeddable** – only standard Go + a few small dependencies.
*   **Human-friendly** – plain TCP protocol, ergonomic CLI, JSON/YAML configuration.

---

## Quick start

### 1. Build (or install) the server

```
# From the project root
$ go build -o kayakdb ./cmd
# Or run directly without building
$ go run ./cmd
```

### 2. Provide (optional) configuration

Out of the box the server starts on port **8080** (client traffic) and **9090** (Raft traffic).  Every setting can be overridden via **`raft.json`** or environment variables – see the [Configuration](#configuration) section.

If you are just playing on your laptop you can skip this step entirely.

### 3. Fire it up

```
$ ./kayakdb             # or: go run ./cmd
INFO	Server is starting
INFO	Server is Listening on {"port": "8080"}
INFO	Raft Server is listening on port: 9090
```

In another terminal use the CLI:

```
$ kayakctl put str:name str:alice
$ kayakctl get str:name
┌──────┬───────┐
│ key  │ value │
├──────┼───────┤
│ name │ alice │
└──────┴───────┘
```

> **Tip** – both tools expose `-h/--help` for all commands and flags.

---

## Configuration

Configuration is loaded in three layers (lowest precedence first):

1. **Defaults** – hard-coded in `config.Configuration` struct tags.
2. **`raft.json`** – optional JSON file placed next to the binary.
3. **Environment variables** – override everything.

Below is a minimal example file:

```json
{
  "kayak_port": "8080",
  "raft_port": "9090",
  "log_level": "debug",
  "worker_pool_size": 4,
  "wait_queue_size": 1000,
  "peer_discovery": false,
  "seed_peers": ["127.0.0.1:9090"]
}
```

| Field | Env var | Default | Description |
|-------|---------|---------|-------------|
| `kayak_port` | `KAYAK_PORT` | `8080` | TCP port for client requests |
| `raft_port`  | `RAFT_PORT`  | `9090` | TCP port for Raft internal RPCs |
| `log_level`  | `LOG_LEVEL`  | `info` | `debug`, `info`, `warn`, `error` |
| `max_log_batch` | `MAX_LOG_BATCH` | `50` | How many log entries are sent in a single replication batch |
| `worker_pool_size` | `WORKER_POOL_SIZE` | `4` | Goroutines processing asynchronous jobs |
| `wait_queue_size`  | `WAIT_QUEUE_SIZE`  | `1000` | Size of the worker-pool queue |
| `peer_discovery`   | `PEER_DISCOVERY`  | `false` | Use DNS-SRV service discovery instead of static peers |
| `service_name` | `SERVICE_NAME` | `kayakdb` | DNS-SRV record when discovery is enabled |
| `seed_peers` | – | – | Array of `host:port` strings for the initial cluster |

---

## The API server

*   Listens on **`kayak_port`** (default **8080**).
*   Accepts raw TCP connections – there is **no HTTP** layer for maximum throughput.
*   Messages are encoded using the custom [`types.Payload`](types/) binary format.  Two endpoints are currently available:
    * **`/put`** – store one or more key/value pairs.
    * **`/get`** – retrieve the current value for a given key.
*   Once a request is received the server delegates the heavy-lifting to the embedded Raft library and eventually sends back a binary response which the CLI converts to a pretty table.

> The API is intentionally minimal at this stage; it will grow as kayakDB matures.

---

## Raft under the hood

The implementation lives in [`raft/`](raft/) and is completely self-contained.  Highlights:

*   **Leader election**, **log replication** and **state machine application** closely follow the Raft paper.
*   Persistence is abstracted behind `raft/storage.Driver` – the default is an in-memory store suitable for tests and local development.
*   **Raft RPCs** (`AppendEntries`, `RequestVote`, `Ping`) are served over Go’s `net/rpc` on the **`raft_port`** (9090 by default).
*   Concurrency is handled via the project’s [worker-pool](#worker-pool); each outgoing RPC is queued as an asynchronous job keeping the critical Raft logic free from goroutine bookkeeping.

If you want to embed kayakDB as a library you can simply:

```go
cfg := &config.Configuration{}
utils.LoadConfigurations(cfg)
logger := utils.InitLogger(cfg.LogLevel)
raft := raft.NewRaft(cfg, logger)
raft.Start()
```

---

## Worker-pool

The generic worker-pool lives in [`utils/workerpool.go`](utils/workerpool.go):

*   Fixed number of workers (`Size`), configurable queue (`WaitQueue`).
*   Jobs are ordinary functions (`func(...any) error`) with optional post-processing hooks.
*   The pool is **only** used by the Raft library at the moment – for example, log replication (`Append` RPC), leader heartbeats (`Ping`), and vote requests are dispatched asynchronously through it.

Feel free to reuse the pool in your own code – it has 100% test coverage.

---

## Command-line client `kayakctl`

The CLI, found under [`cli/`](cli/), offers a friendly way to interact with the server.
It is built with Cobra and distributed as a single binary.

### Installation

```
$ go install github.com/MohammedShetaya/kayakdb/cli@latest
```

### Global flags

```
-d, --hostname   Server hostname (default: "localhost")
-p, --port       Server port     (default: "8080")
```

### Examples

Store a value:

```
$ kayakctl put str:city str:Berlin
```

Fetch it back:

```
$ kayakctl get str:city
┌─────────┬───────────┐
│ key     │ value     │
├─────────┼───────────┤
│ country │ Palestine │
└─────────┴───────────┘
```

Values can be typed explicitly (`str:`, `num:`, `bool:`) or left for auto-detection.


---

## License

Licensed under the MIT license – see the [LICENSE](LICENSE) file for details.

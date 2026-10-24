<div align="center">

# 🔴 GoRedis

### A Redis-compatible server, built from the socket up, in pure Go, with zero dependencies.

[![Go Version](https://img.shields.io/badge/go-1.26.2-00ADD8?logo=go&logoColor=white)](go.mod)
[![Dependencies](https://img.shields.io/badge/dependencies-0-success)](go.mod)
[![License](https://img.shields.io/badge/license-unspecified-lightgrey)](#-license)
[![Build Status](https://img.shields.io/badge/build-add_CI_workflow-inactive)](#-roadmap)
[![Protocol](https://img.shields.io/badge/protocol-RESP2-DC382D?logo=redis&logoColor=white)](https://redis.io/docs/latest/develop/reference/protocol-spec/)

[![GoRedis demo placeholder, replace with a terminal GIF of redis-cli talking to GoRedis](https://via.placeholder.com/900x420.png?text=GoRedis+Demo+%28replace+with+terminal+GIF%29)](#-usage--examples)

*Speaks real RESP2. Any real Redis client, including `redis-cli`, can talk to it without knowing the difference.*

</div>

---

## 💡 The Problem & The Solution

Learning how Redis actually works by reading its C source is a slog: decades of optimizations, legacy code paths, and build complexity stand between you and the core ideas, namely a binary wire protocol, a dispatch table, and a concurrent in-memory store.

**GoRedis strips that away.** It re-implements the pieces that make Redis *Redis*: the RESP2 wire format, inline/array command parsing, pipelining, and a thread-safe keyspace, in under 500 lines of readable, tested, standard-library-only Go. You get a server real Redis tooling can connect to, and a codebase small enough to read end-to-end in one sitting.

> [!NOTE]
> This is **Phase 1 (Foundation)** of a multi-phase build. It is intentionally minimal: one data type (strings), no persistence, no expiry. See the [Roadmap](#-roadmap) for what's next.

---

## ✨ Key Features

- ✅ **Real RESP2 wire protocol**: parsed from scratch, no serialization library involved
- ✅ **Dual input modes**: accepts both RESP arrays (`*2\r\n$3\r\nGET\r\n...`) and inline/telnet-style commands (`PING\r\n`)
- ✅ **Command pipelining**: multiple commands sent back-to-back on one connection, no round-trip stalls
- ✅ **Thread-safe in-memory store**: `map[string]string` guarded by a single `sync.RWMutex`
- ✅ **Concurrent connections**: one goroutine per TCP connection via the standard `net` accept loop
- ✅ **13 commands implemented**: server (`PING`, `ECHO`) plus the full string suite (`GET`/`SET`/`SETNX`/`DEL`/`EXISTS`/`INCR`/`DECR`/`APPEND`/`STRLEN`/`MGET`/`MSET`)
- ✅ **Zero external dependencies**: only `net`, `bufio`, `sync`, `log/slog`, `flag`, `strconv`, `strings`, `errors` from the Go standard library
- ✅ **Compatible with real Redis clients**: `redis-cli`, raw `nc`, or any RESP2-speaking library works unmodified
- ✅ **Fully unit and integration tested**: 37 test cases across protocol, storage, command, and live-socket layers

---

## 🛠️ Tech Stack & Architecture

### Tech Stack

| Layer | Technology |
|---|---|
| Language | ![Go](https://img.shields.io/badge/-Go_1.26.2-00ADD8?logo=go&logoColor=white) |
| Networking | `net` (stdlib TCP sockets) |
| Protocol | RESP2 (hand-rolled encoder/decoder, no SDK) |
| Concurrency | goroutines + `sync.RWMutex` |
| Logging | `log/slog` (structured, stdlib) |
| CLI flags | `flag` (stdlib) |
| Testing | `testing` (stdlib), no mocking framework |
| Dependencies | **None.** `go.mod` declares zero third-party modules. |

### Architecture / How It Works

[![Architecture diagram placeholder, replace with a diagram of the TCP connection, RESP Reader, Dispatch, and Store flow](https://via.placeholder.com/800x360.png?text=Architecture+Diagram+%28placeholder%29)](#)

GoRedis is organized as one Go package per responsibility: no god-objects, no shared mutable globals outside the store.

| Package | Responsibility |
|---|---|
| [`internal/resp`](internal/resp) | Wire format. `Value` (tagged union over `+`/`-`/`:`/`$`/`*`), `Encode()` to bytes, and a buffered `Reader` that parses both RESP arrays and inline commands |
| [`internal/storage`](internal/storage) | The keyspace: a `map[string]string` behind a `sync.RWMutex`, exposing `Get`/`Set`/`SetNX`/`Del`/`Exists`/`Incr`/`Decr`/`Append`/`StrLen`/`MGet`/`MSet` |
| [`internal/commands`](internal/commands) | A `map[string]Handler` dispatch table. `Dispatch()` upper-cases the command name, looks up the handler, and routes `args[1:]` to it |
| [`internal/server`](internal/server) | The TCP accept loop. `Serve(ln)` takes a pre-bound `net.Listener` (so tests can bind ephemeral ports); each accepted connection gets its own goroutine and its own `resp.Reader` |
| [`cmd/server`](cmd/server) | The binary entrypoint. Parses `--port` and calls `server.New(addr).ListenAndServe()` |

**Request lifecycle**, end to end:

```
TCP connection accepted
        |
        v
 resp.Reader.Read()    parses one line/array off the buffered socket
                        (handles RESP arrays AND inline text identically)
        |
        v
 commands.Dispatch()   strings.ToUpper(name) -> table lookup -> handler(store, args)
        |
        v
 storage.Store         RLock for reads, Lock for writes, on a plain Go map
        |
        v
 resp.Value.Encode()   serialized back to the client as RESP bytes
        |
        v
      loop (same goroutine reads the next pipelined command)
```

Because `Read()` is called in a tight loop per-connection and the underlying `bufio.Reader` only blocks when its buffer is empty, a client can fire several commands in one `Write()` and GoRedis will drain and answer them all without waiting on round trips. This is exercised directly by `TestServerPipelinedCommands` and `TestReadMultipleCommandsPipelined`.

---

## 🚀 Getting Started

### Prerequisites

| Requirement | Version | Why |
|---|---|---|
| [Go](https://go.dev/dl/) | `1.26.2` or newer (per `go.mod`) | Build/run the server, no other toolchain needed |
| `redis-cli` *(optional)* | any | Easiest way to exercise the server; install via `redis-tools`/`brew install redis` |
| `nc` / `netcat` *(optional)* | any | Fallback client if you don't want to install `redis-cli` |

> [!TIP]
> There is **no Docker, no database, no external service** to stand up. `go build` and you're done.

### Installation

```bash
# 1. Clone the repository
git clone <your-fork-or-origin-url>
cd "build your redis"

# 2. Download dependencies (there are none, this is a no-op, but confirms go.mod integrity)
go mod download

# 3. Build the server binary
go build -o goredis ./cmd/server

# 4. Run it
./goredis --port 6379
```

You should see a structured log line confirming the listener is up:

```
time=2026-10-24T12:00:00.000-04:00 level=INFO msg="goredis listening" addr=[::]:6379
```

### Environment Variables

**There are none.** GoRedis reads no environment variables and loads no `.env` file; `cmd/server/main.go` only defines a single `flag.String`. All runtime configuration is the CLI flag below:

| Flag | Default | Description |
|---|---|---|
| `--port` | `6379` | TCP port the server listens on |

```bash
# Example: run on a non-default port
./goredis --port 7000
```

---

## 📖 Usage & Examples

### Running the App

```bash
# Development: run directly from source, no build step
go run ./cmd/server --port 6379

# Production-style: build once, run the binary
go build -o goredis ./cmd/server && ./goredis --port 6379
```

### Talking to it with redis-cli

```bash
redis-cli -p 6379 ping                 # PONG
redis-cli -p 6379 set foo bar          # OK
redis-cli -p 6379 get foo              # "bar"
redis-cli -p 6379 incr counter         # (integer) 1
redis-cli -p 6379 mset a 1 b 2         # OK
redis-cli -p 6379 mget a b missing     # "1" "2" (nil)
redis-cli -p 6379 append foo baz       # (integer) 6
redis-cli -p 6379 strlen foo           # (integer) 6
```

### Talking to it over raw RESP (no client library at all)

```bash
printf '*1\r\n$4\r\nPING\r\n' | nc localhost 6379
# +PONG
```

### Talking to it inline (telnet-style, no array framing needed)

```bash
printf 'PING\r\nECHO hello\r\n' | nc localhost 6379
# +PONG
# $5
# hello
```

### Using the command table as a library (Go)

Because `internal/commands.Dispatch` just takes a `*storage.Store` and `[]string`, you can drive it directly in tests or tooling without opening a socket:

```go
store := storage.New()
reply := commands.Dispatch(store, []string{"SET", "foo", "bar"})
// reply == resp.SimpleStringValue("OK")

reply = commands.Dispatch(store, []string{"GET", "foo"})
// reply == resp.BulkStringValue("bar")
```

[![Terminal usage placeholder, replace with a screen recording of a redis-cli session](https://via.placeholder.com/900x300.png?text=Terminal+Session+%28placeholder%29)](#)

---

## ⚙️ Configuration & Customization

GoRedis currently has **no config file, no feature flags, and no plugin system**, by design, at this phase. The entire configurable surface is the `--port` flag documented above.

- To **add a new command**: write a `Handler` (`func(*storage.Store, []string) resp.Value`) and register it in the `table` map in [`internal/commands/commands.go`](internal/commands/commands.go). No framework or registration ceremony beyond that one map entry.
- To **change the listen address beyond a bare port**, edit `cmd/server/main.go`; `server.New` already accepts a full `addr string` (e.g. `"127.0.0.1:6379"`), only the flag is restricted to a port today.
- A config-file system (per the `redis.conf`-style setup Redis itself uses) is **not implemented**, see [Roadmap](#-roadmap).

---

## 🧪 Testing & Development

### Running Tests

```bash
# Full suite
go test ./...

# Verbose, see every test name
go test ./... -v

# With the race detector (recommended, the store and server are concurrent)
go test ./... -race
```

The suite currently covers (37 tests across 4 files):

| Package | File | Focus |
|---|---|---|
| `internal/resp` | `value_test.go` | RESP `Encode()` for every type, including nulls |
| `internal/resp` | `reader_test.go` | Array parsing, inline parsing, null bulk strings, pipelining |
| `internal/storage` | `store_test.go` | Every store method, including the non-integer `INCR` error path |
| `internal/commands` | `commands_test.go` | Dispatch table, argument-count validation, case-insensitivity |
| `internal/server` | `server_test.go` | **Real TCP sockets**: PING, SET/GET, inline commands, pipelining, end-to-end |

### Linting & Formatting

No third-party linter is configured (no `.golangci.yml` present). Use the standard Go toolchain:

```bash
go fmt ./...     # format
go vet ./...     # static analysis / common mistakes
```

### Project Structure

```
build your redis/
├── cmd/
│   └── server/
│       └── main.go          # Entrypoint: flag parsing, server bootstrap
├── internal/
│   ├── resp/                 # RESP2 wire format: Value, Encode(), Reader
│   │   ├── value.go
│   │   ├── reader.go
│   │   └── *_test.go
│   ├── storage/               # Thread-safe map[string]string store
│   │   ├── store.go
│   │   └── store_test.go
│   ├── commands/               # Command name to handler dispatch table
│   │   ├── commands.go         # table, Dispatch(), PING/ECHO
│   │   ├── strings.go          # GET/SET/SETNX/DEL/EXISTS/INCR/DECR/APPEND/STRLEN/MGET/MSET
│   │   └── commands_test.go
│   └── server/
│       ├── server.go           # TCP accept loop, per-connection goroutine
│       └── server_test.go
├── go.mod                      # module goredis, go 1.26.2, zero deps
└── README.md
```

> [!NOTE]
> `.gitignore` already reserves `/goredis-cli`, `*.rdb`, and `*.aof`, a sign of where the roadmap is headed (a bundled CLI client, RDB/AOF persistence) even though those artifacts don't exist yet.

---

## 🩹 Troubleshooting & FAQ

<details>
<summary><strong>The server exits immediately with "address already in use"</strong></summary>

Another process (possibly a previous `goredis` run, or a real Redis instance) is already bound to port `6379`. Either stop it or pick a different port:

```bash
./goredis --port 6380
```
</details>

<details>
<summary><strong>My client hangs forever after sending a command</strong></summary>

`internal/resp/reader.go`'s `readBulkString` calls `io.ReadFull` for exactly the byte count declared in the `$<n>` header, plus the trailing `\r\n`. If your hand-written RESP payload declares the wrong length (e.g. `$3\r\nbar\r\n` but the string is actually 4 bytes), the reader will block waiting for bytes that never arrive. Double-check byte counts, or just use `redis-cli`/a real RESP client so framing is never hand-rolled.
</details>

<details>
<summary><strong>I restarted the server and all my data is gone</strong></summary>

Expected: GoRedis is currently **in-memory only**, with no AOF/RDB persistence (see [Roadmap](#-roadmap)). Every restart starts from an empty `storage.Store`.
</details>

---

## 🤝 Contributing & Roadmap

### Contributing

1. Fork the repo and create a branch: `git checkout -b feature/<short-description>` (or `fix/<short-description>` for bug fixes).
2. Keep packages single-responsibility: a new command goes in `internal/commands`, a new data type gets its own file there (mirroring `strings.go`), not scattered across the dispatch table.
3. Add tests alongside the code you touch. The project has no untested package today; keep it that way.
4. Run `go test ./... -race` and `go vet ./...` before opening a PR.
5. Open a PR describing *what* changed and *why*. The project has no issue/PR template yet, so a clear description in the PR body is enough.

### Roadmap

**Phase 1, Foundation** *(current)*

- [x] RESP2 protocol (arrays + inline commands + pipelining)
- [x] Thread-safe in-memory string store
- [x] Core string commands (`GET`/`SET`/`SETNX`/`DEL`/`EXISTS`/`INCR`/`DECR`/`APPEND`/`STRLEN`/`MGET`/`MSET`)
- [x] Server commands (`PING`/`ECHO`)
- [x] Unit + live-socket integration tests

**Future phases** *(not yet implemented)*

- [ ] Lists, Sets, Hashes + generic keyspace commands (`TYPE`, `KEYS`, `EXPIRE`)
- [ ] Multi-database support (`SELECT`)
- [ ] TTL / active key expiration
- [ ] Persistence: AOF and RDB + config file (`redis.conf`-style)
- [ ] Pub/Sub
- [ ] Transactions (`MULTI`/`EXEC`)
- [ ] Replication
- [ ] Observability: metrics, structured access logging, slowlog
- [ ] Standalone `goredis-cli` client + benchmarking tool
- [ ] Docker image + CI pipeline

---

## 📄 Footer

### License

> [!WARNING]
> **No `LICENSE` file is currently present in this repository.** Until one is added, default copyright law applies and no reuse rights are granted. Add a `LICENSE` file (e.g. MIT, Apache-2.0) to make the terms explicit.

### Acknowledgements

- [Redis](https://redis.io/), the original project this re-implements the wire protocol and command semantics of.
- [RESP2 Protocol Specification](https://redis.io/docs/latest/develop/reference/protocol-spec/), the primary reference for the `internal/resp` implementation.
- The Go standard library: `net`, `bufio`, `sync`, `log/slog`, `flag` did all the heavy lifting; no third-party package was needed.

### Contact / Support

- 💬 Discord: `<placeholder-invite-link>`
- 🐦 Twitter/X: `@<placeholder-handle>`
- ✉️ Email: `<placeholder@example.com>`
- 🐛 Found a bug? Open an issue on this repository.

</div>

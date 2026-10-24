# GoRedis

A Redis-compatible server built from scratch in Go — no external dependencies, just the standard library.

GoRedis speaks RESP2 over TCP, so any real Redis client (`redis-cli` included) can talk to it directly.

## Features

- **RESP2 protocol**, parsed from scratch — arrays, inline (telnet-style) commands, and pipelining all work over a single connection.
- **Thread-safe in-memory store** guarded by `sync.RWMutex`.
- **String commands**: `GET` `SET` `SETNX` `DEL` `EXISTS` `INCR` `DECR` `APPEND` `STRLEN` `MGET` `MSET`
- **Server commands**: `PING` `ECHO`
- **Zero dependencies** — `net`, `bufio`, `sync`, `log/slog`, `flag`, `testing` only.

## Quick Start

```bash
go build -o goredis ./cmd/server
./goredis --port 6379
```

Then, from another terminal:

```bash
redis-cli -p 6379 set foo bar
redis-cli -p 6379 get foo
redis-cli -p 6379 incr counter
redis-cli -p 6379 mset a 1 b 2
redis-cli -p 6379 mget a b missing
```

No `redis-cli`? Any RESP2 or raw TCP client works — GoRedis accepts both wire formats:

```
$ printf '*1\r\n$4\r\nPING\r\n' | nc localhost 6379
+PONG
```

## Architecture

One responsibility per package:

| Package | Responsibility |
|---|---|
| `internal/resp` | Wire format — `Value` types, `Encode()`, and a `Reader` that parses RESP arrays and inline commands |
| `internal/storage` | Thread-safe `map[string]string` store |
| `internal/commands` | Command name → handler dispatch table |
| `internal/server` | TCP accept loop, wires the above together per connection |
| `cmd/server` | Entrypoint (`--port` flag) |

## Testing

```bash
go test ./...
```

## Roadmap

This is Phase 1 (Foundation) of a multi-phase build. Not yet implemented:

Lists/Sets/Hashes + keyspace commands + multi-DB, TTL/active expiration, AOF/RDB persistence + config file, Pub/Sub + Transactions, Replication, observability (metrics/logging/slowlog), CLI client + benchmarks, Docker/CI.

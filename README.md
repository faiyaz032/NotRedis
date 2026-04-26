# NotRedis

NotRedis — a distributed key-value store built from scratch in Go to learn how systems like Redis work internally.

## Project Structure

```text
cmd/
	server/
		main.go              # Application entrypoint only
internal/
	protocol/
		command.go           # Command parsing and execution rules
	server/
		tcp_server.go        # TCP listener and connection lifecycle
	store/
		store.go             # Thread-safe in-memory key-value store
```

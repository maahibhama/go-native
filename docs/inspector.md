# Local runtime inspector

The `runtime/inspector` package exposes the portable diagnostics log and current
UI-tree snapshot as an opt-in, read-only HTTP service. It is intended for local
development tools, not for inclusion in release builds.

```go
service := inspector.New(appRuntime, "") // 127.0.0.1 on an ephemeral port
address, err := service.Start()
// Query http://<address>/v1/tree and /v1/logs.

ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
err = service.Stop(ctx)
```

Endpoints are JSON and send `Cache-Control: no-store`:

- `GET /healthz` reports service availability.
- `GET /v1/tree` returns a detached snapshot of the last rendered tree.
- `GET /v1/logs` returns the bounded structured runtime event history.

`Start` rejects wildcard and non-loopback addresses. There are no mutation,
event-dispatch, source-evaluation, or file endpoints. Callers must explicitly
start and stop the service; merely importing the package never opens a socket.
Applications should guard inspector startup with their own debug-build policy.

Handler IDs and text props may reveal application behavior or user-entered data.
Loopback binding reduces exposure but is not authentication. Do not port-forward
the service, bind it through a proxy, or ship it enabled in production.


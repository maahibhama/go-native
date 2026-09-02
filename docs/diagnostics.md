# Runtime diagnostics and UI-tree inspection

The `runtime` package contains platform-independent diagnostics foundations. They
are deliberately detached from UIKit, Android Views, and the binary mutation
protocol so the same data can feed a CLI, an in-app inspector, or test tooling.

## Structured event log

`runtime.Diagnostics` stores a thread-safe, bounded history of `LogEntry`
values. Entries use stable string kinds such as `runtime.started`,
`event.dispatched`, and `render.batch_applied`. Batch entries can include the
sequence and mutation count; event entries can include a handler ID. The buffer
keeps the newest 1,024 events so enabling diagnostics cannot cause unbounded
growth in a long-running application.

`Entries` always returns a copy, and `Clear` can reset a debugging session.
`Record` supplies a timestamp when its caller does not provide one.

## Tree snapshots

`runtime.SnapshotTree` turns a `ui.Node` tree into a detached `TreeSnapshot`.
The snapshot contains node IDs, readable primitive type names, props, and
children. It intentionally excludes Go callback functions and internal
identity flags. Mutating or replacing the live tree after capture does not
change an existing snapshot.

Snapshots can be marshaled directly to JSON:

```go
snapshot := runtime.SnapshotTree(root)
document, err := json.Marshal(snapshot)
```

The owner of a live tree must synchronize the call with tree replacement. A
future interactive inspector should request the snapshot through `Runtime`,
which already owns the required tree lock, rather than reading runtime state
directly.

## Integration boundary

Runtime integration needs only three hooks:

1. own a `Diagnostics` collector;
2. record lifecycle, event-dispatch, renderer-success, and renderer-error
   entries at their existing transition points;
3. expose a locked method that snapshots the last successfully rendered tree.

No native pointers, native views, or callbacks belong in these records. Native
timing remains in `TimingSample`; log entries may correlate to it through a
batch sequence ID.

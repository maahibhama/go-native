// Package runtime mounts portable UI trees, preserves identity, manages event
// handlers, computes mutations, and submits one ordered batch to a native renderer.
//
// Applications normally construct a Runtime with New or NewContext. Platform
// bridges implement Renderer and forward complete mutation batches to the native
// UI thread. Native code refers back to callbacks only through integer HandlerID
// values owned by EventRegistry.
package runtime

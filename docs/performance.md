# Performance measurements

This project makes no performance claims yet. Benchmarks establish repeatable baselines and expose regressions while the architecture is still small.

## Portable suite

Run:

```bash
make benchmark
```

The suite measures:

- construction and `Build` of a 1,000-child declarative tree;
- initial mutation generation for 1,000 nodes;
- unchanged, one-property, and 100-property reconciliation in a 1,000-node tree;
- binary serialization for one update, 100 updates, and 1,000 creates;
- in-process `HandlerID` registry dispatch.

These are Go-side microbenchmarks. Registry dispatch is not native event latency, and serialization is not native mutation application. Results must always be reported with CPU, Go version, sample count, and allocation data.

## Baseline: 2026-09-02

Environment: Apple M5 Pro, Darwin arm64, Go 1.26.5. Three samples per benchmark at `-benchtime=200ms`; the table reports the middle observed time and the deterministic allocation result.

| Benchmark | Time/op | Bytes/op | Allocs/op |
|---|---:|---:|---:|
| Declarative tree build, 1,000 children | 84.8 µs | 427,918 | 4,906 |
| Initial reconciliation, 1,000 children | 62.8 µs | 917,270 | 14 |
| Unchanged reconciliation, 1,000 children | 178 µs | 82,080 | 11 |
| One property update, 1,000 children | 179 µs | 82,192 | 12 |
| 100 property updates, 1,000 children | 182 µs | 114,608 | 19 |
| Serialize one update | 207 ns | 312 | 17 |
| Serialize 100 updates | 14.6 µs | 22,776 | 1,211 |
| Serialize 1,000-node initial batch | 291 µs | 390,203 | 24,027 |
| Event registry dispatch | 3.38 ns | 0 | 0 |

These numbers are baselines, not user-facing performance claims. They identify two immediate optimization targets: unchanged reconciliation currently pays for full child maps, and the initial serializer performs many small `binary.Write` operations. Changes should be justified with before/after multi-sample results rather than intuition.

## Native instrumentation

Mutation protocol version 2 assigns every runtime batch a sequence number. The iOS renderer measures decode plus UIKit mutation application with `CLOCK_MONOTONIC_RAW`; Android measures Java decode plus View mutation application with `System.nanoTime`. Both acknowledge the sequence back to Go after UI-thread application. The runtime retains bounded per-batch timestamps and can report:

- native decode/application duration;
- bridge submission to native completion, including UI-queue delay;
- button dispatch to native completion for event-triggered batches.

`Runtime.TimingSamples` exposes completed samples for benchmark harnesses. These hooks do not log in the timed path.

## Native measurements still required

The following need platform harnesses rather than approximation by Go benchmarks:

- automated, repeated Objective-C/UIKit and JNI/Android View samples across batch sizes;
- automated distributions for native button event to completed label update latency;
- cold startup, resident memory, and packaged binary contribution;
- frame pacing for larger update batches.

For end-to-end latency, timestamps need a shared monotonic clock or platform-native spans around both bridge directions. Measurements should separate queue delay, decode time, reconciliation, serialization, and native application. Debug logging inside the timed loop is prohibited because it would dominate small batches.

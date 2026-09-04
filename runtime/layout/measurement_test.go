package layout

import (
	"context"
	"testing"

	"github.com/go-native/go-native/ui"
)

type recordingBatchMeasurer struct {
	calls    int
	requests int
}

func (m *recordingBatchMeasurer) MeasureBatch(_ context.Context, requests []MeasurementRequest) ([]MeasurementResult, error) {
	m.calls++
	m.requests += len(requests)
	out := make([]MeasurementResult, len(requests))
	for i, r := range requests {
		out[i] = MeasurementResult{ID: r.ID, Size: ui.Size{Width: float32(len(r.Text) * 8), Height: 20}}
	}
	return out, nil
}

func TestLayoutMeasuredBatchesAndCachesIntrinsicSizes(t *testing.T) {
	root := ui.Column(ui.Text("one"), ui.Text("longer")).Gap(4).Build()
	native := &recordingBatchMeasurer{}
	cache := NewMeasurementCache()
	engine := Engine{}
	box, err := engine.LayoutMeasured(context.Background(), root, Constraints{MaxWidth: 200, MaxHeight: 400}, native, cache)
	if err != nil {
		t.Fatal(err)
	}
	if native.calls != 1 || native.requests != 2 {
		t.Fatalf("calls=%d requests=%d", native.calls, native.requests)
	}
	if box.Frame.Width != 48 || box.Frame.Height != 44 {
		t.Fatalf("frame=%#v", box.Frame)
	}
	if _, err = engine.LayoutMeasured(context.Background(), root, Constraints{MaxWidth: 200, MaxHeight: 400}, native, cache); err != nil {
		t.Fatal(err)
	}
	if native.calls != 1 {
		t.Fatalf("cached layout made %d native calls", native.calls)
	}
}

func TestLayoutMeasuredRejectsIncompleteResults(t *testing.T) {
	bad := batchMeasureFunc(func(context.Context, []MeasurementRequest) ([]MeasurementResult, error) { return nil, nil })
	_, err := (Engine{}).LayoutMeasured(context.Background(), ui.Text("x").Build(), Constraints{MaxWidth: 100, MaxHeight: 100}, bad, nil)
	if err == nil {
		t.Fatal("expected result validation error")
	}
}

type batchMeasureFunc func(context.Context, []MeasurementRequest) ([]MeasurementResult, error)

func (f batchMeasureFunc) MeasureBatch(ctx context.Context, r []MeasurementRequest) ([]MeasurementResult, error) {
	return f(ctx, r)
}

package layout

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-native/go-native/ui"
)

// MeasurementRequest describes one intrinsic native-control measurement.
// Requests contain values only; no Go pointers cross the native boundary.
type MeasurementRequest struct {
	ID          uint64
	NodeType    ui.NodeType
	Text        string
	ImageSource string
	Style       ui.Style
	Constraints Constraints
}

type MeasurementResult struct {
	ID   uint64
	Size ui.Size
	Err  string
}

// BatchMeasurer performs one coarse native call for a set of intrinsic sizes.
type BatchMeasurer interface {
	MeasureBatch(context.Context, []MeasurementRequest) ([]MeasurementResult, error)
}

type measurementKey struct {
	nodeType    ui.NodeType
	text, image string
	style       ui.Style
	constraints Constraints
}

// MeasurementCache deduplicates native measurements by content, style, and constraints.
type MeasurementCache struct {
	mu     sync.RWMutex
	values map[measurementKey]ui.Size
}

func NewMeasurementCache() *MeasurementCache {
	return &MeasurementCache{values: make(map[measurementKey]ui.Size)}
}

func (c *MeasurementCache) Get(key measurementKey) (ui.Size, bool) {
	if c == nil {
		return ui.Size{}, false
	}
	c.mu.RLock()
	value, ok := c.values[key]
	c.mu.RUnlock()
	return value, ok
}
func (c *MeasurementCache) set(key measurementKey, value ui.Size) {
	if c == nil {
		return
	}
	c.mu.Lock()
	c.values[key] = value
	c.mu.Unlock()
}
func (c *MeasurementCache) Clear() {
	if c == nil {
		return
	}
	c.mu.Lock()
	c.values = make(map[measurementKey]ui.Size)
	c.mu.Unlock()
}

// LayoutMeasured batches every uncached intrinsic leaf measurement before layout.
func (e Engine) LayoutMeasured(ctx context.Context, root *ui.Node, constraints Constraints, native BatchMeasurer, cache *MeasurementCache) (*Box, error) {
	if root == nil {
		return nil, nil
	}
	if native == nil {
		return nil, fmt.Errorf("layout: native batch measurer unavailable")
	}
	if cache == nil {
		cache = NewMeasurementCache()
	}
	constraints = normalize(constraints)
	requests := make([]MeasurementRequest, 0)
	keys := make(map[uint64]measurementKey)
	sizes := make(map[measurementKey]ui.Size)
	var nextID uint64
	collectMeasurements(root, constraints, constraints.MaxWidth, constraints.MaxHeight, cache, sizes, &requests, keys, &nextID)
	if len(requests) > 0 {
		results, err := native.MeasureBatch(ctx, requests)
		if err != nil {
			return nil, fmt.Errorf("layout: native measurement batch: %w", err)
		}
		if len(results) != len(requests) {
			return nil, fmt.Errorf("layout: expected %d measurement results, got %d", len(requests), len(results))
		}
		seen := make(map[uint64]bool, len(results))
		for _, result := range results {
			key, ok := keys[result.ID]
			if !ok || seen[result.ID] {
				return nil, fmt.Errorf("layout: invalid measurement result id %d", result.ID)
			}
			if result.Err != "" {
				return nil, fmt.Errorf("layout: measurement %d: %s", result.ID, result.Err)
			}
			seen[result.ID] = true
			sizes[key] = result.Size
			cache.set(key, result.Size)
		}
	}
	measurer := MeasureFunc(func(node *ui.Node, c Constraints) ui.Size {
		return sizes[makeMeasurementKey(node, c)]
	})
	e.Measurer = measurer
	return e.Layout(root, constraints), nil
}

func collectMeasurements(node *ui.Node, c Constraints, parentW, parentH float32, cache *MeasurementCache, sizes map[measurementKey]ui.Size, requests *[]MeasurementRequest, keys map[uint64]measurementKey, nextID *uint64) {
	s := node.Style.Layout
	w, hasW := resolve(s.Width, parentW)
	h, hasH := resolve(s.Height, parentH)
	maxW, maxH := maxZero(c.MaxWidth-s.Padding.Leading-s.Padding.Trailing), maxZero(c.MaxHeight-s.Padding.Top-s.Padding.Bottom)
	if hasW {
		maxW = maxZero(w - s.Padding.Leading - s.Padding.Trailing)
	}
	if hasH {
		maxH = maxZero(h - s.Padding.Top - s.Padding.Bottom)
	}
	childConstraints := Constraints{MaxWidth: maxW, MaxHeight: maxH}
	if len(node.Children) == 0 {
		key := makeMeasurementKey(node, childConstraints)
		if value, ok := cache.Get(key); ok {
			sizes[key] = value
			return
		}
		*nextID++
		id := *nextID
		*requests = append(*requests, MeasurementRequest{ID: id, NodeType: node.Type, Text: node.Props.Text, ImageSource: node.Props.ImageSource, Style: node.Style, Constraints: childConstraints})
		keys[id] = key
		return
	}
	for _, child := range node.Children {
		collectMeasurements(child, childConstraints, maxW, maxH, cache, sizes, requests, keys, nextID)
	}
}

func makeMeasurementKey(node *ui.Node, c Constraints) measurementKey {
	return measurementKey{nodeType: node.Type, text: node.Props.Text, image: node.Props.ImageSource, style: node.Style, constraints: c}
}

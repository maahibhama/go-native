package ui

import (
	"fmt"
	"sync"
)

type Orientation uint8

const (
	OrientationUnknown Orientation = iota
	OrientationPortrait
	OrientationLandscape
)

type LayoutDirection uint8

const (
	DirectionLTR LayoutDirection = iota
	DirectionRTL
)

type Contrast uint8

const (
	ContrastNormal Contrast = iota
	ContrastHigh
)

type PointerCapabilities struct{ Touch, Mouse, Hover bool }

type MediaQuery struct {
	Viewport       Size
	Scale          float32
	TextScale      float32
	Orientation    Orientation
	SafeAreaInsets EdgeInsets
	KeyboardInsets EdgeInsets
	Contrast       Contrast
	ReducedMotion  bool
	Pointers       PointerCapabilities
}

// LayoutRect is a Go-computed rectangle expressed in platform-independent
// logical points. It is value-only so it can safely cross a native boundary.
type LayoutRect struct{ X, Y, Width, Height float32 }

// LayoutFrame associates deterministic geometry with a stable node identity.
type LayoutFrame struct {
	NodeID NodeID
	Rect   LayoutRect
}

type LifecycleState uint8

const (
	LifecycleCreated LifecycleState = iota
	LifecycleForeground
	LifecycleActive
	LifecycleInactive
	LifecycleBackground
	LifecycleMemoryPressure
	LifecycleDestroyed
)

type DependencyKey[T any] struct{ Name string }

type Dependencies struct {
	mu     sync.RWMutex
	values map[string]any
}

func NewDependencies() *Dependencies { return &Dependencies{values: make(map[string]any)} }

func Provide[T any](dependencies *Dependencies, key DependencyKey[T], value T) {
	if dependencies == nil {
		return
	}
	dependencies.mu.Lock()
	dependencies.values[key.Name] = value
	dependencies.mu.Unlock()
}

func Resolve[T any](dependencies *Dependencies, key DependencyKey[T]) (T, bool) {
	var zero T
	if dependencies == nil {
		return zero, false
	}
	dependencies.mu.RLock()
	value, ok := dependencies.values[key.Name]
	dependencies.mu.RUnlock()
	if !ok {
		return zero, false
	}
	typed, ok := value.(T)
	return typed, ok
}

type Environment struct {
	Theme        Theme
	Locale       string
	Direction    LayoutDirection
	MediaQuery   MediaQuery
	Lifecycle    LifecycleState
	Dependencies *Dependencies
	Focus        *FocusManager
}

func DefaultEnvironment() Environment {
	return Environment{Theme: DefaultTheme(), Locale: "en", Direction: DirectionLTR, Lifecycle: LifecycleCreated, Dependencies: NewDependencies(), Focus: NewFocusManager(nil)}
}

// BuildContext is immutable application context scoped to a mounted component.
type BuildContext struct {
	Environment Environment
	Path        string
	hooks       *HookRegistry
}

func NewBuildContext(environment Environment) BuildContext {
	return BuildContext{Environment: environment, Path: "root"}
}
func (c BuildContext) Child(key string) BuildContext {
	if key == "" {
		key = "child"
	}
	c.Path = fmt.Sprintf("%s/%s", c.Path, key)
	return c
}
func (c BuildContext) WithEnvironment(environment Environment) BuildContext {
	c.Environment = environment
	return c
}

// FocusManager returns the application-scoped portable focus owner.
func (c BuildContext) FocusManager() *FocusManager { return c.Environment.Focus }

// WithHooks attaches mounted hook storage. It is used by runtimes and custom
// render hosts; application components normally receive an already configured context.
func (c BuildContext) WithHooks(hooks *HookRegistry) BuildContext { c.hooks = hooks; return c }

// ContextComponent is the production component contract. The legacy Component
// contract remains supported during v0 migration.
type ContextComponent interface{ BuildContext(BuildContext) *Node }

func BuildWithContext(component Component, context BuildContext) *Node {
	if component == nil {
		return nil
	}
	if contextual, ok := component.(ContextComponent); ok {
		return contextual.BuildContext(context)
	}
	return component.Build()
}

# ADR 0005: Declarative navigation and modal contracts

## Status

Accepted as an API foundation. Native controller mounting is deferred.

## Context

Navigation controllers, Android activities and fragments, and modal presentation
APIs have incompatible ownership models. Exposing any one platform abstraction in
Go would prevent a component tree from behaving consistently on both mobile
platforms. Treating push and dismiss as imperative commands would also make native
state difficult to reconcile after an interrupted transition or lifecycle change.

## Decision

Navigation is declared as a complete ordered stack of `Route` values. The first
route is the root and the last is visible. Every route has a stable string key,
an optional title, and component content. A later reconciler compares route keys:
appending means push, removing the last route means pop, and replacing or reordering
routes rebuilds the affected native controller suffix.

A modal is optional metadata attached to base content. Its stable key defines
presentation identity, while its content, automatic/sheet/fullscreen style, and
dismissibility describe desired state. `OnDismiss` is Go-owned and runs only after
the platform reports that dismissal completed, including an allowed interactive
dismissal. Application state must then remove the modal intent.

The current `NavigationStack` fallback builds the last non-nil route, and
`PresentModal` builds its base. This keeps decorated applications renderable before
native controller integration without pretending that push or presentation has
occurred.

## Lifecycle and ownership

The runtime owns route and dismissal handlers. Native code owns controllers only
while their keys remain declared. Backgrounding does not change the desired stack.
Stopping the runtime releases callbacks and dismisses owned presentations without
calling application dismissal handlers.

Native interactive back navigation must be reported to Go and reflected in the
next declarative stack. If state restores the removed route, the renderer presents
it again. The same rule applies when a dismissible modal is swiped away.

## Consequences

The public contract is typed and platform-neutral, but this ADR does not complete
native navigation or modal support. Completion still requires node/wire encoding,
integer callback registration, UIKit and Android mounting, transition
serialization, restoration, accessibility checks, and device tests.

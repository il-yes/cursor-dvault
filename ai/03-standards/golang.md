
# Go Engineering Standards

## Purpose

This document defines the Go development standards for Ankhora.

The objective is to produce Go code that is:

- readable
- testable
- maintainable
- explicit
- aligned with DDD architecture

---

# Core Principle

Go code should favor:

```

Simplicity

*

Explicit Design

*

Clear Ownership

```

Avoid unnecessary abstraction.

---

# Project Structure

Go packages follow bounded context ownership.

Example:

```

internal/

├── vault/
│   ├── domain/
│   ├── application/
│   ├── infrastructure/
│   └── interface/

├── tracecore/
│   ├── domain/
│   ├── application/
│   └── infrastructure/

├── federation/

└── identity/

```

---

# Package Rules

Packages should have a single responsibility.

Good:

```

vault/domain/asset

```

Contains asset business rules.

---

Bad:

```

utils/
helpers/
common/
misc/

````

Avoid generic dumping grounds.

---

# Naming Standards

Go names should be:

- clear
- concise
- meaningful

Prefer:

```go
VaultRepository
AssetService
CreateWorkspace
````

Avoid:

```go
VaultManager
DataProcessor
HelperService
```

---

# Interfaces

Interfaces should belong to the consumer.

Example:

Application defines what it needs:

```go
type AssetRepository interface {
    Save(ctx context.Context, asset Asset) error
    Find(ctx context.Context, id string) (*Asset, error)
}
```

Infrastructure implements:

```go
type PostgresAssetRepository struct {}
```

---

Avoid:

Creating interfaces before a need exists.

Do not create:

```go
IAssetService
```

Go does not use interface prefixes.

---

# Dependency Injection

Dependencies are explicit.

Prefer:

```go
type CreateAssetHandler struct {
    repo AssetRepository
    events EventPublisher
}
```

Avoid:

* global variables
* service locators
* hidden dependencies

---

# Context Usage

Every operation involving:

* database
* network
* filesystem
* external service

must accept:

```go
context.Context
```

Example:

```go
func (r Repository) Save(
    ctx context.Context,
    asset Asset,
) error
```

---

Context rules:

Use context for:

* cancellation
* deadlines
* request scope values

Do not store business data in context.

---

# Error Handling

Errors must preserve meaning.

Bad:

```go
return err
```

when context is lost.

---

Good:

```go
return fmt.Errorf(
    "failed to save asset %w",
    err,
)
```

---

Errors should answer:

```
What failed?

Where?

Why?
```

---

# Domain Errors

Domain errors represent business failures.

Example:

```go
ErrInvalidAssetState

ErrUnauthorizedSharing

ErrInvalidTransition
```

---

Do not use infrastructure errors inside domain logic.

---

# Struct Design

Prefer explicit structures.

Example:

```go
type Asset struct {
    ID string
    CID string
    OwnerID string
}
```

---

Avoid:

* huge structs
* unrelated fields
* hidden state

---

# Constructors

Use constructors when validation is required.

Example:

```go
func NewAsset(
    id string,
    owner string,
) (*Asset, error)
```

The constructor protects invariants.

---

# Methods

Methods should express business behavior.

Good:

```go
asset.Share(groupID)
```

Bad:

```go
asset.SetShared(true)
```

---

# Concurrency Standards

Go concurrency must be intentional.

Before using goroutines ask:

```
Why does this need concurrency?
```

---

Rules:

* always manage goroutine lifecycle
* use context cancellation
* avoid leaks
* protect shared state

---

Example:

Good:

```go
go worker.Run(ctx)
```

because cancellation exists.

---

# Channels

Channels communicate ownership transfer.

Avoid using channels as a replacement for normal function calls.

---

Use channels for:

* asynchronous processing
* pipelines
* event distribution

---

# JSON Standards

External data must be validated.

Do not trust:

* API payloads
* federation messages
* imported files

---

Domain objects should not directly depend on JSON format.

---

# Database Rules

Database code belongs in infrastructure.

Repositories expose domain needs.

Avoid:

```go
domain.Asset

imports postgres package
```

---

# Logging Standards

Logs should provide:

* operation
* identifier
* context
* error

Example:

```
asset reconstruction failed

asset_id=123

reason=missing metadata
```

---

Avoid:

```
error occurred
```

---

# Testing Standards

Every important behavior requires tests.

Prioritize:

1. Domain tests
2. Application tests
3. Integration tests

---

Tests should verify behavior.

Good:

```go
TestCannotShareAssetWithoutPermission
```

Bad:

```go
TestShareFunction
```

---

# Go Tooling

Code should pass:

```
go fmt

go vet

go test

go test -race
```

---

# AI Code Generation Rules

When generating Go code, AI must:

* follow existing package ownership
* avoid unnecessary abstractions
* use context.Context
* return meaningful errors
* create tests with features
* preserve DDD boundaries

---

# Forbidden Patterns

Avoid:

* global state
* god structs
* giant services
* circular dependencies
* hidden side effects
* unnecessary interfaces
* premature optimization

---

# Final Principle

Good Go code is boring.

It is easy to read, easy to test, and obvious to maintain.

The best Go architecture lets future engineers understand the system without fighting the code.

````

---



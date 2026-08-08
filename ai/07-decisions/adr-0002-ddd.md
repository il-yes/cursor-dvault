
# ADR-0002: Domain-Driven Design Architecture

## Status

Accepted

## Date

2026

---

# Decision

Ankhora adopts Domain-Driven Design (DDD) as the primary architectural approach.

The system is organized around:

- bounded contexts
- domain ownership
- explicit business capabilities
- clear dependency boundaries

---

# Context

Ankhora is not a single-purpose application.

It is a platform containing multiple complex domains:

- secure storage
- collaboration
- lifecycle management
- identity
- federation
- commercial capabilities
- industry applications

A traditional layered architecture creates problems:

```

All Models

↓

All Services

↓

All Controllers

```

This causes:

- unclear ownership
- increasing coupling
- difficult evolution
- accidental dependency

---

# Decision Details

Each major capability is represented as a bounded context.

Current contexts:

```

Identity

Vault

C3

TraceCore

Federation

Subscription

Domain Applications

```

---

# Bounded Context Principle

A bounded context owns:

- its domain model
- its business rules
- its language
- its lifecycle

Example:

Vault owns:

```

Assets

Encryption

Keys

Ownership

```

C3 does not.

---

C3 owns:

```

Channels

Threads

Collaboration

Trust Groups

```

Vault does not.

---

TraceCore owns:

```

Commits

Events

History

Validation

```

Domain applications do not.

---

# Dependency Rules

The architecture follows:

```

Interface Layer

```
    ↓
```

Application Layer

```
    ↓
```

Domain Layer

```
    ↓
```

Infrastructure

```

---

# Domain Layer Rules

The domain layer must contain:

- entities
- value objects
- aggregates
- domain events
- business rules

The domain layer must not contain:

- HTTP
- databases
- external APIs
- framework dependencies

---

# Application Layer Rules

The application layer coordinates:

- use cases
- transactions
- workflows
- event publishing

It does not contain complex business rules.

---

# Infrastructure Layer Rules

Infrastructure provides:

- persistence
- external communication
- filesystem
- networking
- third-party integrations

---

# Context Communication

Bounded contexts communicate through:

- contracts
- events
- interfaces
- explicit APIs

Avoid:

```

Direct database access between contexts

```

---

# Domain Events

Important state changes should be represented as events.

Example:

```

AssetCreated

ThreadOpened

CommitValidated

TrustEstablished

```

Events communicate facts.

---

# Alternatives Considered

## Monolithic Layered Architecture

Rejected.

Reason:

Does not scale with domain complexity.

---

## Microservices Everywhere

Rejected.

Reason:

Introduces operational complexity too early.

Ankhora prefers modular boundaries first.

---

## Shared Domain Model

Rejected.

Reason:

Different contexts have different meanings.

---

# Consequences

Positive:

- clear ownership
- easier evolution
- better AI understanding
- reduced coupling
- independent scaling possibilities

Negative:

- more concepts to maintain
- explicit communication required
- initial complexity

---

# AI Rules

When creating new code, AI must ask:

```

Which bounded context owns this?

Does this belong in the domain?

Is this crossing a boundary?

Should this be an event?

```

---

AI must avoid:

- shared global models
- god services
- cross-context database access
- moving logic into convenient locations

---

# Final Principle

Architecture follows ownership.

The system should evolve by strengthening boundaries, not removing them.

DDD allows Ankhora to grow into a platform without becoming a tangled application.
```

---



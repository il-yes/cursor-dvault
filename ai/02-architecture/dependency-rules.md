# Dependency Rules

## Purpose

This document defines the dependency rules of the Ankhora architecture.

The purpose is to maintain:

- clean boundaries
- domain independence
- testability
- maintainability
- architectural integrity

Dependencies are architectural decisions.

Every dependency must have a reason.

---

# Core Dependency Principle

Dependencies must point toward business meaning.

Technical details should depend on business rules.

Business rules should never depend on technical details.

The direction is:
Infrastructure

    ↓

Application

    ↓

Domain


The domain is the center of the system.

---

# Architectural Layers

Each bounded context follows a layered architecture.

internal/context/
domain/

application/

infrastructure/

interfaces/

---

# Domain Layer

## Responsibility

The domain contains the core business logic.

It defines:

- entities
- value objects
- aggregates
- domain services
- domain events
- business rules

---

# Domain Rules

The domain layer:

MUST:

- remain independent
- contain business behavior
- protect invariants
- express domain concepts

---

The domain layer MUST NOT:

- import databases
- import HTTP frameworks
- import message brokers
- know infrastructure details
- depend on external services

---

Example:

Correct:
domain.Asset

contains:

Encrypt()
ValidateOwnership()
ChangeStatus()


---

Incorrect:


domain.Asset

imports:

postgres
redis
http.Client


---

# Application Layer

## Responsibility

The application layer coordinates use cases.

It manages:

- workflows
- transactions
- orchestration
- permissions checks
- domain interaction

---

The application layer:

MUST:

- call domain behavior
- coordinate multiple domain objects
- publish events
- use interfaces

---

The application layer MUST NOT:

- contain low-level infrastructure logic
- contain database queries
- become a second domain layer

---

Example:

Correct:
CreateAssetUseCase

validate request
create Asset aggregate
save through repository interface
publish AssetCreated event

---

Incorrect:


CreateAssetUseCase

contains:

SQL queries
encryption implementation
HTTP calls


---

# Infrastructure Layer

## Responsibility

Infrastructure provides technical implementations.

Examples:

- PostgreSQL repositories
- MongoDB storage
- Redis cache
- message brokers
- external APIs
- cryptography providers

---

Infrastructure:

MUST:

- implement domain/application interfaces
- hide technical details
- provide adapters

---

Infrastructure MUST NOT:

- define business rules
- modify domain decisions
- become the source of truth

---

# Interface Layer

## Responsibility

Interfaces translate external communication into application actions.

Examples:

- HTTP handlers
- WebSocket handlers
- CLI commands
- API controllers

---

Interface layer:

MUST:

- validate input format
- authenticate requests
- call use cases
- return responses

---

Interface layer MUST NOT:

- contain business logic
- access repositories directly
- manipulate domain state

---

Incorrect:
HTTP Handler

    |
    v

Database


---

Correct:


HTTP Handler

    |
    v

UseCase

    |
    v

Repository Interface

    |
    v

Database Implementation

---

# Repository Rules

Repositories represent domain persistence needs.

The domain defines the interface.

Infrastructure implements it.

Example:

Domain:

type AssetRepository interface {
    Save(ctx context.Context, asset Asset) error
    Find(ctx context.Context, id ID) (*Asset, error)
}

Infrastructure:
type PostgresAssetRepository struct {
    db *sql.DB
}
The domain should never know:

PostgreSQL
MongoDB
filesystem
cloud storage
Cross-Bounded Context Rules

Bounded contexts communicate through contracts.

Allowed:

APIs
domain events
messages
shared protocols

Forbidden:

importing another context's internal packages
sharing database tables
sharing domain entities

Incorrect:

tracecore/domain

imports

vault/domain.Asset

Correct:

TraceCore

receives:

AssetReferenceDTO

or

AssetCreatedEvent
Model Ownership Rules

Every context owns its models.

Example:

Vault:

Asset

means:

encrypted protected object.

TraceCore:

AssetReference

means:

lifecycle tracked object.

They may represent related concepts.

They are not the same object.

Dependency Inversion

High-level policies should not depend on low-level details.

Example:

Application:

Storage interface

Infrastructure:

S3 implementation

The application does not know S3 exists.

Event Dependencies

Events reduce coupling.

Prefer:

Context A

publishes event

        ↓

Context B reacts

over:

Context A

directly modifies

Context B
Circular Dependency Prevention

Forbidden:

Vault
 |
 v
TraceCore
 |
 v
Vault

If two contexts need communication:

introduce:

contracts
events
interfaces
shared protocols
Dependency Review Checklist

Before adding a dependency ask:

Which layer owns this responsibility?
Is this dependency necessary?
Does it point toward the domain?
Does it create coupling?
Can an interface replace it?
Does another context already own this concept?
AI Implementation Rule

When generating code, AI must verify:

Is the layer correct?
Is the dependency direction correct?
Does this violate bounded context ownership?
Is business logic placed correctly?
Are interfaces used appropriately?


Final Principle

A clean architecture is not defined by folders.

It is defined by dependency direction.

The system remains maintainable when business meaning stays independent from technical details.

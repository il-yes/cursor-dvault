
# Ankhora Current Engineering State

## Purpose

This document provides the current state of the Ankhora platform.

AI assistants should read this file before modifying the codebase.

This document is a snapshot, not a historical document.

For architectural decisions, refer to:

- ai/07-decisions
- ai/02-architecture

---

# Platform Overview

Ankhora is a secure data sovereignty platform composed of multiple bounded contexts.

The platform combines:

- encrypted vault storage
- collaboration
- identity management
- lifecycle tracking
- compliance and audit capabilities

The architecture follows:

- Domain Driven Design
- Clean Architecture
- Event Driven Architecture
- Zero Trust principles

---

# Current Architecture State

## Core Layers

```

Domain Applications

```
    |
```

Ankhora Vault

```
    |
```

TraceCore

```
    |
```

C3 Collaboration

```
    |
```

Federation

```

Each bounded context owns its business rules.

Cross-context communication happens through:

- domain events
- integration events
- explicit interfaces

---

# Bounded Context Status

## Vault

Status:

ACTIVE

Responsibilities:

- encrypted object storage
- asset management
- encryption lifecycle
- sharing foundations

Current capabilities:

- vault creation
- encrypted payload management
- asset storage
- attachment handling

---

## Vault Desktop Engine

Status:

ACTIVE

Responsibilities:

- local vault operation
- desktop application behavior
- offline-first interaction

Technology:

- Go backend
- Wails desktop application
- React frontend

---

## Vault Cloud Service

Status:

ACTIVE

Responsibilities:

- remote vault operations
- cloud synchronization
- distributed access

---

## C3 Collaboration

Status:

ACTIVE DEVELOPMENT

Responsibilities:

- workspaces
- channels
- threads
- trust groups
- collaboration lifecycle

Current implemented concepts:

- Workspace
- Channel
- Thread
- Asset
- ShareEntry
- TrustGroup
- Federation foundation

---

## Identity

Status:

ACTIVE

Responsibilities:

- user identity
- authentication
- trust relationships
- cryptographic identity

---

## Federation

Status:

FOUNDATION IMPLEMENTED

Responsibilities:

- remote vault trust
- message validation
- cryptographic verification
- remote communication

---

## TraceCore

Status:

ACTIVE DEVELOPMENT

Responsibilities:

- lifecycle history
- commits
- validation
- compliance
- audit trail

---

# Current Development Focus

Main focus:

```

C3 Collaboration completion

```

Current bounded context:

```

internal/channel

```

Current work:

- Channel use cases
- Channel lifecycle
- Event integration

---

# Engineering Standards

All new code must follow:

- DDD boundaries
- dependency inversion
- thin interfaces
- explicit ownership
- domain events for business changes

AI must not introduce:

- shortcuts
- cross-context coupling
- duplicated business rules
- infrastructure logic inside domain

---

# Current AI Collaboration Mode

AI assistants should behave as:

```

Senior developer
+
DDD reviewer
+
Architecture guardian

```

Before coding:

1. Understand ownership.
2. Check existing patterns.
3. Propose design.
4. Implement.
5. Test.
6. Review against architecture.

---

# Last Updated

Update this document when:

- a major bounded context changes
- architecture evolves
- a milestone is completed


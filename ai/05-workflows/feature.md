
# Feature Development Workflow

## Purpose

This document defines the standard workflow for implementing new features in Ankhora.

The goal is to ensure that every feature respects:

- Domain boundaries
- Security principles
- DDD architecture
- Existing design decisions
- Long-term maintainability

---

# Core Principle

A feature is not only code.

A feature is:

```

Business Need

*

Domain Decision

*

Architecture Integration

*

Implementation

*

Validation

```

AI must understand the system before modifying it.

---

# Feature Development Process

Every feature follows:

```

1. Understand

2. Analyze Ownership

3. Design

4. Implement

5. Validate

6. Document

```

---

# Step 1 — Understand The Request

Before writing code, AI must identify:

- What problem is being solved?
- Who is the user?
- Which bounded context owns this?
- What existing behavior is affected?

---

AI should read:

```

00-vision/

01-principles/

02-architecture/

04-contexts/

```

before making architectural decisions.

---

# Step 2 — Identify Ownership

Every feature must have an owner.

Ask:

```

Which bounded context owns this capability?

```

Examples:

## Encryption Feature

Owner:

```

Vault

```

Not:

```

C3
TraceCore

```

---

## Collaboration Feature

Owner:

```

C3

```

Not:

```

Vault

```

---

## Audit Feature

Owner:

```

TraceCore

```

Not:

```

Domain Application

```

---

# Step 3 — Analyze Impact

Before implementation, identify:

## Domain Impact

Questions:

- New entities?
- New aggregates?
- New domain rules?
- New events?

---

## Architecture Impact

Questions:

- New bounded context?
- New dependency?
- New interface?
- New workflow?

---

## Security Impact

Questions:

- Does data protection change?
- Does trust change?
- Does authorization change?
- Does encryption change?

---

# Step 4 — Design Before Coding

The AI should propose:

## Domain Model

Example:

```

Entity

Value Objects

Aggregates

Events

```

---

## Application Flow

Example:

```

Request

↓

Use Case

↓

Domain Operation

↓

Event

↓

Infrastructure

```

---

## Data Flow

Example:

```

User

↓

Desktop

↓

Vault Engine

↓

Synchronization

↓

Cloud

```

---

# Step 5 — Implementation Rules

AI must respect:

## Domain Layer

Contains:

- business rules
- entities
- value objects
- domain events

Must not contain:

- HTTP
- database
- frameworks

---

## Application Layer

Contains:

- use cases
- orchestration
- transactions

---

## Infrastructure Layer

Contains:

- databases
- APIs
- external services

---

## Interface Layer

Contains:

- HTTP handlers
- websocket handlers
- UI adapters

---

# Step 6 — Testing Requirements

Every feature requires:

## Domain Tests

Validate:

- business rules
- invariants
- state transitions

---

## Application Tests

Validate:

- use case behavior
- orchestration

---

## Integration Tests

Validate:

- repositories
- APIs
- external systems

---

# Step 7 — Documentation Update

After implementation, update:

Possible documents:

```

04-contexts/

02-architecture/

07-decisions/

```

depending on impact.

---

# Feature Review Checklist

Before completion:

## Architecture

- [ ] Correct bounded context
- [ ] Dependencies respected
- [ ] No responsibility leakage

---

## Security

- [ ] Ownership preserved
- [ ] Keys protected
- [ ] Permissions validated

---

## DDD

- [ ] Domain logic isolated
- [ ] Aggregates respected
- [ ] Events meaningful

---

## Quality

- [ ] Tests added
- [ ] Errors handled
- [ ] Documentation updated

---

# AI Behavior Rules

When asked to implement a feature, AI must:

1. Explain understanding first
2. Identify affected contexts
3. Propose design
4. Wait for validation if architecture changes
5. Implement incrementally
6. Review the result

---

# Forbidden Behavior

AI must not:

- create code before understanding ownership
- bypass bounded contexts
- duplicate existing logic
- modify architecture casually
- introduce dependencies without justification

---

# Final Principle

A feature is successful when it adds capability without damaging the architecture.

The objective is not the fastest code.

The objective is the fastest sustainable evolution.
```

---



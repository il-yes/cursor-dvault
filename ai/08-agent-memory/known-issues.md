Excellent. `known-issues.md` is the **engineering memory of pain**.

This file is extremely important because AI agents have a natural tendency to "solve" problems they already know nothing about. Without memory, the agent may:

* reintroduce an old bug,
* choose a previously rejected approach,
* break a known constraint,
* waste tokens rediscovering architecture decisions.

The purpose is:

```text
current-state.md
    = Where are we?

active-work.md
    = What are we building?

known-issues.md
    = What must we remember?
```

Create:

```text
ai/
└── 08-agent-memory
    └── known-issues.md
```

Content:

```md
# Ankhora Known Issues

## Purpose

This document contains accumulated engineering knowledge about known problems, constraints, previous discoveries, and important implementation lessons.

AI assistants should read this file before proposing solutions.

A known issue is not necessarily a bug.

It can also represent:

- architectural constraints
- rejected approaches
- technical debt
- future improvements
- operational risks

---

# Usage Rules

Before changing existing systems:

1. Check whether the problem already exists here.
2. Understand previous decisions.
3. Avoid reintroducing rejected solutions.
4. Prefer incremental improvements.

---

# Architecture Constraints

---

## Bounded Context Isolation

Status:

ACTIVE RULE

Problem:

Cross-context coupling creates architectural drift.

Rule:

Bounded contexts must communicate through:

- interfaces
- domain events
- integration events

Never:

- import another context's domain objects
- access another context's repository directly
- share database models

Example:

Wrong:

```

Channel

imports

TraceCore.Commit

```

Correct:

```

Channel

emits

ChannelArchived

↓

TraceCore consumes event

```

---

# Event Architecture

---

## Domain Event vs Integration Event

Status:

IMPORTANT

Rule:

Not every domain event is immediately an integration event.

Example:

Inside C3:

```

ChannelArchived

```

is a domain event.

Later:

```

ChannelArchivedIntegrationEvent

```

may exist for TraceCore or Federation.

Do not create cross-context events prematurely.

---

## Event Publication Reliability

Status:

KNOWN FUTURE IMPROVEMENT

Current pattern:

```

Repository Update

↓

Publish Event

```

Risk:

If event publication fails after persistence:

```

State changed

but

Event missing

```

Future solution:

Transactional outbox pattern.

Possible implementation:

```

Aggregate Change

↓

Database Transaction

↓

Outbox Table

↓

Event Dispatcher

↓

Consumers

```

Current decision:

Not required until distributed workflows become critical.

---

# Vault Payload Evolution

Status:

IMPORTANT

Problem:

Vault data structures evolve over time.

Known risk:

JSON reconstruction failures caused by schema mismatch.

Example:

```

cannot unmarshal object into Go struct field

````

Cause:

Persisted data structure differs from current Go model.

Rules:

Before changing persisted structures:

- define migration strategy
- maintain backward compatibility
- update reconstruction tests

Never silently break stored vault data.

---

# Serialization Issues

Status:

KNOWN

Problem:

Go JSON serialization is sensitive to:

- struct changes
- field type changes
- array/object mismatches

Example:

Changing:

```go
[]Slot
````

into:

```go
Slot
```

breaks existing payloads.

Before modifying persisted models:

Verify:

* old data compatibility
* migration needs
* reconstruction tests

---

# Repository Patterns

Status:

ACTIVE RULE

Repositories are abstractions owned by their bounded context.

Rules:

Repositories must:

* hide persistence details
* expose domain operations
* return meaningful errors

Avoid:

* database logic in use cases
* direct SQL from application layer
* infrastructure leakage

---

# Go Architecture Issues

---

## Dependency Direction

Status:

ACTIVE RULE

Required dependency flow:

```
Interface

↓

Application

↓

Domain

↓

Infrastructure implementation
```

Domain must never import:

* HTTP
* database
* framework code

---

## Context Usage

Status:

ACTIVE RULE

Go methods interacting with:

* repositories
* external systems
* event buses

should receive:

```go
context.Context
```

---

# Testing Lessons

---

## Tests Must Protect Behavior

Avoid testing implementation details.

Prefer:

Testing:

```
Channel cannot be archived when revoked
```

over:

```
Archive() line 42 executed
```

---

## Domain Tests First

For business rules:

Order:

```
Domain behavior

↓

Application workflow

↓

Infrastructure integration
```

---

# C3 Known Issues

---

## Channel Lifecycle Semantics

Status:

UNDER DESIGN

Issue:

Multiple closing states exist:

```
revoked
archived
```

Must maintain clear semantics.

Current definition:

```
Revoked

=
access/trust removal


Archived

=
collaboration paused while preserving history
```

Do not merge these concepts.

---

## TraceCore Integration

Status:

FUTURE

Issue:

C3 owns collaboration state.

TraceCore owns:

* lifecycle history
* audit
* compliance evidence

Future integration should use events.

Avoid:

```
C3 directly calls TraceCore
```

Preferred:

```
C3 Event

↓

TraceCore Handler
```

---

# Federation Known Issues

---

## Remote Trust Management

Status:

FOUNDATION ONLY

Known future topics:

* retry policy
* delivery acknowledgements
* remote vault lifecycle
* synchronization conflicts

Do not implement incomplete federation behavior without a defined protocol.

---

# WebSocket / Realtime

Status:

KNOWN TECHNICAL AREA

Previous issues:

* connection closing unexpectedly
* abnormal closure 1006
* heartbeat handling

Important:

Realtime reliability requires:

* heartbeat
* reconnect strategy
* connection lifecycle management

Do not assume a successful websocket handshake means a stable connection.

---

# Security Constraints

---

## Encryption Boundaries

Status:

ACTIVE RULE

Sensitive data must remain encrypted outside trusted boundaries.

Never:

* log decrypted content
* expose encryption keys
* store secrets in configuration files

---

## Trust Model

Status:

ACTIVE RULE

Access is based on:

* identity
* cryptographic trust
* authorization rules

Never assume:

"same user" = "trusted access"

---

# Rejected Approaches

---

## Direct Database Access From Use Cases

Rejected because:

* breaks DDD boundaries
* couples application to infrastructure

---

## Generic Utility Packages

Rejected because:

* hide ownership
* create unclear dependencies

Prefer:

small explicit domain services.

---

## Premature Abstractions

Rejected because:

* increase complexity
* slow evolution

Prefer:

simple designs that can evolve.

---

# AI Reminder

When proposing a solution:

Ask:

1. Has this problem already been solved?
2. Does this violate a known rule?
3. Is there historical context?
4. Is there a simpler approach?

The goal is not only to fix today's problem.

The goal is to preserve Ankhora's engineering memory.

---

## ChannelCreated Event Missing WorkspaceID

Status:

OPEN

Discovered:

2026-08-08 — ArchiveChannel session review

Issue:

`CreateChannelUsecase` publishes `ChannelCreated` without populating `WorkspaceID`.

The struct has the field. The use case does not set it.

`ArchiveChannelUsecase` correctly populates `WorkspaceID` in `ChannelArchived`.

Impact:

Any future consumer of `ChannelCreated` will receive an empty `WorkspaceID`.

Fix:

Add `WorkspaceID: channel.WorkspaceID` to the event in `create_usecases.go`.

---

## Missing Domain-Level Aggregate Tests

Status:

OPEN

Discovered:

2026-08-08 — ArchiveChannel session review

Issue:

`Channel.Archive()` is tested indirectly through application tests.

No isolated domain test exists for aggregate behavior.

Domain tests should verify:

- state transitions
- timestamp updates
- invariant enforcement

independently from use case orchestration.

Recommendation:

Create `internal/channel/domain/aggregate_test.go`.

---

## Error Strings Declared as var

Status:

OPEN

Discovered:

2026-08-08 — ArchiveChannel session review

Issue:

Domain error messages in `errors.go` use `var` instead of `const`.

These are immutable string constants used with `errors.New()`.

`var` allows accidental mutation.

Recommendation:

Change to `const` block.


## Architectural Implementation Gaps

### Client / Backend / Cloud Integration

**Status:** OPEN — architecture approved, implementation pending.

The following components are required before the encrypted Desktop → Backend → Cloud workflow is complete:

| Gap                                    | Priority | Status |
| -------------------------------------- | -------- | ------ |
| Device domain model                    | HIGH     | OPEN   |
| Multi-device TrustGroup key envelopes  | HIGH     | OPEN   |
| Versioned WrappedKEK model             | HIGH     | OPEN   |
| TrustGroup KEK rotation use case       | HIGH     | OPEN   |
| Eager DEK re-wrapping                  | HIGH     | OPEN   |
| TrustGroup-level ShareEntry workflow   | HIGH     | OPEN   |
| Desktop cryptographic engine           | HIGH     | OPEN   |
| Desktop/Wails crypto bridge            | HIGH     | OPEN   |
| Desktop → Backend command contracts    | HIGH     | OPEN   |
| Backend → Cloud integration validation | HIGH     | OPEN   |

### Architectural Constraint

Do NOT move domain authorization or business invariants to the Desktop.

Do NOT move plaintext asset handling, raw DEK handling, or raw KEK handling to the Backend.

The responsibility boundary is defined by the approved Client / Backend / Cloud ADR.

### Deferred

Lazy DEK re-wrapping is deferred to V2.

Federation-specific cryptographic workflows remain subject to the Federation architecture.

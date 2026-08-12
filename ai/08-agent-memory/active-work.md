
# Ankhora Active Work

## Purpose

This document describes the current engineering focus.

AI assistants should read this file before starting development tasks.

This document represents the active development context.

Completed work should be moved to historical documentation or architecture decisions.

---

# Current Objective

## Complete C3 Collaboration Foundations

Current focus:

```

C3 bounded context stabilization

```

The goal is to establish the core collaboration primitives required by Ankhora:

- workspaces
- channels
- threads
- assets
- sharing
- trust groups
- federation foundations

---

# Current Development Area

Bounded Context:

```

C3 Collaboration

```

Current package:

```

internal/channel

```

Primary objective:

Complete Channel lifecycle and application workflows.

---

# Completed Work

## Channel Aggregate

Implemented:

- Channel domain model
- Channel repository abstraction
- Channel creation workflow
- Channel persistence flow

---

## CreateChannelUsecase

Status:

COMPLETED

Responsibilities:

- validate request
- create Channel aggregate
- persist through repository
- publish ChannelCreated event

Pattern established:

```

Request

↓

Application UseCase

↓

Domain Aggregate

↓

Repository

↓

Domain Event

```

---

## ListChannelUsecase

Status:

COMPLETED

Purpose:

Provide channel querying capability.

Implementation decisions:

- Query operation
- No domain event emission
- Repository read operation only
- No DomainBus dependency

Reason:

Queries do not represent domain state changes.

---

# Current Work

## ArchiveChannelUsecase

Status:

COMPLETED

Objective:

Introduce Channel lifecycle transition:

```

active

|

v

archived

```

---

# ArchiveChannel Design

## Ownership

Owner:

```

C3 bounded context

```

Reason:

C3 owns:

- collaboration structure
- collaboration lifecycle
- participant relationships

---

## Domain Rule

Archive is different from revoke.

```

Archive

=
soft closure

preserve history

pause collaboration

```
```

Revoke

=
remove trust

deny access

```

---

## Proposed Domain Change

Channel aggregate will own:

```

Archive()

```

The aggregate protects:

- valid state transitions
- invariant enforcement
- timestamp update

Expected transition:

Allowed:

```

Active → Archived

```

Forbidden:

```

Pending → Archived

Revoked → Archived

Archived → Archived

```

---

# Event Strategy

New domain event:

```

ChannelArchived

```

Owned by:

```

C3

```

Purpose:

Notify internal consumers that collaboration state changed.

Future possibility:

Promote to integration event for:

```

TraceCore

```

when lifecycle history synchronization is implemented.

---

# Pending Architectural Decisions

## Event Publication Reliability

Current design:

```

Update Channel

↓

Publish Event

```

Potential future improvement:

```

Aggregate Change

↓

Transactional Outbox

↓

Event Dispatcher

↓

Consumers

```

Decision:

Not required for current implementation.

Revisit when:

- federation synchronization increases
- TraceCore integration becomes active
- distributed consistency becomes critical

---

# Next Implementation Steps

## Step 1

Implement domain changes:

- add Archived status
- add ArchivedAt timestamp
- implement Channel.Archive()

---

## Step 2

Implement event:

```

ChannelArchived

```

---

## Step 3

Implement:

```

ArchiveChannelUsecase

```

Flow:

```

Request

↓

Validate

↓

Repository.GetChannel()

↓

Channel.Archive()

↓

Repository.UpdateChannel()

↓

Publish Event

```

---

## Step 4

Testing

Required:

Domain:

- active archive success
- invalid transitions
- already archived protection

Application:

- success
- repository errors
- event errors
- validation errors

---

# Current AI Instructions

When working on this area:

AI must:

1. Read:
   - ai/04-contexts/c3.md
   - ai/02-architecture/event-driven-design.md
   - ai/01-principles/ownership.md

2. Preserve:
   - C3 ownership boundaries
   - domain event patterns
   - DDD aggregate rules

3. Avoid:
   - putting lifecycle rules in use cases
   - coupling C3 directly to TraceCore
   - bypassing repositories

---

# Current Status

Phase:

```

C3 Collaboration Core

```

Current feature:

```

Channel Lifecycle Management

```

Last completed:

```

ArchiveChannelUsecase

```

Next:

```

Thread lifecycle

```

## Current Work

### Client / Backend / Cloud Integration

**Status:**

ARCHITECTURE APPROVED — IMPLEMENTATION PENDING

**Architectural boundary:**

* Desktop is the cryptographic execution authority.
* Backend is the domain and access-control authority.
* Cloud is untrusted persistence and relay infrastructure.
* Backend MUST remain zero-knowledge with respect to plaintext assets, raw DEKs, and raw KEKs.

**Desktop responsibilities:**

* Plaintext asset processing
* Asset encryption/decryption
* DEK generation
* KEK generation during TrustGroup creation/rotation
* WrappedKEK generation
* WrappedDEK generation
* Private-key management
* Local cryptographic state

**Backend responsibilities:**

* Authentication verification
* Authorization
* Domain invariants
* Aggregate mutations
* Persistence orchestration
* Domain events
* Cross-context workflows

**Cloud responsibilities:**

* Persist encrypted assets
* Persist domain records
* Persist WrappedDEK / WrappedKEK envelopes
* Persist public keys
* Persist event/audit data
* Never receive plaintext assets, raw DEKs, or raw KEKs

**Architectural source:**

* ADR: Client / Backend / Cloud Responsibility Boundary

---

## Current Implementation Phase

### Phase: Cryptographic Identity + TrustGroup Integration

Implementation order:

1. Formalize Device entity in `internal/identity`
2. Extend TrustGroupMember for multi-device key envelopes
3. Extend WrappedGroupKey for key versioning
4. Implement TrustGroup KEK rotation
5. Refactor ShareEntry asset sharing to TrustGroup-level sharing
6. Implement Desktop Cryptographic Engine / Wails bridge
7. Define Desktop → Backend command contracts
8. Connect Backend → Cloud persistence
9. Validate end-to-end encrypted asset sharing

---

## Cryptographic Rotation Policy

**V1:** Eager DEK re-wrapping.

When a TrustGroup KEK rotates:

```text
KEK v1
  ↓
generate KEK v2
  ↓
unwrap existing DEKs using KEK v1
  ↓
wrap DEKs using KEK v2
  ↓
create WrappedKEK v2 for authorized devices
  ↓
persist new envelopes
```

Lazy DEK re-wrapping is explicitly deferred to V2.

---

## Completed Work

* Workspace creation
* Channel creation
* Thread creation
* Channel lifecycle
* Channel Slot lifecycle use cases
* TrustGroup creation
* TrustGroup member model
* Asset creation
* ShareEntry creation
* Initial TrustGroup / ShareEntry integration

---

## Next

### Step 1 — Device Model

Completed:

* Created `internal/identity/domain/device.go` (Device aggregate, Status, KeyType, IsActive(), Revoke())
* Implemented Device repositories (`MemoryDeviceRepository`, `GormDeviceRepository`)
* Implemented Device application use cases (`CreateDeviceUseCase`, `GetDeviceUseCase`, `ListDevicesUseCase`, `RevokeDeviceUseCase`)
* Published `DeviceCreated` and `DeviceRevoked` domain/application events

---

### Step 2 — Identity Device ↔ TrustGroup KeyEnvelope Integration

Completed:

* Created `DeviceResolver` port in `internal/trust_group/application/ports/device_resolver.go`
* Created `IdentityDeviceAdapter` in `internal/trust_group/infrastructure/adapters/identity_device_adapter.go`
* Created `AddTrustGroupKeyEnvelopeUseCase` in `internal/trust_group/application/usecases/envelope/add_envelope_usecase.go`
* Implemented cross-context validation (TrustGroup existence, member check, active device check, vault ownership check, KEKVersion check, duplicate active envelope rejection)
* Verified tests in `trust_group/application/usecases/envelope` and `trust_group/infrastructure/adapters`

---

### Step 3 — Desktop Cryptographic Orchestration

Completed:

* Extended `VaultKeyring` and `KeyringService` with `KeyTypeTrustGroupKEK` and `GetTrustGroupKEK(kr, trustGroupID, version)` / `StoreTrustGroupKEK(...)`
* Implemented `TrustGroupCryptoOrchestrator` in `internal/trust_group/application/orchestrator/crypto_orchestrator.go`
* Implemented local KEK/DEK generation, AES-256-GCM payload encryption, DEK wrapping, multi-device asymmetric KEK wrapping (`EncryptPayload`), and opaque DTO envelope building
* Verified end-to-end device unwrapping test, active vs revoked device envelope filtering, multi-device envelope distinction, and zero-knowledge boundary in `crypto_orchestrator_test.go`

---

### Step 4 — KEK Rotation

Next focus:

Create:

`internal/trust_group/application/usecases/rotate_kek.go`

V1 uses eager DEK re-wrapping.



---

### Step 4 — ShareEntry Refactoring

Refactor asset sharing from member-level targeting to TrustGroup-level sharing.

---

### Step 5 — Desktop Cryptographic Engine

Implement the cryptographic execution boundary in the Desktop/Wails layer.

Backend must receive only encrypted payloads and cryptographic envelopes.

---

### Step 6 — Communication Layer

Connect:

```text
Desktop
    ↓
Backend API / Commands
    ↓
Application Usecases
    ↓
Domain
    ↓
Cloud Repository
```

Do not bypass the Backend domain/application layer from the Desktop.



V1 Multi-Device Key Model — APPROVED

- Device is owned by Identity BC.
- TrustGroup owns TrustGroupMember.
- TrustGroup owns TrustGroupKeyEnvelope.
- TrustGroupKeyEnvelope targets a specific Device.
- Multiple devices may belong to one TrustGroupMember.
- Device addition does not rotate KEK.
- Member removal/security breach triggers KEK rotation.
- V1 KEK rotation uses eager DEK re-wrapping.
- ShareEntry tracks KEKVersion.
- Raw KEK/DEK never crosses Desktop → Backend boundary.

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



## Current Architecture State

### Client / Backend / Cloud Boundary

The architecture has been approved.

```text
Desktop
  └── Cryptographic execution authority

Backend
  └── Domain + authorization authority

Cloud
  └── Untrusted encrypted persistence
```

The Desktop may handle plaintext and raw cryptographic material locally.

The Backend MUST NOT receive, process, log, or persist:

* plaintext asset payloads
* raw DEKs
* raw KEKs

The Backend remains authoritative for:

* authentication
* authorization
* aggregate invariants
* state transitions
* persistence
* domain events

### Current C3 State

Completed:

* Workspace creation
* Channel creation
* Thread creation
* Channel lifecycle
* Channel Slot lifecycle
* TrustGroup creation
* TrustGroup membership structures
* Asset creation
* ShareEntry creation
* Initial TrustGroup / ShareEntry integration
* Identity Device aggregate, use cases, and persistence repositories
* Device ↔ TrustGroup KeyEnvelope cross-context application integration (`AddTrustGroupKeyEnvelopeUseCase`)
* Desktop Cryptographic Orchestrator (`TrustGroupCryptoOrchestrator`)
* Collaborative `ShareEntry` / CID Workflow Integration (`CreateCollaborativeShareUseCase`)
* Production-Hardened V1 KEK Rotation (`RotateTrustGroupKEKUseCase`) with atomicity, idempotency, stale client rejection, and concurrency safety

### Current Focus

Move from domain implementation toward:

**Desktop ↔ Backend ↔ Cloud integration**

Immediate priority:

**API & Wails Handler Integration**

## Channel Desktop → Cloud Status — 2026-08-15

### Completed (live-verified against the running Cloud backend)

Channel lifecycle, participants, and invitations are operational end-to-end:

```text
React
→ Wails
→ App (Wails bindings)
→ ChannelHandler
→ ChannelUseCase
→ TracecoreClient (internal/tracecore/channel_client.go)
→ Cloud REST API
```

All HTTP lives exclusively in `internal/tracecore/channel_client.go`. Cloud wire
DTOs use default Go JSON encoding (capitalized field names) and are mapped
explicitly (`channel_cloud_dto.go` → domain) — never unmarshaled into domain
directly.

#### Channel lifecycle

- Create (POST /channels) → List (GET /channels/workspace/{id}) → Activate
  (POST /channels/{id}/activate) → Revoke (POST /channels/{id}/revoke).
- Activation is Cloud-authoritative: a gated slot is fulfilled when a
  participant matches its VaultID (or Role) or an assignment targets it.

#### Participants

- AddParticipant (POST /channels/{id}/participants) → ListParticipants
  (GET /channels/{id}/participants). Cloud derives/validates role from the
  optional slot_id; the Desktop never decides locally.

#### Invitations (this session)

- InviteToChannel (POST /channels/{id}/invitations) body
  `{channel_id, inviter_vault_id, invitee_vault_id}` → returns pending
  Invitation. Cloud dedupes pending invitations per channel+invitee.
- AcceptChannelInvitation (POST /channels/invitations/{invitation_id}/accept —
  note: nested under the /api/channels group, NOT /api/invitations/{id}/accept)
  body `{invitation_id, invitee_vault_id, invitee_public_key}` → returns the
  accepted Invitation and creates the participant (Direction bidirectional,
  empty Role, JoinedAt set). Accepting twice is idempotent (no duplicate
  participant). Wrong invitee → 400 `invitation not for you`.
- Invitation carries NO slot/role info; an invitation-accepted participant has
  empty Role, so a gated slot can only be satisfied via VaultID matching.
- Acceptance identity: the Desktop user's own vault
  (`vault_runtime_context.VaultID` + `UserConfig.stellar_account.public_key`)
  acts as the invitee.

### Known limitations

- Desktop DeleteChannel is a stub; temporary live-test channels/invitations are
  cleaned by reporting their IDs, not by destructive API calls.
- Frontend has no invitation-list endpoint (Cloud has none); the Invitations
  panel lists session-created invitations only.
- Pre-existing `RequestChallenge` Wails-binding warning (unrelated) and
  `TestDirectAddToIPFSCall` failure (needs live local IPFS) remain.

# 3. Update `ai/08-agent-memory/current-state.md`

This is particularly important because this is the file an AI engineer should consult before touching the code.

Add:

```md
## Channel Desktop → Cloud Status — 2026-08-15

### Completed

Channel Create/List/Display is operational.

Desktop → Cloud path:

```text
React
→ Wails
→ App
→ ChannelHandler
→ ChannelUseCase
→ TracecoreClient
→ Cloud REST API
````

Dedicated client:

```text
internal/tracecore/channel_client.go
```

Cloud wire-format DTO:

```text
internal/tracecore/types/channel_cloud_dto.go
```

The Cloud Channel response uses Go-default JSON field names. It MUST be mapped through TraceCore DTOs rather than unmarshaled directly into the domain.

### Create

Working against Cloud.

### List

Working against Cloud.

Endpoint:

```text
GET /channels/workspace/{workspace_id}
```

### Display

Working.

The C3 ledger and Channel detail consume the real Channel store populated from Cloud.

Mocks are not part of the production Channel list/detail workflow.

### Channel persistence

Cloud persistence was corrected on 2026-08-15.

Slots:

```text
PRIMARY KEY (channel_id, id)
```

Assignments:

```text
PRIMARY KEY (channel_id, slot_id)
```

This prevents cross-channel corruption when multiple Channels use identical template slot IDs.

Cloud Create now returns the persisted representation.

### Federation

`channel.created` with zero recipients is a valid no-op.

Do not treat:

```text
federation: no target recipients found
```

as a Channel creation failure when no remote participants exist.

### Next lifecycle work

Before implementing activation on Desktop:

1. Persist Channel wizard slots.
2. Persist assignments.
3. Verify the persisted representation through GET.
4. Then implement Desktop activation.

Do not implement activation as an isolated UI feature while the prerequisites are not persisted.

````

---




### Key Management

TrustGroup owns the conceptual KEK.

Devices receive versioned WrappedKEK envelopes.

V1 uses **eager DEK re-wrapping** during KEK rotation.

Lazy re-wrapping is deferred to V2.



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

---

## Desktop ConnectVault Protocol & Cloud Workspace DTO Integration — 2026-08-18

### Completed

1. **Authentication vs Vault Delegation Boundary**:
   - `SignIn()` authenticates user credentials and retrieves Cloud bearer token (JWT).
   - Vault operation access requires an active user-vault delegation in Cloud (`UserVaultIdentity`).
   - `ConnectVault()` flow executes: `POST /api/identity/challenge` → local Ed25519 Stellar signature of challenge payload → `POST /api/identity/` registration.
   - `SignIn()` invokes `ConnectVault()` immediately after opening the vault, ensuring delegation is established before accessing Ledger.

2. **Workspace Infrastructure DTO Mapping**:
   - Cloud serializes `Workspace` aggregates using PascalCase field names (`ID`, `VaultID`, `Name`, `Description`, `Status`, `OwnerID`, `CreatedAt`, `UpdatedAt`).
   - Implemented `CloudWorkspaceDTO` in `internal/tracecore/types/workspace_cloud_dto.go` and explicit mapper `MapCloudWorkspaceToTypes()`.
   - All fields (`VaultID`, `OwnerID`, `CreatedAt`, `UpdatedAt`) are strictly preserved when decoding Cloud responses into Desktop application/domain models.
   - Unit tests added in `internal/tracecore/workspace_dto_test.go` to prevent regression.
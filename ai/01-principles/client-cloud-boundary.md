Client Responsibility Principle

Desktop/mobile clients are non-authoritative clients.

Clients may:
- handle presentation and interaction
- perform local validation for UX
- maintain local UI/session state
- perform cryptographic operations that must remain client-side
- optimistically update UI when explicitly supported

Clients must not:
- implement authoritative domain invariants
- orchestrate cross-bounded-context workflows
- determine authorization
- directly coordinate cloud persistence
- reproduce cloud business workflows
- assume local state is authoritative

Authoritative business decisions belong to the backend/cloud application layer.


Command-Oriented Client Boundary

The client communicates intent, not implementation steps.

Client:
    "Share Asset X with Member Y"

Backend:
    determines how that intent is fulfilled.


## Cloud Wire Format Is an Integration Concern

Cloud persistence/domain serialization is not necessarily identical to Desktop domain serialization.

The Desktop MUST NOT assume that:

```text
Cloud JSON == Desktop domain JSON
````

Instead:

```text
Cloud API
    ↓
Cloud DTO
    ↓
TraceCore mapper
    ↓
Desktop domain
```

The TraceCore boundary is responsible for absorbing transport-format differences.

If Cloud uses Go-default JSON serialization while Desktop uses explicit snake_case JSON tags, introduce a dedicated Cloud DTO rather than modifying the domain solely to match the wire format.

### Rule

> Never contaminate the domain model with transport-specific serialization requirements.

````

---

Connect Vault is not an authentication operation and must never become one.

So:

SignIn()
    → local identity
    → Cloud authentication


ConnectVault()
    → vault cryptographic proof
    → Cloud delegation


Ledger()
    → consume already-established authorization

That separation is worth protecting.

In particular, do not put Connect Vault inside ListWorkspaces() as a fallback. If the vault isn't connected, the application should know that explicitly.

----

## Authentication vs Vault Delegation Invariants

1. **Authentication and Vault Authorization Are Separate Concerns**:
   - `SignIn()` returns a Cloud JWT bearer token (`401` class concern if invalid/missing).
   - Having a valid Cloud bearer token does **not** grant access to a vault.
   - Access to a vault requires an active `UserVaultIdentity` delegation record in Cloud (`403` class concern if non-delegated).

2. **Explicit ConnectVault Protocol**:
   - A Desktop-created vault is initially local and is **not** automatically registered/delegated in Cloud.
   - Connecting a vault is an explicit `ConnectVault()` protocol: `POST /api/identity/challenge` → local Ed25519 signature of challenge payload → `POST /api/identity/` registration.
   - The Desktop owns the vault private key and performs the local signature. The private key MUST NEVER be transmitted to Cloud.

3. **Infrastructure DTO Mapping & Field Integrity**:
   - External Cloud DTOs (e.g. `CloudWorkspaceDTO`) must be explicitly mapped at the TraceCore infrastructure boundary and must NOT leak into domain models.
   - Silent loss of security-relevant or identity-relevant DTO fields (e.g. `VaultID`, `OwnerID`, `CreatedAt`, `UpdatedAt`) is unacceptable; regression tests must verify field preservation.

4. **CreateVault vs RegisterVault vs DelegateVault Distinction**:
   - `ConnectVault` is a trust-establishment operation, not a vault-creation operation.
   - The Desktop creates/owns the local vault. Cloud registration establishes Cloud's awareness of that vault and creates caller delegation (`UserVaultIdentity`).
   - `CreateVault != RegisterVault != DelegateVault` — these represent three distinct responsibility domains.

5. **Architectural Boundary Protection Rule**:
   - **Do NOT modify these boundaries merely to make a test pass**. Treat ADR-0006 and the client/cloud boundary rules as immutable architectural constraints. If implementation code appears to conflict with them, stop and report the conflict before modifying the architecture.

----


TraceCore is an internal Ankhora Cloud service. The Desktop communicates exclusively with Ankhora Cloud for remote operational services.

-----

Your phrase from 3.26 is becoming the core architectural statement:

The Desktop can cryptographically establish that this sovereign vault belongs to this authenticated user, and Cloud subsequently authorizes access to that vault.

Phase 3.27 adds:

C3 authorizes collaboration without duplicating vault authentication or delegation.

Phase 3.28 adds the third piece:

Authorized C3 operations become distributable domain events without coupling the originating bounded context to Federation.

Together, those three statements describe something much bigger than a working CRUD application.

Sovereign Identity
        ↓
Cloud Authorization
        ↓
Collaborative Authority
        ↓
Domain Events
        ↓
Federated Distribution
        ↓
Remote Projection

That's the platform architecture you've been trying to reach.
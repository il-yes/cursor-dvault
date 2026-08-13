# Vault ↔ Cloud Contract

## Purpose

This document defines the application-level relationship between Ankhora Vault Desktop and Ankhora Cloud.

The objective is to keep the two repositories independently maintainable while allowing them to participate in the same business workflows.

---

# 1. Contract Over Transport

The architectural contract is independent from the transport mechanism.

The current Desktop application may invoke backend functionality through direct application/Wails methods rather than a browser HTTP API.

This does not make the backend part of the Desktop domain.

The conceptual boundary remains:

```text
Vault Application
       │
       │ contract
       ▼
Cloud Application
```

Transport can evolve independently.

---

# 2. Request Ownership

Desktop initiates user-driven operations such as:

* creating collaborative assets;
* sharing assets;
* adding TrustGroup devices;
* rotating local cryptographic material;
* preparing encrypted payloads.

Cloud validates and commits operations requiring authoritative server state.

---

# 3. Cloud Response

Cloud responses should expose:

* identifiers;
* authorization results;
* state transitions;
* versions;
* metadata;
* opaque encrypted values;
* errors.

Cloud responses must not expose:

* private keys;
* raw DEKs;
* raw KEKs;
* plaintext payloads;
* unlocked keyrings.

---

# 4. Application Boundary

The backend must expose explicit application contracts.

Examples include:

```text
CreateCollaborativeShare
AddTrustGroupKeyEnvelope
RotateTrustGroupKEK
ResolveDevice
ListActiveDevices
```

These are application operations.

They are not permissions for another bounded context to directly manipulate another aggregate.

---

# 5. Identity Boundary

Identity owns Device state.

Collaboration consumes Device information through a port:

```text
DeviceResolver
```

The collaboration domain must not import Identity domain aggregates.

The Cloud implementation may provide an adapter:

```text
Identity Device Repository
        │
        ▼
IdentityDeviceAdapter
        │
        ▼
DeviceResolver
        │
        ▼
Collaboration Application
```

This preserves bounded-context independence.

---

# 6. TrustGroup Boundary

TrustGroup owns:

* membership;
* KEK version;
* key-envelope lifecycle;
* TrustGroup-specific authorization rules.

Collaboration orchestrates workflows involving TrustGroup and ShareEntry.

The TrustGroup domain must not own ShareEntry persistence.

---

# 7. ShareEntry Boundary

ShareEntry owns collaborative asset-sharing metadata.

At minimum:

```text
AssetCID
TrustGroupID
WrappedDEK
KEKVersion
```

The actual encrypted asset payload is not a plaintext domain object.

---

# 8. Error Semantics

Cross-repository operations should preserve meaningful domain/application errors.

Examples:

```text
DeviceNotFound
DeviceRevoked
DeviceMemberMismatch
MemberNotInTrustGroup
TrustGroupNotFound
StaleKEKVersion
DuplicateKeyEnvelope
InvalidKEKVersionIncrement
```

The transport layer may serialize these errors differently, but their business meaning must remain stable.

---

# 9. Versioning

Versioned cryptographic state must be explicit.

The Cloud side must never infer the client's intended KEK version from unrelated metadata.

The request must explicitly identify the expected version.

Example:

```text
OldVersion = 1
NewVersion = 2
```

Cloud validates the transition.

---

# 10. Request Idempotency

Operations susceptible to network retries should carry a RequestID.

This is especially important for:

* KEK rotation;
* collaborative share creation;
* federation delivery;
* remote commands.

Idempotency belongs to the authoritative side that commits the state.

---

# 11. No Hidden Cross-Repository Coupling

Neither repository may depend on:

```text
file paths
internal packages
private structs
database tables
implementation-specific services
```

from the other repository.

Only explicit contracts may cross the repository boundary.

---

# 12. Contract Change Procedure

When changing a shared contract:

1. Identify the owning repository.
2. Update the cross-repository contract.
3. Update the implementation in the owning repository.
4. Update the consuming repository.
5. Add compatibility tests where appropriate.
6. Run both repositories' relevant test suites.
7. Update the relevant ADR if the boundary itself changes.

A contract change is not complete until both sides agree.

---

# 13. Current Platform Model

```text
┌───────────────────────────────────────────────────────────┐
│                    ANKHORA PLATFORM                       │
│                                                           │
│  ┌────────────────────┐       ┌────────────────────────┐ │
│  │   ANKHORA VAULT    │       │     ANKHORA CLOUD      │ │
│  │                    │       │                        │ │
│  │ Wails/Desktop      │       │ Backend                │ │
│  │ Vault              │       │ Identity               │ │
│  │ Keyring            │◄─────►│ TrustGroup             │ │
│  │ Crypto             │       │ Collaboration           │ │
│  │ Plaintext          │       │ Persistence             │ │
│  │ Private Keys       │       │ Federation              │ │
│  └────────────────────┘       └────────────────────────┘ │
│             │                           │                 │
│             └────── explicit contract ──┘                 │
└───────────────────────────────────────────────────────────┘
```

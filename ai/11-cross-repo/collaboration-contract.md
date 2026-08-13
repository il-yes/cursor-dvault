# Collaboration Contract

## Purpose

This document defines the shared contract between Desktop and Cloud for TrustGroup-based collaborative sharing.

---

# 1. Core Entities

The collaboration contract uses the following identifiers:

```text
TrustGroupID
MemberID
DeviceID
AssetCID
KEKVersion
```

These identifiers must not be conflated.

### TrustGroupID

Identifies the collaboration TrustGroup.

### MemberID

Identifies the authorized Vault/member participating in the TrustGroup.

### DeviceID

Identifies a concrete device belonging to a member.

A member may have multiple active devices.

### AssetCID

Identifies the encrypted asset payload.

### KEKVersion

Identifies the TrustGroup key version protecting collaborative asset DEKs.

---

# 2. Device Model

Identity owns the authoritative Device aggregate.

The relevant contract is:

```text
Device
 ├── ID
 ├── VaultID
 ├── PublicKey
 ├── KeyType
 ├── Status
 ├── CreatedAt
 └── RevokedAt
```

Cloud collaboration logic resolves devices through an application-level port.

TrustGroup must not directly depend on Identity's internal implementation.

---

# 3. Device Authorization

A device may receive a TrustGroup KEK only if:

```text
Device exists
AND
Device is active
AND
Device.VaultID == MemberID
AND
Member belongs to TrustGroup
AND
KEKVersion matches current TrustGroup version
```

A revoked device must never receive a new KEK envelope.

---

# 4. Key Envelope

A TrustGroup key envelope represents:

```text
TrustGroupID
MemberID
DeviceID
KEKVersion
WrappedKEK
```

`WrappedKEK` is opaque.

Cloud must not attempt to decrypt or inspect it.

Each active device receives its own envelope.

Therefore:

```text
Member A
 ├── Device A1 → WrappedKEK
 └── Device A2 → WrappedKEK

Member B
 └── Device B1 → WrappedKEK
```

---

# 5. Collaborative ShareEntry

A collaborative ShareEntry references:

```text
AssetCID
TrustGroupID
WrappedDEK
KEKVersion
```

The ShareEntry does not contain raw DEK or KEK material.

---

# 6. Asset Storage

The asset flow is:

```text
Desktop
  │
  ├── plaintext
  ├── generate DEK
  ├── encrypt payload
  └── produce encrypted bytes
          │
          ▼
       AssetCID
          │
          ▼
       Cloud / IPFS
```

The CID must represent the encrypted payload.

The Cloud backend may persist metadata referencing the CID without knowing the plaintext.

---

# 7. Create Collaborative Share

The conceptual operation is:

```text
Desktop
   │
   ├── Encrypt payload
   ├── Generate/obtain DEK
   ├── Wrap DEK with KEK
   ├── Wrap KEK for active devices
   │
   ▼
Cloud
   │
   ├── Validate TrustGroup
   ├── Validate members
   ├── Validate devices
   ├── Validate KEKVersion
   ├── Persist ShareEntry
   └── Persist KeyEnvelopes
```

Cloud never regenerates the cryptographic values.

---

# 8. KEK Rotation

KEK rotation is split between two authorities.

### Desktop

Desktop performs:

```text
Generate KEK N+1
Unwrap DEKs with KEK N
Re-wrap DEKs with KEK N+1
Wrap KEK N+1 for remaining devices
```

### Cloud

Cloud performs:

```text
Validate current KEK version
Validate membership
Validate devices
Transition TrustGroup version
Revoke old envelopes
Persist new envelopes
Persist new WrappedDEKs
Commit transaction
```

---

# 9. Rotation Invariants

A successful rotation must guarantee:

```text
newKEK != oldKEK

newVersion = oldVersion + 1

AssetCID unchanged

EncryptedPayload unchanged

DEK unchanged

WrappedDEK changed

WrappedKEK changed

removed devices receive no new envelope
```

---

# 10. Transaction Boundary

Cloud owns persistence atomicity.

A rotation must not expose partial state.

Conceptually:

```text
BEGIN TRANSACTION

Update TrustGroup
Update KeyEnvelopes
Update ShareEntries

COMMIT
```

If any persistence operation fails:

```text
ROLLBACK
```

The Desktop performs cryptography before submission; Cloud controls persistence atomicity.

---

# 11. Idempotency

Network retries must not produce duplicate rotations.

Rotation requests therefore require a unique `RequestID`.

For a previously committed request:

```text
RequestID = X
```

Cloud returns the previously committed result rather than performing another rotation.

This prevents:

```text
v1 → v2
retry
v2 → v3   ❌
```

The desired result is:

```text
v1 → v2
retry
return previous v2 result
```

---

# 12. Concurrency

Cloud owns optimistic concurrency.

A rotation request must include the client's expected old version.

Example:

```text
OldVersion = 1
```

If Cloud is already at:

```text
KEKVersion = 2
```

the request is stale and must be rejected.

Only one concurrent rotation may successfully transition a given version.

---

# 13. Backend Contract

The Cloud-facing collaboration contract must contain only:

```text
identifiers
metadata
opaque encrypted values
versions
authorization information
```

It must not contain:

```text
raw DEK
raw KEK
private key
plaintext
```

---

# 14. Evolution Rule

Changes to the collaboration cryptographic contract must be coordinated across both repositories.

If a field changes in:

```text
TrustGroup
Device
ShareEntry
KeyEnvelope
KEKVersion
AssetCID
```

the engineer must verify both:

```text
ankhora-vault
ankhora-cloud
```

before considering the change complete.

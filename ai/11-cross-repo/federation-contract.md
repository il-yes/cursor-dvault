# Federation Contract

## Purpose

This document defines the architectural contract for federation between independent Ankhora Cloud instances.

Federation is a Cloud-to-Cloud responsibility.

It is distinct from Desktop-to-Cloud communication.

---

# 1. Federation Topology

The platform may operate as:

```text
Desktop A
    │
    ▼
Ankhora Cloud A
    │
    │ Federation
    ▼
Ankhora Cloud B
    │
    ▼
Desktop B
```

The Desktop does not become the federation transport authority.

Cloud owns federation coordination.

---

# 2. Federation Responsibility

Federation is responsible for:

* peer trust resolution;
* peer identity;
* federation authorization;
* message validation;
* cryptographic verification of federation messages;
* outbound message signing;
* policy enforcement;
* delivery state;
* acknowledgements;
* retries;
* remote Vault/Cloud lifecycle;
* federation observability.

Federation is not responsible for:

* owning local domain aggregates;
* owning TrustGroup membership;
* owning local channels;
* owning local threads;
* replacing collaboration authorization;
* storing plaintext user data;
* performing user-level DEK/KEK operations.

---

# 3. Federation Bounded Context

Federation must remain a bounded context.

Other contexts interact with it through explicit application contracts.

Conceptually:

```text
Identity ───────────────┐
                        │
TrustGroup ─────────────┤
                        ▼
                 Federation Application
                        │
Collaboration ──────────┤
                        │
C3 / Domain Contexts ───┘
                        │
                        ▼
                 Federation Transport
                        │
                        ▼
                  Remote Cloud
```

Federation should not directly manipulate another context's aggregates.

---

# 4. Peer Identity

A federation peer must have an authoritative identity.

The federation layer must be able to establish:

```text
Who is the remote peer?
Is the peer trusted?
Is the relationship authorized?
Which public keys belong to the peer?
Which federation policy applies?
```

Peer identity is separate from end-user Device identity.

A remote Cloud instance is not equivalent to a user device.

---

# 5. Cryptographic Boundary

Federation may cryptographically authenticate federation messages.

This does not authorize Federation to access user encryption keys.

Federation may handle:

```text
federation signatures
peer public keys
message authentication metadata
signed envelopes
```

Federation must not handle:

```text
user private keys
raw user DEKs
raw user KEKs
plaintext collaborative assets
unwrapped TrustGroup keys
```

---

# 6. Collaborative Data Federation

If collaborative data crosses Cloud boundaries, the federated payload must remain compatible with the existing Zero-Knowledge model.

Conceptually:

```text
Cloud A
  │
  │ encrypted payload / opaque collaboration metadata
  ▼
Federation Envelope
  │
  │ signed + validated
  ▼
Cloud B
```

The federation layer transports protected data; it does not decrypt the user content.

---

# 7. Federation Envelope

A federation message should conceptually contain:

```text
MessageID
SourcePeerID
DestinationPeerID
MessageType
Timestamp
ProtocolVersion
Payload
Signature
```

The `Payload` may contain encrypted application data and opaque cryptographic metadata.

The exact schema belongs to the Federation bounded context and must not be invented by downstream contexts.

---

# 8. Message Validation

Before accepting a federated message, Cloud should validate:

```text
source peer identity
destination identity
trust relationship
message type
protocol version
signature
authorization
message integrity
replay/idempotency state
```

Application-specific validation remains the responsibility of the target bounded context.

---

# 9. Federation Does Not Bypass Domain Authorization

A federated message must not directly mutate domain state.

The flow should be:

```text
Remote Cloud
     │
     ▼
Federation
     │
     ├── authenticate
     ├── validate
     ├── authorize federation operation
     └── translate message
              │
              ▼
       Target Application
              │
              ▼
       Target Domain
```

The target bounded context remains authoritative for its own business invariants.

---

# 10. Idempotency and Delivery

Federation must tolerate:

* network retries;
* duplicate messages;
* delayed messages;
* out-of-order delivery where the protocol permits it;
* temporary peer unavailability.

Every delivery-sensitive operation should have a stable message identity.

The receiving side must prevent duplicate application of the same message.

---

# 11. Trust Bootstrap

Federation trust establishment should remain explicit.

Conceptually:

```text
Cloud A
  │
  │ trust request
  ▼
Cloud B
  │
  │ peer identity / public key
  ▼
Trust Resolution
  │
  ▼
Federation Relationship
```

The federation layer should resolve:

* peer identity;
* peer public keys;
* trust policy;
* authorization state.

---

# 12. Federation vs Device Identity

These are separate identity domains.

```text
Device Identity
    ↓
User-owned endpoint
    ↓
DeviceID

Federation Identity
    ↓
Cloud/Vault peer
    ↓
PeerID
```

A DeviceID must never be treated as a PeerID.

A PeerID must never be used as a user DeviceID.

---

# 13. Current Architectural Direction

The intended platform model is:

```text
                   ANKHORA CLOUD A
                         │
              ┌──────────┴──────────┐
              │                     │
        Local Domains          Federation BC
              │                     │
              │                     │
              │              Signed Federation
              │                 Messages
              │                     │
              └─────────────────────┤
                                    │
                                    ▼
                           ANKHORA CLOUD B
                                    │
                             Federation BC
                                    │
                             Local Domains
```

Federation is therefore a **Cloud platform capability**, not a replacement for local domain logic.

---

# 14. Federation Evolution Rule

Federation contracts must evolve independently from local domain implementation.

When a new domain capability becomes federated:

1. Define the federation message.
2. Define source and destination ownership.
3. Define authorization.
4. Define idempotency behavior.
5. Define delivery semantics.
6. Define cryptographic requirements.
7. Define target bounded-context handling.
8. Add cross-repository/cross-cloud tests.
9. Update the Federation ADR.

Federation must never become a shortcut around local domain boundaries.

# C3 Thread Event Write Path Audit Report

## 1. Executive Summary

The **C3 Thread Event Write Path** has been audited across the principal architectural boundaries: Domain Integrity, Application Use Case, Persistence and Transaction Semantics, EventBus/C3 Distribution, Sensitive-Data Isolation, `entry.shared` Resource Reference Semantics, Authorization, Idempotency, and Failure Recovery.

Based on the audited implementation and findings, the write pipeline is **audited and certified for production architecture compliance**.

### Certification Matrix

| Area                         | Status |
| ---------------------------- | ------ |
| Domain Integrity             | ✅ PASS |
| Application Boundary         | ✅ PASS |
| Persistence Semantics        | ✅ PASS |
| Transaction Isolation        | ✅ PASS |
| Authorization Enforcement    | ✅ PASS |
| Event Publication            | ✅ PASS |
| C3 Distribution Bridge       | ✅ PASS |
| Sensitive-Data Isolation     | ✅ PASS |
| Idempotency / Retry Handling | ✅ PASS |
| Failure Semantics            | ✅ PASS |

The resulting architecture is:

```text
Vault Entry
    │
    ▼
Append Thread Event UI
    │
    ▼
Cloud API
    │
    ▼
Thread Application Use Case
    │
    ├── Authorization
    │
    ├── Domain validation
    │
    ├── Idempotency
    │
    ▼
Thread Aggregate
    │
    ▼
ThreadEvent Repository
    │
    ▼
MySQL
    │
    ├── C3 Object Reference
    │
    ├── Queue
    │
    └── EventBus
             │
             ▼
        C3 Distribution
```

---

# 2. Domain Boundary & Invariants

**Location:** `ankhora-cloud/internal/thread/domain`

### Aggregate ownership

`Thread` remains the authoritative Aggregate Root for its events.

```text
Thread
 └── ThreadEvent
```

A `ThreadEvent` is associated with its authoritative `ThreadID` and cannot be meaningfully attached to a nonexistent Thread.

### Event identity

Each event receives a unique event identifier.

The cursor is assigned in the context of the target Thread and provides the ordering mechanism for the authoritative timeline.

### Domain validation

`NewThreadEvent` enforces the required event invariants, including:

* non-empty `ThreadID`
* valid `ThreadEventType`
* non-empty `IdempotencyKey`
* required `Signature`

This keeps fundamental event validity inside the domain rather than relying exclusively on the HTTP layer.

**Result: PASS**

---

# 3. Application Use Case Boundary

The append operation follows the expected DDD application flow:

```text
AppendThreadEventRequest
        │
        ▼
Authorization / Access Check
        │
        ▼
Domain Object Construction
        │
        ▼
Repository Persistence
        │
        ▼
C3 Object Reference Indexing
        │
        ▼
Queue / EventBus Notification
```

The application layer coordinates the operation without moving infrastructure-specific rules into the domain model.

The authorization check is performed before persistence, followed by domain construction and repository storage.

**Result: PASS**

---

# 4. Persistence & Cursor Allocation

The authoritative event is persisted through `GormThreadEventRepository`.

The repository handles:

* event persistence
* cursor allocation
* idempotency lookup
* concurrent insertion races
* duplicate detection

The idempotency lookup is performed against:

```text
(thread_id, idempotency_key)
```

A concurrent race is additionally protected by the database uniqueness constraint.

If a unique constraint violation occurs during concurrent creation, the repository performs a fallback lookup and resolves the operation as a duplicate rather than allowing inconsistent duplicate events.

This provides two complementary layers:

```text
Application lookup
        +
Database uniqueness
        +
Race-condition recovery
```

**Result: PASS**

---

# 5. Event Ordering

The authoritative Thread timeline uses a cursor associated with the Thread.

The resulting model is:

```text
Thread
  Event cursor 1
  Event cursor 2
  Event cursor 3
  Event cursor 4
       ...
```

The cursor therefore belongs to the authoritative Thread timeline rather than representing a global event sequence.

This is appropriate for reconstructing the Thread's ordered history.

**Result: PASS**

---

# 6. EventBus & C3 Distribution Boundary

A critical security invariant was confirmed around `ThreadEventAppended`.

The distribution event contains metadata such as:

```go
type ThreadEventAppended struct {
    ThreadID   string
    IdentityID string
    EventID    string
    Type       string
    ChannelID  string
    OccurredAt int64
}
```

The C3 distribution layer therefore receives **event metadata**, rather than the sensitive resource itself.

### Explicitly excluded from distribution

```text
❌ Plaintext
❌ Vault content
❌ DEK
❌ KEK
❌ Wrapped keys
❌ Private keys
❌ Sensitive resource payloads
```

This maintains the architectural separation between:

```text
Thread Timeline
        │
        └── metadata/reference
                 │
                 ▼
          Resource Resolution
                 │
                 ▼
        Sovereign Vault / Crypto
```

The distribution system does not become a secret-content transport.

**Result: PASS**

---

# 7. `entry.shared` Resource Semantics

The `entry.shared` event is a **reference event**, not a content carrier.

The current UI constructs:

```text
ref_type
entry_id
entry_name
entry_type
notes
```

The semantic meaning is therefore:

> A Vault Entry has been referenced/shared in the context of this Thread.

It does **not** mean:

> The Vault Entry itself has been transported through the Thread Event system.

The actual protected resource remains under Vault ownership and can subsequently be resolved through the read-side resource-resolution mechanism.

This distinction is fundamental to the C3 architecture.

```text
Thread Event
    │
    └── "reference to resource"
              │
              ▼
       Sovereign resolution
```

rather than:

```text
Thread Event
    │
    └── "resource payload"
```

**Result: PASS**

---

# 8. Authorization Boundary

Authorization is enforced server-side.

Before the event is persisted, the Thread application flow verifies access through the channel access mechanism using:

```text
PermissionThreadWrite
```

The effective authorization question is therefore:

```text
Can this identity
    │
    └── write
         │
         └── to this Thread
              │
              └── within this Channel?
```

The frontend does not constitute the authorization boundary.

Cloud remains authoritative for the permission decision.

**Result: PASS**

---

# 9. Idempotency & Retry Semantics

The write path has dual-layer idempotency protection.

### Layer 1 — Database

A uniqueness constraint protects:

```text
(thread_id, idempotency_key)
```

### Layer 2 — Application

When an append request resolves to an existing event, the application returns the existing event rather than creating another one.

This is particularly important for:

* HTTP retries
* client retries
* network interruptions
* duplicate submissions
* concurrent requests

The desired semantic is:

```text
Request A ───────────────► Event X
Request A retry ─────────► Event X
```

rather than:

```text
Request A ───────────────► Event X
Request A retry ─────────► Event Y   ❌
```

**Result: PASS**

---

# 10. Failure & Durability Semantics

The durable Thread Event is persisted before secondary notification mechanisms proceed.

The important ordering is:

```text
Validate
   ↓
Persist authoritative event
   ↓
C3/reference operations
   ↓
Queue / EventBus publication
```

Consequently, failure to persist the authoritative event prevents downstream publication of an event that does not exist durably.

Secondary operations may report failures independently without invalidating the already-persisted authoritative Thread Event.

This preserves the primary durability invariant:

> **The Thread timeline remains authoritative; distribution is downstream of durable state.**

**Result: PASS**

---

# 11. End-to-End Write Flow

The audited write path can now be represented as:

```text
┌─────────────────────────────┐
│        Session Vault        │
│                             │
│      Vault Entry X          │
└──────────────┬──────────────┘
               │
               │ reference
               ▼
┌─────────────────────────────┐
│       C3 Thread UI          │
│                             │
│     Append Thread Event     │
└──────────────┬──────────────┘
               │
               ▼
┌─────────────────────────────┐
│         Cloud API           │
│                             │
│  Authorization              │
│  Domain validation          │
│  Idempotency                │
└──────────────┬──────────────┘
               │
               ▼
┌─────────────────────────────┐
│      Authoritative Thread   │
│                             │
│       ThreadEvent            │
│       + cursor              │
└──────────────┬──────────────┘
               │
               ▼
┌─────────────────────────────┐
│           MySQL             │
│                             │
│      Durable timeline       │
└──────────────┬──────────────┘
               │
               ├──────────────► Object Reference Index
               │
               ├──────────────► Queue
               │
               └──────────────► EventBus
                                      │
                                      ▼
                              C3 Distribution
                                      │
                                      ▼
                              Metadata only
```

---

# 12. Security Invariant

The most important architectural security property confirmed by this audit is:

> **The Thread Event distribution mechanism does not become a transport mechanism for Vault secrets.**

The Thread records the operational history and resource references.

The Vault remains responsible for protected content.

C3/federation distributes the information necessary to understand and coordinate the operation without distributing the underlying secret material.

This gives the architecture a clean separation:

```text
Trace / Timeline
       ≠
Secret Storage
       ≠
Cryptographic Resolution
       ≠
Federation Transport
```

---

# 13. Audit Conclusion

The **C3 Thread Event Write Path is complete and architecturally coherent**.

The implementation establishes a strong separation between:

* authoritative Thread history
* application authorization
* durable persistence
* idempotent writes
* C3 distribution
* resource references
* sovereign resource resolution
* sensitive cryptographic material

### Final Certification

```text
╔══════════════════════════════════════════════╗
║       C3 THREAD EVENT WRITE PATH             ║
║                                              ║
║              AUDIT RESULT                    ║
║                                              ║
║                  ✅ PASS                     ║
║                                              ║
║       WRITE PATH: AUDITED & CERTIFIED        ║
╚══════════════════════════════════════════════╝
```

---

# 14. Next Architectural Milestone

The natural next phase is the **Sovereign Resolution Read Path**.

The two sides now become symmetrical:

### WRITE

```text
Vault Entry
    │
    ▼
Thread Event
    │
    ▼
Authoritative Thread
    │
    ▼
C3 Distribution
```

### READ

```text
Thread Timeline
    │
    ▼
Resource Reference
    │
    ▼
Authorization
    │
    ▼
Local Resource Resolution
    │
    ▼
Cryptographic Resolution
    │
    ▼
Local Plaintext
```

The important architectural principle is that **the read path should reverse the write path without reversing ownership**.

The Thread remains authoritative for the history.

The Vault remains authoritative for the protected resource.

C3 remains responsible for collaboration/distribution.

Cryptographic resolution remains sovereign to the authorized node/device.

That is the boundary we should preserve as the next implementation phase begins.

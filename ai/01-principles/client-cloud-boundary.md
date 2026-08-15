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
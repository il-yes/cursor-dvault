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
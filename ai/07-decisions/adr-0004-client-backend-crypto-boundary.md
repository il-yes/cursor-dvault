ARCHITECTURAL DECISION

Status:
APPROVED

Desktop:
- Authoritative for local UI rendering, local session management, plaintext payload processing, raw key management, AES-256-GCM asset encryption/decryption, DEK generation, KEK generation (upon TG creation/rotation), WrappedKEK generation, and WrappedDEK generation.
- Non-authoritative for domain invariants, authorization, authentication, or state persistence.

Backend:
- Authoritative for authentication context verification, authorization checks, domain invariant validation, aggregate lifecycle mutations, persistence orchestration, domain event publication, and cross-context workflow execution.
- Strictly Zero-Knowledge: Must NEVER receive, process, log, or persist plaintext asset payloads, unencrypted DEKs, or raw KEKs.

Cloud:
- Untrusted persistence and relay infrastructure storing encrypted binary blobs, structural entity records, wrapped key envelopes (WrappedDEK, WrappedKEK), public keys, and event logs.

Cryptographic Authority:
- Split Boundary: Desktop Client holds Cryptographic Execution Authority (possesses private keys and performs cipher operations); Backend Application holds Domain & Access Control Authority (enforces authorization and state transitions).

Required Changes:
- Update TrustGroupMember and WrappedGroupKey domain structs to support multi-device/versioned key envelopes.
- Refactor ShareAssetWithTrustGroupMemberUsecase to target TrustGroup-level asset sharing.
- Formalize Device entity in internal/identity.

Existing Components To Keep:
- Workspace aggregate & use cases
- Channel aggregate, slots, policies, and use cases
- Thread aggregate & use cases
- Asset & ShareEntry aggregate models

Existing Components To Modify:
- internal/trust_group/domain/aggregate.go
- internal/collaboration/application/usecases/share_asset_with_member.go

New Components:
- internal/identity/domain/device.go (Device Aggregate)
- internal/trust_group/application/usecases/rotate_kek.go (KEK Rotation Use Case)
- Desktop Cryptographic Engine / Wails Bridge (Client-side cipher & key wrapping)


Device ↔ TrustGroupMember Relationship:
- Decoupled across Bounded Context boundaries. Identity owns Device. TrustGroup owns TrustGroupMember and TrustGroupKeyEnvelope.
- TrustGroupKeyEnvelope maps (TrustGroupID, MemberCID, DeviceID, KEKVersion) -> WrappedKEK.

Key Envelope Model:
- Represented by TrustGroupKeyEnvelope entity inside the TrustGroup Aggregate Root.
- Supports multi-device per member, explicit KEK versioning, and per-device public-key targeting.

Rotation & Eager Re-wrapping Policy (V1):
- Member Removal / Security Breach triggers KEK rotation (KEK v1 -> KEK v2).
- Eager DEK re-wrapping: Admin Desktop client unwraps DEKs using KEK v1, re-wraps DEKs using KEK v2, and updates ShareEntry.KEKVersion to 2 in a single atomic transaction.
- Device addition DOES NOT rotate KEK (issues envelope for KEK vCurrent).

ShareEntry Model:
- ShareEntry explicitly tracks KEKVersion matching the TrustGroup KEK version used to wrap the DEK.

Next Implementation Step:
- Update domain aggregates in internal/trust_group/domain and internal/c3_asset/domain to reflect TrustGroupKeyEnvelope and KEKVersion fields.

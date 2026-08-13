# Cross-Repository Architecture Contracts

## Purpose

The `09-cross-repo` directory defines the architectural contracts shared between the Ankhora Vault Desktop project and the Ankhora Cloud backend project.

These documents exist because the platform is implemented across two independent repositories with different responsibilities:

* **Ankhora Vault** — Desktop / Wails client and local cryptographic authority.
* **Ankhora Cloud** — Backend / cloud application, persistence, authorization, collaboration, and federation authority.

The repositories must remain independently buildable and internally decoupled while respecting these shared contracts.

---

## Repositories

```text
                    ANKHORA PLATFORM
                           │
              ┌────────────┴────────────┐
              │                         │
       ankhora-vault              ankhora-cloud
       Desktop / Wails             Cloud Backend
              │                         │
       Local Crypto Authority      Cloud Authority
       Local Keyring               Persistence
       Private Keys                Authorization
       Plaintext                   Collaboration
       DEK / KEK                   Federation
              │                         │
              └──────── CONTRACT ───────┘
```

---

## Core Principle

The cross-repository layer defines **contracts, boundaries, and invariants**.

It does not define implementation details belonging to either repository.

A Cloud engineer must not assume that a Desktop implementation can be imported or executed by Cloud.

A Desktop engineer must not assume that Cloud owns cryptographic secrets.

---

## Contract Documents

### `platform-boundaries.md`

Defines ownership boundaries between Desktop and Cloud.

### `vault-cloud-contract.md`

Defines the communication and application-level contract between Desktop and Cloud.

### `cryptographic-boundary.md`

Defines the Zero-Knowledge cryptographic boundary and explicitly identifies which cryptographic operations and materials are local versus cloud-visible.

### `collaboration-contract.md`

Defines the shared contract for TrustGroup collaboration, devices, key envelopes, ShareEntries, AssetCID, KEK versions, and collaborative assets.

### `federation-contract.md`

Defines the cross-cloud federation contract and the relationship between a local Cloud instance and a remote Cloud/Vault peer.

---

## Rules for AI Engineers

Before implementing a cross-repository feature:

1. Read this directory.
2. Identify which repository owns each responsibility.
3. Never move responsibility across the boundary merely because implementation would be easier.
4. Treat opaque cryptographic values as opaque.
5. Never introduce raw private keys, raw DEKs, or raw KEKs into Cloud contracts.
6. Never make a bounded context depend directly on another repository's internal implementation.
7. Prefer explicit application ports, DTOs, adapters, and events.
8. Update these contracts when a deliberate architectural boundary changes.
9. Record major boundary changes through an ADR in the owning repository.
10. Cross-repository documentation must describe the agreed platform contract, not temporary implementation details.

---

## Source of Truth

Repository-specific implementation is authoritative within its repository.

Cross-repository contracts are authoritative for:

* ownership;
* data crossing the repository boundary;
* security guarantees;
* cryptographic boundaries;
* identifiers and versions;
* collaboration semantics;
* federation semantics.

When implementation and a cross-repository contract disagree, the engineer must stop and resolve the architectural discrepancy rather than silently changing one side.

---

## Current Architectural Invariant

The most important invariant is:

> Ankhora Cloud may authorize, validate, persist, route, and federate encrypted data, but it must never require plaintext application data or private cryptographic material to perform those responsibilities.

The Desktop remains the cryptographic execution authority for user-controlled encryption and key operations.

# ADR-0006 — Desktop Vault Connection & Cloud Delegation Protocol

## Status

Accepted

## Date

2026-08-18

## Context

During Cloud workspace and channel integration, a fundamental distinction was established between **User Authentication** and **Vault Delegation**:

* Cloud user authentication (`SignIn`) verifies user credentials (email/password) and returns a Cloud JWT bearer token.
* Vault operation (`/api/workspaces`, `/api/channels`) requires an active `UserVaultIdentity` delegation record in Ankhora Cloud.
* A Desktop-created vault is sovereign and initially local; it is not automatically delegated in Cloud upon creation.

Attempting to perform vault operations with a valid user JWT on an un-delegated vault results in `403 Forbidden` (`vault not delegated`).

## Decision

1. **Strict Separation of Concerns**:
   - Authentication (`401` class concern) and Vault Authorization (`403` class concern) are separate independent trust boundaries.
   - A valid bearer token authorizes a User to call Cloud, but does NOT grant access to a specific Vault until delegation is established.

2. **Explicit ConnectVault Protocol**:
   - Vault registration/delegation is executed via an explicit `ConnectVault(userID, vaultID)` flow:
     - `POST /api/identity/challenge` with `vault_id` → receives `challenge_id` + `signing_payload`.
     - The Desktop signs `signing_payload` locally using the vault's Stellar Ed25519 private key.
     - `POST /api/identity/` with `{challenge_id, signature, vault_id, public_key}` → Cloud creates `identity_vaults` and `user_vault_identities` delegation.
   - `ConnectVault` MUST be executed prior to entering Ledger or making vault-scoped API calls.

3. **Key Isolation**:
   - The vault private key is owned exclusively by the Desktop and MUST NEVER be sent to Cloud.

4. **Infrastructure DTO Mapping**:
   - Cloud response structures (which use default PascalCase serialization) MUST be unmarshaled into dedicated infrastructure DTOs (e.g., `CloudWorkspaceDTO`) and explicitly mapped to application/domain models.
   - Direct leakage of wire DTOs or silent field drop (e.g. `VaultID`, `OwnerID`, `CreatedAt`, `UpdatedAt`) is forbidden.

5. **Responsibility Distinction**:
   - `ConnectVault` is a trust-establishment operation, not a vault-creation operation.
   - The conceptual model is: `CreateVault != RegisterVault != DelegateVault`. Desktop creates and owns local vault; `ConnectVault` proves key possession and registers caller delegation in Cloud.

## Consequences

* The Desktop cleanly establishes sovereign vault ownership with Cloud without compromising private key security.
* Authorization failures (`403`) are prevented before rendering the Ledger UI.
* Domain models remain decoupled from transport-layer serialization details.

## Architectural Constraint Notice — MUST HOLD

> **Do NOT modify these boundaries merely to make a test pass. Treat ADR-0006 and the client/cloud boundary rules as immutable architectural constraints. If implementation appears to conflict with them, stop and report the conflict before changing the architecture.**


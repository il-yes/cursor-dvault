# Session 2026-08-18 — Desktop Authentication, ConnectVault Delegation & Workspace DTO Integration

## Summary

Debugging and root-cause resolution of Desktop → Cloud vault authorization failures and workspace DTO field deserialization.

## Chronological Discovery Path

1. **Authentication Worked**:
   - `SignIn()` successfully logged into Ankhora Cloud via `/authenticate` and returned a 26-character JWT.

2. **Token Propagation Worked**:
   - Trace logs confirmed `TracecoreClient.SetToken()` executed and the JWT was supplied in the `Authorization: Bearer <token>` HTTP header.

3. **Duplicate Client Ruled Out**:
   - Trace logging of `TracecoreClient` pointer identities confirmed that all handlers shared the exact same client instance.

4. **Delegation Missing (`403 Forbidden`)**:
   - Cloud log inspection revealed the backend query `SELECT * FROM user_vault_identities WHERE user_id = 70 AND vault_id = 'fbffb9ad-d4de-4581-b7af-7b8e160f63bb' AND revoked_at IS NULL` returned zero rows.
   - Root cause: Cloud user authentication (`SignIn`) does NOT automatically delegate sovereign local vaults.

5. **ConnectVault Protocol Missing**:
   - The Desktop was missing the active execution of the vault challenge/registration protocol (`ConnectVault`).

6. **Protocol Implementation & Endpoint Alignment**:
   - Implemented `ConnectVault(userID, vaultID)`:
     - `POST /api/identity/challenge` (returns `challenge_id` + `signing_payload`).
     - Local Ed25519 signature of `signing_payload` using vault's Stellar private key.
     - `POST /api/identity/` (registers delegation in Cloud).
   - Wired `ConnectVault` into `SignIn()` immediately following `OpenVault()`.

7. **Delegation Succeeded**:
   - Cloud returned successful registration (`identity_vaults` and `user_vault_identities` created).
   - Subsequent `GET /api/workspaces?vault_id=...` returned `HTTP 200 OK`.

8. **Workspace DTO Field Loss Resolution**:
   - Inspection of workspace list logs showed that `ID`, `Name`, `Description`, and `Status` survived, but `VaultID`, `OwnerID`, `CreatedAt`, and `UpdatedAt` were empty/defaulted.
   - Root cause: Cloud serializes Workspace JSON with PascalCase field names (`VaultID`, `OwnerID`, `CreatedAt`, `UpdatedAt`).
   - Fix: Created `CloudWorkspaceDTO` in `internal/tracecore/types/workspace_cloud_dto.go` and explicit mapping `MapCloudWorkspaceToTypes()`.
   - Added unit test `internal/tracecore/workspace_dto_test.go` confirming full field survival across the infrastructure boundary.

## Decisive Result

* `SignIn()` → `ConnectVault()` → Ledger → `GET /workspaces` operates end-to-end with `200 OK`.
* Complete field preservation verified by regression unit tests.
* Architectural rules captured in ADR-0006 and documented in core principles.

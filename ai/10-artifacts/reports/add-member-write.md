Viewed app.go:3708-3865
Viewed app.go:3860-3920
Searched for "AddTrustGroupMember"
Searched for "AddTrustGroupMember"
Searched for "trustGroupId"
Viewed api.ts:1705-1735
Searched for "addTrustGroupMember"
Viewed TrustGroupDetailView.tsx:35-60
Searched for "AttachRuntime"
Viewed auth_handler.go:170-260
Viewed auth_handler.go:261-330
Viewed auth_handler.go:331-390
Viewed auth_handler.go:391-430
Searched for "func (s *KeyringService) LoadHybrid"
Searched for "LoadHybrid"
Searched for "type KeyringService"
Searched for "LoadHybrid"
Listed directory security
Viewed keyring_security.go:1-339
Edited app.go
Ran command: `/usr/local/go/bin/go test -count=1 . -run TestC3_FullEndToEndUserAcceptanceFlow -v`
Viewed app.go:3710-3730
Edited app.go
Ran command: `/usr/local/go/bin/go test -count=1 . -run TestC3_FullEndToEndUserAcceptanceFlow -v`

### Production Execution Trace & Failure Analysis for `App.AddTrustGroupMember`

#### 1. Exact Execution Call Chain

```text
Frontend UI (TrustGroupDetailView.tsx)
  │
  ├──> addTrustGroupMember(trustGroupId, memberId, "member") [api.ts]
  │       │
  │       └──> AppAPI.AddTrustGroupMember(jwtToken, trustGroupId, memberId, "member")
  │
  ▼
app.go: App.AddTrustGroupMember(JwtToken, trustGroupID, vaultID, role)
  │
  ├──> 1. RequireAuth(JwtToken)
  │       ├── Decodes claims.UserID (caller's vault ID, e.g. Alice)
  │       └── Returns unauthorized error if token is invalid
  │
  ├──> 2. Cloud Member Addition (First Cloud Mutation)
  │       └── addTrustGroupMemberUC.Execute(ctx, AddMemberToTrustGroupRequest{...})
  │             └── HTTP POST /api/trustgroups/{trustGroupID}/members
  │             └── Cloud adds vaultID to trustGroup.member_cids and returns HTTP 200
  │             └── STATUS: SUCCESS ✅ (Observed: Cloud mutation persists Bob in Cloud DB)
  │
  ├──> 3. Active Device Resolution for Invitee (Bob)
  │       └── identityDeviceAdapter.ListActiveDevices(ctx, vaultID)
  │             └── SELECT * FROM identity_devices WHERE vault_id = vaultID AND status = 'active'
  │             └── STATUS EVALUATION:
  │                   ├── IF devices == 0 ──> Branch B (Fallback Resolution)
  │                   └── IF devices > 0  ──> Branch A (Envelope Provisioning)
  │
  ├──> 4A. Branch A (len(devices) > 0)
  │       ├── a. Session Secrets Lookup:
  │       │      a.Vault.GetSession(claims.UserID)
  │       │      Extracts pass = SessionSecrets["password"], secret = SessionSecrets["device_seed"]
  │       ├── b. Keyring Unwrapping:
  │       │      a.Vault.KeyringService.LoadHybrid(claims.UserID, pass, secret)
  │       │      └── Fails if session pass is empty ("failed to unlock keyring")
  │       └── c. Envelope Provisioning & Cloud Sync:
  │              a.provisionEnvelopeUC.Execute(ctx, req, keyring)
  │              └── Wraps Group KEK v1 for Bob's active device public key
  │              └── Calls AddTrustGroupKeyEnvelopeUseCase -> HTTP PUT /api/trustgroups/{id}
  │
  └──> 4B. Branch B (len(devices) == 0) [Fallback Path]
          ├── a. Identity.FindUserById(ctx, vaultID) & tracecoreClient.GetUserByEmail(ctx, email)
          └── b. If targetPubKey == "":
                 Returns Error: "no active device registered for member <vaultID>" ❌
```

---

#### 2. Identification of the First Failing Operation

The observed UI anomaly (*Cloud mutation succeeds → UI displays error → reload shows Bob present*) occurs because **`App.AddTrustGroupMember` is a 2-phase non-atomic operation**:

1. **Phase 1 (Cloud Mutation)**: `addTrustGroupMemberUC.Execute(...)` runs first and immediately mutates state on Ankhora Cloud backend (`POST /api/trustgroups/{id}/members`). This request returns HTTP 200 and permanently persists Bob's `vaultID` in Cloud `member_cids`.
2. **Phase 2 (Client-Side Provisioning)**: Executed immediately after Phase 1 succeeds in the same function:
   - **Failure Point 1 (`ListActiveDevices`)**: If Bob's account was created prior to atomic onboarding fixes, `ListActiveDevices(BobVaultID)` returns `0` active devices. The function enters Branch B, fails to resolve a public key, and returns:
     ```text
     "no active device registered for member <BobVaultID>"
     ```
   - **Failure Point 2 (`LoadHybrid`)**: If Alice's in-memory session `VaultRuntimeContext.SessionSecrets["password"]` is unpopulated or lost upon app restart, `LoadHybrid(AliceVaultID, "", "")` returns:
     ```text
     "failed to load caller keyring <AliceVaultID>: failed to unlock keyring: tried 1 wrappers, none succeeded"
     ```

Because Phase 1 already succeeded on Cloud, returning an error from Phase 2 causes Wails to return a rejected promise to the UI. The UI catches the rejection and displays an error toast. However, upon app reload/re-login, the UI fetches the TrustGroup from Cloud (`GET /api/trustgroups/{id}`), which already contains Bob from Phase 1.

---

#### 3. State Comparison: Real UI Session vs. `TestC3_FullEndToEndUserAcceptanceFlow`

| Phase / State Component | `TestC3_FullEndToEndUserAcceptanceFlow` State | Real UI Runtime / Session State |
| :--- | :--- | :--- |
| **Invitee Device State (`ListActiveDevices`)** | Bob is created via `CreateAccountUseCase` with atomic `OnCreateDevice`. `IdentityDevice` exists in local SQLite DB with `vault_id = BobVaultID` and `status = 'active'`. `ListActiveDevices` returns `1` device. | In real UI accounts onboarded before atomic provisioning, `identity_devices` table had 0 rows for `BobVaultID`. `ListActiveDevices` returns `0` devices, triggering Phase 2 error. |
| **Caller Session Credentials (`LoadHybrid`)** | `AuthHandler.Login` starts session and populates `VaultRuntimeContext.SessionSecrets["password"]`. Test preserves `sessionMgr` runtime so `GetSession(AliceVaultID)` returns valid password. | If Alice's session was restored from disk or initialized without `password` in `SessionSecrets`, `GetSession(AliceVaultID)` returns empty password. `LoadHybrid` fails, triggering Phase 2 error. |
| **Phase 1 Cloud Mutation** | Invokes `addTrustGroupMemberUC.Execute(...)`. Cloud updates `member_cids`. | Invokes `addTrustGroupMemberUC.Execute(...)`. Cloud updates `member_cids` (Succeeds). |
| **Phase 2 Envelope Provisioning** | `provisionEnvelopeUC.Execute(...)` wraps `KEK v1` using Bob's device public key and posts envelope to Cloud. Returns `*TrustGroup, nil`. | Fails at Step 3 or Step 4 due to missing device or missing session password. Returns `nil, error` to Wails frontend. |
| **UI Outcome** | Test asserts `require.NoError(t, err)`. Complete pass. | Wails Promise rejects -> UI displays Error Toast -> Cloud already has Bob -> App reload shows Bob present. |
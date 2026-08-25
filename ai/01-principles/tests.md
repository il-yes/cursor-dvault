> **A passing unit/integration test is not evidence that the real desktop → HTTP → Cloud → DB path works.**

We were testing pieces of the architecture while the actual boundary had several mismatches:

* `/api/api/...` URL construction
* `identity_users` vs Cloud's actual identity storage
* empty `public_key` in issued tokens
* `vault_id` wire-format mismatch
* real delegation state in MySQL vs test fixtures

None of those necessarily showed up in isolated tests.

### So from now on, I'd use this verification hierarchy

**1. Domain/unit tests**

> Does the component behave correctly in isolation?

Useful, but **not sufficient**.

**2. Repository/integration tests**

> Does the component work against its actual persistence mechanism?

Still useful, but not sufficient for cross-BC behavior.

**3. Wire-contract tests**

> Does the desktop actually serialize exactly what Cloud expects?

This catches things like `VaultID` vs `vault_id`.

**4. Live boundary tests**

> Does the real desktop client communicate with the real running Cloud?

This is the critical layer for Ankhora.

```text
Desktop
   │
   │ real HTTP
   ▼
TracecoreClient
   │
   ▼
Cloud API
   │
   ├── Identity
   ├── Federation / Delegation
   ├── Workspace
   └── ...
   │
   ▼
Real MySQL
```

**5. Manual smoke test**

> Can I actually perform the operation from the application?

This is the final sanity check.

And now we have an extremely valuable pattern:

```text
TEST PASS
   ↓
WIRE CONTRACT PASS
   ↓
LIVE CLOUD BOUNDARY PASS
   ↓
MANUAL UI PASS
   ↓
CONFIDENT FEATURE
```

That's much stronger than saying *"I have 30 tests and therefore it's done."*

### For the current feature

We now have something materially different:

**Stellar authentication**
→ works against real Cloud

**Bearer token**
→ correctly identifies the Stellar user

**Vault delegation**
→ exists and is resolved by Cloud

**Workspace request**
→ correct `vault_id` crosses the HTTP boundary

**Cloud authorization**
→ accepts the delegated vault

**Workspace persistence**
→ real MySQL row created

**Desktop**
→ successfully creates the workspace without error

That's a **real vertical slice**, not just a test artifact.

And I would absolutely keep the new `live_cloud_boundary_test.go` methodology. For every important cross-boundary feature, we should add a focused live acceptance test rather than accumulating dozens of mocked tests that can give us false confidence.

**The next time an agent tells you "all tests pass," the question should be:**

> **"Show me the live boundary test and the actual HTTP request/response."**

That should become our standard for Ankhora.

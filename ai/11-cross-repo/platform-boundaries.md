# Platform Boundaries

## Purpose

This document defines ownership between Ankhora Vault Desktop and Ankhora Cloud.

The objective is to prevent responsibility leakage between repositories and preserve the Zero-Knowledge architecture.

---

# 1. Ankhora Vault — Desktop Authority

Ankhora Vault is the local execution environment for user-controlled cryptographic operations.

It owns:

* local Vault state;
* local VaultKeyring;
* private keys;
* unlocked cryptographic material;
* plaintext application data;
* DEK generation;
* KEK generation;
* DEK wrapping and unwrapping;
* KEK wrapping and unwrapping;
* AES-256-GCM encryption/decryption;
* asymmetric key operations;
* local cryptographic orchestration;
* preparation of encrypted payloads;
* preparation of opaque key envelopes.

The Desktop may temporarily hold raw cryptographic material in memory when performing authorized local operations.

Raw cryptographic material must never cross into the Cloud application boundary.

---

# 2. Ankhora Cloud — Cloud Authority

Ankhora Cloud owns:

* authentication;
* authorization;
* user/device identity metadata;
* TrustGroup state;
* TrustGroup membership;
* Device authorization state;
* ShareEntry metadata;
* Asset metadata;
* encrypted payload persistence;
* opaque key-envelope persistence;
* KEK version metadata;
* collaboration orchestration;
* transactional consistency;
* idempotency;
* concurrency control;
* auditability;
* federation;
* remote-peer trust;
* message validation;
* outbound federation signing;
* delivery state;
* retry and acknowledgement handling.

Cloud is a **policy and coordination authority**, not a cryptographic secret authority.

---

# 3. Responsibilities That Must Not Move to Cloud

The Cloud must not:

* generate user DEKs;
* generate user KEKs;
* decrypt application payloads;
* unwrap user DEKs;
* unwrap TrustGroup KEKs;
* access user private keys;
* derive private keys;
* receive raw private keys;
* require plaintext asset payloads;
* inspect encrypted payload contents;
* replace local cryptographic orchestration.

The Cloud may validate metadata associated with these operations.

---

# 4. Responsibilities That Must Not Move to Desktop

The Desktop must not become the authoritative source for:

* global TrustGroup membership;
* authoritative Device status;
* Cloud persistence;
* federation state;
* remote peer state;
* server-side authorization decisions;
* transaction commit semantics;
* idempotency guarantees across network retries;
* server-side concurrency control.

The Desktop can request and consume these decisions, but Cloud remains authoritative.

---

# 5. Communication Model

The current Desktop application communicates with the backend through the application's direct Wails/backend integration layer rather than assuming a browser-style HTTP API between the two local layers.

This implementation detail must not alter the architectural boundary.

Conceptually:

```text
Desktop Application Layer
        │
        │ Application Contract
        ▼
Cloud Application Boundary
```

The transport mechanism is an implementation concern.

The contract remains explicit.

---

# 6. Cryptographic Data Boundary

Values crossing from Desktop to Cloud may include:

* encrypted payloads;
* AssetCID;
* WrappedDEK;
* WrappedKEK;
* DeviceID;
* MemberID;
* TrustGroupID;
* KEKVersion;
* cryptographic metadata required for validation.

Values that must not cross the boundary:

* plaintext payload;
* raw DEK;
* raw KEK;
* private key;
* private key seed;
* decrypted VaultKeyring;
* user secrets used to unlock local key material.

---

# 7. Authority Model

```text
                  ┌──────────────────────┐
                  │      DESKTOP         │
                  │                      │
                  │ Crypto Authority     │
                  │ Plaintext Authority  │
                  │ Keyring Authority    │
                  └──────────┬───────────┘
                             │
                       opaque contract
                             │
                  ┌──────────▼───────────┐
                  │       CLOUD          │
                  │                      │
                  │ Policy Authority     │
                  │ Persistence Authority│
                  │ Collaboration        │
                  │ Federation Authority │
                  └──────────────────────┘
```

Neither side should silently assume the other's authority.

---

# 8. Architectural Decision Rule

When a new feature crosses the repositories, ask:

1. Who owns the state?
2. Who owns the secret?
3. Who owns the policy?
4. Who owns persistence?
5. Who owns the transaction?
6. Who owns the cryptographic operation?
7. What exact data crosses the boundary?

If these questions cannot be answered explicitly, the feature is not ready for implementation.

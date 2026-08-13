# Cryptographic Boundary

## Purpose

This document defines the cryptographic Zero-Knowledge boundary between Ankhora Vault Desktop and Ankhora Cloud.

The fundamental rule is:

> Cryptographic execution happens locally. Cryptographic authorization and opaque cryptographic state may be coordinated by Cloud.

---

# 1. Local Cryptographic Authority

Ankhora Vault Desktop owns cryptographic execution.

This includes:

* key generation;
* key derivation;
* encryption;
* decryption;
* key wrapping;
* key unwrapping;
* private-key operations;
* VaultKeyring access;
* plaintext access.

The Desktop may use existing cryptographic primitives such as:

* AES-256-GCM;
* Curve25519/Ed25519-compatible key wrapping;
* cryptographically secure random generation.

The exact primitive is implementation-owned by the Desktop repository.

---

# 2. Cloud Cryptographic Responsibility

Cloud may store and validate opaque cryptographic material.

For example:

```text
WrappedKEK
WrappedDEK
EncryptedPayload
AssetCID
KEKVersion
```

Cloud may validate:

* envelope ownership;
* Device existence;
* Device status;
* TrustGroup membership;
* Member ownership;
* KEK version;
* duplicate envelopes;
* stale versions;
* authorization;
* lifecycle state.

Cloud must not attempt to interpret the cryptographic contents of opaque values.

---

# 3. Forbidden Cloud Data

The following must never be required by a Cloud use case:

```text
RawPrivateKey
PrivateKeySeed
RawDEK
RawKEK
PlaintextAsset
UnlockedVaultKeyring
UserPassword
StellarSecret
```

These values must never be persisted in Cloud repositories or emitted through Cloud application events.

---

# 4. Collaborative Asset Encryption

The collaborative asset flow is:

```text
Plaintext
   │
   │ Desktop
   ▼
DEK
   │
   ▼
AES-256-GCM
   │
   ▼
Encrypted Payload
   │
   ▼
AssetCID
```

The DEK is protected by the TrustGroup KEK:

```text
DEK
 │
 └── KEK ──► WrappedDEK
```

The TrustGroup KEK is protected for each authorized device:

```text
KEK
 │
 ├── Device A PublicKey ──► WrappedKEK A
 ├── Device B PublicKey ──► WrappedKEK B
 └── Device C PublicKey ──► WrappedKEK C
```

Cloud stores the resulting opaque values.

---

# 5. AssetCID Invariant

For collaborative assets:

> AssetCID identifies the encrypted payload, not plaintext.

Therefore:

```text
Plaintext
   ↓
DEK
   ↓
Encrypted Payload
   ↓
AssetCID
```

Key rotation must not require re-encrypting the payload.

If only the key wrapping layer changes:

```text
Encrypted Payload = unchanged
AssetCID          = unchanged
DEK               = unchanged
WrappedDEK        = changed
WrappedKEK        = changed
```

---

# 6. KEK Rotation

When a TrustGroup KEK changes from version N to N+1:

```text
KEK N
 │
 ├── unwrap DEK
 │
 ▼
DEK
 │
 └── wrap with KEK N+1
        │
        ▼
     WrappedDEK N+1
```

The encrypted payload does not change.

The new KEK is then wrapped for remaining authorized devices:

```text
KEK N+1
 │
 ├── Device B → WrappedKEK N+1
 ├── Device C → WrappedKEK N+1
 └── Device D → WrappedKEK N+1
```

Removed or revoked devices receive no new envelope.

---

# 7. Cloud Must Remain Cryptographically Blind

Cloud may know:

```text
DeviceID = D1
TrustGroupID = TG1
KEKVersion = 2
WrappedKEK = opaque value
WrappedDEK = opaque value
AssetCID = CID
```

Cloud must not know:

```text
KEK bytes
DEK bytes
private key
plaintext
```

This distinction is a core security invariant.

---

# 8. Logging Rule

Raw cryptographic material must never be logged.

This includes:

* private keys;
* private key seeds;
* raw DEKs;
* raw KEKs;
* decrypted payloads;
* decrypted VaultKeyring contents.

Logs may contain identifiers and safe metadata such as:

```text
TrustGroupID
DeviceID
MemberID
KEKVersion
AssetCID
RequestID
operation status
```

---

# 9. Boundary Test Requirement

Every new cryptographic cross-repository operation must include a test proving that the Cloud-facing DTO contains no raw cryptographic secret.

At minimum:

```text
assert DTO contains no raw KEK
assert DTO contains no raw DEK
assert DTO contains no private key
assert DTO contains no plaintext payload
```

The strongest test is an end-to-end local round-trip:

```text
Plaintext
 → Encrypt
 → CID
 → WrappedDEK
 → WrappedKEK
 → Device private key
 → KEK
 → DEK
 → Plaintext
```

while Cloud only observes opaque values.

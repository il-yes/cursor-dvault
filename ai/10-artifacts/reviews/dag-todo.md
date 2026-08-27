Yes — and your first question is important.

### 1. How did I know the commit message?

Because **you explicitly included it in the message you pasted**:

> `Ran command: gcmsg "fix: session boudary - attachments"`

So I didn't infer it from Git, inspect your repository, or retrieve it from somewhere else. I simply read it from the command output you provided.

And I should distinguish that clearly when discussing your repo: **“you showed me the commit message”**, not “I know your commit is…”.

---

# 2. Yes — the terminology should absolutely go into the AI architecture docs

This is more than documentation polish. You've just uncovered an important architectural distinction that will prevent a lot of confusion as the platform grows.

The terminology should become **canonical vocabulary for the project**.

I'd formalize it roughly like this:

## Vault Cryptographic Key Model

### 1. Master Password

**Definition:**
The user's authentication/unlock credential.

**Role:**
It is **not the payload encryption key**.

It is used to unlock the persistent encrypted keyring:

```text
Master Password
      ↓
KDF / key derivation
      ↓
unwrap Stored Keyring
      ↓
Vault DEK
```

Therefore:

> **Master Password = credential used to unlock key material.**

It should normally be required at the **authentication/unlock boundary**, not for every vault operation.

---

### 2. Keyring Container

The persistent encrypted container:

```text
data/keyrings/<UserOnboardingID>.json
```

Its responsibility is to protect persistent cryptographic material.

Conceptually:

```text
Encrypted Keyring Container
        │
        └── contains Vault DEK
```

The keyring is therefore **key storage**, not the vault data encryption mechanism itself.

---

### 3. Vault DEK

This is the really important one.

**Vault DEK = Data Encryption Key**

It is the actual **32-byte symmetric key** used by AES-GCM to encrypt/decrypt private vault payloads.

```text
Vault DEK
   ↓
AES-GCM
   ↓
Encrypted vault payload
```

So the precise terminology should be:

> **Vault DEK is the actual private vault payload encryption key.**

This is the key your recent refactor exposed.

---

### 4. Session VaultKey

`Session.VaultKey []byte`

This is **not a second cryptographic key**.

It is the **in-memory runtime representation of the Vault DEK after the vault has been unlocked**.

That's an important distinction.

```text
Persistent world                    Runtime world

Master Password
       │
       ▼
Encrypted Keyring
       │
       ▼
   Vault DEK ──────────────────────► Session.VaultKey
                                      │
                                      ▼
                                   AES-GCM
```

Therefore:

> **Session VaultKey = unlocked in-memory Vault DEK.**

And your new tags:

```go
VaultKey []byte `json:"-" gorm:"-"`
```

make the architectural boundary explicit:

**session-only / ephemeral / non-persistent.**

---

### 5. Share / Collaboration Key Material

This needs to remain conceptually separate.

Your C3 collaboration architecture has another cryptographic layer involving:

* Trust Groups
* KEK versions
* key envelopes
* DEKs
* wrapped DEKs
* device keys

So I would **not call all of that “ShareKey” generically** in the architecture documentation.

I'd establish a vocabulary such as:

```text
PRIVATE VAULT CRYPTOGRAPHY

Master Password
      │
      ▼
Keyring Container
      │
      ▼
Vault DEK
      │
      ▼
Session VaultKey
      │
      ▼
AES-GCM
      │
      ▼
Private Vault Payload


COLLABORATIVE CRYPTOGRAPHY

Device Key
      │
      ▼
Trust Group Key / KEK
      │
      ▼
Wrapped DEK / Key Envelope
      │
      ▼
Collaborative Asset DEK
      │
      ▼
AES-GCM
      │
      ▼
Collaborative Payload
```

That separation is **very valuable** because otherwise six months from now an AI assistant could easily confuse:

> Vault DEK
> TrustGroup KEK
> Asset DEK
> Session VaultKey
> Master Password

and make a disastrous architectural change.

---

# The most important invariant

I'd put this in the AI documentation almost verbatim:

> **The Master Password does not encrypt vault payloads. It unlocks the persistent keyring from which the Vault DEK is obtained. The Vault DEK is the actual symmetric encryption key used for private vault payload encryption. Once unlocked, the Vault DEK may be retained as `Session.VaultKey` for the lifetime of the authenticated session, eliminating repeated password/keyring-unlock operations.**

And another invariant:

> **Password handling belongs to the authentication/unlock boundary. Payload encryption belongs to the Vault DEK boundary. These concerns must not be coupled.**

That's probably the architectural insight you were looking for when you initially noticed the legacy `"password"` requirement.

Your latest test is particularly strong evidence:

```text
Password = ""
UnlockVaultHandler = nil
Session.VaultKey = 32 bytes
        ↓
AES-GCM encryption
        ↓
SUCCESS
```

That isn't merely a convenience optimization. It demonstrates that **the payload encryption architecture does not intrinsically require the master password**.

The password was previously being carried too far down the stack because the old implementation was repeatedly reopening the keyring.

---

## I would also add one naming rule

For the AI docs:

| Term                        | Never confuse with   |
| --------------------------- | -------------------- |
| **Master Password**         | Vault DEK            |
| **Keyring Container**       | Vault itself         |
| **Vault DEK**               | Master Password      |
| **Session VaultKey**        | A new/generated key  |
| **Trust Group KEK**         | Vault DEK            |
| **Collaborative Asset DEK** | Private Vault DEK    |
| **Key Envelope**            | Actual plaintext key |

And specifically:

> **`Session.VaultKey` is an alias/reference to the unlocked Vault DEK, not an independently generated session encryption key.**

That's an excellent canonical distinction for your AI assistants.

If you're having your coding assistant modify the AI documentation now, I would have it create a dedicated **“Cryptographic Terminology & Key Hierarchy”** section rather than burying this inside the Vault README. This terminology is foundational enough that **every future BC touching encryption should be able to reference it.**

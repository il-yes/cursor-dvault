
# Domain-Driven Design Standards

## Purpose

This document defines the implementation standards for applying Domain-Driven Design in Ankhora.

The objective is to maintain:

- clear ownership
- strong boundaries
- independent evolution
- understandable business logic

---

# Core Rule

Code organization follows business ownership.

The primary question is:

> "Who owns this concept?"

Not:

> "Where is this technically convenient?"

---

# Bounded Context Structure

Every bounded context follows:

```

context/

├── domain/
├── application/
├── infrastructure/
└── interface/

```

Example:

```

vault/

├── domain/
├── application/
├── infrastructure/
└── interface/

```

---

# Domain Layer Standards

The domain layer contains:

- entities
- aggregates
- value objects
- domain services
- domain events
- business rules

---

The domain layer must NOT contain:

- HTTP handlers
- database code
- JSON transport logic
- framework dependencies

---

Example:

Good:

```

Asset.Share()

validates sharing rules

```

Bad:

```

AssetRepository.Save()

inside domain object

```

---

# Entity Standards

Entities:

- have identity
- contain behavior
- protect invariants

Example:

```

Asset

ID

Owner

EncryptionState

SharingRules

```

---

Avoid:

Anemic models:

```

Asset struct

with only fields

```

and all logic elsewhere.

---

# Value Object Standards

Value Objects represent concepts without identity.

Examples:

```

CID

Hash

Email

Version

EncryptionMetadata

```

Rules:

- immutable
- validated at creation
- compared by value

---

# Aggregate Standards

Aggregates protect consistency boundaries.

Rules:

- one aggregate root
- internal state controlled by root
- external code cannot directly modify internals

Example:

```

Workspace

```
|
+-- Channels

+-- Members
```

```

Workspace controls changes.

---

# Repository Standards

Repositories belong to domain abstractions.

Example:

Domain:

```

type AssetRepository interface {

Save(asset Asset)

Find(id ID)

}

```

Infrastructure:

```

PostgresAssetRepository

```

---

The domain knows WHAT.

Infrastructure decides HOW.

---

# Use Case Standards

Application layer contains use cases.

Example:

```

CreateAssetUseCase

ShareAssetUseCase

RestoreVaultUseCase

```

Use cases:

- coordinate actions
- call domain behavior
- publish events

They do not contain database logic.

---

# Domain Event Standards

Events represent facts.

Good:

```

AssetCreated

TrustEstablished

CommitValidated

```

Bad:

```

CreateAssetExecuted

```

Events must:

- be immutable
- contain meaningful data
- have clear ownership

---

# Context Communication

Contexts communicate through:

- events
- interfaces
- APIs

Never:

- shared database tables
- importing another context's domain package

---

# Dependency Rules

Allowed:

```

Interface

↓

Application

↓

Domain

```

Infrastructure connects from outside.

---

Forbidden:

```

Domain

↓

Database

↓

HTTP

```

---

# Naming Standards

Use business language.

Good:

```

TrustGroup

VaultAsset

ValidationCommit

```

Avoid technical names:

```

DataObject

Manager

HandlerService

```

---

# AI Review Checklist

Before creating code, AI should verify:

- [ ] Correct bounded context
- [ ] Domain ownership clear
- [ ] Business rules inside domain
- [ ] Infrastructure isolated
- [ ] Events considered
- [ ] No cross-context leakage

---

# Final Principle

DDD is not a folder structure.

DDD is ownership made explicit in code.

A good Ankhora implementation makes it obvious:

- who owns a concept,
- who can change it,
- how it communicates,
- why it exists.
```

---


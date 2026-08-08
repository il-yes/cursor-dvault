
# Ankhora Glossary

## Purpose

This glossary defines the official terminology used throughout the Ankhora platform.

The objective is to create a shared language between:

- engineers
- architects
- domain experts
- AI assistants

Terms defined here have a specific Ankhora meaning.

AI assistants must prefer these definitions over generic interpretations.

---

# Core Concepts

## Ankhora

Ankhora is a sovereign digital trust platform designed to protect, manage, collaborate on, and verify digital information throughout its lifecycle.

Core capabilities:

- secure ownership
- encrypted storage
- trusted collaboration
- lifecycle tracking
- federation
- verification

---

# Vault

## Definition

A Vault is the secure ownership layer of Ankhora.

It manages:

- encrypted data
- ownership
- keys
- assets
- protected metadata

A Vault is not simply storage.

It represents:

```

Ownership

*

Protection

*

Control

```

---

# Vault Engine Desktop

## Definition

The Vault Engine Desktop is the local execution environment responsible for:

- local vault operations
- encryption
- secure user interaction
- offline capabilities
- desktop resource management

It represents the user's sovereign environment.

---

# Vault Cloud Service

## Definition

The Vault Cloud Service provides cloud capabilities around Vault without replacing ownership.

Responsibilities:

- synchronization
- availability
- collaboration support
- remote operations

It must preserve zero-knowledge principles.

---

# Asset

## Definition

An Asset is a protected unit of information managed by a Vault.

An Asset may represent:

- documents
- files
- credentials
- identities
- structured data
- domain information

An Asset contains:

- ownership information
- metadata
- protection state

An Asset is not simply a file.

---

# Identity

## Definition

Identity represents the authenticated actor participating in the system.

Identity answers:

```

Who are you?

```

Identity does not answer:

```

What are you allowed to do?

```

Authorization remains separate.

---

# Trust

## Definition

Trust represents the confidence relationship between entities.

Trust is established through:

- identity verification
- cryptographic validation
- policies
- permissions

Trust is not assumed.

---

# Trust Group

## Definition

A Trust Group defines a controlled group of identities allowed to collaborate around protected resources.

Responsibilities:

- membership
- access relationships
- sharing boundaries

A Trust Group is not a generic user group.

---

# C3

## Definition

C3 is the collaborative communication and coordination layer of Ankhora.

It manages:

- channels
- threads
- collaboration flows
- shared activities

C3 enables interaction without owning the underlying assets.

---

# Channel

## Definition

A Channel is a collaboration space where participants exchange information.

A Channel belongs to C3.

It does not own:

- encryption keys
- asset lifecycle
- historical validation

---

# Thread

## Definition

A Thread represents a collaboration lifecycle around a topic, asset, or activity.

A Thread has:

- identity
- participants
- state
- lifecycle

A Thread is not a programming execution thread.

---

# TraceCore

## Definition

TraceCore is the lifecycle memory and verification engine of Ankhora.

It manages:

- commits
- history
- validation
- branches
- workflows
- audit trails

TraceCore answers:

```

What happened?

When?

By whom?

Was it valid?

```

---

# Commit

## Definition

A Commit is an immutable recorded state change inside TraceCore.

A Commit represents:

- a validated change
- a historical checkpoint
- an auditable event

A Commit is inspired by Git but is domain-oriented.

---

# Event

## Definition

An Event represents something meaningful that happened.

Examples:

```

AssetCreated

TrustEstablished

CommitValidated

```

Events are:

- immutable
- owned
- observable

---

# Domain Application

## Definition

A Domain Application represents a specialized business capability built on top of Ankhora.

Examples:

- construction workflows
- pharmaceutical traceability
- supply chain processes
- enterprise operations

Domain applications own business rules.

They do not replace the platform foundations.

---

# Federation

## Definition

Federation enables trusted interaction between independent Ankhora environments.

Responsibilities:

- trust resolution
- message validation
- secure exchange
- remote synchronization

Federation does not own domain data.

---

# Subscription

## Definition

Subscription manages commercial access to Ankhora capabilities.

Responsibilities:

- plans
- tiers
- billing lifecycle
- entitlements

Subscription does not define technical permissions.

---

# Bounded Context

## Definition

A bounded context is an independent business capability with:

- its own model
- ownership
- rules
- language

Examples:

```

Vault

C3

TraceCore

Identity

```

---

# Aggregate

## Definition

An Aggregate is a consistency boundary inside a domain.

It protects:

- invariants
- lifecycle rules
- internal state

---

# Domain Event

## Definition

A Domain Event represents a business fact produced by a bounded context.

Example:

```

AssetShared

```

The owning context creates it.

Other contexts may react.

---

# Zero-Knowledge

## Definition

Zero-knowledge means the system minimizes knowledge outside the owner's control.

The platform should not require access to user data to provide functionality.

---

# Sovereignty

## Definition

Sovereignty means maintaining ownership and control over digital information.

A sovereign system preserves:

- ownership
- portability
- control
- verification

---

# AI Rule

When an AI assistant encounters an Ankhora term:

1. Check this glossary first.
2. Use the Ankhora definition.
3. Do not replace it with generic industry assumptions.

---

# Final Principle

A shared language creates a shared understanding.

The glossary is the semantic foundation of the Ankhora engineering system.
```

---

After adding this file, your AI folder reaches:

```text
AI Engineering Platform

Vision              ✅
Principles          ✅
Architecture        ✅
Standards           ✅
Contexts            ✅
Workflows           ✅
AI Roles            ✅
Decisions           ✅
Glossary            ✅
```


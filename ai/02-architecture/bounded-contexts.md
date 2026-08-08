# Ankhora Bounded Contexts

## Purpose

This document defines the bounded contexts of the Ankhora platform.

Each bounded context represents a specific business capability with:

- its own responsibility
- its own domain model
- its own rules
- its own lifecycle
- its own language

A bounded context is an ownership boundary.

---

# Bounded Context Philosophy

Ankhora is not designed as one large application.

It is designed as a collaboration of specialized domains.

Each context:

- owns its concepts
- protects its invariants
- exposes explicit contracts
- communicates through controlled interactions

No context should become a "god module."

---

# Context Overview

Ankhora currently contains:
Identity

Onboarding

Vault

ShareEntries

C3 Collaboration

TraceCore

Federation

Subscription

Domain Applications


Each context has a distinct purpose.

---

# Identity Context

## Mission

Identity establishes trusted participants in the system.

Identity answers:

> "Who is this entity?"

---

# Responsibilities

Identity manages:

- user identity
- authentication
- identity verification
- cryptographic identity
- trust bootstrap

---

# Core Concepts

## Identity

Represents a verified participant.

Possible attributes:

- identifier
- public information
- verification status
- credentials

---

## Trust Identity

Represents the ability to establish trust with another entity.

---

# Identity Owns

- identity lifecycle
- authentication mechanisms
- identity verification

---

# Identity Does Not Own

Identity does not own:

- business permissions
- files
- workflows
- collaboration rules

Authorization decisions belong to the context owning the resource.

---

# Vault Context

## Mission

Vault provides secure ownership and protection of information.

Vault answers:

> "How do we protect and manage owned information?"

---

# Responsibilities

Vault manages:

- encrypted storage
- assets
- keys
- secure sharing
- data lifecycle

---

# Core Concepts

## Vault

Root container of protected information.

---

## Asset

A protected object.

Examples:

- documents
- files
- encrypted objects

---

## Share Entry

Represents controlled sharing of an asset.

---

## Encryption Keys

Protect ownership and confidentiality.

---

# Vault Owns

- asset protection
- encryption lifecycle
- secure storage

---

# Vault Does Not Own

Vault does not own:

- business workflows
- approval processes
- domain rules
- collaboration semantics

---

# C3 Collaboration Context

## Mission

C3 enables secure collaboration between trusted participants.

C3 answers:

> "How do participants collaborate?"

---

# Responsibilities

C3 manages:

- workspaces
- channels
- threads
- trust groups
- collaboration assets

---

# Core Concepts

## Workspace

A collaboration environment.

---

## Channel

A communication and organization space.

---

## Thread

A focused collaboration flow.

---

## Trust Group

A group of participants sharing controlled access.

---

# C3 Owns

- collaboration structure
- collaboration lifecycle
- participant relationships inside collaboration

---

# C3 Does Not Own

C3 does not own:

- user identity creation
- encryption implementation
- business workflows

---

# TraceCore Context

## Mission

TraceCore provides verifiable lifecycle management.

TraceCore answers:

> "What happened and why?"

---

# Responsibilities

TraceCore manages:

- commits
- versions
- branches
- workflows
- validation
- approvals
- audit history

---

# Core Concepts

## Commit

A recorded change in lifecycle history.

---

## Repository

A controlled lifecycle space.

---

## Workflow

A defined process with validation rules.

---

## Validation Rule

A rule used to verify correctness.

---

# TraceCore Owns

- lifecycle state
- history
- validation
- compliance evidence

---

# TraceCore Does Not Own

TraceCore does not own:

- encrypted storage
- identity management
- collaboration interfaces

---

# Federation Context

## Mission

Federation enables trusted communication between independent environments.

Federation answers:

> "Can two independent systems collaborate safely?"

---

# Responsibilities

Federation manages:

- remote trust relationships
- message validation
- cryptographic verification
- synchronization agreements

---

# Core Concepts

## Remote Vault

A trusted external vault relationship.

---

## Federation Message

A validated communication envelope.

---

## Trust Relationship

An established relationship between independent entities.

---

# Federation Owns

- trust resolution
- message validation
- federation policies

---

# Federation Does Not Own

Federation does not own:

- transport infrastructure
- business objects
- collaboration rules

---

# Subscription Context

## Mission

Subscription manages commercial access to platform capabilities.

---

# Responsibilities

Subscription manages:

- plans
- tiers
- billing lifecycle
- entitlement rules

---

# Subscription Owns

- subscription state
- pricing rules
- feature access rules

---

# Subscription Does Not Own

Subscription does not own:

- identity
- payment infrastructure
- business domain data

---

# Domain Application Contexts

## Mission

Domain applications implement specific business capabilities.

Examples:

- construction management
- pharmaceutical workflows
- supply chain systems

---

# Responsibilities

Domain applications own:

- business rules
- domain entities
- operational processes

---

# Domain Applications Use

Platform capabilities:

- Identity
- Vault
- C3
- TraceCore
- Federation

---

# Context Interaction Rules

When a new feature is designed:

Ask:

## Ownership

Which context owns this behavior?

---

## Data

Which context owns this information?

---

## Rules

Where does the business rule live?

---

## Communication

Should this interaction be:

- API call?
- domain event?
- message?

---

# Anti-Patterns

## Shared Database Model

Wrong:

Every context accesses every table


---

## Shared Domain Objects

Wrong:


Vault.Asset
used directly by TraceCore


---

## Context Responsibility Leakage

Wrong:


Identity decides business permissions

---

# Final Principle

A bounded context is not a technical module.

It is a boundary of responsibility, ownership, and meaning.

The architecture remains healthy when every concept has exactly one rightful home.
# Ankhora System Map

## Purpose

This document describes the global system structure of Ankhora.

It provides a high-level map of:

- major bounded contexts
- architectural relationships
- communication flows
- ownership boundaries

This document should be read before modifying interactions between subsystems.

---

# Global Architecture

Ankhora is composed of independent but connected bounded contexts.

The platform follows a principle:

> Each system owns its responsibility and collaborates through explicit contracts.

High-level view:
                     Users
                       |
                       |
                 Identity Context
                       |
                       |
    +------------------+------------------+
    |                                     |
    |                                     |
        Secure Vault C3 Collaboration
    |                   |
    |                   |
    +------------------+------------------+
                        |
                        |
            TraceCore Context
                        |
                        |
                Domain Applications



---

# Core Bounded Contexts

## Identity Context

### Responsibility

Identity establishes trusted participants.

Identity answers:

> Who is interacting with the platform?

---

### Owns

- user identity
- authentication
- identity verification
- public identity information
- trust establishment

---

### Does Not Own

- business resources
- files
- workflows
- collaboration permissions

---

### Provides

- verified identity
- identity claims
- trust information

---

# Vault Context

## Responsibility

Secure ownership and protection of information.

Vault answers:

> How do we protect and manage owned data?

---

## Owns

- encrypted objects
- assets
- keys
- secure storage lifecycle
- sharing primitives

---

## Does Not Own

- business workflows
- domain validation
- collaboration semantics

---

## Provides

- secure storage
- encrypted assets
- protected data access

---

# C3 Collaboration Context

## Responsibility

Secure collaboration between participants.

C3 answers:

> How do trusted entities collaborate?

---

## Owns

- workspaces
- channels
- threads
- collaboration assets
- trust groups
- collaboration lifecycle

---

## Does Not Own

- business domain meaning
- raw storage implementation
- identity creation

---

## Provides

- collaboration structures
- secure communication primitives
- collaboration events

---

# TraceCore Context

## Responsibility

Lifecycle, history, and verification.

TraceCore answers:

> What happened, when, and under which rules?

---

## Owns

- commits
- history
- branches
- workflows
- validations
- approvals
- audit records

---

## Does Not Own

- encryption
- user identity
- raw asset storage

---

## Provides

- verifiable lifecycle
- operational history
- compliance evidence

---

# Federation Context

## Responsibility

Trusted communication between independent systems.

Federation answers:

> Can two independent environments collaborate safely?

---

## Owns

- remote vault relationships
- trust resolution
- message validation
- federation protocols
- outbound signing

---

## Does Not Own

- domain aggregates
- transport infrastructure
- collaboration rules

---

## Provides

- trusted exchange
- remote verification
- secure synchronization

---

# Domain Application Contexts

## Responsibility

Business-specific capabilities.

Examples:

- Construction
- Healthcare
- Supply Chain
- Banking
- Manufacturing

---

## Owns

- business processes
- business rules
- domain entities

---

## Uses

Platform capabilities:

- Identity
- Vault
- C3
- TraceCore
- Federation

---

# Communication Rules

Bounded contexts communicate through:

## Synchronous communication

Used when immediate information is required.

Examples:

- identity verification
- permission checks
- validation requests

---

## Asynchronous communication

Used for domain events.

Examples:
AssetCreated

TrustGroupUpdated

CommitCreated

WorkflowApproved

VaultShared


---

# Dependency Direction

The preferred dependency direction:
Domain Applications

    |
    v

TraceCore / C3

    |
    v

Vault / Identity / Federation

    |
    v

Infrastructure

Lower-level systems provide capabilities.

Higher-level systems use those capabilities.

---

# Forbidden Dependencies

The following patterns are forbidden.

## Direct database coupling

Example:
C3 ---> Vault database

Incorrect.

Contexts communicate through contracts.

---

## Cross-context model sharing

Example:
TraceCore uses Vault internal structs


Incorrect.

Each context owns its model.

---

## Infrastructure leaking into domain

Example:


Domain imports PostgreSQL package


Incorrect.

---

# Data Ownership Map

| Data | Owner |
|---|---|
| Identity information | Identity |
| Encrypted assets | Vault |
| Sharing relationships | Vault / C3 |
| Channels | C3 |
| Threads | C3 |
| Lifecycle history | TraceCore |
| Validation rules | TraceCore |
| Business objects | Domain Applications |
| Remote trust relationships | Federation |

---

# Security Boundary Map

Each context represents a security boundary.
Identity
|
|
Trust established
|
|
Vault
|
|
Protected information
|
|
C3
|
|
Collaboration
|
|
TraceCore
|
|
Verifiable history


Every boundary requires explicit authorization and validation.

---

# Architectural Mental Model

When designing a new feature, ask:

1. Which bounded context owns this responsibility?
2. Who owns the data?
3. Which trust relationship is involved?
4. Should this be synchronous or event-driven?
5. Does this introduce coupling?
6. Does this preserve security boundaries?

---

# Final Principle

Ankhora is a network of specialized systems.

The power of the platform comes not from one component doing everything.

The power comes from independent components collaborating through trust, ownership, and explicit contracts.











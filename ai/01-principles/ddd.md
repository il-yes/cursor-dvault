# Domain-Driven Design Principle

## Purpose

This document defines how Domain-Driven Design guides the architecture of Ankhora.

DDD is not a folder structure.

DDD is a way of organizing software around business reality.

---

# Domain First Philosophy

The business domain is the source of truth.

Technology exists to support domain needs.

The system should reflect real-world concepts rather than technical abstractions.

---

# Principle 1 — Bounded Contexts Define Responsibility

Each bounded context represents a specific business capability.

A bounded context owns:

- its model
- its rules
- its decisions
- its language

No context should become a universal database for the entire system.

---

# Principle 2 — Shared Language Matters

Every domain concept should have a clear meaning.

Terms must be consistent.

Examples:

- Vault
- Asset
- Workspace
- Channel
- Thread
- Trust Group
- Commit
- Federation

A shared language prevents architectural confusion.

---

# Principle 3 — Aggregates Protect Invariants

Aggregates define consistency boundaries.

An aggregate is responsible for protecting its own rules.

External systems should not directly modify internal state.

Changes should happen through controlled operations.

---

# Principle 4 — Domain Logic Belongs In The Domain

Business rules should not live in:

- HTTP handlers
- database code
- infrastructure services

The domain should express business behavior.

---

# Principle 5 — Infrastructure Serves The Domain

Infrastructure provides capabilities:

- databases
- messaging
- storage
- networking
- external services

Infrastructure should not define business decisions.

---

# Principle 6 — Events Represent Important Changes

Domain events communicate meaningful business changes.

Examples:

- AssetCreated
- TrustGroupUpdated
- CommitCreated
- VaultShared

Events allow systems to evolve independently.

---

# Principle 7 — Dependencies Must Point Inward

The core domain should not depend on external technology.

Preferred direction:
Infrastructure

    ↓

Application

    ↓

Domain


The domain remains independent.

---

# Principle 8 — Models Should Reflect Reality

Avoid designing around:

- database tables
- API endpoints
- frameworks

Design around:

- business concepts
- workflows
- responsibilities

---

# DDD In Ankhora

## Vault Context

Responsible for secure storage concepts.

## TraceCore Context

Responsible for lifecycle, history, validation, and audit.

## Federation Context

Responsible for trusted communication between independent systems.

## Identity Context

Responsible for identity and trust establishment.

## C3 Context

Responsible for collaborative interactions.

---

# Final Principle

Good architecture is not about organizing code.

It is about organizing understanding.

The structure of the software should reveal the structure of the business.
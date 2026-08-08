# Ankhora Architecture Overview

## Purpose

This document provides the global mental model of the Ankhora platform.

It explains how the major systems interact and how responsibilities are distributed.

Before modifying any component, engineers and AI assistants should understand this architecture.

---

# Architectural Vision

Ankhora is a modular, domain-driven, event-oriented platform designed to enable secure collaboration through cryptographic trust.

The architecture is built around several fundamental ideas:

- bounded contexts
- explicit ownership
- cryptographic trust
- domain-driven design
- event-driven communication
- modular evolution

---

# High-Level Model

Ankhora can be understood as a trusted collaboration stack.

             Business Applications

                     ▲

          Domain Use Cases

                     ▲

              TraceCore

    Lifecycle • History • Validation
    Workflow • Compliance • Audit

                     ▲

                 C3 Layer

   Collaboration • Channels • Threads
   Assets • Sharing • Federation

                     ▲

              Secure Vault

   Encryption • Storage • Keys
   Assets • Access Control

                     ▲

                Identity

   Users • Trust • Authentication


---

# Core Architectural Layers

## Identity Layer

Identity establishes who participants are.

Responsibilities:

- identity management
- authentication
- cryptographic identity
- trust establishment

Identity does not own business data.

---

## Vault Layer

The Vault provides secure information management.

Responsibilities:

- encrypted storage
- key management
- assets
- sharing
- recovery
- secure ownership

The Vault is responsible for protecting information.

It is not responsible for business workflows.

---

## C3 Collaboration Layer

C3 provides secure collaboration primitives.

Responsibilities:

- workspaces
- channels
- threads
- collaboration assets
- trust groups
- federation relationships

C3 enables interaction without owning business meaning.

---

## TraceCore Layer

TraceCore manages lifecycle and verifiable history.

Responsibilities:

- commits
- versions
- branches
- validation
- workflows
- approvals
- audit trail

TraceCore provides operational truth.

---

## Domain Application Layer

Business domains use the platform.

Examples:

- construction
- healthcare
- supply chain
- banking
- manufacturing

These domains should not rebuild:

- security
- storage
- collaboration
- lifecycle management

They consume platform capabilities.

---

# Architectural Style

Ankhora follows:

## Domain-Driven Design

Business concepts define boundaries.

---

## Clean Architecture

Dependencies point toward business logic.

---

## Event-Driven Architecture

Important state changes can produce domain events.

---

## Modular Monolith Evolution

The platform is designed with strong boundaries while allowing future distribution.

---

# Separation Of Responsibilities

Each layer answers different questions.

Identity:

"Who are you?"

Vault:

"What information do you own?"

C3:

"How do you collaborate?"

TraceCore:

"What happened and why?"

Domain Applications:

"What business process are you executing?"

---

# Core Design Rule

No component should become the owner of concepts outside its responsibility.

Examples:

Vault does not own business workflows.

TraceCore does not store raw application data.

Federation does not manage domain entities.

Identity does not define permissions for every business case.

---

# Communication Principles

Contexts communicate through:

- explicit interfaces
- domain events
- validated messages
- contracts

Avoid:

- direct database access
- hidden coupling
- shared internal models

---

# Architecture Evolution

The architecture supports future growth.

Possible future extensions:

- additional business domains
- new storage technologies
- new federation protocols
- distributed deployments
- sovereign infrastructure

The principles remain stable while implementations evolve.

---

# Final Principle

Ankhora is not a collection of features.

It is a trust-oriented platform where specialized systems collaborate while preserving clear ownership and responsibility.


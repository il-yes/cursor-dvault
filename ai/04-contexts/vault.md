# Vault Bounded Context

## Purpose

This document defines the Vault bounded context of Ankhora.

Vault is the secure ownership and protection layer of the platform.

Vault answers:

> "How does an entity own, protect, and manage information?"

---

# Mission

The Vault context provides:

- secure ownership of information
- encrypted asset management
- controlled sharing
- key lifecycle management
- protected data access

Vault is the foundation of data sovereignty in Ankhora.

---

# Core Principle

A Vault is not a storage bucket.

A Vault represents:

- ownership
- protection
- trust
- lifecycle
- controlled access

The purpose is not simply storing data.

The purpose is preserving control over data.

---

# Vault Responsibilities

Vault owns:

- vault lifecycle
- asset ownership
- encryption lifecycle
- sharing primitives
- access boundaries
- protected object management

---

# Vault Does Not Own

Vault does not own:

- business workflows
- collaboration semantics
- lifecycle approvals
- identity creation
- billing

---

# Core Domain Model

The main concepts are:
Vault

|

+-- Asset

|

+-- Share Entry

|

+-- Encryption Keys

|

+-- Access Rules

---

# Vault Aggregate

## Vault

The Vault is the root aggregate representing a protected environment.

A Vault contains:

- protected assets
- ownership information
- security metadata
- lifecycle state

---

Possible attributes:
Vault ID

Owner ID

Name

Version

Created At

Updated At

Security Metadata

---

# Vault Invariants

A Vault must ensure:

- ownership is always known
- protected assets have controlled access
- security metadata remains consistent
- lifecycle transitions are valid

---

# Asset

## Definition

An Asset represents protected information managed by the Vault.

Examples:

- document
- encrypted file
- secret
- credential
- structured data object

---

An Asset is not the raw file.

It is the managed protected object.

---

Possible attributes:
Asset ID

Content Reference

Content Hash

Size

Metadata

Owner

Created At

---

# Asset Lifecycle

Example:
Created

|

Encrypted

|

Stored

|

Shared

|

Updated

|

Archived


---

# Asset Ownership

Assets always have an owner.

Sharing access does not transfer ownership.

Example:


Owner

|

Asset

|

Shared Access


---

# Share Entry

## Definition

A Share Entry represents controlled access to an asset.

It describes:

- who can access
- what permissions exist
- how access is protected

---

Possible attributes:


Share ID

Asset Reference

Recipient

Permission

Created At

Expiration


---

# Sharing Principle

Sharing is permission delegation.

It is not ownership transfer.

---

Example:

Incorrect:


Share Asset

=

Transfer Asset


---

Correct:


Share Asset

=

Grant Controlled Access


---

# Encryption Model

Encryption is a core Vault responsibility.

The Vault protects:

- confidentiality
- integrity
- ownership

---

High-level model:


Content

|

v

Encryption

|

v

Protected Asset

|

v

Storage

---

# Key Management

Vault manages key lifecycle concepts:

- creation
- rotation
- protection
- revocation

---

Keys must follow:

- ownership rules
- access rules
- recovery rules

---

# Vault Security Boundary

Vault is a security boundary.

External systems should not directly manipulate Vault internals.

Access occurs through:

- application services
- contracts
- controlled APIs

---

# Vault Events

Vault owns events related to protected information.

Examples:
VaultCreated

AssetCreated

AssetEncrypted

AssetShared

AssetRevoked

KeyRotated


---

# Integration With Other Contexts

---

# Identity Integration

Identity provides:

- authenticated actor
- identity claims

Vault decides:

- ownership
- permissions
- access

Flow:
Identity

Who are you?

    |

    v

Vault

Can you access this resource?


---

# C3 Integration

C3 uses Vault capabilities for collaboration.

Example:


Vault

AssetShared

    |

    v

C3

Collaboration Reference


C3 does not become the owner of the asset.

---

# TraceCore Integration

TraceCore records lifecycle evidence.

Example:


Vault

AssetCreated

    |

    v

TraceCore

Lifecycle Event


TraceCore does not replace Vault storage.

---

# Federation Integration

Federation allows controlled exchange.

Example:


Local Vault

    |

Encrypted Exchange

    |

Remote Vault


Trust must be established before sharing.

---

# Domain Events

Vault domain events represent meaningful changes.

Examples:


VaultInitialized

AssetCreated

AssetUpdated

AssetShared

AccessRevoked

KeyRotationCompleted


---

# Forbidden Patterns

## Vault As Generic Database

Wrong:


Everything stores data in Vault


Vault protects information.

It does not replace every domain.

---

## Business Logic Inside Vault

Wrong:


Vault decides construction workflow approval


---

## Removing Ownership

Wrong:


Asset exists without owner


---

## Global Decryption

Wrong:


Every service receives plaintext assets

---

# AI Implementation Rules

Before implementing Vault functionality, AI must ask:

1. Who owns this information?
2. Is this an Asset or domain data?
3. Does ownership change?
4. Does encryption happen here?
5. Should this be shared or referenced?
6. Should another context own this logic?
7. Does this action require an audit event?

---

# Final Principle

Vault is the sovereignty layer of Ankhora.

It does not simply store information.

It preserves ownership, protection, and trust throughout the entire lifecycle of data.



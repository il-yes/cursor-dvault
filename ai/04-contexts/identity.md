# Identity Bounded Context

## Purpose

This document defines the Identity bounded context of Ankhora.

Identity is responsible for establishing trusted participants in the platform.

Identity answers:

> "Who is this entity?"

It does not answer:

> "What can this entity do?"

Authorization belongs to the context owning the protected resource.

---

# Mission

Identity provides the foundation for trust across Ankhora.

It establishes:

- users
- organizations
- devices
- authentication
- identity verification
- trust bootstrap

Identity creates the foundation upon which other contexts build their own authorization decisions.

---

# Core Principle

Identity is about recognition.

Ownership and permissions are separate concerns.

Example:

Identity says:
This is Alice.


Vault decides:


Can Alice access this asset?


C3 decides:


Can Alice join this workspace?


TraceCore decides:


Can Alice approve this workflow?

---

# Responsibilities

Identity owns:

- identity lifecycle
- authentication
- credential management
- identity verification
- identity claims
- trust establishment

---

# Identity Does Not Own

Identity does not own:

- business permissions
- vault access rules
- collaboration permissions
- workflow approvals
- domain roles

---

# Core Domain Concepts

## Identity

Represents a verified participant.

An Identity may represent:

- person
- organization
- service
- device

---

Possible attributes:
Identity ID

Public Information

Verification Status

Created At

Updated At

---

# Identity Types

Ankhora may support different identity categories.

## Human Identity

Represents a person.

Examples:

- user
- administrator
- collaborator

---

## Organization Identity

Represents an organization.

Examples:

- company
- institution
- enterprise

---

## Device Identity

Represents a trusted device.

Examples:

- desktop installation
- secure endpoint
- service node

---

## Service Identity

Represents automated systems.

Examples:

- integration service
- federation node
- automation agent

---

# Authentication

Authentication proves identity ownership.

Possible mechanisms:

- credentials
- cryptographic keys
- device verification
- tokens
- signatures

---

# Authentication Flow

Example:
User

|

v

Authentication Request

|

v

Identity Context

|

v

Identity Verification

|

v

Authenticated Session


---

# Identity Claims

Identity may provide verified claims.

Examples:


Identity ID

Organization

Verification Status

Trust Level

Capabilities

Other contexts consume claims.

They do not modify identity state.

---

# Device Trust

Because Ankhora is desktop-first, devices are important trust objects.

A device may have:

- identity
- cryptographic material
- registration state
- trust status

---

Example:
User Identity

    |

    |

Trusted Desktop Device

    |

    |

Vault Access

---

# Trust Bootstrap

Identity participates in establishing trust relationships.

Examples:

- first device registration
- organization invitation
- federation handshake

---

# Identity Events

Identity owns identity-related events.

Examples:
IdentityCreated

IdentityVerified

DeviceRegistered

DeviceRevoked

TrustEstablished

---

# Integration With Other Contexts

---

# Identity → Vault

Purpose:

Authenticate ownership and access requests.

Flow:
Identity

Provides:

Verified Identity

    |

    v

Vault

Evaluates:

Ownership + Permission


Identity does not grant asset access.

---

# Identity → C3

Purpose:

Identify collaboration participants.

Flow:


Identity

Participant Verification

    |

    v

C3

Workspace / Channel Authorization


---

# Identity → TraceCore

Purpose:

Attribute actions.

Flow:


Identity

User Identity

    |

    v

TraceCore

Commit Author

Approval Actor


---

# Identity → Federation

Purpose:

Establish trusted communication.

Flow:


Identity

Verification

    |

    v

Federation

Trust Relationship

---

# Security Principles

Identity must protect:

- authentication secrets
- credentials
- identity metadata
- cryptographic material

---

# Identity Boundary Rules

Identity must not become a global permission engine.

Incorrect:
Identity

decides all access everywhere


---

Correct:


Identity

proves who you are

Resource Context

decides what you can do

---

# Failure Handling

Identity must support:

- revoked credentials
- lost devices
- compromised identities
- recovery flows
- trust expiration

---

# AI Implementation Rules

When implementing Identity features, AI should ask:

1. Is this authentication or authorization?
2. Does Identity own this concept?
3. Should another context decide?
4. Is cryptographic verification required?
5. Does this create a new trust relationship?
6. Is the action auditable?

---

# Final Principle

Identity is the foundation of trust.

Identity proves existence.

Other contexts decide responsibility.

A trustworthy system knows both:
who someone is,
and who controls each resource.




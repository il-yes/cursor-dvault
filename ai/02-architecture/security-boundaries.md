# Ankhora Security Boundaries

## Purpose

This document defines the security boundaries of the Ankhora platform.

Security boundaries describe where:

- trust changes
- permissions are evaluated
- encryption is applied
- ownership is enforced
- verification is required

A security boundary is not only a technical boundary.

It is an ownership and trust boundary.

---

# Security Philosophy

Ankhora assumes that no environment should automatically be fully trusted.

Security is achieved through:

- explicit trust relationships
- cryptographic verification
- minimal access
- ownership enforcement
- observable actions

The goal is not to create a single secure zone.

The goal is to create controlled interactions between independent zones.

---

# Security Boundary Model

High-level view:
+------------------------------------------------+

            User Environment

+------------------------------------------------+

                |
                |
                v

+------------------------------------------------+

          Desktop Security Boundary
Private Keys
Plain Data
Vault Engine
Local Operations

+------------------------------------------------+

                |
      Encrypted Communication

                |

                v

+------------------------------------------------+

          Cloud Security Boundary
Identity
Coordination
Synchronization
Collaboration Services

+------------------------------------------------+

                |

      Federation Trust Boundary

                |

                v

+------------------------------------------------+

      External Organization Boundary

+------------------------------------------------+


---

# Boundary 1 — User Sovereignty Boundary

## Purpose

Protect the user's ownership and control.

This is the strongest personal security boundary.

---

## Contains

- user-controlled devices
- private keys
- local vault state
- decrypted information

---

## Trust Level

Highest trust because the user controls the environment.

---

## Security Responsibilities

The desktop must protect:

- credentials
- cryptographic material
- local data
- user actions

---

# Boundary 2 — Desktop Application Boundary

## Purpose

Protect the internal components of the Ankhora desktop environment.

---

## Contains

- Vault Engine
- local storage
- local domain services
- local event processing

---

## Rules

Internal components communicate through:

- application services
- domain interfaces
- controlled contracts

---

## Forbidden

Direct uncontrolled access between modules.

Example:
UI

directly modifies

Vault database

Incorrect.

---

# Boundary 3 — Cryptographic Boundary

## Purpose

Separate protected information from external systems.

This is one of Ankhora's most important boundaries.

---

## Principle

Encrypted information may leave the ownership boundary.

Plain information should not.

---

Example:
Plain Document

    |

    v

Encryption Boundary

    |

    v

Encrypted Asset

    |

    v

Cloud Storage / Coordination

---

# Boundary 4 — Cloud Coordination Boundary

## Purpose

Provide collaboration capabilities without becoming the owner of user data.

---

## Contains

- authentication services
- collaboration services
- synchronization
- notifications
- subscription services

---

## Cloud Can Know

- identities
- permissions required for coordination
- encrypted references
- operational metadata

---

## Cloud Should Not Automatically Know

- plaintext user content
- private encryption keys
- protected documents

---

# Boundary 5 — Context Boundary

Every bounded context is a security boundary.

Examples:

Vault

owns:

encrypted assets

TraceCore

owns:

history and lifecycle

Federation

owns:

trusted exchange


---

# Rules

A context must protect its internal state.

External contexts interact through:

- APIs
- events
- contracts

---

# Boundary 6 — Federation Boundary

## Purpose

Protect communication between independent Ankhora environments.

Example:


Organization A

    |

    v

Federation Layer

    |

    v

Organization B

---

## Federation Requires

- identity verification
- trust establishment
- message validation
- policy agreement
- cryptographic verification

---

## Federation Does Not Assume

- remote systems are safe
- remote users are authorized
- remote data is correct

---

# Boundary 7 — Enterprise Boundary

Organizations introduce their own security domains.

Examples:

- companies
- institutions
- government entities

Each organization may define:

- policies
- access rules
- compliance requirements

Ankhora enables collaboration without removing organizational sovereignty.

---

# Authentication vs Authorization Boundary

Important distinction:

Authentication:

> Who are you?

Owned primarily by Identity.

---

Authorization:

> What are you allowed to do?

Owned by the resource context.

---

Example:

Identity says:
User = Alice


Vault decides:


Can Alice access Asset X?


C3 decides:


Can Alice join Channel Y?


TraceCore decides:


Can Alice approve Workflow Z?

---

# Encryption Boundary Rules

## Data At Rest

Protected by:

- encryption
- key management
- ownership rules

---

## Data In Transit

Protected by:

- secure communication
- authenticated channels
- message validation

---

## Data In Use

Protected by:

- minimal exposure
- controlled decryption
- authorized execution boundary

---

# Security Event Boundary

Important security events should be observable.

Examples:
LoginAttempted

KeyRotated

AccessGranted

AccessRevoked

TrustEstablished

FederationRejected

---

# Security Boundary Violations

The following patterns are forbidden.

---

## Cloud Access To Private Keys

Wrong:
Cloud Service

    |

    v

User Private Keys


---

## Shared Plain Data Everywhere

Wrong:


Every Service

    |

    v

Plain User Data


---

## Context Bypass

Wrong:


TraceCore

directly accesses

Vault database


---

## Invisible Trust

Wrong:


Remote System

automatically trusted

---

# AI Security Review Questions

Before implementing a feature, AI should ask:

1. What security boundary does this cross?
2. Does ownership change?
3. Does encryption need to happen here?
4. Who owns the keys?
5. Who validates access?
6. Should this be local or cloud?
7. Should this interaction be audited?

---

# Final Principle

Ankhora security is built through boundaries.

Every boundary exists to protect:

- ownership
- privacy
- integrity
- trust

A secure system is not one where everything is protected by one wall.

A secure system is one where every interaction crosses a controlled boundary.





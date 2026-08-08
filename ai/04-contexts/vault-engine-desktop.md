
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
# Vault Engine Desktop Context

## Purpose

This document defines the desktop implementation of the Vault bounded context.

The Vault Engine Desktop is the sovereign execution environment responsible for protecting user-owned information.

It provides:

- local vault management
- encryption operations
- key ownership
- local persistence
- offline capability
- secure reconstruction

---

# Mission

The Vault Engine Desktop answers:

> "How does a user locally own and control protected information?"

The desktop is the primary trust boundary for user data.

---

# Core Principle

The desktop is not a cloud client.

The desktop is a sovereign Ankhora node.

It can:

- own data
- protect data
- transform data
- operate offline
- synchronize later

---

# Desktop Architecture Position

---

# Responsibilities

The Vault Engine Desktop owns:

- local vault lifecycle
- vault reconstruction
- encryption/decryption
- local asset management
- local indexing
- local state management
- offline operations

---

# Does Not Own

The desktop Vault Engine does not own:

- global identity management
- cloud billing
- remote user discovery
- federation routing
- collaboration presence

Those belong to other contexts.

---

# Core Components

## VaultPayload

VaultPayload represents the reconstructed local vault state.

It contains:

- vault metadata
- personal data
- collaborative data
- indexes
- synchronization information

Example conceptual model:
VaultPayload

|
+-- Personal Content

|
+-- Collaborative Content

|
+-- Index

|
+-- Metadata

---

# Local Vault Lifecycle

## Initialization

Example:
Create Vault

    |

Generate Local Metadata

    |

Initialize Encryption Context

    |

Create Empty VaultPayload


---

## Loading

Example:


Encrypted Storage

    |

Decrypt Authorized Data

    |

Reconstruct VaultPayload

    |

Open User Session


---

## Modification

Example:


User Action

    |

Domain Operation

    |

Update Local State

    |

Create Event

    |

Synchronize Later


---

# Encryption Boundary

The desktop owns cryptographic operations.

Flow:


Plain Content

    |

    v

Encryption Service

    |

    v

Encrypted Asset

    |

    v

Local / Remote Storage

---

# Key Management

Private keys should remain under desktop control whenever possible.

Responsibilities:

- key creation
- key usage
- key rotation
- key protection

---

# Local Storage

The desktop may use:

- local database
- encrypted filesystem
- embedded storage

The storage technology is an implementation detail.

The Vault domain does not depend on it.

---

# Offline-First Principle

The desktop should support:

- creating data offline
- modifying data offline
- validating local operations
- synchronizing later

---

Example:
Offline User Action

    |

Local Commit

    |

Pending Synchronization

    |

Cloud Sync When Available


---

# Synchronization Model

Synchronization is not database replication.

It is controlled exchange.

Example:


Desktop

Local Change

    |

    v

Synchronization Layer

    |

    v

Cloud Coordination


---

# Events Produced

The desktop may produce Vault events:

Examples:


VaultCreated

AssetCreated

AssetEncrypted

AssetModified

AssetShared

KeyRotated

---

# Interaction With Cloud

The desktop sends:

- encrypted assets
- metadata
- references
- events
- synchronization requests

The desktop does not send:

- private keys
- unnecessary plaintext
- uncontrolled local state

---

# Interaction With Identity

Identity provides:

- authenticated user
- device trust
- identity claims

Vault Engine uses this information.

It does not create identities.

---

# Interaction With C3

C3 uses Vault capabilities.

Example:
User Shares Asset

    |

Vault Engine

Encrypts + Creates Permission

    |

C3

Coordinates Collaboration


---

# Interaction With TraceCore

The desktop may create lifecycle information.

Example:


Local Change

    |

Commit Candidate

    |

TraceCore Synchronization


---

# Security Rules

The Vault Engine must:

- protect private keys
- minimize plaintext lifetime
- validate ownership
- isolate sensitive operations

---

# Forbidden Patterns

## Cloud-First Storage

Wrong:


User Data

    |

    v

Cloud Database

    |

    v

Desktop Cache


---

## Remote Key Ownership

Wrong:


Cloud

owns encryption keys


---

## Plain Synchronization

Wrong:


Desktop

uploads plaintext objects


---

## UI Controls Domain Logic

Wrong:


React Component

directly modifies Vault state

---

# AI Implementation Rules

When generating desktop Vault code, AI must verify:

1. Does this operation belong locally?
2. Is plaintext exposure minimized?
3. Who owns the encryption key?
4. Can this work offline?
5. Is synchronization separated from domain logic?
6. Is the cloud being treated as coordinator only?
7. Is the operation auditable?

---

# Final Principle

The Vault Engine Desktop is the user's sovereign node.

The cloud may help coordinate information.

The desktop preserves ownership.

The user's data remains under the user's control.










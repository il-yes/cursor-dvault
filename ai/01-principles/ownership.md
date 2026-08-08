# Ownership Principle

## Purpose

This document defines how Ankhora understands ownership.

Ownership is a fundamental concept that determines:

- control
- responsibility
- permissions
- trust boundaries
- lifecycle management

---

# Ownership Philosophy

Data is not simply stored.

Data belongs to someone.

Every important resource should have:

- an owner
- a lifecycle
- defined responsibilities
- controlled access

---

# Principle 1 — Ownership Must Be Explicit

Ownership should never be ambiguous.

Every important object should answer:

- Who owns this?
- Who controls it?
- Who can modify it?
- Who can delegate access?
- Who can revoke access?

---

# Principle 2 — Ownership Creates Boundaries

Ownership defines boundaries.

Examples:

A user owns personal information.

An organization owns organizational resources.

A workspace owns collaboration structures.

A vault protects owned assets.

Clear ownership prevents uncontrolled access.

---

# Principle 3 — Ownership Is Different From Access

Ownership and access are not the same.

Ownership defines control.

Access defines permission.

An owner may delegate access without transferring ownership.

---

# Principle 4 — Ownership Must Be Transferable When Appropriate

Some resources require controlled transfer.

Examples:

- organizational changes
- succession
- recovery scenarios
- business transitions

Transfers must be:

- explicit
- authorized
- traceable

---

# Principle 5 — Ownership Requires Accountability

Ownership means responsibility.

An owner is responsible for:

- protecting resources
- managing permissions
- maintaining integrity
- defining lifecycle rules

---

# Principle 6 — Cryptographic Ownership

For sensitive resources, ownership should be reinforced through cryptographic mechanisms.

Examples:

- encryption keys
- signatures
- proofs of control
- verification mechanisms

Ownership should not rely only on database records.

---

# Principle 7 — Ownership Exists Across Boundaries

Federated systems require ownership clarity.

When information crosses systems:

- original ownership must remain visible
- permissions must be respected
- responsibilities must remain clear

---

# Ownership In Ankhora

## Vault

Ownership defines who controls protected information.

## Assets

Ownership defines who controls encrypted objects.

## Trust Groups

Ownership defines who manages collaboration relationships.

## TraceCore

Ownership defines who controls lifecycle and validation.

## Federation

Ownership defines responsibility between independent systems.

---

# Final Principle

A trustworthy system must always know:

Who owns this?

Who controls this?

Who is responsible for this?
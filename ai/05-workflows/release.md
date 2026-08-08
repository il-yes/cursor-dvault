
# Bug Fix Workflow

## Purpose

This document defines the standard workflow for diagnosing and fixing bugs in Ankhora.

The goal is not only to remove errors.

The goal is to restore correctness while preserving architecture.

---

# Core Principle

A bug is not only a coding mistake.

A bug is usually:

- a broken invariant
- an incorrect assumption
- a missing validation
- a violated boundary
- an unexpected interaction

---

A bug fix must answer:

```

What failed?

Why did it fail?

Where should the correction live?

How do we prevent recurrence?

```

---

# Bug Fix Process

Every bug follows:

```

1. Reproduce

2. Understand

3. Locate Ownership

4. Identify Root Cause

5. Implement Fix

6. Validate

7. Prevent Regression

```

---

# Step 1 — Reproduce The Problem

Before changing code:

Confirm:

- exact error
- reproduction steps
- affected environment
- expected behavior
- actual behavior

---

Example:

Bad:

```

WebSocket broken

```

Good:

```

WebSocket connection closes with code 1006
after authentication handshake.

```

---

# Step 2 — Understand The Failure

Classify the problem.

Possible categories:

## Domain Failure

Example:

```

Invalid state transition

```

Owner:

```

Domain Layer

```

---

## Application Failure

Example:

```

Incorrect use case orchestration

```

Owner:

```

Application Layer

```

---

## Infrastructure Failure

Example:

```

Database timeout

API unavailable

Storage failure

```

Owner:

```

Infrastructure Layer

```

---

## Integration Failure

Example:

```

Federation message rejected

Invalid event payload

```

Owner:

```

Integration Boundary

```

---

# Step 3 — Identify Ownership

Every bug must have a responsible context.

Ask:

```

Which bounded context owns the failed rule?

```

---

Examples:

Encryption failure:

```

Vault

```

not:

```

C3

```

---

Invalid collaboration permission:

```

C3

```

not:

```

Identity

```

---

Incorrect lifecycle validation:

```

TraceCore

```

not:

```

Domain Application

```

---

# Step 4 — Root Cause Analysis

AI must avoid symptom fixes.

Example:

Problem:

```

Asset cannot be decrypted

```

Possible symptoms:

- wrong key
- missing metadata
- corrupted payload
- incorrect version

The fix must target the cause.

---

Useful questions:

```

What assumption failed?

What invariant was violated?

What boundary allowed invalid data?

Why was this not detected earlier?

```

---

# Step 5 — Design The Fix

Before coding:

Explain:

## Cause

Example:

```

Vault reconstruction accepted incomplete metadata.

```

---

## Correction

Example:

```

VaultPayload validation now rejects incomplete assets.

```

---

## Location

Example:

```

Vault Domain Layer

```

---

## Risk

Example:

```

Existing corrupted vaults require migration handling.

```

---

# Step 6 — Implement The Fix

Fix rules:

- minimal change
- correct ownership
- no unrelated refactoring
- preserve architecture

---

Preferred:

```

Small Correct Fix

*

Regression Test

```

---

Avoid:

```

Large Rewrite

```

---

# Testing Requirements

Every bug fix requires:

## Reproduction Test

A test that fails before the fix.

---

## Correction Test

A test that passes after the fix.

---

## Regression Protection

Future changes must not recreate the problem.

---

# Debugging Boundaries

Many bugs occur between contexts.

AI should inspect:

## Vault ↔ C3

Check:

- ownership
- permissions
- asset references

---

## Desktop ↔ Cloud

Check:

- synchronization state
- encryption
- versions

---

## Identity ↔ Federation

Check:

- trust
- signatures
- authentication

---

## TraceCore ↔ Domain Applications

Check:

- events
- lifecycle state
- validation rules

---

# Logging Principles

Logs should help answer:

- what happened?
- where?
- by whom?
- with which correlation ID?

---

Avoid:

```

Something failed

```

Prefer:

```

Asset reconstruction failed

vault_id=X

asset_id=Y

reason=missing encryption metadata

```

---

# Security Bug Rules

Security bugs receive priority.

Examples:

- unauthorized access
- key exposure
- plaintext leakage
- invalid trust acceptance

AI must:

1. Identify impact
2. Stop unsafe behavior
3. Add validation
4. Add regression coverage

---

# Forbidden Patterns

AI must not:

- patch errors blindly
- remove validation to hide failures
- ignore failing tests
- move logic into the wrong context
- fix symptoms instead of causes

---

# AI Debugging Checklist

Before proposing a fix:

- [ ] Can the bug be reproduced?
- [ ] Is the root cause understood?
- [ ] Is ownership clear?
- [ ] Is the fix placed in the correct layer?
- [ ] Is a regression test added?
- [ ] Are security implications checked?

---

# Final Principle

A bug fix is successful when:

```

The problem disappears

*

The reason it happened is understood

*

The architecture becomes stronger

```

A mature system does not only remove failures.

It learns from them.
```

---


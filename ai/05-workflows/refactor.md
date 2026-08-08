
# Refactoring Workflow

## Purpose

This document defines the standard workflow for refactoring Ankhora.

The goal of refactoring is to improve:

- clarity
- maintainability
- performance
- architecture alignment

without changing intended behavior.

---

# Core Principle

Refactoring changes the structure of the system.

It does not change the meaning of the system.

```

Same Behavior

*

Better Structure

=

Successful Refactor

```

---

# Refactoring Objectives

A refactor may improve:

- duplicated logic
- unclear responsibilities
- dependency direction
- naming
- performance bottlenecks
- testability
- technical debt

---

# Refactoring Process

Every refactor follows:

```

1. Identify Problem

2. Understand Current Design

3. Define Target State

4. Execute Incrementally

5. Validate Behavior

6. Document Decision

```

---

# Step 1 — Identify The Problem

Before changing code, identify:

- What is wrong?
- Why is it a problem?
- What behavior must remain unchanged?

Examples:

Bad:

```

This code looks ugly.

```

Good:

```

Vault encryption logic is duplicated
across three services causing inconsistent behavior.

```

---

# Step 2 — Understand Ownership

Before moving code, determine:

```

Who owns this responsibility?

```

Examples:

Encryption logic:

```

Vault

```

not:

```

C3
TraceCore

```

---

Lifecycle validation:

```

TraceCore

```

not:

```

Vault

```

---

Collaboration rules:

```

C3

```

not:

```

Identity

```

---

# Step 3 — Define Target Architecture

Before implementation, describe:

## Current State

Example:

```

Handler

|

Business Logic

|

Repository

```

---

## Target State

Example:

```

Handler

|

Use Case

|

Domain

|

Repository Interface

|

Infrastructure

```

---

# Refactoring Categories

## Structural Refactor

Examples:

- moving packages
- separating responsibilities
- improving dependency direction

---

## Domain Refactor

Examples:

- improving aggregates
- extracting value objects
- clarifying invariants

---

## Performance Refactor

Examples:

- reducing allocations
- improving queries
- optimizing synchronization

---

## Security Refactor

Examples:

- reducing plaintext exposure
- improving key handling
- strengthening validation

---

# DDD Refactoring Rules

AI must preserve:

## Bounded Context Ownership

Do not move logic into the wrong context.

---

## Dependency Direction

Allowed:

```

Interface

↓

Application

↓

Domain

```

Infrastructure stays outside.

---

Forbidden:

```

Domain

imports

Database package

```

---

# Testing Requirements

Before refactoring:

- existing tests must pass

During refactoring:

- add missing tests

After refactoring:

- behavior must remain verified

---

# Safe Refactoring Strategy

Prefer:

```

Small Change

↓

Run Tests

↓

Review

↓

Next Change

```

over:

```

Large Rewrite

```

---

# AI Refactoring Rules

When asked to refactor, AI must:

1. Explain the current problem
2. Identify affected context
3. Explain risks
4. Propose migration steps
5. Keep behavior unchanged
6. Add or update tests

---

# Forbidden Behavior

AI must not:

- rewrite working systems without reason
- optimize prematurely
- remove abstractions without analysis
- merge bounded contexts
- introduce hidden coupling

---

# Final Principle

A good refactor makes the architecture clearer.

The best refactor is one where future engineers understand the system faster.
```

---

Now:

Create:

```text
ai/05-workflows/release.md
```

---

```md
# Release Workflow

## Purpose

This document defines the standard workflow for releasing Ankhora versions.

A release represents a trusted transition from development state to delivered product capability.

---

# Core Principle

A release is not just deployment.

A release is:

```

Validated Code

*

Verified Architecture

*

Documented Change

*

Operational Confidence

```

---

# Release Objectives

A release must provide:

- stability
- traceability
- reproducibility
- security confidence

---

# Release Process

Every release follows:

```

1. Prepare

2. Validate

3. Review

4. Package

5. Deploy

6. Monitor

```

---

# Step 1 — Prepare Release

Verify:

- completed features
- completed fixes
- updated documentation
- version changes

---

Required review:

```

What changed?

Why?

Which contexts were affected?

```

---

# Step 2 — Validation

Run:

## Tests

- unit tests
- domain tests
- integration tests

---

## Quality Checks

Verify:

- formatting
- static analysis
- dependency health

---

## Security Checks

Verify:

- secrets handling
- permissions
- encryption boundaries

---

# Step 3 — Architecture Review

Before release, verify:

## Context Boundaries

Questions:

- Did any responsibility leak?
- Did a context become coupled?

---

## Dependencies

Questions:

- Are dependencies intentional?
- Did infrastructure enter domain code?

---

## Data Protection

Questions:

- Is ownership preserved?
- Are encryption rules respected?

---

# Step 4 — Versioning

Versions should communicate change.

Example:

```

Major.Minor.Patch

```

---

Major:

Breaking architectural or API changes.

---

Minor:

New capabilities.

---

Patch:

Bug fixes and small improvements.

---

# Step 5 — Release Artifacts

A release may include:

- desktop application
- cloud services
- migration scripts
- documentation
- configuration updates

---

# Desktop Release

Verify:

- application packaging
- local storage compatibility
- encryption compatibility
- migration safety

---

# Cloud Release

Verify:

- API compatibility
- database migrations
- deployment health
- rollback strategy

---

# Step 6 — Traceability

Every release should be connected to:

- commits
- decisions
- issues
- documentation

TraceCore may record:

```

ReleaseCreated

VersionPublished

MigrationExecuted

```

---

# Rollback Strategy

Every release must define:

- rollback possibility
- data migration safety
- recovery procedure

---

# AI Release Checklist

Before recommending release:

## Code

- [ ] Tests passing
- [ ] No unresolved critical errors
- [ ] Dependencies reviewed

---

## Architecture

- [ ] Context boundaries preserved
- [ ] ADRs updated if needed

---

## Security

- [ ] Secrets protected
- [ ] Access rules verified

---

## Documentation

- [ ] User changes documented
- [ ] Technical changes documented

---

# Forbidden Patterns

AI must not:

- release untested features
- skip migrations
- ignore security impact
- deploy architectural experiments as stable versions

---

# Final Principle

A release is a promise.

It tells users:

"This version has been understood, validated, and can be trusted."
```

---


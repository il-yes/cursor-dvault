
# Testing Standards

## Purpose

This document defines the testing standards for Ankhora.

The objective is to ensure the system remains:

- correct
- secure
- evolvable
- reliable

Testing is part of design, not only validation.

---

# Core Principle

A test should protect a behavior.

The primary question is:

> "What rule or expectation does this test protect?"

Not:

> "What line of code does this test cover?"

---

# Testing Strategy

Ankhora follows a layered testing approach:

```

```
          End-to-End Tests

                ↑

      Integration Tests

                ↑

      Application Tests

                ↑

         Domain Tests
```

```

Each layer has a different responsibility.

---

# Domain Tests

## Purpose

Validate business rules and invariants.

Domain tests must be:

- fast
- isolated
- deterministic

---

Examples:

Vault:

```

Asset cannot be shared without permission.

```

---

C3:

```

Closed thread cannot receive new messages.

```

---

TraceCore:

```

Validated commit cannot be modified.

````

---

# Domain Test Rules

Domain tests should:

- create domain objects directly
- avoid databases
- avoid HTTP
- avoid external services

---

Good:

```go
asset.Share(trustGroup)
````

Verify:

```
sharing rules are respected
```

---

Bad:

```
HTTP request

↓

Database

↓

Check result
```

for a domain rule.

---

# Application Tests

## Purpose

Validate use cases.

Application tests verify:

* orchestration
* workflows
* event publishing
* transaction behavior

---

Example:

```
CreateAssetUseCase

↓

Validate

↓

Encrypt

↓

Persist

↓

Publish AssetCreated
```

---

Verify:

* correct services called
* correct events produced
* failures handled

---

# Integration Tests

## Purpose

Validate boundaries.

Important Ankhora integrations:

---

## Vault ↔ Storage

Verify:

* encryption
* persistence
* reconstruction

---

## Vault ↔ C3

Verify:

* sharing flow
* permissions
* collaboration state

---

## C3 ↔ WebSocket Gateway

Verify:

* event delivery
* realtime updates

---

## Federation ↔ Remote Systems

Verify:

* message validation
* trust verification
* synchronization

---

## TraceCore ↔ Domain Applications

Verify:

* lifecycle events
* validation flows

---

# End-to-End Tests

## Purpose

Validate complete user scenarios.

Examples:

```
User creates vault

↓

Stores asset

↓

Shares asset

↓

Collaborator receives access

↓

History recorded
```

---

E2E tests should represent real workflows.

---

# Test Naming

Tests must describe behavior.

Preferred:

```go
TestCannotShareAssetWithoutTrustPermission
```

Avoid:

```go
TestAsset1
```

---

# Test Structure

Prefer:

```
Arrange

Act

Assert
```

Example:

```go
func TestAssetCannotBeSharedWithoutOwnerApproval(t *testing.T) {

    // Arrange

    // Act

    // Assert

}
```

---

# Edge Case Testing

Always test:

## Empty State

Examples:

* empty vault
* empty workspace
* missing metadata

---

## Invalid State

Examples:

* invalid transition
* corrupted payload
* unauthorized actor

---

## Failure State

Examples:

* database unavailable
* network failure
* invalid signature

---

## Concurrent State

Examples:

* simultaneous updates
* synchronization conflicts
* multiple collaborators

---

# Security Testing

Security scenarios are mandatory.

Verify:

## Authorization

```
User A cannot access User B data.
```

---

## Encryption

```
Stored data remains protected.
```

---

## Federation

```
Invalid trust messages are rejected.
```

---

## Input Validation

```
Malformed external data cannot enter the domain.
```

---

# Regression Testing

Every bug fix requires:

```
Bug discovered

        ↓

Regression test created

        ↓

Future protection
```

---

Example:

Problem:

```
Vault reconstruction failed with missing metadata.
```

Required:

```
TestVaultRejectsIncompletePayload
```

---

# Event Testing

Events must be tested.

Verify:

* correct event created
* correct payload
* correct ownership
* correct consumers

---

Example:

```
AssetCreated
```

must contain:

* asset identity
* owner information
* creation metadata

---

# Repository Testing

Repositories should verify:

* persistence correctness
* mapping correctness
* error handling

Avoid testing SQL implementation details.

---

# Performance Testing

Performance tests are required for:

* large vaults
* large histories
* synchronization workloads
* event processing

Measure:

* execution time
* memory usage
* throughput

---

# Go Testing Standards

Required commands:

```
go test ./...

go test -race ./...

go vet ./...
```

---

# Test Quality Rules

A good test is:

* readable
* deterministic
* isolated
* meaningful

---

Avoid:

* flaky tests
* excessive mocking
* testing implementation details
* giant integration tests for simple rules

---

# AI Testing Rules

When implementing a feature, AI should automatically ask:

```
What behavior needs protection?

What can fail?

What invalid states exist?

What regression risk exists?
```

---

AI should generate tests together with implementation.

---

# Testing Checklist

Before accepting code:

* [ ] Domain behavior tested
* [ ] Invalid cases tested
* [ ] Failure cases tested
* [ ] Security cases considered
* [ ] Regression coverage added
* [ ] Integration boundaries validated

---

# Final Principle

Tests are executable documentation.

They preserve the rules of the system.

A mature system is not one without bugs.

A mature system is one where mistakes cannot silently return.

````

---


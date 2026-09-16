
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

A test is valid only when it exercises the same application boundary that the real user exercises.
That means we explicitly distinguish:
UNIT TEST
    ↓
tests isolated behavior

INTEGRATION TEST
    ↓
tests real component boundaries

END-TO-END ACCEPTANCE TEST
    ↓
starts the actual Desktop application
    ↓
uses the actual frontend
    ↓
performs actual user actions
    ↓
exercises Wails
    ↓
exercises real backend services
    ↓
verifies observable result
And fake tests cannot be presented as E2E tests.

_____


## LOOP Testing Principles

### 1. Test the behavior before the integration

Use pragmatic TDD to define and verify the expected behavior at the appropriate unit, domain, application, and integration boundaries before relying on the full application flow.

### 2. Never confuse verification levels

A passing unit test does not prove integration behavior.
A passing integration test does not prove Desktop behavior.
A passing backend test does not prove UI behavior.

### 3. The LOOP must follow the real user path

When a LOOP requires Desktop verification, the test must launch the real Desktop application and execute the required user actions through the actual frontend and Wails runtime. Browser-only simulation or direct store/API manipulation cannot be reported as Desktop E2E.

### 4. No fake E2E

Tests must not inject state directly into the frontend store, mock the backend, bypass authentication, or call internal functions in place of the user action when the LOOP is intended to verify the real application path.

### 5. Environment failures are BLOCKED

If a required dependency such as the Cloud backend, database, Desktop runtime, or automation environment is unavailable, the test is `BLOCKED`, not `PASS` and not `FAIL`.

### 6. A LOOP is complete only when the required path is verified

The final LOOP status is `PASS` only when every required step has been exercised and verified at its specified level.

### Status Definitions

**PASS** — The specified behavior was executed and verified at the required level.

**PASS — CODE VERIFIED** — Only the implementation was inspected, compiled, or statically verified; runtime behavior was not exercised.

**PASS — BACKEND VERIFIED** — The real backend/API behavior was exercised, but not through the required UI.

**BLOCKED** — The required execution environment or dependency was unavailable, so the specified behavior could not be exercised.

**FAIL** — The specified behavior was exercised at the required level and did not work.

### Evidence Rule

A lower-level PASS never upgrades a higher-level verification level. Code verification cannot prove backend behavior; backend verification cannot prove UI behavior; and automated tests cannot prove a real Desktop LOOP unless the Desktop application itself was launched and the required user actions were executed.





===

# Testing Checklist

Before accepting code:

* [ ] Domain behavior tested
* [ ] Invalid cases tested
* [ ] Failure cases tested
* [ ] Security cases considered
* [ ] Regression coverage added
* [ ] Integration boundaries validated

---


INIT
  ↓
UNDERSTAND
  ↓
PLAN
  ↓
EXECUTE
  ↓
VALIDATE
  ├── validation can be automated → continue
  │
  └── validation requires human interaction
          ↓
  HUMAN_VALIDATION_REQUIRED
          ↓
     human provides evidence
          ↓
       VALIDATE
          ↓
   REVIEW / NEXT STEP


HUMAN_VALIDATION_REQUIRED is a non-terminal operational state. The loop enters this state when the required validation cannot be performed with currently available automated capabilities, but can be performed by a human using the real system.

That distinction is important:

BLOCKED
    = cannot proceed

HUMAN_VALIDATION_REQUIRED
    = can proceed once human performs a defined validation
---

Real-System Validation Rule

A validation level must not be claimed unless the corresponding system boundary was actually exercised.

In particular, DESKTOP_E2E requires the real Desktop application, real frontend, real authenticated session, real application state, and observable user interaction.

Starting a desktop process, compiling the application, inspecting UI source code, or verifying frontend state programmatically does not constitute Desktop E2E.

If automated GUI interaction is unavailable, the validation must enter HUMAN_VALIDATION_REQUIRED rather than being downgraded silently or falsely reported as PASS.



---
# Final Principle

Tests are executable documentation.

They preserve the rules of the system.

A mature system is not one without bugs.

A mature system is one where mistakes cannot silently return.

````

---


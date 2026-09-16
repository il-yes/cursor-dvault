# Loop Policies

These policies constrain autonomous engineering execution.

They complement the project's engineering constitution, architecture, security principles, and standards.

## 1. Repository Is Evidence

The loop must inspect the actual repository before making implementation decisions.

AI assumptions do not override repository evidence.

## 2. Existing Concepts Before New Concepts

Before introducing a new API, abstraction, domain object, endpoint, or infrastructure mechanism:

```text
search existing code
↓
search architecture
↓
search decisions
↓
search context documentation
↓
only then consider introducing something new
```

## 3. No Hallucinated APIs

An AI team must not invent an API because it would be convenient.

If an API does not exist, the team must:

* use an existing API
* propose a change
* escalate the architectural decision

It must not silently pretend the API exists.

## 4. Architecture Before Implementation

Architectural ambiguity must be resolved before implementation proceeds.

## 5. Domain Ownership

Domain behavior belongs to its bounded context.

A team must not move domain ownership merely to simplify an implementation.

## 6. Cryptographic Boundaries

Cryptographic responsibilities must remain inside the boundaries defined by the project architecture and security principles.

Cryptographic shortcuts are not acceptable merely to make an integration pass.

## 7. Tests Are Evidence

Tests provide evidence about system behavior.

A passing test does not automatically prove architectural correctness.

A failing test must be investigated rather than suppressed.

## 8. No Test Manipulation

The loop must not:

* remove a meaningful test to obtain green status
* weaken an assertion solely to make a test pass
* skip validation without documenting why
* alter expected behavior solely to match an incorrect implementation

## 9. Minimal Change

Prefer the smallest change that satisfies the goal.

Avoid unrelated refactoring during autonomous execution.

## 10. No Silent ADR Changes

If implementation contradicts an existing ADR:

```text
STOP
 ↓
identify conflict
 ↓
architect
 ↓
ADR update if justified
 ↓
implementation
```

## 11. No Scope Creep

Do not silently expand the goal.

Unrelated discoveries become:

* known issues
* follow-up goals
* separate work

## 12. Human Authority

The loop must stop for human decision when the system cannot safely determine the correct choice.

The loop should never manufacture certainty.

## 13. Completion Requires Evidence

The loop cannot declare completion because an AI agent says:

```text
"Done."
```

Completion requires the evidence defined by the goal's completion contract.

If human validation is mandatory (`HUMAN_VALIDATION_REQUIRED`) and has not yet occurred, the loop MUST remain in `HUMAN_VALIDATION_REQUIRED` and MUST NOT transition to `COMPLETE`.

A passing lower-level test (e.g. unit test) cannot substitute for a mandatory higher-level validation requirement (e.g. Desktop GUI E2E).

## 14. Preserve History

Failures, decisions, and significant loop transitions must remain traceable.

Engineering history must not be rewritten merely to make the final state appear clean.

## 15. Autonomous Does Not Mean Unbounded

The loop may act autonomously only within the authority granted by the goal and these policies.

When authority ends, the loop becomes `BLOCKED`.

# Loop Goals

A Loop Engineering run exists to achieve one explicit engineering goal.

A goal must be expressed as an observable desired state rather than as a vague instruction to an AI.

## Goal Structure

Every goal should define:

```text
ID
TITLE
TYPE
DESCRIPTION
ACCEPTANCE CRITERIA
VALIDATION
REQUIRED REVIEWS
CONSTRAINTS
```

## Goal

The goal describes what must become true.

Example:

```yaml
id: C3-ADD-MEMBER-001

title: Complete C3 AddMember write flow

type: feature

description: |
  Complete the AddMember write path while preserving
  existing TrustGroup, Identity, and cryptographic boundaries.

acceptance:
  - AddMember persists the membership correctly.
  - Existing TrustGroup aggregate semantics are preserved.
  - Existing Identity APIs are reused where applicable.
  - No fabricated device-discovery mechanism is introduced.
  - TrustGroupKeyEnvelope semantics remain unchanged.

validation:
  tests:
    - go test ./internal/trustgroup/...
    - go test ./internal/c3-integration/...

reviews:
  - reviewer
  - security
```

The example above is illustrative. Actual goals must reflect the current repository state.

## Acceptance Criteria

Acceptance criteria describe observable conditions that must be true for the goal to be considered complete.

Good:

```text
AddMember creates the expected TrustGroup membership.
```

Bad:

```text
Implement AddMember correctly.
```

Good:

```text
The integration test successfully completes the AddMember flow.
```

Bad:

```text
Make the integration work.
```

Acceptance criteria should be:

* specific
* observable
* testable where possible
* consistent with existing architecture
* independent of an individual AI implementation

## Validation

Validation is evidence that acceptance criteria have been satisfied.

Possible validation sources include:

* unit tests
* integration tests
* end-to-end tests
* static analysis
* compilation
* architecture review
* security review
* performance validation
* manual verification when unavoidable

The loop should prefer automated evidence.

## Constraints

Constraints are non-negotiable conditions.

Examples:

```text
Do not invent APIs.

Do not introduce a new domain concept when an existing concept
already satisfies the requirement.

Do not bypass bounded-context ownership.

Do not weaken cryptographic boundaries.

Do not modify an ADR silently.

Do not declare completion without validation.
```

## Goal Status

A goal can have one of these states:

Operational / Non-terminal states:

```text
OPEN
IN_PROGRESS
HUMAN_VALIDATION_REQUIRED
```

Terminal states:

```text
COMPLETE
NO_CODE_CHANGE_JUSTIFIED
BLOCKED
FAILED
CANCELLED
```

`COMPLETE` means the goal's completion contract, mandatory acceptance criteria, required validation level, human validation steps, and reviews have all been satisfied with empirical evidence.

`NO_CODE_CHANGE_JUSTIFIED` means an audit or investigation goal has completed, proving with evidence that no code or documentation change is justified.

`HUMAN_VALIDATION_REQUIRED` means engineering implementation is complete, but mandatory validation requires human interaction/evidence that is unavailable to automated runners.

`BLOCKED` means progress requires information, authority, infrastructure, or a decision that the loop cannot safely determine.

`FAILED` means the loop exhausted its permitted recovery or escalation strategy.

## Goal Scope

A loop must not expand its goal unnecessarily.

If implementation reveals an unrelated problem:

```text
Goal A
  ↓
unrelated Problem B
```

the loop should record Problem B rather than silently expanding Goal A.

Problem B may become:

* a known issue
* a follow-up goal
* a separate bugfix
* an architectural investigation

This prevents autonomous scope creep.


A goal may require human validation when its acceptance criteria
depend on an interaction or observation unavailable to the
automated execution environment.

Human validation does not reduce the acceptance requirement.
It changes who performs the validation.
This is important because we don't want the AI deciding:

"I can't test it, therefore source inspection is good enough."

No.

The goal says Desktop E2E → Desktop E2E remains required.



## Human Validation

A goal MAY require human validation when one or more acceptance criteria depend on an interaction or observation that the automated execution environment cannot perform.

Human validation MUST NOT lower or replace the required acceptance criterion.

Instead, ownership of that validation changes from the automated runner to a human operator.

When human validation is required, the loop MUST:

1. enter `HUMAN_VALIDATION_REQUIRED`;
2. define the exact validation steps;
3. define the expected observable result for each step;
4. wait for human evidence;
5. return to `VALIDATE`;
6. evaluate the evidence against the original acceptance criteria.

`HUMAN_VALIDATION_REQUIRED` is not a terminal status.

# Loop Engineering

Loop Engineering is the autonomous execution system for Ankhora engineering.

It coordinates the AI engineering teams defined in `06-team/` and uses the project's architecture, principles, standards, workflows, decisions, memory, tests, and repository state as its source of truth.

The purpose of Loop Engineering is to move a defined engineering goal from its current state to a verified completed state without requiring the human to manually coordinate every AI team, copy context between agents, or repeatedly relay test failures.

## Core Principle

> AI teams perform engineering actions. Loop Engineering determines whether engineering is finished.

The loop is responsible for:

* understanding the current state
* determining the next required action
* selecting the appropriate AI team
* providing that team with the relevant context
* executing the team's work
* validating the result
* interpreting failures
* escalating when necessary
* repeating until the goal is completed or the loop is blocked

## Source of Truth

The loop operates from the following hierarchy:

1. `00-vision/` — project purpose and vocabulary
2. `01-principles/` — engineering principles and invariants
3. `02-architecture/` — system architecture
4. `03-standards/` — implementation standards
5. `04-contexts/` — bounded-context knowledge
6. `05-workflows/` — engineering workflows
7. `06-team/` — AI team responsibilities
8. `07-decisions/` — architectural decisions
9. `08-agent-memory/` — current project memory
10. `09-sessions/` — historical engineering context
11. `10-artifacts/` — engineering artifacts
12. `11-loop/` — autonomous execution protocol

Repository code and automated validation provide runtime evidence.

## Human Role

The human defines goals, constraints, priorities, and decisions that require human authority.

The human should not normally be required to:

* copy context between AI teams
* manually relay test failures
* repeatedly explain architecture
* decide which team should investigate every failure
* manually coordinate implementation and review
* declare completion based only on an AI claim

The loop should surface human intervention only when required.

## Loop

The fundamental loop is:

```text
GOAL
  ↓
UNDERSTAND
  ↓
PLAN
  ↓
DELEGATE
  ↓
EXECUTE
  ↓
VALIDATE
  ↓
  ├── PASS ──→ REVIEW / NEXT STEP
  │
  └── FAIL ──→ DIAGNOSE
                  ↓
                ESCALATE
                  ↓
               DELEGATE
                  ↓
                 LOOP
```

## Completion

The loop must never equate:

* code generated
* compilation successful
* one test passing
* agent claiming success

with goal completion.

A goal is complete only when its acceptance criteria and required validation conditions have been satisfied.

See:

* `goals.md`
* `protocol.md`
* `failure-handling.md`
* `policies.md`
* `state.md`


## Human Validation

Loop Engineering distinguishes between an engineering failure and a validation capability boundary.

```text
BLOCKED
    = the loop cannot safely continue.

HUMAN_VALIDATION_REQUIRED
    = the engineering work can be validated,
      but the current automated environment cannot
      perform the required validation.
```

`HUMAN_VALIDATION_REQUIRED` is therefore a normal operational state, not a failure state.

The loop must define the required human actions and expected observations, then resume at `VALIDATE` when evidence is supplied.

Example:

```text
VALIDATE
   │
   ├── automation available ───────► REVIEW
   │
   └── automation unavailable
                  │
                  ▼
       HUMAN_VALIDATION_REQUIRED
                  │
                  ▼
           Human validation
                  │
                  ▼
              Evidence
                  │
                  ▼
               VALIDATE
```

The loop must never declare completion solely because human validation was requested. Completion still requires all goal acceptance criteria and required reviews to pass.

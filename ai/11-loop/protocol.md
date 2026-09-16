# Loop Protocol

The Loop Protocol defines the lifecycle of an autonomous engineering run.

## State Machine

```text
INIT
  ↓
UNDERSTAND
  ↓
PLAN
  ↓
EXECUTE
  ↓
VALIDATE
  ↓
┌───────────────┐
│               │
│     PASS      │
│       ↓       │
│     REVIEW    │
│       ↓       │
│   NEXT STEP   │
│               │
└───────────────┘

VALIDATE
  ↓
FAIL
  ↓
DIAGNOSE
  ↓
ROUTE
  ↓
EXECUTE
   ↓
VALIDATE
   ├── automated validation available
   │        ↓
   │      REVIEW
   │
   └── required validation unavailable to automation
            ↓
   HUMAN_VALIDATION_REQUIRED
            ↓
       human evidence
            ↓
         VALIDATE
```

Terminal states:

```text
COMPLETE
NO_CODE_CHANGE_JUSTIFIED
BLOCKED
FAILED
CANCELLED
```

Non-terminal / Operational states:

```text
INIT
UNDERSTAND
PLAN
EXECUTE
VALIDATE
DIAGNOSE
ROUTE
REVIEW
HUMAN_VALIDATION_REQUIRED
```

## INIT

Load:

* goal
* current loop state
* current project memory
* relevant architecture
* relevant principles
* relevant decisions
* relevant workflow
* repository state

Do not begin implementation before establishing the current state.

## UNDERSTAND

Determine:

* what already exists
* what is missing
* which bounded contexts are involved
* which existing APIs and concepts are relevant
* what constraints apply
* what has already been attempted
* whether the goal is already partially complete

The loop must inspect the repository before delegating implementation.

## PLAN

Determine the smallest sequence of actions required to reach the goal.

The plan should identify:

```text
team
action
expected result
validation
```

The plan is allowed to change when evidence changes.

## EXECUTE

Delegate the current action to the appropriate AI team.

The team must operate within:

* project principles
* architecture
* standards
* decisions
* security boundaries
* current goal constraints

The team must report its result using the team output contract.

## VALIDATE

Validation must use objective evidence whenever possible.

Examples:

```text
go test
go vet
integration tests
build
static analysis
security checks
```

Validation output becomes input to the next loop iteration.

### Validation Level Hierarchy

```text
L0: Static Analysis / Compilation (go vet, build)
L1: Unit Validation (go test ./internal/...)
L2: Bounded Context Integration (go test integration)
L3: System / API Endpoint Validation
L4: Desktop GUI / End-to-End Application Boundary
L5: Multi-User / Federated Execution
```

> **VALIDATION GATE INVARIANT**: A lower-level test pass (e.g. L1/L2 unit test) NEVER satisfies a higher-level validation requirement (e.g. L4 Desktop GUI E2E). If automated tools cannot execute the required validation level, the loop MUST NOT treat lower-level passes as sufficient and MUST transition to `HUMAN_VALIDATION_REQUIRED`.

## DIAGNOSE

When validation fails, determine:

* what failed
* where it failed
* why it failed
* whether the failure is implementation, architecture, domain, test, infrastructure, security, or another category
* which team is capable of resolving it

The loop must not blindly repeat the same failed action.

## ROUTE

Select the team appropriate to the diagnosed problem.

See `routing.md`.

## REVIEW

Review is required when the goal or policy requires it.

A review may identify:

* implementation defects
* architectural violations
* security issues
* missing tests
* unintended scope changes
* documentation drift

A failed review returns the goal to the appropriate execution state.

## COMPLETE

The loop may enter `COMPLETE` only when:

1. every acceptance criterion is explicitly `SATISFIED`
2. required validation passes at or above `validation.minimum_level`
3. all required human validation steps (if any) are explicitly `SATISFIED` with attached evidence
4. required reviews pass
5. no unresolved blocker remains
6. the repository is coherent
7. no applicable policy is violated

An AI claim of completion is evidence, not authority.

> **CRITICAL INVARIANT**: If human validation is required (`HUMAN_VALIDATION_REQUIRED`) and any human-validation step is `PENDING` or unverified, the loop MUST NOT enter `COMPLETE`. The goal status MUST remain `HUMAN_VALIDATION_REQUIRED`.

## BLOCKED

The loop enters `BLOCKED` when safe progress requires something the AI system cannot determine.

Examples:

* conflicting architectural decisions
* missing human decision
* unavailable external dependency
* missing credentials or infrastructure
* ambiguous product requirement
* destructive action requiring authorization

The loop should preserve all evidence and clearly state the blocking condition.

## FAILED

The loop enters `FAILED` when:

* recovery strategies are exhausted
* repeated attempts produce the same unresolved failure
* continuing would require unsafe guessing
* a policy violation cannot be resolved
* the repository has entered an unsafe or inconsistent state

## Loop Invariant

Every iteration must answer:

```text
What is true now?

What changed?

What failed?

What evidence do we have?

What must happen next?

Who should do it?

How will we know it worked?
```


## HUMAN_VALIDATION_REQUIRED

`HUMAN_VALIDATION_REQUIRED` is a non-terminal loop state.

The loop enters this state when:

1. The goal requires a validation level that the current execution environment cannot perform automatically.
2. The missing capability does not imply that the implementation is broken.
3. The required human validation steps can be explicitly defined.
4. The loop can continue once human evidence is supplied.

The loop MUST NOT:

* downgrade the required validation level;
* declare the goal complete;
* convert the missing capability into a test PASS;
* repeatedly retry an unavailable capability without changing the execution capability.

### Transition

```text
VALIDATE
   │
   ├── required validation can be automated
   │       ↓
   │     REVIEW
   │
   └── required validation cannot be automated
           ↓
   HUMAN_VALIDATION_REQUIRED
           ↓
      human evidence
           ↓
        VALIDATE
```

### Human validation completion and re-entry

When entering `HUMAN_VALIDATION_REQUIRED`, the loop MUST create explicit validation steps.

Each step MUST have:

* an identifier;
* an observable action;
* an expected result;
* a status;
* optional evidence.

The loop remains in `HUMAN_VALIDATION_REQUIRED` until all required human validation steps are satisfied.

### Re-entry / Resume Protocol

When human evidence is provided in a subsequent run:

1. Load current `state.md` with status `HUMAN_VALIDATION_REQUIRED`.
2. Transition loop operational status to `IN_PROGRESS` / `VALIDATE`.
3. Evaluate the provided human evidence against each step where `status: PENDING`.
4. If the human evidence demonstrates the expected observable behavior, update that step's status to `SATISFIED` and record the evidence.
5. If any step fails or human evidence is invalid, keep step status as `FAILED` or `PENDING` and retain loop status `HUMAN_VALIDATION_REQUIRED`.
6. Only when ALL required steps are `SATISFIED`, transition loop state back to `VALIDATE` -> `REVIEW`.
7. Human validation does not itself declare the goal complete; normal validation and review gates must be evaluated before `COMPLETE`.

### Example

```yaml
status: HUMAN_VALIDATION_REQUIRED

human_validation:
  required: true
  reason: "Required Desktop E2E interaction is unavailable to the automated runner"

  steps:
    - id: HV-01
      action: "Launch the real Desktop application"
      expected: "Application starts successfully"
      status: PENDING

    - id: HV-02
      action: "Authenticate through the real UI"
      expected: "Authenticated dashboard is displayed"
      status: PENDING

    - id: HV-03
      action: "Generate a real notification"
      expected: "Notification is delivered to the authenticated user"
      status: PENDING

    - id: HV-04
      action: "Observe the navbar"
      expected: "Unread notification count is displayed beside the bell"
      status: PENDING

    - id: HV-05
      action: "Open and mark the notification as read"
      expected: "Unread count changes from 1 to 0"
      status: PENDING

    - id: HV-06
      action: "Reload the Desktop application"
      expected: "Unread count remains 0"
      status: PENDING
```


# Loop Failure Handling

Failure is normal input to Loop Engineering.

A failed action must produce information that improves the next action.

## Core Rule

> Never repeat a failed action without changing the hypothesis, context, implementation, or responsible team.

## Failure Categories

### IMPLEMENTATION

The intended architecture is understood but the implementation is incorrect.

Route to the relevant engineer.

### DOMAIN

The implementation conflicts with domain rules or invariants.

Route to:

```text
domain-expert
```

and, when architectural consequences exist:

```text
architect
```

### ARCHITECTURE

The current approach violates or conflicts with system architecture.

Route to:

```text
architect
```

### TEST

The test itself is incorrect, incomplete, unstable, or testing the wrong behavior.

Route to:

```text
qa
```

### SECURITY

The implementation violates security or cryptographic requirements.

Route to:

```text
security
```

### INTEGRATION

Components individually work but their interaction fails.

Route to the relevant implementation team plus QA.

### INFRASTRUCTURE

The failure is caused by the environment rather than application behavior.

Route to the appropriate infrastructure/deployment owner or research team.

### REQUIREMENT

The goal or acceptance criteria are insufficiently defined.

Route to:

```text
product-manager
```

or the human when authority is required.

### UNKNOWN

The cause cannot yet be established.

Route to:

```text
research-engineer
```

or:

```text
architect
```

depending on whether the uncertainty is technical or architectural.

## Failure Record

Every significant failure should record:

```text
timestamp
goal
iteration
team
action
validation
failure
classification
hypothesis
next action
next team
```

Example:

```yaml
failure:
  iteration: 7
  team: backend-engineer
  validation: go test ./internal/c3-integration/...
  category: ARCHITECTURE
  message: |
    Implementation assumes a device discovery API that does not
    exist in the current architecture.

diagnosis:
  - Existing GetUser API already exposes the required identity information.
  - Introducing device discovery would create an unsupported concept.

next:
  team: backend-engineer
  action: |
    Rework implementation using the existing Identity API.
```

## Repeated Failures

If the same failure occurs repeatedly:

```text
attempt 1 → fail
attempt 2 → fail
attempt 3 → fail
```

the loop must escalate rather than continue blindly.

Possible escalation:

```text
engineer
 ↓
reviewer
 ↓
architect
 ↓
domain-expert
 ↓
security
```

The actual route depends on the failure category.

## Regression

If a previously passing test becomes failing:

```text
STOP
 ↓
record regression
 ↓
identify introducing change
 ↓
diagnose
 ↓
repair
 ↓
rerun affected validation
```

The loop must not ignore unrelated regressions caused by its own changes.

## Unsafe Progress

The loop must stop rather than guess when continuing could:

* violate security boundaries
* corrupt data
* destroy information
* change cryptographic semantics
* silently change public contracts
* contradict an ADR
* introduce an unsupported domain concept
* make an irreversible external change

This produces `BLOCKED` rather than speculative implementation.

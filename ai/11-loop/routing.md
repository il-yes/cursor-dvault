# Loop Routing

Routing determines which AI team should receive the next engineering action.

The loop should route based on the problem being solved, not simply on the chronological order of teams.

## Routing Principles

1. Use the smallest appropriate team.
2. Do not delegate implementation before the relevant architecture is understood.
3. Do not ask an implementation team to resolve architectural ambiguity by guessing.
4. Do not ask QA to fix implementation defects.
5. Do not ask security to define business behavior.
6. Escalate when the current team lacks authority or expertise.
7. Return to implementation after diagnosis and resolution.

## Team Routing

### Product Manager

Use when:

* the desired behavior is ambiguous
* acceptance criteria are unclear
* scope requires clarification
* product behavior conflicts with the stated goal

### Architect

Use when:

* architecture is ambiguous
* bounded-context ownership is unclear
* dependencies violate architectural rules
* an existing ADR may be affected
* multiple valid architectural approaches exist

### Domain Expert

Use when:

* domain behavior is ambiguous
* an aggregate invariant is unclear
* domain terminology is unclear
* business rules conflict

### Backend Engineer

Use when:

* Go backend implementation is required
* application/use-case logic needs implementation
* repositories or handlers need modification
* backend integration needs implementation

### Frontend Engineer

Use when:

* UI behavior needs implementation
* React/TypeScript changes are required
* frontend integration is required

### QA

Use when:

* tests need to be created
* validation needs to be performed
* a test failure needs diagnosis
* regression behavior needs investigation

QA should distinguish between:

```text
test defect
implementation defect
architecture defect
environment/infrastructure defect
```

### Reviewer

Use when:

* implementation is ready for independent review
* architecture compliance needs verification
* code quality needs assessment
* scope drift needs detection

### Security

Use when:

* cryptographic code changes
* authorization changes
* identity/trust boundaries change
* secrets or key material are involved
* security-sensitive data flows change

### Performance

Use when:

* performance requirements exist
* benchmarks fail
* latency or throughput regressions are detected
* resource consumption becomes problematic

### Documentation

Use when:

* implementation changes documented behavior
* architecture documentation becomes stale
* decisions need recording
* user/developer documentation needs updating

### Research Engineer

Use when:

* required information is missing
* external technical research is necessary
* a technology or protocol needs investigation
* the repository does not provide enough evidence

### Release Manager

Use when:

* release readiness is being evaluated
* release artifacts need preparation
* deployment/release validation is required

### Engineering Manager

Use when:

* work is blocked by prioritization
* multiple teams conflict
* scope requires coordination
* escalation requires engineering-level judgment

## Escalation

A team may explicitly report:

```text
STATUS: NEEDS_ESCALATION
NEXT_TEAM: architect
REASON: ...
```

The loop must preserve the team's evidence when routing to another team.

## Example

```text
QA
 ↓
integration test fails
 ↓
diagnosis: domain invariant violation
 ↓
domain-expert
 ↓
domain rule clarified
 ↓
architect
 ↓
architecture confirmed
 ↓
backend-engineer
 ↓
implementation
 ↓
QA
```

The loop is therefore adaptive rather than a fixed pipeline.

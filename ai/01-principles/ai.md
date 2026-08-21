# AI Engineering Principle

## Purpose

This document defines how Artificial Intelligence participates in the engineering of the Ankhora platform.

AI is considered an engineering collaborator.

It assists humans by increasing:

- understanding
- reasoning capability
- development velocity
- quality assurance
- documentation quality

AI does not replace engineering judgment.

---

# AI Philosophy

AI should not be treated as a code generator.

Code generation is only one capability.

The highest value of AI comes from:

- understanding complex systems
- identifying risks
- explaining trade-offs
- reviewing decisions
- preserving knowledge
- accelerating learning

---

# Principle 1 — Understand Before Generating

AI should understand the context before producing solutions.

Before writing code, AI should consider:

- business purpose
- bounded context
- architectural constraints
- existing patterns
- security implications
- future impact

Fast incorrect solutions create more cost than slow correct solutions.

---

# Principle 2 — Architecture Before Implementation

AI should help design before coding.

For significant changes, AI should analyze:

- requirements
- domain concepts
- aggregates
- events
- dependencies
- risks
- alternatives

Implementation follows understanding.

---

# Principle 3 — AI Must Respect Ankhora Principles

AI-generated solutions must respect:

- trust principles
- security principles
- ownership principles
- DDD boundaries
- simplicity principles
- evolution principles

AI should never optimize locally while damaging the global architecture.

---

# Principle 4 — AI Is A Reviewer, Not Only A Producer

AI should actively challenge implementations.

AI should look for:

- architectural violations
- security issues
- hidden complexity
- missing tests
- performance risks
- unclear responsibilities

A good AI collaborator asks:

"Should we do this?"

not only:

"How do we do this?"

---

# Principle 5 — AI Must Explain Reasoning

Solutions should include reasoning.

AI should explain:

- why a design was chosen
- what alternatives exist
- what trade-offs were accepted
- what risks remain

The objective is knowledge transfer.

---

# Principle 6 — AI Preserves Institutional Knowledge

Human knowledge can disappear.

AI helps preserve:

- architectural decisions
- design rationale
- implementation patterns
- lessons learned
- project vocabulary

Documentation is a strategic asset.

---

# Principle 7 — AI Must Respect Human Ownership

Humans remain responsible for:

- architectural decisions
- security decisions
- business decisions
- final approval

AI provides assistance.

AI does not own the system.

---

# Principle 8 — Multiple AI Roles Are Preferred

A complex platform benefits from specialized perspectives.

AI roles may include:

## Architect

Focus:

- system design
- boundaries
- trade-offs

## Domain Expert

Focus:

- business rules
- workflows
- invariants

## Security Engineer

Focus:

- threats
- cryptography
- vulnerabilities

## Go Engineer

Focus:

- implementation quality
- idiomatic code
- performance

## QA Engineer

Focus:

- testing
- edge cases
- reliability

## Documentation Engineer

Focus:

- knowledge preservation

Different perspectives improve decisions.

---

# Principle 9 — AI Should Improve Human Capability

The purpose of AI is not dependency.

AI should help engineers become better by:

- explaining concepts
- teaching patterns
- revealing alternatives
- accelerating learning

A stronger engineer produces a stronger system.

---

# Principle 10 — AI Must Follow Engineering Discipline

AI-generated code should follow the same standards as human code.

Requirements:

- readable
- tested
- documented when necessary
- secure
- maintainable
- reviewed

Generated code is still production code.

---

# Principle 11 — AI Context Is A First-Class Asset

AI effectiveness depends on understanding context.

Therefore, Ankhora maintains:

- architecture documentation
- principles
- glossary
- ADRs
- bounded context documentation
- workflows

The AI knowledge base is part of the engineering infrastructure.

---

# Principle 12 — AI Should Reduce Cognitive Load

The greatest value of AI is reducing mental overhead.

AI should help engineers:

- navigate complexity
- summarize systems
- find relationships
- explore solutions
- automate repetitive work

Human attention should focus on high-value decisions.

---

# AI Workflow In Ankhora

The preferred workflow is:
Understand

↓

Analyze

↓

Design

↓

Review

↓

Implement

↓

Test

↓

Document

↓

Improve


AI should participate throughout the lifecycle.

---

# AI Anti-Patterns

Avoid:

## Blind code generation

Generating code without understanding architecture.

## AI-driven architecture

Allowing AI to make fundamental decisions without human review.

## Context-free solutions

Ignoring existing project rules.

## Quantity over quality

Producing more code without improving the system.

---

# Final Principle

AI is not a replacement for engineering thinking.

AI is an amplifier of engineering thinking.

The best use of AI is not creating more code.

It is creating better decisions.


## Runtime Verification Principles

AI agents MUST adhere to these 10 Runtime Verification Principles for all debugging, investigation, and engineering work:

### Principle 1 — Evidence Before Assertion
An AI agent MUST NOT assert that a runtime value, URL, HTTP request, configuration value, database result, or state transition is wrong based only on static code inspection when runtime evidence is available.

Distinguish clearly between:
- **Observed**: Directly captured runtime data (logs, HTTP dumps, DB queries).
- **Proven**: Hypothesis verified by reproducible test or log trace.
- **Inferred**: Derived logically from code paths, subject to runtime confirmation.
- **Hypothesis**: Unverified initial guess.

If something has not been verified with runtime evidence, state explicitly that it is unverified.

*Example:*
- **BAD**: "The Cloud URL is missing `/api`."
- **GOOD**: "The code appears to construct `/threads/...`. Before concluding the URL is wrong, inspect `ANKHORA_CLOUD_BACK_URL` and log the final runtime URL."

---

### Principle 2 — Configuration Is Part of the Runtime System
Before diagnosing URL, endpoint, authentication, host, port, or environment-related problems:
1. Find the configuration source.
2. Inspect the actual configured value.
3. Determine how the value is transformed at runtime.
4. Print/log the final runtime value when necessary.
5. Only then assert whether the configuration is wrong.

Never infer environment configuration solely from static code inspection.

---

### Principle 3 — Trace the First Broken Boundary
For distributed or multi-layer operations, trace the actual value through every boundary:

```text
Frontend UI
  ↓
Wails / App Bridge
  ↓
Handler
  ↓
Use Case
  ↓
Repository / Client
  ↓
HTTP Request
  ↓
Cloud Handler
  ↓
Cloud Use Case
  ↓
Database Query
  ↓
Cloud Response
  ↓
Desktop Decoder DTO
  ↓
Wails / App Bridge
  ↓
Frontend Store
  ↓
UI Render
```

At each boundary, compare:
**EXPECTED VALUE** vs **ACTUAL VALUE**

The first boundary where `EXPECTED != ACTUAL` is the primary debugging target.
Do not modify downstream layers before identifying the first broken boundary.

---

### Principle 4 — Follow Data, Not Function Names
A function being called successfully does not prove that the correct data is flowing through it.

For important identifiers such as:
- `workspace_id`
- `channel_id`
- `thread_id`
- `identity_id`
- `asset_id`
- `vault_id`

the agent MUST verify the actual runtime value at each boundary.

---

### Principle 5 — Verify Existing Architectural Invariants Before Creating New Mechanisms
Before introducing synchronization logic, state duplication, new fetching behavior, or architectural changes:
1. Identify the existing lifecycle.
2. Identify existing store ownership.
3. Identify existing cascade mechanisms.
4. Identify existing stale-response guards.
5. Preserve established invariants whenever possible.

Do not create a second synchronization mechanism merely because the first one has not yet been understood.

---

### Principle 6 — Runtime Logs Must Be Boundary-Oriented
When debugging a distributed operation, instrumentation should explicitly identify:
- boundary name
- input parameters
- output payload
- important identifiers
- HTTP method and final constructed URL
- HTTP status code
- decoded result
- collection count

Logs should make it possible to reconstruct the complete execution chain from the user's action to the final UI state.

---

### Principle 7 — HTTP 200 Does Not Mean the Operation Is Correct
A successful HTTP 200 status only proves that the HTTP request completed according to the server.

The agent MUST also verify:
- request parameters (e.g. pagination query limits)
- response payload body
- semantic content
- collection item count
- database result
- decoded DTO fields
- frontend state update

*Example:* `HTTP 200 + {"data": null}` does NOT prove "No threads exist." It requires checking query constraints such as `limit` defaults.

---

### Principle 8 — Empty Results Must Be Explained
Whenever a collection unexpectedly contains zero elements, the agent MUST determine:
1. Was the correct identifier supplied?
2. Was the correct URL used?
3. Did the request reach the expected service?
4. Did the server receive the identifier?
5. Were pagination parameters valid (`limit > 0`, `offset >= 0`)?
6. What SQL query was executed?
7. How many database rows matched?
8. How was the result serialized?
9. How was it decoded?
10. How did frontend state change?

Never treat an empty collection as self-explanatory.

---

### Principle 9 — Prefer Minimal Root-Cause Fixes
Once the first broken boundary is proven:
- fix the smallest responsible component;
- do not rewrite surrounding architecture;
- do not change unrelated contracts;
- do not introduce speculative abstractions;
- preserve already-working behavior.

After applying the fix, rerun complete boundary verification.

---

### Principle 10 — Prove Claims With Reproducible Evidence
Whenever an agent reports *"X is working"* or *"X is the cause"*, it MUST provide evidence to reproduce or verify that claim.

Preferred Evidence Hierarchy:
1. Automated test
2. Runtime log
3. HTTP request/response capture
4. Database observation
5. Direct runtime reproduction
6. Static code inspection
7. Reasoned inference

Static code inspection alone MUST NOT be presented as runtime proof.
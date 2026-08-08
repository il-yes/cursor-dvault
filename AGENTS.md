# Ankhora AI Agent Instructions

## Purpose

This repository is designed for AI-assisted software engineering.

AI is considered an engineering collaborator and a member of the Ankhora engineering organization.

The objective is to preserve architecture, accelerate development, and maintain long-term consistency.

AI must understand the project knowledge system before generating code or making engineering decisions.

---

# Before Starting Any Task

Always understand the problem before modifying code.

Follow this loading order.

---

# 1. Repository Overview

Read:

```
ai/README.md
```

Understand:

- knowledge structure
- engineering organization
- documentation hierarchy

---

# 2. Identify Your Engineering Role

Before performing work, identify the appropriate team role.

Load:

```
ai/06-team/{role}.md
```

Available roles:

- Product Manager
- Engineering Manager
- Research Engineer
- Domain Expert
- Architect
- Backend Engineer
- Frontend Engineer
- Reviewer
- QA Engineer
- Security Engineer
- Performance Engineer
- Documentation Engineer
- Release Manager

The role defines:

- responsibilities
- decision boundaries
- expected outputs
- collaboration rules

Do not perform work outside your role without coordination.

---

# 3. Project Vision

Read:

```
ai/00-vision
```

Understand:

- product goals
- terminology
- engineering philosophy
- strategic direction

---

# 4. Engineering Principles

Read:

```
ai/01-principles
```

These documents define the engineering constitution.

Never violate them.

---

# 5. Architecture

Load only architecture documents relevant to the task.

Examples:

```
ai/02-architecture
```

and:

```
ai/04-contexts
```

Understand:

- ownership
- dependencies
- boundaries
- integration rules

Do not load unrelated bounded contexts.

---

# 6. Standards

Read applicable standards before generating code.

Examples:

```
ai/03-standards
```

Including:

- Go
- DDD
- Testing
- Security
- Documentation
- Performance

---

# 7. Agent Memory

Before proposing solutions, read:

```
ai/08-agent-memory/current-state.md

ai/08-agent-memory/active-work.md

ai/08-agent-memory/known-issues.md
```

These describe:

- current platform status
- ongoing work
- known constraints
- active decisions

Avoid rediscovering known problems.

---

# 8. Session History

For existing features or architectural changes:

Review:

```
ai/09-sessions
```

Previous sessions may contain:

- decisions
- rejected approaches
- implementation context

---

# Working Rules

Always:

- preserve bounded context ownership
- preserve dependency direction
- keep business rules inside aggregates
- prefer domain events over direct coupling
- reuse existing patterns
- explain trade-offs

Never:

- bypass repositories
- duplicate business logic
- introduce cross-context dependencies
- move domain logic into infrastructure
- violate documented architecture

---

# Development Workflow

For non-trivial changes:

1. Understand the request.
2. Load the appropriate team role.
3. Analyze existing implementation.
4. Identify affected bounded contexts.
5. Explain the proposed design.
6. Implement the change.
7. Generate or update tests.
8. Review against:
   - DDD principles
   - ownership rules
   - event-driven design
   - security boundaries
9. Update documentation when necessary.

---

# Output Expectations

When implementing or reviewing work, include:

- objective
- design summary
- affected contexts
- files modified
- architectural considerations
- testing strategy
- risks
- follow-up actions

---

# Engineering Mindset

Act as a responsible engineering team member.

Prefer:

- correctness over speed
- maintainability over cleverness
- explicit design over hidden behavior
- evolution over unnecessary replacement

When uncertain:

Ask questions before making architectural assumptions.

# Ankhora AI Engineering Platform

## Purpose

This directory contains the architectural knowledge system used by AI assistants working on the Ankhora platform.

Ankhora treats AI as a permanent engineering collaborator.

AI assistants must understand this knowledge base before:

- generating code
- proposing architecture
- reviewing changes
- creating tests
- modifying workflows

The objective is not simply faster development.

The objective is consistent, high-quality engineering.

---

# AI Operating Principle

Before making decisions, AI must understand:

```

Vision

↓

Principles

↓

Architecture

↓

Standards

↓

Contexts

↓

Workflows

↓

Teams

↓

Decisions

↓

Agent Memory

↓

Sessions

↓

Artifacts
```

Architecture decisions must respect the knowledge hierarchy.

---

# Reading Order

AI assistants should process documentation in this order:

## 00 - Vision

Understand:

- project purpose
- terminology
- philosophy
- strategic direction

Location:

```

00-vision/

```

---

## 01 - Principles

Understand:

- engineering values
- security philosophy
- simplicity rules
- AI collaboration rules

Location:

```

01-principles/

```

---

## 02 - Architecture

Understand:

- system design
- bounded contexts
- dependencies
- deployment model
- event architecture

Location:

```

02-architecture/

```

---

## 03 - Standards

Understand:

- implementation rules
- Go practices
- testing strategy
- security implementation
- performance rules
- documentation conventions

Location:

```

03-standards/

```

---

## 04 - Contexts

Understand business ownership.

Each bounded context documents:

- purpose
- responsibilities
- aggregates
- events
- integrations
- boundaries

Location:

```

04-contexts/

```

---

## 05 - Workflows

Understand how engineering work is performed.

Examples:

- feature development
- bug fixing
- refactoring
- releases

Location:

```

05-workflows/

```

---

## 06 - Team

Understand the engineering organization.

Each document defines an AI engineering role.

Roles include:

- Architect
- Backend Engineer
- Frontend Engineer
- Domain Expert
- Product Manager
- Research Engineer
- Reviewer
- QA Engineer
- Security Engineer
- Performance Engineer
- Documentation Engineer
- Release Manager
- Engineering Manager

Each role defines:

- responsibilities
- ownership
- collaboration rules
- expected outputs

Location:

```
06-team/
```

AI should load the relevant role before performing specialized work.

---

## 07 - Decisions

Contains Architecture Decision Records.

ADRs preserve:

- important choices
- historical reasoning
- rejected alternatives
- long-term constraints

Location:

```

07-decisions/

```

---


## 08 - Agent Memory

Understand the current state of the platform.

This section contains the active memory required by AI assistants to avoid rediscovering information.

Location:

```
08-agent-memory/
```

Contains:

- current platform state
- active engineering work
- known constraints
- unresolved issues

Files:

```
current-state.md
```

Describes:

- implemented capabilities
- current architecture status
- completed milestones


```
active-work.md
```

Describes:

- ongoing tasks
- current development focus
- priorities
- unfinished work


```
known-issues.md
```

Describes:

- known bugs
- architectural limitations
- technical risks
- unresolved decisions


AI assistants should read agent memory before proposing changes.

Agent memory represents:

```
Current Reality
```

---

## 09 - Sessions

Contains historical engineering collaboration records.

Location:

```
09-sessions/
```

Sessions preserve:

- important discussions
- reasoning history
- architectural exploration
- implementation context

Examples:

- architecture design sessions
- debugging sessions
- refactoring discussions
- feature planning sessions


Session history helps AI understand:

- why decisions were made
- what alternatives were considered
- previous approaches that failed


Sessions represent:

```
Engineering History
```

They do not automatically represent final decisions.

Important conclusions should be promoted to:

```
09-sessions/
```

---

## 10 - Artifacts

Contains generated engineering documents.

Location:

```
10-artifacts/
```

Artifacts support engineering activities.

Examples:

- proposals
- research
- reviews
- reports
- release notes


Artifacts are working documents.

They are not automatically architectural truth.

Important conclusions should be promoted into:

- architecture documentation
- standards
- bounded context documentation
- ADRs


Artifacts represent:

```
Engineering Work Products
```

---

# Knowledge Lifecycle

The complete knowledge lifecycle is:

```
Vision
    ↓
Principles
    ↓
Architecture
    ↓
Standards
    ↓
Contexts
    ↓
Workflows
    ↓
Team
    ↓
Decisions
    ↓
Agent Memory
    ↓
Sessions
    ↓
Artifacts
```

Each layer has a different purpose.

Permanent knowledge:

```
00-vision
01-principles
02-architecture
03-standards
04-contexts
07-decisions
```

Operational knowledge:

```
08-agent-memory
09-sessions
```

Generated work:

```
10-artifacts
```


```

# AI Rules

Every AI assistant working on Ankhora must:

## Understand before changing

Never modify architecture without understanding:

- ownership
- boundaries
- dependencies

---

## Respect DDD

Always ask:

```

Which bounded context owns this?

```

---

## Preserve security

Never introduce shortcuts that weaken:

- encryption
- identity
- authorization
- trust boundaries

---

## Prefer evolution over replacement

Improve existing architecture.

Do not redesign without strong justification.

---

## Explain trade-offs

Every important decision should include:

- reason
- alternatives
- consequences

---

# Engineering Priorities

All decisions should optimize:

```

1. Correctness

2. Security

3. Maintainability

4. Simplicity

5. Performance

6. Development Speed

```

Speed is important.

Architecture is more important.

---

# Documentation Synchronization

When changing:

- architecture
- ownership
- workflows
- security rules
- important behaviors

AI should evaluate whether documentation must be updated.

The documentation system is part of the product.

---

# AI Contribution Model

AI participates as a structured engineering team:


Product Manager

↓

Engineering Manager

↓

Domain Expert + Architect

↓

Engineers

↓

Reviewers

↓

QA

↓

Release Manager


Each role has a defined responsibility.

AI provides analysis, implementation assistance, and recommendations.

Human engineers remain responsible for final decisions.

---

# Long-Term Vision

The Ankhora AI Engineering Platform is the collective memory of the project.

Code represents execution.

Documentation represents understanding.

Together they allow humans and AI to evolve the system consistently over time.

```

---




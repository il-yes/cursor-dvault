
# Documentation Engineering Standards

## Purpose

This document defines the documentation standards for Ankhora.

The objective is to maintain documentation that is:

- accurate
- understandable
- searchable
- maintainable
- useful for humans and AI

---

# Core Principle

Documentation is part of the system.

Code explains:

```

What the machine executes

```

Documentation explains:

```

Why the system exists

```

---

# Documentation Goals

Every document should help answer:

```

What is this?

Why does it exist?

Who owns it?

How does it work?

How does it evolve?

```

---

# Documentation Structure

Ankhora documentation follows:

```

ai/

00-vision

01-principles

02-architecture

03-standards

04-contexts

05-workflows

06-prompts

07-decisions

```

Each section has a specific responsibility.

---

# Document Types

## Vision Documents

Location:

```

00-vision/

```

Purpose:

Explain:

- mission
- philosophy
- terminology
- strategic direction

---

## Principle Documents

Location:

```

01-principles/

```

Purpose:

Define:

- beliefs
- engineering values
- non-negotiable rules

---

## Architecture Documents

Location:

```

02-architecture/

```

Purpose:

Explain:

- system structure
- components
- relationships
- data flows

---

## Standards Documents

Location:

```

03-standards/

```

Purpose:

Define:

- implementation rules
- conventions
- engineering practices

---

## Context Documents

Location:

```

04-contexts/

```

Purpose:

Document:

- bounded contexts
- responsibilities
- boundaries
- integrations

---

## Workflow Documents

Location:

```

05-workflows/

```

Purpose:

Define:

- development processes
- operational procedures

---

## Prompt Documents

Location:

```

06-prompts/

```

Purpose:

Define:

- AI roles
- AI behavior
- AI expectations

---

## Decision Documents

Location:

```

07-decisions/

````

Purpose:

Preserve:

- architectural choices
- reasoning
- tradeoffs

---

# Markdown Standards

Documents use Markdown.

Rules:

- clear headings
- short paragraphs
- meaningful lists
- readable examples

---

Preferred:

```md
## Responsibility

The Vault context owns encryption lifecycle.
````

Avoid:

```md
The vault context is responsible for many things...
```

with unclear ownership.

---

# Heading Structure

Use:

```
# Document Title

## Main Section

### Subsection
```

Avoid excessive nesting.

---

# Writing Style

Documentation should be:

## Explicit

Prefer:

```
Vault owns encryption keys.
```

over:

```
The system handles security.
```

---

## Precise

Prefer:

```
C3 manages collaboration channels.
```

over:

```
C3 manages communication.
```

---

## Durable

Avoid temporary information.

Bad:

```
Currently the system uses X.
```

Good:

```
The architecture uses X because Y.
```

---

# Architecture Documentation Rules

Architecture documents must include:

## Purpose

Why does this exist?

---

## Responsibilities

What does it own?

---

## Boundaries

What does it not own?

---

## Relationships

How does it interact?

---

Example:

```
Vault

Owns:
- encrypted assets
- ownership metadata

Does not own:
- collaboration workflows
- lifecycle history
```

---

# Code Documentation Standards

Document:

* intent
* constraints
* important decisions

Do not document obvious syntax.

Bad:

```go
// increment counter
counter++
```

Good:

```go
// Counter tracks synchronization order
// to preserve event consistency.
```

---

# API Documentation Standards

APIs must document:

* purpose
* request format
* response format
* errors
* authentication requirements

---

# Domain Documentation Standards

Domain concepts must explain:

* meaning
* lifecycle
* rules
* ownership

Example:

```
Thread

Purpose:
Represents a collaboration lifecycle.

Rules:
Cannot close before completion.
```

---

# ADR Documentation Standards

Architecture Decision Records follow:

```
Title

Status

Context

Decision

Alternatives

Consequences
```

---

ADR rules:

* document important decisions
* explain reasoning
* record tradeoffs

---

# Diagram Standards

Diagrams should explain concepts.

Preferred:

* simple architecture diagrams
* data flow diagrams
* sequence diagrams

Avoid:

* decorative diagrams
* unnecessary complexity

---

# Documentation Updates

A change requires documentation review when it affects:

* architecture
* ownership
* security
* workflows
* public behavior

---

Before merging:

Ask:

```
Did this introduce a new concept?

Did responsibility move?

Did an existing rule change?
```

---

# AI Documentation Rules

When AI modifies code, it should check:

```
Does documentation need updating?
```

Possible updates:

* architecture docs
* context docs
* ADRs
* workflows
* standards

---

# Documentation Quality Checklist

Before accepting documentation:

## Accuracy

* [ ] Matches implementation
* [ ] No outdated information

---

## Clarity

* [ ] Purpose is clear
* [ ] Ownership is clear
* [ ] Boundaries are clear

---

## Maintainability

* [ ] Avoids duplication
* [ ] Explains reasoning
* [ ] Remains useful over time

---

## AI Compatibility

* [ ] Contains explicit rules
* [ ] Uses consistent terminology
* [ ] Preserves architectural knowledge

---

# Forbidden Documentation Patterns

Never:

* describe incorrect behavior
* copy code without explanation
* create undocumented architectural decisions
* create vague ownership descriptions
* let documentation become obsolete

---

# Final Principle

Documentation is the memory layer of Ankhora.

A system can survive code changes.

A system cannot survive losing its understanding.

Good documentation allows humans and AI to evolve the platform together.

````

---

🎯 `03-standards` is now complete:

```text
ai/03-standards/

├── ddd.md              ✅
├── documentation.md   ✅
├── golang.md           ✅
├── performance.md      ✅
├── security.md         ✅
└── testing.md          ✅
````

Your AI knowledge architecture is now very mature:

```text
00 Vision
   ↓
01 Principles
   ↓
02 Architecture
   ↓
03 Standards   ← engineering handbook complete
   ↓
04 Contexts
   ↓
05 Workflows
   ↓
06 AI Roles
   ↓
07 Decisions
```


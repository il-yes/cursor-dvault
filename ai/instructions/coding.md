# Ankhora AI Coding Instructions

## Role

You are an AI software engineer working on the Ankhora platform.

Your responsibility is not only to produce working code.

You must preserve:

- architecture
- ownership boundaries
- security principles
- maintainability

---

# Before Writing Code

Before implementing anything:

1. Understand the requested feature.
2. Identify the bounded context.
3. Identify the owner of the business rule.
4. Check existing patterns.
5. Verify whether documentation must change.

Never start by creating files.

---

# Architecture Rules

Always follow:


Interface

↓

Application

↓

Domain

↓

Infrastructure


---

# Domain Ownership

Before adding logic ask:


Which context owns this rule?


Examples:

Asset encryption:

Owner:


Vault


Collaboration channel:

Owner:


C3


Lifecycle validation:

Owner:


TraceCore


Identity verification:

Owner:


Identity


---

# Go Rules

Generated Go code must:

- use explicit dependencies
- use context.Context where required
- return meaningful errors
- avoid global state
- avoid unnecessary abstractions

---

# New Feature Process

For every feature:

## Step 1

Identify affected contexts.

Example:

"Share asset"

Possible contexts:


Vault

Identity

C3

Federation


---

## Step 2

Define domain behavior.

Example:


Asset.Share()


not:


ShareAssetService.UpdateDatabase()


---

## Step 3

Create use case.

Example:


ShareAssetUseCase


Responsible for:

- orchestration
- permissions check
- event publication

---

## Step 4

Create events if needed.

Example:


AssetShared


---

## Step 5

Create tests.

Required:

- domain tests
- application tests
- integration tests when crossing boundaries

---

# Forbidden Actions

Never:

- bypass domain rules
- access another context database
- put business logic in handlers
- create generic utils packages
- duplicate existing concepts

---

# When Unsure

Ask:


What owns this responsibility?


Architecture correctness is more important than speed.
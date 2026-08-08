
# Ankhora Active Work

## Purpose

This document describes the current engineering focus.

AI assistants should read this file before starting development tasks.

This document represents the active development context.

Completed work should be moved to historical documentation or architecture decisions.

---

# Current Objective

## Complete C3 Collaboration Foundations

Current focus:

```

C3 bounded context stabilization

```

The goal is to establish the core collaboration primitives required by Ankhora:

- workspaces
- channels
- threads
- assets
- sharing
- trust groups
- federation foundations

---

# Current Development Area

Bounded Context:

```

C3 Collaboration

```

Current package:

```

internal/channel

```

Primary objective:

Complete Channel lifecycle and application workflows.

---

# Completed Work

## Channel Aggregate

Implemented:

- Channel domain model
- Channel repository abstraction
- Channel creation workflow
- Channel persistence flow

---

## CreateChannelUsecase

Status:

COMPLETED

Responsibilities:

- validate request
- create Channel aggregate
- persist through repository
- publish ChannelCreated event

Pattern established:

```

Request

↓

Application UseCase

↓

Domain Aggregate

↓

Repository

↓

Domain Event

```

---

## ListChannelUsecase

Status:

COMPLETED

Purpose:

Provide channel querying capability.

Implementation decisions:

- Query operation
- No domain event emission
- Repository read operation only
- No DomainBus dependency

Reason:

Queries do not represent domain state changes.

---

# Current Work

## ArchiveChannelUsecase

Status:

COMPLETED

Objective:

Introduce Channel lifecycle transition:

```

active

|

v

archived

```

---

# ArchiveChannel Design

## Ownership

Owner:

```

C3 bounded context

```

Reason:

C3 owns:

- collaboration structure
- collaboration lifecycle
- participant relationships

---

## Domain Rule

Archive is different from revoke.

```

Archive

=
soft closure

preserve history

pause collaboration

```
```

Revoke

=
remove trust

deny access

```

---

## Proposed Domain Change

Channel aggregate will own:

```

Archive()

```

The aggregate protects:

- valid state transitions
- invariant enforcement
- timestamp update

Expected transition:

Allowed:

```

Active → Archived

```

Forbidden:

```

Pending → Archived

Revoked → Archived

Archived → Archived

```

---

# Event Strategy

New domain event:

```

ChannelArchived

```

Owned by:

```

C3

```

Purpose:

Notify internal consumers that collaboration state changed.

Future possibility:

Promote to integration event for:

```

TraceCore

```

when lifecycle history synchronization is implemented.

---

# Pending Architectural Decisions

## Event Publication Reliability

Current design:

```

Update Channel

↓

Publish Event

```

Potential future improvement:

```

Aggregate Change

↓

Transactional Outbox

↓

Event Dispatcher

↓

Consumers

```

Decision:

Not required for current implementation.

Revisit when:

- federation synchronization increases
- TraceCore integration becomes active
- distributed consistency becomes critical

---

# Next Implementation Steps

## Step 1

Implement domain changes:

- add Archived status
- add ArchivedAt timestamp
- implement Channel.Archive()

---

## Step 2

Implement event:

```

ChannelArchived

```

---

## Step 3

Implement:

```

ArchiveChannelUsecase

```

Flow:

```

Request

↓

Validate

↓

Repository.GetChannel()

↓

Channel.Archive()

↓

Repository.UpdateChannel()

↓

Publish Event

```

---

## Step 4

Testing

Required:

Domain:

- active archive success
- invalid transitions
- already archived protection

Application:

- success
- repository errors
- event errors
- validation errors

---

# Current AI Instructions

When working on this area:

AI must:

1. Read:
   - ai/04-contexts/c3.md
   - ai/02-architecture/event-driven-design.md
   - ai/01-principles/ownership.md

2. Preserve:
   - C3 ownership boundaries
   - domain event patterns
   - DDD aggregate rules

3. Avoid:
   - putting lifecycle rules in use cases
   - coupling C3 directly to TraceCore
   - bypassing repositories

---

# Current Status

Phase:

```

C3 Collaboration Core

```

Current feature:

```

Channel Lifecycle Management

```

Last completed:

```

ArchiveChannelUsecase

```

Next:

```

Thread lifecycle

```


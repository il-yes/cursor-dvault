# Session: ArchiveChannelUsecase

Date:

2026-08-08


# Context

The C3 Collaboration bounded context required a lifecycle transition for Channels.

The Channel aggregate already supported creation (CreateChannelUsecase) and querying (ListChannelUsecase).

The next requirement was enabling controlled closure of collaboration spaces.


# Problem

How should Channel closure be modeled and implemented?

Questions:

- What status transition is needed?
- Who owns the transition?
- What event should represent the change?
- How does this affect TraceCore?


# Analysis

## Explored Options

### Delete

Rejected.

Destroys collaboration history. Conflicts with audit requirements.

### Revoke

Existing concept. Represents trust removal and access denial.

Not suitable. Revoke is security state, not collaboration lifecycle.

### Archive

Selected.

A collaboration space is closed operationally while preserving history.


## Ownership

C3 owns:

- collaboration structure
- collaboration lifecycle
- participant relationships

TraceCore owns:

- lifecycle history
- audit trail

Therefore:

C3 changes the state.

TraceCore records the history (future integration via events).


# Decisions

## Channel aggregate owns Archive()

The Channel aggregate protects business invariants.

Only active channels can be archived.

Reason:

Business rules belong inside the domain model, not the use case.

## Archive and Revoke remain separate

Archive:

- collaboration paused
- history preserved

Revoke:

- access removed
- trust relationship affected

Reason:

Different business semantics. Different triggers. Different consequences.

## ChannelArchived is a C3 domain event

ChannelArchived represents the fact that a collaboration space was closed.

It is owned by C3.

Future: promote to integration event when TraceCore lifecycle synchronization is implemented.

Reason:

Events should represent completed facts. Cross-context integration should use events, not direct coupling.

## State transitions are guarded

Allowed:

Active → Archived

Forbidden:

Pending → Archived
Revoked → Archived
Archived → Archived

Reason:

Invalid state transitions must be prevented at the aggregate level.

## Event publication without outbox

Current design: persist then publish.

Outbox pattern deferred until federation synchronization or TraceCore integration requires it.

Reason:

Pragmatic decision. Matches all existing use cases. Avoids premature infrastructure.


# Rejected Approaches

## Archive logic inside UseCase

Rejected because:

The use case coordinates behavior. The aggregate protects business rules.

Placing lifecycle rules in the use case creates an anemic domain model.

## Direct C3 to TraceCore dependency

Rejected because:

Creates bounded context coupling. Preferred: C3 publishes event, TraceCore consumes.

## Channel deletion instead of archive

Rejected because:

Destroys collaboration history. Conflicts with auditability.


# Architectural Impact

Affected contexts:

- C3 (primary — owns Channel lifecycle)

Affected principles:

- DDD aggregate invariants
- Event-driven design
- Dependency direction


# Implementation Summary

## Files Modified

Domain:

- aggregate.go — StatusArchived, ArchivedAt, Archive()
- events.go — EventChannelArchived, ChannelArchived struct
- errors.go — ErrChannelNotArchivable

Application:

- events/events.go — PublishChannelArchived, SubscribeToChannelArchived
- dtos.go — ArchiveChannelRequest
- usecases/archive_usecase.go — NEW

Tests:

- tests/archive_channel_test.go — NEW (13 tests)
- tests/list_channel_test.go — Updated mocks

## Test Results

21/21 PASS (13 archive + 8 list)

0 regressions


# Open Questions

- Should archived channels be restorable?
- Should Archive require additional authorization rules?
- When should transactional outbox be introduced?
- Should TraceCore consume ChannelArchived immediately?


# Participants

- Engineering Manager
- Domain Expert
- Architect
- Backend Engineer
- Reviewer (pending)
- QA Engineer (pending)


# Next Actions

- Reviewer: validate implementation against DDD, ownership, events
- QA: confirm test coverage and edge cases
- Update agent memory: mark ArchiveChannel as completed
- Next feature: Thread lifecycle or channel restore

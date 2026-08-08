
# Session: Channel Lifecycle Architecture

Date:

2026-08-04


# Context

The C3 Collaboration bounded context required additional lifecycle management capabilities.

The Channel aggregate already supported creation and querying.

The next requirement was introducing controlled closure of collaboration spaces.


# Problem

How should Channel closure be modeled?

Questions:

- Should closing a channel be deletion?
- Should it be a status transition?
- Who owns lifecycle decisions?
- How should TraceCore receive lifecycle information?


# Analysis

Several concepts were evaluated:

## Delete

Rejected.

Reason:

Deletion destroys collaboration history and conflicts with audit requirements.


## Revoke

Existing concept.

Meaning:

Trust removal and access denial.

Not suitable because it represents security state, not collaboration lifecycle.


## Archive

Selected.

Meaning:

A collaboration space is closed operationally while preserving history.


# Ownership Analysis

C3 owns:

- collaboration structure
- collaboration lifecycle
- participant relationships


TraceCore owns:

- lifecycle history
- audit trail
- compliance evidence


Therefore:

C3 changes the state.

TraceCore records the history.


# Decisions


## Channel owns archive state

Decision:

The Channel aggregate owns:

```

Archive()

```

Reason:

Business invariants belong inside the domain model.


## Archive and revoke remain separate

Archive:

- collaboration paused
- history preserved


Revoke:

- access removed
- trust relationship affected


## ChannelArchived starts as domain event

Decision:

Create:

```

ChannelArchived

```

as a C3 domain event.

Future:

Promote to integration event when TraceCore lifecycle synchronization is implemented.


## Query use cases do not publish events

Decision:

ListChannelUsecase does not depend on DomainBus.

Reason:

Queries do not represent domain state changes.


# Rejected Approaches


## Direct C3 to TraceCore dependency

Rejected.

Reason:

Creates bounded context coupling.


Preferred:

```

C3

↓

Event

↓

TraceCore consumer

```


## Archive logic inside UseCase

Rejected.

Reason:

The use case coordinates behavior.

The aggregate protects business rules.


# Open Questions

- Should archived channels be restorable?
- Should Archive require additional authorization rules?
- When should transactional outbox be introduced?
- Should TraceCore consume ChannelArchived immediately?


# Next Actions

- Implement Channel.Archive()
- Add Archived status
- Add ChannelArchived event
- Implement ArchiveChannelUsecase
- Add domain tests
- Add application tests
```

---

After this, your AI workflow becomes:

```
New task arrives

        ↓

Read:
00-vision

        ↓

Read:
01-principles

        ↓

Read:
02-architecture

        ↓

Read:
04-context

        ↓

Read:
08-agent-memory

        ↓

Check:
09-sessions

        ↓

Design

        ↓

Implement

        ↓

Create new session summary
```



# ADR-0003: Event-Driven Architecture

## Status

Accepted

## Date

2026

---

# Decision

Ankhora adopts an event-driven architecture for communication between bounded contexts.

Important business state changes are represented as events.

Contexts communicate through:

- domain events
- application events
- integration events

rather than direct internal coupling.

---

# Context

Ankhora contains multiple independent bounded contexts:

```

Identity

Vault

C3

TraceCore

Federation

Subscription

```

Each context has its own responsibilities.

Direct dependencies create problems:

```

Vault

calls

C3 internals

```

This creates:

- coupling
- fragile evolution
- unclear ownership

---

# Decision Details

Events represent facts.

An event describes something that already happened.

Example:

```

AssetCreated

ThreadOpened

CommitValidated

TrustEstablished

```

The event does not command another system.

It communicates reality.

---

# Event Types

## Domain Events

Owned by a bounded context.

Represent meaningful business changes.

Examples:

Vault:

```

AssetCreated

AssetShared

KeyRotated

```

---

C3:

```

ChannelCreated

MemberAdded

ThreadClosed

```

---

TraceCore:

```

CommitCreated

ValidationCompleted

BranchMerged

```

---

## Application Events

Represent workflow execution.

Examples:

```

SynchronizationStarted

NotificationSent

```

---

## Integration Events

Cross system boundaries.

Examples:

```

FederatedAssetShared

RemoteTrustAccepted

ExternalSynchronizationCompleted

```

---

# Event Flow Model

Example:

```

User Creates Asset

```
    |

    v
```

Vault Domain

```
    |

    v
```

AssetCreated Event

```
    |

    +----------------+

    |                |

    v                v
```

TraceCore          C3

History Record     Collaboration Update

```

---

# TraceCore Relationship

TraceCore acts as lifecycle memory.

Important events may become historical records.

Example:

```

Business Event

```
    |

    v
```

TraceCore Commit

```
    |

    v
```

Immutable History

```

---

# C3 Relationship

C3 uses events for collaboration.

Example:

```

ThreadMessageCreated

```
    |

    v
```

Realtime Gateway

```
    |

    v
```

Connected Users

```

---

# Federation Relationship

Federation uses events for trusted exchange.

Example:

```

TrustEstablished

```
    |

    v
```

Federation Message

```
    |

    v
```

Remote Node

```

---

# Event Design Rules

Events must be:

## Immutable

After creation, an event does not change.

---

## Meaningful

Events describe business facts.

Good:

```

AssetShared

```

Bad:

```

UpdateAssetFunctionExecuted

```

---

## Self-Contained

An event should contain enough information for consumers.

---

## Versioned

Events may evolve.

Compatibility must be considered.

---

# Event Ownership

Every event has an owner.

Example:

```

AssetCreated

Owner:

Vault Context

```

Consumers may react.

They do not own the event.

---

# Event Processing Rules

Consumers should:

- validate input
- handle duplicates
- be resilient to failures
- process asynchronously when appropriate

---

# Reliability Requirements

Event systems must consider:

## Delivery

Questions:

- Was the event delivered?
- Was it processed?

---

## Ordering

Questions:

- Does sequence matter?
- How is ordering preserved?

---

## Retry

Questions:

- What happens after failure?
- Is retry safe?

---

## Idempotency

Processing the same event twice should not corrupt state.

---

# Alternatives Considered

## Direct Service Calls

Rejected.

Reason:

Creates tight coupling between bounded contexts.

---

## Shared Database Communication

Rejected.

Reason:

Breaks ownership boundaries.

---

## Message Bus For Everything

Rejected.

Reason:

Not every internal operation requires asynchronous communication.

Use events where they represent meaningful facts.

---

# Consequences

Positive:

- loose coupling
- better scalability
- auditability
- easier federation
- stronger lifecycle tracking

Negative:

- more distributed complexity
- event version management required
- debugging requires better tooling

---

# AI Rules

When implementing a new feature, AI must ask:

```

Did something meaningful happen?

Should this create an event?

Who owns this event?

Who needs to react?

```

---

AI must avoid:

- hidden side effects
- direct database coupling
- events without meaning
- commands disguised as events

---

# Final Principle

Events are the memory and communication system of Ankhora.

Bounded contexts own their decisions.

Events allow the entire platform to understand what happened without destroying independence.

The system evolves through facts, not hidden dependencies.
```

---

🎯 The ADR foundation is now complete:

```text
07-decisions/

├── adr-0001-project-vision.md   ✅
├── adr-0002-ddd.md              ✅
└── adr-0003-events.md           ✅
```

Your complete AI architecture memory now looks like:

```text
ai/

00-vision
    ↓
01-principles
    ↓
02-architecture
    ↓
03-standards
    ↓
04-contexts
    ↓
05-workflows
    ↓
06-prompts
    ↓
07-decisions
```



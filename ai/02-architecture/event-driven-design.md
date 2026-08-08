# Event-Driven Design Principle

## Purpose

This document defines how events are designed, produced, consumed, and evolved inside the Ankhora platform.

Events are a communication mechanism between bounded contexts.

They allow systems to collaborate while preserving ownership boundaries.

---

# Event Philosophy

An event represents something meaningful that happened.

An event is not:

- a database change
- a technical notification
- an internal implementation detail

An event represents a business fact.

Examples:

Good:
AssetShared

TrustGroupCreated

CommitValidated

WorkflowApproved

RemoteVaultConnected


Bad:


DatabaseRowUpdated

FunctionExecuted

ObjectSaved



---

# Principle 1 — Events Represent Facts

Events are written in the past tense.

They describe something that already happened.

Examples:
VaultCreated

AssetEncrypted

ShareEntryCreated

CommitCreated


An event does not command another component.

Incorrect:


CreateAsset


This is a command.

Correct:


AssetCreated


This is a fact.

---

# Principle 2 — Event Ownership Belongs To The Domain

The bounded context that owns the business concept owns the event.

Example:

Vault owns:


AssetCreated
AssetShared
KeyRotated


TraceCore owns:


CommitCreated
WorkflowApproved
ValidationCompleted


C3 owns:


ChannelCreated
ThreadOpened
TrustGroupUpdated

---

# Principle 3 — Events Preserve Boundaries

Events allow communication without direct ownership sharing.

Example:

Incorrect:
TraceCore

imports

Vault internal package


Correct:


Vault

publishes:

AssetCreatedEvent

TraceCore

consumes:

AssetCreatedEvent


Each context keeps its own model.

---

# Principle 4 — Domain Events vs Integration Events

Ankhora distinguishes two types of events.

---

# Domain Events

Internal events inside one bounded context.

Purpose:

Coordinate behavior inside the domain.

Example:
AssetEncrypted

Used by:

Vault internally.

---

# Integration Events

Events crossing bounded contexts.

Purpose:

Notify external systems about meaningful changes.

Example:
AssetSharedIntegrationEvent

Consumed by:

- C3
- TraceCore
- Federation

---

# Principle 5 — Events Are Not Remote Procedure Calls

Do not use events to simulate commands.

Incorrect:
UserCreated

Listener:

CreateEverything()


This creates hidden coupling.

---

Better:


UserCreated

Identity event

↓

Other contexts decide independently
whether they need this information.

---

# Principle 6 — Commands and Events Are Different

Commands:

Ask something to happen.

Examples:
CreateVault

ShareAsset

ApproveWorkflow


Events:

Describe something that happened.

Examples:


VaultCreated

AssetShared

WorkflowApproved

---

# Principle 7 — Events Must Have Clear Contracts

Every integration event should define:

- event name
- owner
- version
- timestamp
- identifier
- payload
- metadata

Example:

```json
{
  "event_id": "uuid",
  "event_type": "AssetShared",
  "version": "1.0",
  "occurred_at": "timestamp",
  "source": "vault",
  "payload": {}
}
Principle 8 — Events Should Be Immutable

Once published, an event represents history.

Do not modify existing events.

If the meaning changes:

Create a new version.

Example:

AssetShared.v1

AssetShared.v2
Principle 9 — Consumers Must Be Resilient

Event consumers should handle:

duplicate delivery
delayed delivery
missing events
ordering issues

Consumers should be:

idempotent
observable
recoverable
Principle 10 — Event Flow Must Remain Understandable

Avoid uncontrolled subscriptions.

Every event consumer should have a documented reason.

Ask:

Why does this context need this event?

Ankhora Event Flow Examples
Asset Lifecycle
Vault

AssetCreated

        |

        v

TraceCore

Creates lifecycle reference
Collaboration Sharing
Vault

AssetShared

        |

        v

C3

Creates collaboration visibility
Federation
Federation

RemoteVaultConnected

        |

        v

C3

Allows trusted collaboration
Compliance Workflow
TraceCore

WorkflowApproved

        |

        v

Domain Application

Updates business state
Event Storage

Important business events may be persisted.

Reasons:

audit
recovery
replay
compliance
debugging

However:

Event storage responsibility belongs to the context that owns the events.

Events And TraceCore

TraceCore provides historical visibility.

Events can contribute to:

lifecycle reconstruction
audit trails
compliance evidence
operational history

TraceCore does not automatically own every event.

It records events relevant to lifecycle management.

Events And Realtime

WebSocket communication should not replace domain events.

Correct flow:
Business action

↓

Domain Event

↓

Application Handler

↓

Realtime Notification
Not:

WebSocket message

↓

Business change
Event Anti-Patterns
Event Everything

Not every method call needs an event.

Hidden Business Logic

Do not hide important decisions inside consumers.

Event Chains Without Ownership

Avoid:
Event A

↓

Event B

↓

Event C

↓

Unknown behavior

Shared Event Models

Contexts should not depend on another context's internal events.

AI Event Review Checklist

Before creating an event, AI should ask:

Is this a business fact?
Who owns this concept?
Is it domain or integration level?
Who needs this information?
Can a direct dependency be avoided?
Is the contract stable?
How is failure handled?
Final Principle

Events are not messages between objects.

Events are the language through which independent domains collaborate.

A healthy event-driven architecture preserves autonomy while enabling cooperation.

---

Our `02-architecture` layer is now:

```text
02-architecture/

├── overview.md
├── system-map.md
├── bounded-contexts.md
├── dependency-rules.md
└── event-driven-design.md












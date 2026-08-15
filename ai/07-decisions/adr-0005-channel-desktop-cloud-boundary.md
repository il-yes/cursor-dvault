
# ADR-0005 — Channel Desktop → Cloud Boundary

## Status

Accepted

## Date

2026-08-15

## Context

Ankhora Desktop and Ankhora Cloud implement the Channel bounded context across two applications.

The Desktop is not authoritative for Channel persistence.

The canonical architecture is:

```text
Desktop React UI
    ↓
frontend/src/services/api.ts
    ↓ Wails IPC
App
    ↓
ChannelHandler
    ↓
Channel Use Case
    ↓
TracecoreClient
    ↓ HTTP + Bearer JWT
Ankhora Cloud
    ↓
Channel bounded context
    ↓
Persistence
````

This architecture follows the existing Workspace implementation and MUST be used for Channel operations.

## Decision

Channel communication between Desktop and Cloud MUST go through a dedicated TraceCore client.

The Desktop TraceCore package MUST expose a dedicated:

```text
internal/tracecore/channel_client.go
```

The Channel client is responsible for:

* HTTP request construction
* Cloud endpoint routing
* authentication headers
* Cloud response decoding
* Cloud wire-format → domain mapping
* strict validation of Cloud responses

The Channel domain, application layer, handlers and UI MUST NOT construct HTTP requests.

## Canonical Client Pattern

Workspace established the reference pattern:

```text
internal/tracecore/workspace_client.go
```

Channel MUST follow the same architectural pattern.

Do not place Channel HTTP implementation inside:

```text
c3_repository.go
```

Do not create ad-hoc `CreateChannelDirect` / `ListChannelsDirect` helpers.

There must be one authoritative Channel HTTP boundary:

```text
internal/tracecore/channel_client.go
```

## Cloud Wire Format

An important boundary rule was discovered during Channel integration.

Cloud currently serializes Channel aggregates using Go's default JSON field names:

```json
{
  "ID": "...",
  "TemplateID": "...",
  "Title": "...",
  "Status": "pending",
  "Federation": {
    "VaultAID": "...",
    "VaultBID": "...",
    "AllowedEventTypes": ["entry.shared"],
    "AllowedPaths": null,
    "AllowedDirections": "bidirectional"
  },
  "CreatedAt": "...",
  "UpdatedAt": "...",
  "WorkspaceID": "..."
}
```

Desktop domain structures use snake_case JSON tags.

Therefore the TraceCore boundary MUST NOT blindly unmarshal Cloud Channel JSON directly into the domain aggregate.

Dedicated Cloud DTOs are required:

```text
internal/tracecore/types/channel_cloud_dto.go
```

The client maps:

```text
CloudChannelDTO
        ↓
mapCloudChannelDTO
        ↓
channel_domain.Channel
```

This mapping is intentionally owned by TraceCore.

## Response Validation

A successful HTTP status is not sufficient.

The client MUST validate:

* response envelope
* presence of data
* required Channel identity
* element validity
* expected response shape

Do NOT fabricate a successful empty response when Cloud data cannot be decoded.

Bad:

```go
return &CloudResponse[[]Channel]{
    Status: 200,
    Data: []Channel{},
    Success: true,
}, nil
```

A malformed or unexpected Cloud response MUST remain an error.

An explicitly valid empty list is different:

```json
{
  "status": 200,
  "data": []
}
```

That is a valid result.

## Persistence Ownership

Channel persistence is owned by Cloud.

Desktop MUST NOT attempt to compensate for Cloud persistence defects.

When a Cloud persistence defect is discovered, fix the Cloud bounded context first, then resume Desktop integration.

## Channel Child Identity

Channel slots have semantic IDs such as:

```text
draft
finance
signature
```

These IDs are meaningful inside a Channel but are NOT globally unique.

Therefore Cloud persistence uses:

```text
PRIMARY KEY (channel_id, id)
```

for Channel slots.

Assignments use:

```text
PRIMARY KEY (channel_id, slot_id)
```

This allows multiple Channels to use the same semantic slot IDs without cross-channel corruption.

Desktop MUST preserve semantic slot IDs while Cloud persistence scopes them by Channel.

## Create Response

Cloud Channel creation returns the persisted representation.

Cloud repository Create MUST return the committed state rather than the pre-persistence in-memory aggregate.

This prevents the API from claiming that child entities were persisted when they were not.

## Lifecycle Boundary

Channel creation, participant assignment and activation are distinct lifecycle operations.

Creation establishes the Channel.

Participant/assignment operations establish collaboration topology.

Activation validates the Channel's gated slots and transitions:

```text
pending → active
```

Do not add VaultAID/VaultBID to ChannelCreated merely to make federation tests pass.

Federation recipients are resolved from the collaboration topology.

## Federation

`channel.created` may legitimately have zero federation recipients.

An empty audience is a successful distribution no-op.

It MUST NOT be treated as a delivery error.

Actual federation failures remain errors.

## Consequences

This architecture provides:

* a single Desktop → Cloud HTTP boundary
* consistent Workspace/Channel client architecture
* protection against Cloud wire-format differences
* strict response validation
* clean DDD boundaries
* no HTTP logic in domain/application/UI
* safe Channel child persistence
* explicit separation between Channel lifecycle and federation

````

---

# 2. Update `ai/04-contexts/tracecore.md`

Add a section like this:

```md
## Desktop TraceCore Client Pattern

TraceCore is the Desktop-side integration boundary for Cloud APIs.

For every Cloud-backed bounded context, prefer a dedicated client:

```text
internal/tracecore/
├── workspace_client.go
├── channel_client.go
└── ...
````

The Workspace client is the canonical reference implementation.

Channel MUST follow the same structure.

### Channel

```text
Channel UI
    ↓
Channel application use case
    ↓
TracecoreClient
    ↓
channel_client.go
    ↓
Ankhora Cloud
```

The client owns:

* HTTP
* authentication
* Cloud DTO decoding
* wire-format mapping
* response validation

The Channel domain owns:

* Channel invariants
* lifecycle
* slots
* assignments
* domain events

The Cloud owns:

* authoritative persistence
* server-side lifecycle validation
* distribution/federation integration

````

---



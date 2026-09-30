Exactly. `ConstructionAction` is the natural complement to `ConstructionIssue`.

The important distinction is:

```text
Issue    = something requires attention
Action   = something someone must do about it
Decision = what was decided
Event    = what actually happened
```

So I would keep `Action` relatively lightweight. It should represent **work/obligation**, not become another workflow engine.

### `construction.action.v1`

```json id="n5c2qm"
{
  "template_id": "construction.action.v1",
  "record_type": "construction_action",
  "schema_version": 1,
  "fields": {
    "action_reference": "string",
    "project_id": "string",
    "type": "string",
    "resource_type": "string",
    "resource_id": "string",
    "assigned_to": "string",
    "status": "string",
    "priority": "string",
    "description": "string",
    "due_date": "date",
    "completed_at": "timestamp",
    "result": "string",
    "evidence_references": ["string"],
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### Why these fields?

| Field                       | Purpose                                                               |
| --------------------------- | --------------------------------------------------------------------- |
| `action_reference`          | Human/system reference such as `ACT-1042`                             |
| `project_id`                | Project context                                                       |
| `type`                      | Approval, rejection, inspection, review, confirmation, transfer, etc. |
| `resource_type`             | Type of object the action concerns                                    |
| `resource_id`               | Specific object concerned                                             |
| `assigned_to`               | Person or organization responsible                                    |
| `status`                    | Open, closed, etc.                                                    |
| `priority`                  | Operational importance                                                |
| `description`               | What is actually expected                                             |
| `due_date`                  | When the action should be completed                                   |
| `completed_at`              | When it actually happened                                             |
| `result`                    | Outcome of the action                                                 |
| `evidence_references`       | Supporting documents/evidence                                         |
| `created_at` / `updated_at` | Record lifecycle                                                      |

### One change from your original model

I would replace:

```json
"resource_id": "string"
```

with:

```json
"resource_type": "string",
"resource_id": "string"
```

for the same reason as with `ConstructionIssue`.

Your original comment says:

```go
ResourceID string // material / document / issue
```

But the ID alone doesn't tell us what it refers to.

So:

```json
{
  "resource_type": "construction_issue",
  "resource_id": "ISS-1042"
}
```

is considerably more expressive than:

```json
{
  "resource_id": "ISS-1042"
}
```

And it keeps the Action model generic.

---

## The interesting part: Action vs Decision

This distinction will matter a lot in TraceCore.

Imagine the delivery problem:

```text
ISS-1042
"Delivery D-1042 delayed"
```

The project manager creates:

```text
ACT-1042
type: review
resource: ISS-1042
assigned_to: engineer-23
status: open
```

The engineer investigates.

Then a decision is made:

```text
DEC-1042
decision:
"Approve alternative transport route"
```

That decision can generate several actions:

```text
ACT-1043
type: confirmation
assigned_to: supplier

ACT-1044
type: transfer
assigned_to: transporter

ACT-1045
type: approval
assigned_to: project_manager
```

And **TraceCore records the actual evolution**:

```text
construction.issue.reported
construction.action.created
construction.action.assigned
construction.action.completed
construction.decision.created
construction.decision.approved
construction.delivery.updated
construction.delivery.received
construction.inspection.completed
```

So the domain state remains clean.

```text
┌──────────────────────┐
│ Construction Domain  │
│                      │
│ Issue                │
│ Action               │
│ Decision             │
│ Delivery             │
│ Inspection           │
└──────────┬───────────┘
           │
           │ state changes
           ▼
┌──────────────────────┐
│      TraceCore       │
│                      │
│ What happened?       │
│ When?                │
│ Who?                 │
│ What changed?        │
└──────────┬───────────┘
           │
           ▼
┌──────────────────────┐
│         C3           │
│                      │
│ Who needs to know?   │
│ Who needs to act?    │
└──────────────────────┘
```

That separation is becoming very strong in this construction model.

And `Action` gives us the missing piece between **coordination** and **execution**: C3 can coordinate the actors, but the construction domain knows what operational action exists and whether it is still open.

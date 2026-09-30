Yes. For `ConstructionIssue`, I would keep the model **deliberately centered on the problem itself**.

It sits at an important point in the construction model:

```text
Construction state
       │
       ▼
     Issue
       │
       ├── investigation
       ├── coordination
       ├── decision
       └── resolution
```

The key distinction remains:

* **Issue** = something is wrong / requires attention.
* **Decision** = what the stakeholders decide to do about it.
* **TraceCore** = the history of how the issue evolved.
* **C3** = how the relevant participants coordinate.
* **Ankhora** = sensitive evidence attached to the issue.

### `construction.issue.v1`

I would evolve your original struct to this:

```json
{
  "template_id": "construction.issue.v1",
  "record_type": "construction_issue",
  "schema_version": 1,
  "fields": {
    "issue_reference": "string",
    "project_id": "string",
    "title": "string",
    "description": "string",
    "type": "string",
    "severity": "string",
    "status": "string",
    "reported_by": "string",
    "reported_at": "timestamp",
    "affected_resource_type": "string",
    "affected_resource_id": "string",
    "site_id": "string",
    "impact": "string",
    "evidence_references": ["string"],
    "assigned_to": "string",
    "resolution": "string",
    "resolved_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### Why these fields?

| Field                       | Purpose                                            |
| --------------------------- | -------------------------------------------------- |
| `issue_reference`           | Human/system reference such as `ISSUE-1042`        |
| `project_id`                | Project context                                    |
| `title`                     | Short operational description                      |
| `description`               | Detailed problem description                       |
| `type`                      | Delay, quality, safety, logistics, technical, etc. |
| `severity`                  | Importance/urgency of the issue                    |
| `status`                    | Open, investigating, resolved, closed, etc.        |
| `reported_by`               | Actor who identified/reported it                   |
| `reported_at`               | When the issue was reported                        |
| `affected_resource_type`    | What kind of resource is affected                  |
| `affected_resource_id`      | Which specific resource                            |
| `site_id`                   | Physical location when applicable                  |
| `impact`                    | Known consequence on the project                   |
| `evidence_references`       | References to supporting evidence                  |
| `assigned_to`               | Actor responsible for handling it                  |
| `resolution`                | Current/final resolution description               |
| `resolved_at`               | When resolution occurred                           |
| `created_at` / `updated_at` | Record lifecycle                                   |

### One important choice: don't put `delivery_id` here

For our current delivery-delay scenario, we **could** have:

```json
"delivery_id": "D-1042"
```

But I would resist adding it to the generic Issue template.

Instead:

```json
"affected_resource_type": "construction_delivery",
"affected_resource_id": "D-1042"
```

That keeps the issue model reusable.

The same issue mechanism can therefore represent:

```text
Issue
 ├── construction_delivery
 ├── construction_material
 ├── construction_inspection
 ├── construction_document
 └── construction_requirement
```

without turning `ConstructionIssue` into a collection of foreign keys for every possible domain object.

### And this becomes very interesting for the simulator

Our delivery scenario can now produce:

```text
Delivery D-1042
    │
    │ ETA changes
    ▼
Issue ISS-1042
    │
    │ reported by transporter
    ▼
TraceCore event
construction.issue.reported
    │
    ▼
C3 coordination
    │
    ├── Supplier
    ├── Transporter
    ├── Contractor
    └── Project Manager
    │
    ▼
Decision DEC-1042
    │
    │ alternative route approved
    ▼
TraceCore event
construction.issue.resolved
    │
    ▼
Delivery D-1042
    │
    ▼
Inspection INS-1042
```

This gives us something much more powerful than a CRUD demo.

We can later ask the simulator:

> **Why is this project currently blocked?**

and reconstruct the answer from **domain state + TraceCore history + evidence references + coordination events**.

That's exactly the kind of semantics we want to establish **before touching the frontend**.

The next model I'd take is **`ConstructionAction`**, because once an issue exists, we need to represent the concrete work that must happen to resolve it.

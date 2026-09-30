Yes — and this is an important model because **Decision is where TraceCore's historical value becomes very visible**.

Your original definition is exactly the right conceptual basis:

```text
Decision
   │
   ├── requested by
   ├── evidence
   ├── technical assessment
   ├── risks identified
   ├── participants consulted
   ├── decision
   └── consequence
```

I would turn that into a domain template without embedding the entire evidence or participant objects. They should remain references to other records.

### `construction.decision.v1`

```json id="48215"
{
  "template_id": "construction.decision.v1",
  "record_type": "construction_decision",
  "schema_version": 1,
  "fields": {
    "decision_reference": "string",
    "project_id": "string",
    "type": "string",
    "requested_by": "string",
    "request_date": "date",
    "subject": "string",
    "context": "string",
    "evidence_references": ["string"],
    "technical_assessment": "string",
    "risks_identified": ["string"],
    "participants_consulted": ["string"],
    "options_considered": ["string"],
    "decision": "string",
    "decided_by": "string",
    "decision_date": "date",
    "consequence": "string",
    "status": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### Why these fields

| Field                    | Purpose                                                                           |
| ------------------------ | --------------------------------------------------------------------------------- |
| `decision_reference`     | Stable reference such as `DEC-DELIVERY-001`                                       |
| `project_id`             | Project to which the decision belongs                                             |
| `type`                   | Delivery change, supplier selection, technical change, inspection, approval, etc. |
| `requested_by`           | Person/organization requesting the decision                                       |
| `request_date`           | When the decision was requested                                                   |
| `subject`                | What the decision concerns                                                        |
| `context`                | Why a decision is necessary                                                       |
| `evidence_references`    | References to documents, reports, events, communications, etc.                    |
| `technical_assessment`   | Technical analysis supporting the decision                                        |
| `risks_identified`       | Known risks associated with the options                                           |
| `participants_consulted` | Stakeholders whose input was considered                                           |
| `options_considered`     | Alternatives evaluated                                                            |
| `decision`               | The actual decision taken                                                         |
| `decided_by`             | Authority/person responsible for the decision                                     |
| `decision_date`          | When the decision was made                                                        |
| `consequence`            | Expected/recorded consequence of the decision                                     |
| `status`                 | Requested, under_review, approved, rejected, implemented, etc.                    |
| timestamps               | Record lifecycle                                                                  |

### One deliberate addition: `options_considered`

I added this because it is extremely useful for our demo.

Imagine the material delivery is delayed:

```text id="1f2a7c"
DEC-DELIVERY-001

Subject:
Structural material delivery delay

Context:
Supplier reports road restriction.

Options considered:
- Wait for original route
- Use alternative route
- Use alternative transporter
- Source replacement material locally

Technical assessment:
Alternative route remains compatible
with vehicle constraints.

Risks:
Additional transport cost
Possible 6h delay

Decision:
Approve alternative route.

Consequence:
Expected delivery restored to Aug 16 07:30.
```

That gives us something much richer than:

```text
status = approved
```

We're preserving **why the organization acted**.

---

## And this is where Ankhora + TraceCore separate beautifully

The decision may reference evidence:

```text id="k6s4nz"
evidence_references:
[
  "transport-report-1042",
  "site-communication-882",
  "supplier-confirmation-551"
]
```

But the actual sensitive documents can live in Ankhora.

So:

```text id="b2r9qd"
                 DECISION
                    │
          ┌─────────┼─────────┐
          │         │         │
       context    evidence   decision
                    │
                    ▼
                ANKHORA
             encrypted evidence
                    │
                    │ reference
                    ▼
                TRACECORE
              verifiable history
```

TraceCore doesn't need to possess the confidential report.

It needs to be able to say:

> **This decision was made, by whom, based on which referenced evidence, after consulting whom, and with what consequence.**

That is exactly the kind of project memory we're aiming for.

### And the simulator is now getting interesting

We have enough models to simulate:

```text id="c7m3w8"
PROJECT
   ↓
REQUIREMENT
   ↓
MATERIAL
   ↓
SUPPLIER
   ↓
OFFER
   ↓
DECISION
```

Then:

```text
DECISION
   ↓
DELIVERY
   ↓
TRANSPORT
   ↓
SITE
   ↓
INSPECTION
```

And the **event stream** connects the whole thing.


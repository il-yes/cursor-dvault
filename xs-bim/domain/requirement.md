

{
  "template_id": "construction.requirement.v1",
  "record_type": "construction_requirement",
  "schema_version": 1,
  "fields": {
    "requirement_reference": "string",
    "project_id": "string",
    "type": "string",
    "description": "string",
    "specification": "string",
    "material_id": "string",
    "quantity": "decimal",
    "unit": "string",
    "required_date": "date",
    "delivery_window": {
      "start": "date",
      "end": "date"
    },
    "site_id": "string",
    "project_phase": "string",
    "priority": "string",
    "status": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}

| Field                   | Purpose in our demo                                                |
| ----------------------- | ------------------------------------------------------------------ |
| `requirement_reference` | Stable business identifier, e.g. `REQ-STRUCT-001`                  |
| `project_id`            | Links requirement to the construction project                      |
| `type`                  | Material/service/etc.                                              |
| `description`           | Human-readable requirement                                         |
| `specification`         | Technical requirement that suppliers must satisfy                  |
| `material_id`           | Links to the material definition                                   |
| `quantity`              | Required quantity                                                  |
| `unit`                  | `ton`, `m3`, `unit`, etc.                                          |
| `required_date`         | Date by which the material is needed                               |
| `delivery_window`       | Operational delivery constraint                                    |
| `site_id`               | Where it must arrive                                               |
| `project_phase`         | Why it is needed at that point in the project                      |
| `priority`              | Important because our demo involves a critical delayed requirement |
| `status`                | Requirement lifecycle                                              |
| timestamps              | Domain record lifecycle                                            |


construction.project
        │
        │ requirement_ids
        ▼
construction.requirement
        │
        ├── material_id
        ├── specification
        ├── quantity
        ├── required_date
        ├── delivery_window
        ├── site_id
        └── project_phase
                 │
                 ▼
          SUPPLIER PROCESS


Tracecore events:
construction.requirement.created
            ↓
construction.supplier.confirmed
            ↓
construction.transport.assigned
            ↓
construction.delivery.delayed
            ↓
construction.decision.proposed
            ↓
construction.decision.approved
            ↓
construction.delivery.received
            ↓
construction.inspection.completed
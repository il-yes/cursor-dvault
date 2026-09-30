{
  "template_id": "construction.project.v1",
  "record_type": "construction_project",
  "schema_version": 1,
  "fields": {
    "project_reference": "string",
    "project_name": "string",
    "project_type": "string",
    "status": "string",

    "location": {
      "address": "string",
      "city": "string",
      "country": "string",
      "coordinates": {
        "lat": "number",
        "lng": "number"
      }
    },

    "start_date": "date",
    "planned_end_date": "date",
    "actual_end_date": "date",

    "budget": {
      "currency": "string (ISO 4217)",
      "total": "decimal",
      "allocated": "decimal"
    },

    "progress_percentage": "number(0-100)",

    "current_phase": "string",

    "milestones": ["string"],

    "stakeholder_ids": ["string"],

    "requirement_ids": ["string"],

    "site_id": "string",

    "status_summary": "string",

    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}

| Field                 | Meaning                            |
| --------------------- | ---------------------------------- |
| `project_reference`   | Stable business reference          |
| `project_name`        | Human-readable project name        |
| `project_type`        | Residential/commercial/etc.        |
| `status`              | Overall project lifecycle          |
| `location`            | Physical project location          |
| `start_date`          | Planned/actual project start       |
| `planned_end_date`    | Target completion                  |
| `actual_end_date`     | Actual completion                  |
| `budget`              | Project financial envelope         |
| `progress_percentage` | Current project progress           |
| `current_phase`       | Current construction phase         |
| `milestones`          | Important project milestones       |
| `stakeholder_ids`     | References to stakeholder records  |
| `requirement_ids`     | References to project requirements |
| `site_id`             | Reference to the construction site |
| `status_summary`      | Human-readable current state       |

Project
 │
 ├── Stakeholders
 │
 ├── Requirements
 │
 ├── Site
 │
 ├── Milestones
 │
 ├── Materials
 │
 ├── Deliveries
 │
 ├── Decisions
 │
 └── Evidence


construction.project.v1
        │
        ├── requirement_ids ──► construction.requirement.v1
        │
        ├── stakeholder_ids ─► construction.stakeholder.v1
        │
        └── site_id ──────────► construction.site.v1

construction.requirement.v1
        │
        ├── material_id ──────► construction.material.v1
        └── supplier_id ──────► construction.supplier.v1

construction.delivery.v1
        │
        ├── material_id
        ├── supplier_id
        ├── transport_id ─────► construction.transport.v1
        └── site_id
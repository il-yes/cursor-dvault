
Exactly. The next one should be **`construction.material.v1`**.

This is the object that the requirement refers to through `material_id`. For our demo, it needs enough information to answer:

> **What material/component is required, what technical characteristics must it satisfy, and where did it come from?**

We should **not** put supplier availability or delivery information here. Those belong to `supplier` and `delivery`.

### `construction.material.v1`

```json id="58321"
{
  "template_id": "construction.material.v1",
  "record_type": "construction_material",
  "schema_version": 1,
  "fields": {
    "material_reference": "string",
    "name": "string",
    "category": "string",
    "description": "string",
    "standard": "string",
    "specification": "string",
    "unit": "string",
    "origin": "string",
    "production_date": "date",
    "batch_reference": "string",
    "certification_references": ["string"],
    "status": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### Why these fields

| Field                      | Role in the demo                               |
| -------------------------- | ---------------------------------------------- |
| `material_reference`       | Stable reference such as `MAT-STRUCT-001`      |
| `name`                     | Human-readable material/component name         |
| `category`                 | Structural steel, concrete, timber, etc.       |
| `description`              | General description                            |
| `standard`                 | Applicable construction/technical standard     |
| `specification`            | Technical characteristics                      |
| `unit`                     | `ton`, `m3`, `unit`, etc.                      |
| `origin`                   | Material origin                                |
| `production_date`          | Useful for provenance/quality                  |
| `batch_reference`          | Links a physical batch to the digital record   |
| `certification_references` | Certifications/test evidence references        |
| `status`                   | Available, reserved, delivered, rejected, etc. |
| timestamps                 | Record lifecycle                               |

The important distinction is:

```text id="b4oscm"
REQUIREMENT
"What do we need?"
       │
       │ material_id
       ▼
MATERIAL
"What exactly is it?"
       │
       ├── specification
       ├── standard
       ├── batch
       ├── origin
       └── certifications
```

Then later:

```text id="0jgjoc"
MATERIAL
    │
    ▼
SUPPLIER
"Who can provide it?"
    │
    ▼
DELIVERY
"When/how is it getting here?"
    │
    ▼
TRANSPORT
"How does it get there?"
    │
    ▼
SITE
"Can we receive it?"
```

And that gives us a very clean separation between **what is required**, **what the material is**, and **what happens to it**.

For the demo, this is particularly important because when the delivery becomes late, TraceCore can reconstruct:

> Requirement `REQ-STRUCT-001` required material `MAT-STRUCT-001`, which had specification X, was confirmed by supplier Y, was transported under delivery D-1042, encountered a road restriction, triggered an alternative-route decision, and was eventually received and inspected.

That's already starting to look like the **project intelligence graph** we want rather than a conventional CRUD construction app.

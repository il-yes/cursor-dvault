Yes. And here I would make one important modeling distinction: **Site is the relatively stable physical/operational context**, while `Inspection` and `Acceptance` are activities/results that happen at the site. We should reference them rather than embed them.

The Site should answer:

> **Where does construction activity happen, what are the site's operational constraints, and what conditions must be satisfied to receive work/materials?**

### `construction.site.v1`

```json id="86134"
{
  "template_id": "construction.site.v1",
  "record_type": "construction_site",
  "schema_version": 1,
  "fields": {
    "site_reference": "string",
    "project_id": "string",
    "name": "string",
    "location": {
      "address": "string",
      "city": "string",
      "country": "string",
      "coordinates": {
        "lat": "number",
        "lng": "number"
      }
    },
    "access_window": {
      "start": "string",
      "end": "string"
    },
    "access_constraints": ["string"],
    "storage_capacity": {
      "value": "decimal",
      "unit": "string"
    },
    "storage_constraints": ["string"],
    "weather_constraints": ["string"],
    "receiving_requirements": ["string"],
    "inspection_required": "boolean",
    "acceptance_required": "boolean",
    "status": "string",
    "notes": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### Why these fields

| Field                    | Purpose                                                       |
| ------------------------ | ------------------------------------------------------------- |
| `site_reference`         | Stable site identifier, e.g. `SITE-001`                       |
| `project_id`             | Project to which the site belongs                             |
| `name`                   | Human-readable site name                                      |
| `location`               | Physical location                                             |
| `access_window`          | When deliveries can enter the site                            |
| `access_constraints`     | Gate restrictions, vehicle restrictions, access permits, etc. |
| `storage_capacity`       | Available capacity for incoming materials                     |
| `storage_constraints`    | Conditions such as covered storage or weight limitations      |
| `weather_constraints`    | Conditions that can affect access or construction             |
| `receiving_requirements` | Requirements for accepting a delivery                         |
| `inspection_required`    | Whether incoming material must be inspected                   |
| `acceptance_required`    | Whether formal site acceptance is required                    |
| `status`                 | Active, inactive, closed, etc.                                |
| `notes`                  | Additional operational information                            |
| timestamps               | Record lifecycle                                              |

### Why `access_window` matters

Our delivery scenario can now become much more realistic.

Suppose:

```text id="f3k8p2"
Site access:
07:00 → 17:00
```

The transport arrives at:

```text id="v7m2q9"
ETA: 18:30
```

The problem is no longer simply:

> "Truck is late."

It becomes:

> **The transport ETA falls outside the site's receiving window.**

That can trigger another coordination event.

```text id="n4c6r1"
Transport ETA changed
        ↓
ETA conflicts with site access window
        ↓
Delivery impact identified
        ↓
Decision required
```

That's exactly the kind of **cross-domain reasoning** our simulator should eventually exercise.

### Storage capacity is equally interesting

Imagine:

```text id="a9r5w3"
Site storage capacity:
100 tonnes

Incoming delivery:
80 tonnes

Current occupied:
60 tonnes
```

The delivery cannot simply be accepted because:

```text
60 + 80 > 100
```

That could trigger:

```text id="p2d7k4"
Storage constraint identified
        ↓
Alternative delivery schedule
        ↓
Decision
        ↓
Updated delivery plan
```

Again, **the Site doesn't make the decision**. It provides the constraint.

TraceCore records what happened.

C3 allows the stakeholders to coordinate.

Ankhora can hold the sensitive evidence.

That's the architecture working exactly as intended.

---

## One deliberate choice

I did **not** put this directly into Site:

```text
inspection_id
acceptance_id
```

because inspection and acceptance are **events/activities around the site**, not intrinsic properties of the site.

The site says:

```text id="z6m1q8"
inspection_required = true
acceptance_required = true
```

Then later:

```text id="s8v3d5"
Delivery
   ↓
Inspection
   ↓
Acceptance
```


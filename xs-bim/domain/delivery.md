Yes. Your original `Delivery` model is the right starting point, and now that we have the `Requirement`, `Supplier`, `Material`, `Offer`, and `Decision` models, we can make the delivery precise without making it overloaded.

The delivery should answer exactly what you stated:

> **What is being delivered, for which project/requirement, by whom, when, where, and in what state?**

### `construction.delivery.v1`

```json id="62418"
{
  "template_id": "construction.delivery.v1",
  "record_type": "construction_delivery",
  "schema_version": 1,
  "fields": {
    "delivery_reference": "string",
    "project_id": "string",
    "requirement_id": "string",
    "material_id": "string",
    "supplier_id": "string",
    "offer_id": "string",
    "quantity": "decimal",
    "unit": "string",
    "status": "string",
    "planned_delivery_date": "date",
    "eta": "timestamp",
    "actual_delivery_date": "date",
    "site_id": "string",
    "transport_id": "string",
    "driver_id": "string",
    "reference": "string",
    "delivery_notes": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### Why these fields

| Field                   | Purpose                                                                    |
| ----------------------- | -------------------------------------------------------------------------- |
| `delivery_reference`    | Stable delivery identifier, e.g. `DEL-1042`                                |
| `project_id`            | Construction project concerned                                             |
| `requirement_id`        | The requirement being fulfilled                                            |
| `material_id`           | Material being delivered                                                   |
| `supplier_id`           | Organization responsible for supplying it                                  |
| `offer_id`              | Accepted offer that resulted in this delivery                              |
| `quantity`              | Quantity being delivered                                                   |
| `unit`                  | Unit of measurement                                                        |
| `status`                | Planned, scheduled, in_transit, delayed, arrived, received, rejected, etc. |
| `planned_delivery_date` | Original committed delivery date                                           |
| `eta`                   | Current estimated arrival                                                  |
| `actual_delivery_date`  | Actual arrival, once completed                                             |
| `site_id`               | Destination construction site                                              |
| `transport_id`          | Transport operation associated with the delivery                           |
| `driver_id`             | Driver responsible, if applicable                                          |
| `reference`             | External logistics/reference number                                        |
| `delivery_notes`        | Operational information                                                    |
| timestamps              | Record lifecycle                                                           |

### One important distinction: planned date vs ETA

This is **very important for our demo**.

We need both:

```text id="x2c9n4"
planned_delivery_date = 2026-08-15
eta                   = 2026-08-15T10:00
```

Then the delay happens:

```text id="m8v4q1"
planned_delivery_date = 2026-08-15
eta                   = 2026-08-16T07:30
status                = delayed
```

The original commitment isn't overwritten.

That's exactly the sort of thing TraceCore can make historically visible.

### The delivery now connects the graph

```text id="r5k7s2"
PROJECT
   │
   ▼
REQUIREMENT
   │
   ├──────────────► MATERIAL
   │
   ▼
OFFER
   │
   ▼
DECISION
   │
   ▼
DELIVERY
   │
   ├── supplier
   ├── material
   ├── transport
   └── site
```

And the **transport itself should remain a separate model**.

We don't put:

```text
vehicle
route
road restrictions
ETA changes
driver status
```

inside Delivery.

Instead:

```text id="q6p3a1"
DELIVERY
   │
   └── transport_id
          │
          ▼
      TRANSPORT
```

This is going to be particularly useful for our actual demo.

The delivery can say:

> **DEL-1042 is delayed.**

The transport model will eventually be able to explain:

> **The vehicle encountered a road restriction on the planned route, causing the ETA to change.**

And TraceCore connects those facts chronologically:

```text id="z7f4m2"
08/15 08:03  transport accepted
08/15 11:41  road restriction reported
08/15 12:05  ETA changed
08/15 12:17  site acknowledged
08/15 13:02  alternative route proposed
08/15 14:20  decision approved
08/16 07:51  delivery received
08/16 09:20  inspection completed
```

That is exactly why I think the delivery model should remain relatively simple.

**Delivery = state of the physical delivery.**

**Transport = movement/logistics.**

**TraceCore = history of what happened.**

**Decision = human/organizational judgment.**

**Ankhora = sensitive evidence.**

That separation is becoming very clean.

### Current model family

We now have:

```text
construction.project.v1
        │
        ▼
construction.requirement.v1
        │
        ├────────► construction.material.v1
        │
        └────────► construction.supplier.v1
                         │
                         ▼
                  construction.offer.v1
                         │
                         ▼
                  construction.decision.v1
                         │
                         ▼
                  construction.delivery.v1
                         │
                         ▼
                  construction.transport.v1
                         │
                         ▼
                    construction.site.v1
```

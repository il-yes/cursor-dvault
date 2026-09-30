Exactly. **`construction.transport.v1`** should describe the movement operation, not the delivery itself.

The key question is:

> **How is the material being transported, by whom, along which route, under what constraints, and what is the current transport state?**

For our delayed-delivery scenario, this is the model that will allow us to represent the road restriction that causes the ETA change.

### `construction.transport.v1`

```json id="52763"
{
  "template_id": "construction.transport.v1",
  "record_type": "construction_transport",
  "schema_version": 1,
  "fields": {
    "transport_reference": "string",
    "delivery_id": "string",
    "vehicle_reference": "string",
    "vehicle_type": "string",
    "driver_id": "string",
    "origin": {
      "address": "string",
      "city": "string",
      "country": "string"
    },
    "destination_site_id": "string",
    "route_reference": "string",
    "planned_departure": "timestamp",
    "actual_departure": "timestamp",
    "planned_arrival": "timestamp",
    "eta": "timestamp",
    "actual_arrival": "timestamp",
    "constraints": ["string"],
    "status": "string",
    "delay_reason": "string",
    "notes": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### Why these fields

| Field                 | Purpose                                                            |
| --------------------- | ------------------------------------------------------------------ |
| `transport_reference` | Stable logistics identifier, e.g. `TR-1042`                        |
| `delivery_id`         | Delivery being transported                                         |
| `vehicle_reference`   | Vehicle identification/reference                                   |
| `vehicle_type`        | Truck, heavy truck, specialized carrier, etc.                      |
| `driver_id`           | Responsible driver                                                 |
| `origin`              | Where the transport starts                                         |
| `destination_site_id` | Destination construction site                                      |
| `route_reference`     | Identifies the planned route                                       |
| `planned_departure`   | Expected departure                                                 |
| `actual_departure`    | Actual departure                                                   |
| `planned_arrival`     | Original expected arrival                                          |
| `eta`                 | Current estimated arrival                                          |
| `actual_arrival`      | Actual arrival                                                     |
| `constraints`         | Known operational constraints                                      |
| `status`              | Planned, assigned, in_transit, delayed, completed, cancelled, etc. |
| `delay_reason`        | Reason for a delay when applicable                                 |
| `notes`               | Additional operational information                                 |
| timestamps            | Record lifecycle                                                   |

### The important part for our demo

We deliberately have:

```text id="8c4k2m"
planned_arrival
        ≠
eta
```

because the whole scenario depends on the difference.

Initially:

```text id="k5d1p8"
planned_arrival = Aug 15 10:00
eta             = Aug 15 10:00
status          = in_transit
```

Then:

```text id="q9w3r6"
constraints = [
  "road restriction"
]

status = delayed

delay_reason = "Road restriction on planned route"

eta = Aug 16 07:30
```

But again, **we don't need to mutate history into the transport record**.

TraceCore records:

```text id="m2v7x4"
transport.accepted
        ↓
transport.departed
        ↓
transport.constraint.reported
        ↓
transport.eta.updated
        ↓
delivery.delay.reported
```

So the current transport state might say:

```text id="u4s8n1"
status = delayed
eta = Aug 16 07:30
```

while TraceCore can still reconstruct:

```text id="r3c6p9"
Aug 15 08:03  accepted
Aug 15 08:20  departed
Aug 15 11:41  road restriction reported
Aug 15 12:05  ETA changed
```

That's an important architectural principle we're establishing:

> **Current domain state tells us where things are. TraceCore tells us how they got there.**

---

### One other useful separation

Notice that `destination_site_id` points to the **Site**, rather than embedding the site.

Likewise:

```text id="f7n2k5"
Delivery
    └── transport_id
             ↓
         Transport
             └── destination_site_id
                         ↓
                       Site
```

This will let the site independently contain things such as:

* access windows
* storage capacity
* receiving constraints
* inspection requirements
* site conditions

which can themselves become causes of coordination problems.


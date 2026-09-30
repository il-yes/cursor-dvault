Yes. This is a good next model because **documents are the bridge between the operational construction model and Ankhora's evidence layer**.

The document should answer:

> **What document is this, what does it relate to, where is its content stored, what is its status, and who provided it?**

### `construction.document.v1`

```json id="61583"
{
  "template_id": "construction.document.v1",
  "record_type": "construction_document",
  "schema_version": 1,
  "fields": {
    "document_reference": "string",
    "project_id": "string",
    "name": "string",
    "category": "string",
    "description": "string",
    "cid": "string",
    "content_hash": "string",
    "version": "string",
    "status": "string",
    "uploaded_by": "string",
    "issued_by": "string",
    "issued_at": "date",
    "effective_from": "date",
    "expires_at": "date",
    "related_resource_type": "string",
    "related_resource_id": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### Why these fields

| Field                   | Purpose                                                                        |
| ----------------------- | ------------------------------------------------------------------------------ |
| `document_reference`    | Stable business identifier, e.g. `DOC-CERT-882`                                |
| `project_id`            | Project to which the document belongs                                          |
| `name`                  | Human-readable document name                                                   |
| `category`              | Drawing, specification, permit, report, blueprint, contract, certificate, etc. |
| `description`           | What the document represents                                                   |
| `cid`                   | IPFS content identifier                                                        |
| `content_hash`          | Integrity verification                                                         |
| `version`               | Document revision/version                                                      |
| `status`                | Draft, active, superseded, archived, etc.                                      |
| `uploaded_by`           | Identity that registered the document                                          |
| `issued_by`             | Organization/person that issued it                                             |
| `issued_at`             | Date of issuance                                                               |
| `effective_from`        | Date from which the document applies                                           |
| `expires_at`            | Optional expiration                                                            |
| `related_resource_type` | Material, delivery, inspection, decision, requirement, etc.                    |
| `related_resource_id`   | Identifier of the related domain object                                        |
| timestamps              | Record lifecycle                                                               |

### The `CID` is particularly important

This is where our existing architecture becomes visible again:

```text id="n4w8p2"
ConstructionDocument
       │
       ├── CID
       └── content_hash
              │
              ▼
             IPFS
              │
              ▼
        actual document
```

The construction domain doesn't need to store the document bytes.

And for sensitive documents:

```text id="r6m3k9"
ConstructionDocument
        │
        └── CID / reference
               │
               ▼
            Ankhora
               │
        encrypted evidence
```

So we maintain the same principle:

> **The domain knows about evidence; it doesn't have to own the sensitive evidence itself.**

### `related_resource_*` is also important

Instead of creating fields like:

```text id="z1q5v7"
delivery_document_id
inspection_report_id
material_certificate_id
decision_document_id
```

we have a generic relationship:

```text id="c8k2m5"
related_resource_type = "inspection"
related_resource_id   = "INSP-1042"
```

or:

```text id="h3v9s1"
related_resource_type = "material"
related_resource_id   = "MAT-STRUCT-001"
```

That allows the same document model to support the entire construction domain.

For example:

```text id="p7d4x2"
Material
   │
   └── certificate
          ↓
ConstructionDocument
          │
          └── CID → IPFS / Ankhora


Delivery
   │
   └── transport report
             ↓
ConstructionDocument


Inspection
   │
   └── inspection report
             ↓
ConstructionDocument


Decision
   │
   └── supporting evidence
             ↓
ConstructionDocument
```

### This gives our simulator an evidence chain

Our delayed delivery scenario can now produce:

```text id="v5r8m1"
DELIVERY DEL-1042
      │
      ├── supplier confirmation
      ├── transport report
      ├── road restriction evidence
      ├── alternative route decision
      ├── site communication
      ├── inspection report
      └── material certificate
```

And TraceCore can preserve the relationships between those records without becoming the document repository itself.

---

At this point, our **core construction scenario model is becoming quite complete**:

```text id="q6s2n8"
PROJECT
   │
   ├── REQUIREMENT
   │       │
   │       └── MATERIAL
   │
   ├── SUPPLIER
   │       │
   │       └── OFFER
   │
   ├── DECISION
   │
   ├── DELIVERY
   │       │
   │       └── TRANSPORT
   │
   ├── SITE
   │       │
   │       └── INSPECTION
   │
   └── DOCUMENTS
```


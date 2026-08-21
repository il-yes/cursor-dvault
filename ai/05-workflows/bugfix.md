# Bugfix Workflow — Runtime Verification & Boundary Tracing

## Purpose

This document defines the standard bugfixing workflow for Ankhora.

The goal of this workflow is to prevent speculative debugging, unnecessary architecture rewrites, and unverified assumptions. AI agents and human engineers MUST follow this procedure when investigating and resolving bugs across multi-layer or distributed components.

---

# Core Mindset

> **Do not guess. Trace. Measure. Compare expected vs actual. Find the first broken boundary. Prove the diagnosis. Make the smallest fix.**

---

# Boundary-Tracing Execution Model

When investigating any issue involving data flow across boundaries (e.g. Frontend → Desktop App → TraceCore → Cloud → DB → Frontend UI), trace the actual data step by step across every layer:

```text
1. Frontend UI Action
     ↓
2. Wails / App Bridge
     ↓
3. Application Handler
     ↓
4. Use Case
     ↓
5. Repository / Client
     ↓
6. HTTP Request
     ↓
7. Cloud HTTP Handler
     ↓
8. Cloud Use Case
     ↓
9. Database Query & Persistence
     ↓
10. Cloud Response Serialization
     ↓
11. Desktop Decoder DTO
     ↓
12. App Bridge Return
     ↓
13. Frontend Store State Update
     ↓
14. UI Component Re-render
```

At every single boundary, compare:

```text
EXPECTED VALUE  vs  ACTUAL VALUE
```

The **first boundary** where `EXPECTED != ACTUAL` is the **primary debugging target**.

---

# Step-by-Step Investigation Workflow

## Step 1 — Verify Environment & Configuration
Before concluding that an endpoint URL, host, port, authentication token, or environment parameter is broken:
- Inspect the actual configuration source (`.env`, config struct, flag).
- Print or log the final, constructed runtime value.
- Do NOT infer environment configuration solely from static code paths.

## Step 2 — Trace Identifiers & Payloads Across Boundaries
Verify that key identifiers (`workspace_id`, `channel_id`, `thread_id`, `identity_id`, `asset_id`, `vault_id`) remain intact and correctly formatted across every boundary:
- Check naming conventions (e.g., Go default `ChannelID` vs JSON tag `channel_id`).
- Ensure identifiers are not lost, truncated, or defaulted to empty strings during serialization/deserialization.

## Step 3 — Analyze Empty Results Systematically
When a collection unexpectedly returns empty (`[]` or `null`), execute the 10-step inspection check:
1. Was the correct identifier supplied by the caller?
2. Was the correct HTTP URL used?
3. Did the request reach the expected server/handler?
4. Did the server handler extract the parameter correctly?
5. Were pagination parameters valid (`limit > 0`, `offset >= 0`)?
6. What actual SQL query was executed by GORM / DB driver?
7. How many database rows matched the query?
8. How was the result serialized into HTTP response JSON?
9. How was the HTTP response decoded into client DTOs?
10. How did frontend store update component state?

## Step 4 — Check Beyond HTTP Status Codes
Do not assume `HTTP 200` means success. Verify payload structure, item counts, error fields, and domain semantics before concluding an operation succeeded.

## Step 5 — Formulate Hypotheses with Evidence
Categorize every claim:
- **Observed**: Data directly captured from logs, HTTP payloads, or DB queries.
- **Proven**: Verified via automated test or log trace.
- **Inferred**: Logically derived from code, requiring runtime confirmation.
- **Hypothesis**: Unverified initial guess.

Never present a hypothesis as a proven finding without evidence.

## Step 6 — Implement the Minimal Fix
- Fix only the smallest responsible component at the first broken boundary.
- Do NOT rewrite surrounding architecture or change unrelated contracts.
- Do NOT add redundant synchronization mechanisms if existing mechanisms can fulfill the requirement.
- Preserve established invariants.

## Step 7 — Verify with Evidence Hierarchy
Validate the fix using the highest available evidence tier:
1. Automated unit / integration test
2. Runtime boundary log
3. HTTP request/response capture
4. Database state observation
5. Direct runtime reproduction
6. Static code inspection

---

# Verification Checklist Before Closing a Bug

- [ ] Configuration and environment verified against runtime values.
- [ ] First broken boundary identified and proven with logs/tests.
- [ ] Identifier transformation verified across all DTOs and serializers.
- [ ] Empty collection results fully explained (pagination, query parameters, DB rows).
- [ ] Smallest root-cause fix applied without architecture scope creep.
- [ ] Automated regression test added or executed.

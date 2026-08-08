# Session: <Title>

Status: Planning

Feature: <Feature>

Owner: Engineering Manager

Started: YYYY-MM-DD

Affected Contexts:
- <Context>

Participants:
- Engineering Manager

Related ADRs:
- <ADR>

Related Memory:
- active-work.md

Outcome:
- Pending

---

Date:

YYYY-MM-DD

---

# Context

Describe the situation that triggered this session.

Example:

A new capability was required in the Channel bounded context.

---


# Problem

What question or challenge needed resolution?

---


# Analysis

Describe:

- explored options
- technical considerations
- architectural impact

---


# Decisions

## Decision 1

---

Description:

Reason:


---

## Decision 2

Description:

Reason:


---

# Rejected Approaches

## <Approach>

Rejected because:


# Architectural Impact

Affected contexts:

- 

Affected principles:

- 


# Open Questions

- 

# Next Actions

- 


# AI Consumption

This section records how AI resources were consumed during the engineering session.

It is used to measure the cost and effectiveness of the engineering workflow over time.

## Thinking Mode

| Role                | Status      | Notes |
| ------------------- | ----------- | ----- |
| Engineering Manager | not started |       |
| Domain Expert       | not started |       |
| Architect           | not started |       |

## Execution Mode

| Role             | Status      | Notes |
| ---------------- | ----------- | ----- |
| Backend Engineer | not started |       |
| Reviewer         | not started |       |
| QA Engineer      | not started |       |

## Quota

| Metric                       | Value       |
| ---------------------------- | ----------- |
| Thinking quota boundary      | not reached |
| Execution quota boundary     | not reached |
| Session interrupted by quota | no          |

### Quota Notes

Record any useful observations about quota consumption, interruptions, or unexpected behavior.

---

## Implementation Consumption

| Metric                         | Value |
| ------------------------------ | ----- |
| Existing implementation reused | no    |
| New implementation required    | no    |
| Tests required                 | no    |
| Documentation updates required | no    |

### Reuse

Describe whether previous implementation, patterns, decisions, session history, or agent memory allowed the AI to avoid rediscovering or rewriting work.

---

## Session Efficiency

| Metric                                           | Value |
| ------------------------------------------------ | ----- |
| Planned thinking stages                          | 3     |
| Executed thinking stages                         | 0     |
| Planned execution stages                         | 3     |
| Executed execution stages                        | 0     |
| Stages skipped                                   | 0     |
| Stages stopped because work was already complete | 0     |

### Efficiency Notes

Record observations such as:

* an agent discovered that previous work was already complete
* a role was unnecessary
* a role repeated work already performed
* session history prevented rediscovery
* an architectural decision eliminated implementation work
* quota prevented a later stage

---

# Future Measurement

Do not optimize the workflow from a single session.

After approximately **10–20 sessions**, aggregate session data to identify actual project-specific patterns.

Initial categories:

| Task Type             | Avg Thinking Stages | Avg Execution Stages |
| --------------------- | ------------------: | -------------------: |
| CRUD                  |                 1–2 |                  1–2 |
| Domain behavior       |                 2–3 |                  2–3 |
| Cross-context feature |                 3–4 |                  3–4 |
| Architecture change   |                 3–5 |                  3–5 |
| Security change       |                 3–5 |                  3–5 |

These are **initial hypotheses, not established measurements**.

The values should eventually be replaced with observed project data.

The purpose of this measurement is to determine:

* which roles are actually necessary for each task type
* where quota is consumed
* where AI work is duplicated
* where previous sessions eliminate work
* which workflows can safely become conditional
* which roles provide the highest engineering value

The workflow should not be optimized until sufficient session data exists.


# Example: Your Channel List Use Case

Let's replay what happened, but using your new system.

## Step 1 — Engineering Manager

The Engineering Manager opens the session.

File:

```
09-sessions/

session-2026-08-06-channel-list-usecase.md
```

```md

# Session

Status: Planning

Feature: ListChannelUsecase

Owner: Engineering Manager

Started: 2026-08-06

Affected Contexts:
- C3

Participants:
- Engineering Manager
- Domain Expert
- Architect
- Backend Engineer
- Reviewer

Related ADRs:
- adr-0002-ddd.md
- adr-0003-events.md

Related Memory:
- active-work.md

Outcome:
- Pending
Topic

ListChannelUsecase

---

Objective

Complete the Channel CRUD implementation.

---

Reason

The Create use case exists.

The remaining CRUD operations should follow existing patterns.

---

Participants

- Engineering Manager
- Domain Expert
- Architect
- Backend Engineer
- Reviewer

---

Expected Deliverables

- approved design
- implementation
- tests
- review
```

Nothing technical yet.

This is simply defining the work.

---

# Step 2 — Domain Expert

Now the Domain Expert answers only one question:

> What does the business want?

Output:

```
Business capability:

List all channels inside a workspace.

Business rules:

• workspace must exist

• only channels belonging to workspace

• read-only

Ownership:

C3
```

Notice:

No Go.

No repository.

No events.

Only business.

---

# Step 3 — Architect

The Architect now asks:

```
Where does this belong?

Who owns it?

Does it already exist?

Can we reuse patterns?

Does it introduce coupling?
```

Output:

```
Decision:

New ListChannelUsecase

Repository already supports it.

No new aggregate.

No domain events.

No ADR required.
```

Still no code.

---

# Step 4 — Backend Engineer

Only now:

```
Implement.
```

The engineer already knows:

* DDD rules
* ownership
* architecture
* standards

They simply write the code.

---

# Step 5 — Reviewer

Now switch roles.

Ask:

```
Review against:

DDD

Ownership

Events

Security

Performance
```

The Reviewer never rewrites the code.

They only evaluate it.

Exactly what happened with ArchiveChannel.

---

# Step 6 — QA

QA asks:

```
What could fail?

What tests are missing?

What edge cases exist?
```

Result:

```
Need tests:

✓ empty list

✓ invalid workspace

✓ nil request

✓ repo failure
```

---

# Step 7 — Documentation

Documentation engineer asks:

```
Does documentation change?
```

Answer:

```
No

No public behavior changed.
```

or

```
Update:

04-contexts/c3.md
```

---

# Step 8 — Engineering Manager

Engineering Manager closes the session.

```
Decision:

Approved

Implementation merged.

Next task:

ArchiveChannel.
```

---

# Step 9 — Agent Memory

Update:

```
08-agent-memory
```

Example:

```
Completed

ListChannelUsecase
```

```
Current Focus

ArchiveChannel
```

---

# The AI conversation becomes

Instead of this:

```
Implement ListChannel.
```

It becomes:

```
Engineering Manager:

Open session.

↓

Domain Expert:

Explain the business capability.

↓

Architect:

Validate architecture.

↓

Backend Engineer:

Implement.

↓

Reviewer:

Review.

↓

QA:

Validate.

↓

Engineering Manager:

Close session.

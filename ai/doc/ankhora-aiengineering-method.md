Ankhora AI Engineering Methodology
Version 1.0


### Testing Workflow ###
I wouldn't think of the session file as something you "run." Think of it as the **agenda and context** for an engineering meeting.

The AI terminal becomes the meeting room.

## Workflow

### 1. Create the session

Create:

```text
ai/09-sessions/session-2026-08-06-channel-list-usecase.md
```

It contains the objective, participants, context, decisions, etc.

---



### 2. Open your AI terminal

Instead of a generic prompt like:

> Implement ListChannelUsecase

you give the AI a role and the session.

For example:

```text
You are the Engineering Manager.

Read:

- ai/README.md
- ai/09-sessions/session-2026-08-06-channel-list-usecase.md

Follow AGENTS.md.

Lead this engineering session.

Do not implement code yet.

Your objective is to produce an implementation plan and identify which engineering roles should act next.
```

The Engineering Manager responds with a plan.

---

### 3. Continue with the next role

Then:

```text
You are the Domain Expert.

Continue the engineering session:

ai/09-sessions/session-2026-08-06-channel-list-usecase.md

Read only the documentation required by AGENTS.md.

Analyze the business capability.

Do not discuss implementation.
```

---

Then:

```text
You are the Architect.

Continue the same session.

Validate the architecture.

Identify ownership, dependencies and risks.

Do not write code.
```

---

Then:

```text
You are the Backend Engineer.

Continue the approved session.

Implement the approved design.

Follow all Ankhora standards.

Generate tests.

Update documentation if necessary.
```

---

Then:

```text
You are the Reviewer.

Review the implementation against:

- DDD
- Ownership
- Security
- Event-driven architecture

Produce findings only.
```

---

Finally:

```text
You are the Engineering Manager.

Close the session.

Summarize:

- decisions
- completed work
- follow-up tasks

Suggest updates for:

- 08-agent-memory
- 07-decisions (if needed)
```

---

## Why this is powerful

Your session file becomes the **single source of truth** for that piece of work.

Instead of every AI call starting from scratch, every role works from the same context:

```text
Session
      │
      ├── Engineering Manager
      ├── Domain Expert
      ├── Architect
      ├── Backend Engineer
      ├── Reviewer
      ├── QA
      └── Release Manager
```

Everyone is discussing the same problem from their own responsibility.

---

## One improvement I'd make

I would add a small metadata block at the top of every session:

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
```

Then, as the work progresses, the Engineering Manager updates:

```md
Status: In Progress
```

and eventually:

```md
Status: Completed
```

This turns each session into a lightweight engineering record that anyone—human or AI—can revisit later to understand not just **what** was built, but **why** it was built that way.













### Conclusion ###
Yes. And I think **freezing now is exactly the right move**.

You've done the important experiment: you didn't just design the workflow theoretically—you ran it against a real Channel feature, with a real strong model, and observed the quota/cost behavior.

The result is quite clear.

## What the experiment proved

Your current pipeline works:

```text
ENGINEERING SESSION
        │
        ▼
┌─────────────────────┐
│ Engineering Manager │
└──────────┬──────────┘
           ▼
┌─────────────────────┐
│    Domain Expert    │
└──────────┬──────────┘
           ▼
┌─────────────────────┐
│      Architect      │
└──────────┬──────────┘
           │
           │ APPROVED DESIGN
           ▼
┌─────────────────────┐
│  Backend Engineer   │
└──────────┬──────────┘
           ▼
┌─────────────────────┐
│      Reviewer       │
└──────────┬──────────┘
           ▼
┌─────────────────────┐
│         QA          │
└──────────┬──────────┘
           ▼
┌─────────────────────┐
│ Engineering Manager │
│      CLOSES         │
└─────────────────────┘
```

And you discovered an important distinction:

### Thinking Mode

Expensive reasoning.

```text
Engineering Manager
        ↓
Domain Expert
        ↓
Architect
```

### Execution Mode

More repetitive, implementation-oriented work.

```text
Backend Engineer
        ↓
Reviewer
        ↓
QA
```

That distinction should become part of your methodology.

---

# The most important discovery: don't make the pipeline unconditional

Your current workflow is conceptually:

```text
A → B → C → D → E → F
```

But your real workflow should eventually be:

```text
A
 ↓
B
 ↓
C
 ↓
[Decision Gate]
 ↓
D
 ↓
[Validation Gate]
 ↓
E
 ↓
[Quality Gate]
 ↓
F
```

For example:

### Gate 1 — Architecture

If the Architect says:

```text
REJECTED
```

stop.

Don't spend tokens on Backend Engineer, Reviewer and QA.

If:

```text
APPROVED
```

continue.

### Gate 2 — Implementation

If Backend Engineer reports:

```text
Implementation already exists.
Build: PASS
Tests: PASS
```

don't ask another agent to implement it again.

Go directly to Reviewer.

### Gate 3 — Review

If Reviewer says:

```text
APPROVED
```

→ QA.

If:

```text
CHANGES REQUIRED
```

→ Backend Engineer gets **only the findings**, not another complete architectural discussion.

That could substantially reduce your AI consumption.

---

# Your Archive experiment already demonstrated this

The Backend Engineer started and discovered:

> The implementation is already complete.

That's an important signal.

A better execution workflow would have recognized:

```text
Backend Engineer
       ↓
Inspect state
       ↓
Already implemented?
       │
      YES
       ↓
Verify
       ↓
Reviewer
```

rather than treating every session as:

```text
Backend Engineer
       ↓
Implement everything
```

This is exactly where your `08-agent-memory` becomes valuable.

---

# And the memory system is working

This is probably the second-biggest success of the experiment.

You now have:

```text
08-agent-memory/

├── current-state.md
├── active-work.md
└── known-issues.md
```

And the execution agent actually updated them.

That means future sessions don't need to rediscover:

* what has already been implemented
* what's currently being worked on
* what architectural problems are known
* what remains open
* what has already been validated

That's a major token optimization **and** a major engineering-quality improvement.

---

# I also really like the specialized execution-model idea

But I would **not implement it yet**.

Freeze the methodology first.

Later, you can evolve:

```text
THINKING MODEL
────────────────────────
Strong reasoning model

Engineering Manager
Domain Expert
Architect


EXECUTION MODEL
────────────────────────
Fast/efficient model

Backend Engineer
Documentation
Refactoring


VERIFICATION MODEL
────────────────────────
Independent model

Reviewer
QA
Security
```

And potentially:

```text
SPECIALIST MODELS
────────────────────────
Vision model       → UI / screenshots
Coding model       → implementation
Research model     → external research
Documentation      → documentation
Security model     → security analysis
```

The architecture of your AI Engineering Platform should therefore be **model-agnostic**.

Your role definitions should describe **responsibility**, not "use Claude."

That's extremely important.

Today:

```text
Architect → Claude Opus
```

Tomorrow:

```text
Architect → GPT
```

or:

```text
Architect → Gemini
```

The engineering process shouldn't care.

---

# And yes: freeze the prompts

I strongly recommend you create a kind of **v1 baseline** now.

Something like:

```text
Ankhora AI Engineering Methodology
Version 1.0
```

Freeze:

* `06-team/*`
* `05-workflows/*`
* `AGENTS.md`
* `ai/README.md`
* session template
* memory structure

Then don't constantly optimize them.

Instead, use them.

After maybe **10–20 real engineering sessions**, you'll have evidence about what actually needs changing.

Otherwise you risk entering an infinite loop:

```text
Improve AI process
       ↓
Improve prompts
       ↓
Improve workflow
       ↓
Improve prompts
       ↓
Improve workflow
       ↓
Never build product
```

😄

You've reached the point where **usage should generate the next evolution**.

---

# One thing I'd record from this experiment

I would preserve this lesson somewhere in the engineering methodology:

> **AI engineering sessions should separate architectural reasoning from implementation. Expensive reasoning should establish a stable design before execution begins. Execution agents should inspect the current repository state and avoid repeating completed work. Each stage should have explicit decision gates so unnecessary agents are not invoked.**

That's probably the most important operational lesson you've learned so far.

---

## So I would freeze here

Your current system is already strong enough:

```text
01 Principles
02 Architecture
03 Standards
04 Contexts
05 Workflows
06 Team
07 Decisions
08 Memory
09 Sessions
10 Artifacts
```

And now:

```text
                ANKHORA
                   │
          AI ENGINEERING OS
                   │
       ┌───────────┴───────────┐
       │                       │
   KNOWLEDGE                EXECUTION
       │                       │
 Vision / Principles       Sessions
 Architecture              Team
 Standards                 Workflows
 Contexts                  Gates
 Decisions                 Memory
       │                       │
       └───────────┬───────────┘
                   │
             REAL ENGINEERING
```

**Freeze it. Use it. Build with it.**

Then the next version of the system should be based on actual engineering evidence rather than speculation.

And honestly, the fact that you reached the quota limit is useful data too. You now know that the methodology isn't just an abstract organizational exercise—it has a measurable **AI operating cost**, which means optimization of the workflow itself becomes an engineering problem. That's a very interesting direction for Ankhora.

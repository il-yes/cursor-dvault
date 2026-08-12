# Engineering Session: Channel Aggregate Lifecycle

**Role:** Engineering Manager

**Status:** Thinking Mode — Planning

**Bounded Context:** C3 Collaboration

---

## Objective

Define the correct application-layer architecture for modifying the lifecycle and child components of the `Channel` aggregate.

The `Channel` aggregate currently contains several mutable components, including:

* Slots
* Assignments
* Policies
* Channel properties
* Federation configuration

The immediate concrete example is **Slot lifecycle management**.

The aggregate already exposes domain operations such as:

* `AddSlot`
* `UpdateSlot`
* `RemoveSlot`
* `GetSlotByID`
* `GetSlotByVaultID`
* `GetSlotsByRole`
* `GetGatedSlots`
* `ListSlots`

The question is how these aggregate mutations should be exposed through the application layer and persisted to the cloud repository.

---

# Problem

We do **not** want to automatically create one application use case for every aggregate mutation.

For example, we want to avoid immediately creating:

```text
AddSlotUsecase
UpdateSlotUsecase
RemoveSlotUsecase

AddAssignmentUsecase
UpdateAssignmentUsecase
RemoveAssignmentUsecase

AddPolicyUsecase
UpdatePolicyUsecase
...
```

The current hypothesis is that Channel lifecycle mutations could be handled through a dedicated lifecycle application service/use case responsible primarily for:

```text
load Channel
    ↓
invoke aggregate behavior
    ↓
persist Channel
```

The aggregate remains responsible for all business rules and invariants.

This is a **proposal only** and must be validated before implementation.

---

# Architectural Questions

The Engineering Manager must determine:

### 1. Application Boundary

Should Channel child-component mutations be exposed through:

* a single `ChannelLifecycleUsecase`
* multiple use cases grouped by responsibility
* another application-layer abstraction

Do not assume the proposed lifecycle use case is correct without evaluating the alternatives.

---

### 2. Aggregate Ownership

Determine which responsibilities belong to:

```text
Channel aggregate
```

and which belong to:

```text
Application layer
```

The application layer should orchestrate.

The aggregate should decide.

---

### 3. Commands vs Queries

Explicitly distinguish mutation operations from read operations.

For example:

```text
Mutations:
AddSlot
UpdateSlot
RemoveSlot

Queries:
GetSlotByID
GetSlotByVaultID
GetSlotsByRole
GetGatedSlots
ListSlots
```

Determine whether queries should remain separate from the lifecycle mutation abstraction.

Do not introduce a lifecycle use case merely because the aggregate contains both commands and queries.

---

### 4. Persistence Boundary

Determine whether the expected flow should remain:

```text
Application
    ↓
Repository.GetChannel()
    ↓
Channel aggregate mutation
    ↓
Repository.UpdateChannel()
```

or whether another persistence strategy is justified.

The repository should persist the aggregate as a whole unless the existing architecture explicitly establishes another ownership boundary.

---

### 5. Domain Events

Determine whether mutations such as:

```text
SlotAdded
SlotUpdated
SlotRemoved
```

should produce domain events.

Do not automatically introduce events for every mutation.

Evaluate whether each mutation represents a meaningful business fact that other components may need to observe.

---

### 6. Child Component Boundaries

Determine whether:

```text
Slot
Assignment
Policy
ChannelProperty
Federation
```

are merely internal components of the `Channel` aggregate or whether any of them represent independent domain concepts requiring their own lifecycle/application boundary.

This decision is important before creating generic lifecycle abstractions.

---

### 7. Authorization

Determine where authorization belongs for Channel lifecycle mutations.

The current use cases validate requests but authorization responsibilities are not yet fully defined.

Do not invent an authorization mechanism.

Identify the architectural boundary where authorization should eventually occur.

---

### 8. Existing Patterns

Inspect existing Channel use cases, repository interfaces, domain methods, event contracts, and relevant architecture documentation.

Compare the proposed lifecycle approach against existing patterns before introducing a new abstraction.

---

# Required Thinking Mode Output

The Engineering Manager must produce:

## 1. Problem Understanding

Clearly describe the architectural problem being solved.

## 2. Ownership Analysis

Identify:

* Channel responsibilities
* Application responsibilities
* Repository responsibilities
* Event responsibilities

## 3. Options Considered

Evaluate at least:

```text
Option A — One lifecycle use case

Option B — Separate use cases per mutation

Option C — Alternative application service abstraction
```

## 4. Recommendation

Select the most appropriate architecture and explain why.

## 5. Slot Example

Show the recommended application flow for:

```text
AddSlot
UpdateSlot
RemoveSlot
```

without implementing it.

## 6. Future Scalability

Explain how the chosen architecture would behave when adding:

```text
Assignments
Policies
ChannelProperties
Federation
```

## 7. Risks

Identify architectural risks, especially:

* anemic domain
* god use case
* generic CRUD abstraction
* hidden business rules
* excessive event generation
* cross-context coupling
* authorization ambiguity

## 8. Decision

End with a clear recommendation:

```text
APPROVE
```

or

```text
REJECT
```

or

```text
APPROVE WITH CONDITIONS
```

No implementation should be performed in this step.

---

# Constraints

Follow the Ankhora AI Engineering Platform rules.

Read and respect:

```text
ai/README.md
ai/00-vision
ai/01-principles
ai/02-architecture
ai/03-standards
ai/04-contexts/c3.md
ai/08-agent-memory
ai/09-sessions
```

Load only documents relevant to this architectural question.

Do not modify code.

Do not create implementation files.

Do not create tests.

Do not redesign unrelated parts of C3.

This is a **Thinking Mode architectural session**.

The purpose is to reach a deliberate engineering decision before execution begins.









Yes. Since the Engineering Manager has produced the architectural proposal, the **next Thinking Mode specialist should challenge the domain model**, not start implementing.

Paste this as the next prompt to the Domain Expert:

````text
# Engineering Session: Channel Aggregate Lifecycle

**Role:** Domain Expert

**Status:** Thinking Mode — Domain Validation

**Bounded Context:** C3 Collaboration

**Previous Specialist:** Engineering Manager

---

# Objective

Validate the proposed application-layer architecture for Channel aggregate lifecycle management from a **domain perspective**.

The Engineering Manager recommends:

- Keep `Channel` as the aggregate root
- Keep Slots, Assignments, Properties, and Policy inside Channel
- Group application mutations by responsibility domain
- Introduce `ChannelSlotUsecase` with:
  - `AddSlot()`
  - `UpdateSlot()`
  - `RemoveSlot()`
- Keep business rules and invariants inside `Channel`
- Persist the complete aggregate through `UpdateChannel()`
- Defer child-component domain events until concrete consumers exist
- Keep `FederationSnapshot` outside the initial Channel lifecycle implementation

Do not implement anything.

Your responsibility is to determine whether this model is **actually correct from the domain perspective**.

---

# Domain Questions

## 1. Slot Semantics

Analyze what a `Slot` represents in the C3 Collaboration domain.

Determine:

- Is a Slot an entity or value object?
- Is its identity meaningful only inside Channel?
- What makes a Slot valid?
- Can two Slots have the same `VaultID`?
- Can two Slots have the same role?
- Can multiple Slots have the same role?
- Does `Gated` have domain implications?
- Does changing `Role` affect existing Assignments?
- Does changing `Gated` affect existing Assignments?
- Can a Slot be modified while the Channel is archived?
- Can a Slot be removed while the Channel is archived?
- Can a Slot be removed while it has an Assignment?

Do not assume these rules from CRUD structure.

Derive them from the domain model and existing documentation/code.

---

# 2. Slot Lifecycle

Validate the proposed operations:

```text
AddSlot
UpdateSlot
RemoveSlot
````

For each operation determine:

* Preconditions
* Invariants
* Side effects
* Possible domain errors
* Whether the operation should be allowed in every Channel state

Produce a table:

| Operation | Preconditions | Domain Rules | Possible Errors |
| --------- | ------------- | ------------ | --------------- |

---

# 3. Channel State Interaction

The Channel already has lifecycle states such as:

```text
Pending
Active
Archived
Revoked
```

Determine which Channel states permit:

```text
AddSlot
UpdateSlot
RemoveSlot
AddAssignment
UpdateAssignment
RemoveAssignment
SetPolicy
UpdatePolicy
RemovePolicy
SetProperty
UpdateProperty
RemoveProperty
```

Produce a matrix:

| Operation           | Pending | Active | Archived | Revoked |
| ------------------- | ------: | -----: | -------: | ------: |
| AddSlot             |       ? |      ? |        ? |       ? |
| UpdateSlot          |       ? |      ? |        ? |       ? |
| RemoveSlot          |       ? |      ? |        ? |       ? |
| Assignment mutation |       ? |      ? |        ? |       ? |
| Policy mutation     |       ? |      ? |        ? |       ? |
| Property mutation   |       ? |      ? |        ? |       ? |

Explain every non-obvious decision.

---

# 4. Slot ↔ Assignment Relationship

This is particularly important.

Determine whether:

```text
Slot
   ↓
Assignment
```

is a domain dependency.

Specifically:

* Can a Slot exist without an Assignment?
* Can an Assignment exist without its Slot?
* What happens when a Slot is removed?
* Should its Assignment automatically be removed?
* Should removal fail when an Assignment exists?
* Can an Assignment reference a Slot that no longer exists?
* Can one Slot have multiple Assignments?
* Can one Vault occupy multiple Slots?
* Are these invariants currently enforced?

If the current aggregate does not enforce required invariants, identify them explicitly.

Do not implement them yet.

---

# 5. AddSlot / UpdateSlot Semantics

Inspect the existing implementation:

```go
func (c *Channel) AddSlot(slot Slot) Channel

func (c *Channel) UpdateSlot(updated Slot) bool
```

Determine whether these methods are currently too CRUD-oriented.

For example, evaluate whether:

```go
AddSlot()
```

should simply append a Slot or should validate domain invariants.

Likewise determine whether:

```go
UpdateSlot()
```

should blindly replace the existing Slot or enforce controlled changes.

Pay particular attention to:

* identity changes
* VaultID changes
* role changes
* gating changes
* assignment compatibility
* duplicate Slots
* Channel state

---

# 6. RemoveSlot Semantics

Evaluate:

```go
RemoveSlotByID()
RemoveSlotByVaultID()
```

Determine whether both are legitimate domain operations.

Questions:

* Is `VaultID` a unique identity for a Slot?
* Should removal by VaultID exist?
* Does removing by VaultID create ambiguity?
* Should the aggregate expose both methods?
* Should removal be based only on Slot identity?
* What should happen if multiple Slots belong to the same Vault?

Recommend the correct domain-level API.

---

# 7. Queries

Validate the existing query methods:

```go
GetSlotByID
GetSlotByVaultID
GetSlotsByRole
GetGatedSlots
ListSlots
```

Determine:

* Which are legitimate domain queries?
* Which expose internal representation?
* Whether returning `[]Slot` is correct
* Whether returning pointers is appropriate
* Whether callers can accidentally mutate aggregate state
* Whether copies are required
* Whether these queries belong on the aggregate

Also pay attention to the recently identified mismatch:

```text
GetSlotByID currently returns Slot / value semantics
but existing tests expect nil.
```

Determine the correct domain API.

Do not simply fix the test because the test is failing.

First determine what the domain API should be.

---

# 8. Aggregate Boundary

Validate the Engineering Manager's conclusion that:

```text
Channel
 ├── Slot
 ├── Assignment
 ├── Property
 ├── Policy
 └── FederationSnapshot
```

Determine whether this is a valid aggregate boundary.

For each component answer:

| Component          | Entity / Value Object | Must belong to Channel? | Independent lifecycle? | Candidate Aggregate Root? |
| ------------------ | --------------------- | ----------------------- | ---------------------- | ------------------------- |
| Slot               | ?                     | ?                       | ?                      | ?                         |
| Assignment         | ?                     | ?                       | ?                      | ?                         |
| Property           | ?                     | ?                       | ?                      | ?                         |
| Policy             | ?                     | ?                       | ?                      | ?                         |
| FederationSnapshot | ?                     | ?                       | ?                      | ?                         |

Pay special attention to `FederationSnapshot`.

Do not automatically accept the Engineering Manager's statement that it may become an aggregate root.

Determine what the domain actually requires.

---

# 9. Application Layer Recommendation

Evaluate the proposed:

```text
ChannelSlotUsecase
    AddSlot()
    UpdateSlot()
    RemoveSlot()
```

from the domain perspective.

Answer:

* Does this grouping represent a real domain responsibility?
* Is "Slot Management" meaningful in C3?
* Is the application layer correctly reflecting the aggregate boundary?
* Does method-level dispatch make sense?
* Would separate use cases better represent meaningful business operations?
* Is the proposed abstraction merely hiding CRUD?

Provide a clear recommendation.

---

# 10. Events

The Engineering Manager recommends initially avoiding events for:

```text
SlotAdded
SlotUpdated
SlotRemoved
```

Validate this decision.

Determine whether these changes represent:

* internal aggregate state changes only
* meaningful domain facts
* collaboration lifecycle facts
* facts that should eventually be consumed by TraceCore
* facts relevant to realtime/federation

Do not recommend events merely because they technically could exist.

Determine whether the domain requires them.

---

# 11. Domain Invariants to Preserve

Identify the complete set of invariants that should hold after:

```text
AddSlot
UpdateSlot
RemoveSlot
```

Separate them into:

### Existing invariants

Already enforced by the current model.

### Missing invariants

Required by the domain but currently not enforced.

### Unknown invariants

Require a business decision before implementation.

This distinction is important.

Do not invent business rules.

---

# 12. Required Output

Produce a structured Domain Expert report containing:

## Domain Understanding

What a Channel Slot actually represents.

## Validated Decisions

Which Engineering Manager conclusions are correct.

## Challenged Decisions

Which conclusions require modification.

## Slot Invariants

Explicit domain rules.

## Channel State Matrix

Allowed/forbidden mutations by Channel status.

## Assignment Relationship

Precise Slot ↔ Assignment semantics.

## Aggregate Boundary

Final recommendation for Channel and its children.

## Application Layer Guidance

Whether `ChannelSlotUsecase` is appropriate.

## Event Guidance

Whether child mutations should produce domain events.

## Open Domain Questions

Only questions that genuinely require human/business decisions.

## Final Verdict

Choose one:

```text
APPROVED
APPROVED WITH CONDITIONS
REJECTED
```

Do not implement code.

Do not create use cases.

Do not modify files.

Do not make architectural decisions outside the domain expert's responsibility.

The purpose of this stage is to validate the **domain model before execution begins**.




# Session — Channel Desktop → Cloud Integration

## Date

2026-08-15

## Objective

Complete the Desktop → Cloud Channel Create/List/Display workflow.

## Reference Architecture

Workspace was used as the canonical implementation pattern.

## Problems Discovered

### 1. Channel client was incorrectly embedded in c3_repository.go

Fixed by creating:

internal/tracecore/channel_client.go

### 2. Cloud wire format differed from Desktop domain JSON

Cloud used Go-default field names.

Fixed with:

internal/tracecore/types/channel_cloud_dto.go

### 3. Desktop accepted malformed Cloud responses

Fixed with strict response validation.

### 4. Federation reported zero recipients as an error

Fixed Cloud FederationAdapter so zero recipients is a successful no-op.

### 5. Cloud child persistence corrupted data across Channels

Slots and assignments used globally incorrect primary keys.

Fixed with composite keys scoped by channel.

## Validation

Create:
PASS

List:
PASS

Display:
PASS

Cloud persistence:
PASS

Same template used by multiple Channels:
PASS

Slots:
PASS

Assignments:
PASS

Federation empty audience:
PASS

## Architectural Lessons

1. Workspace is the canonical Desktop → Cloud client pattern.
2. Each Cloud-backed bounded context should have a dedicated TraceCore client.
3. Never bypass the dedicated client with direct HTTP helpers.
4. Never blindly deserialize Cloud aggregates into Desktop domain objects.
5. Cloud wire formats belong behind the TraceCore boundary.
6. Never fabricate successful empty responses on decode failure.
7. Empty federation audience can be a valid no-op.
8. Child persistence identity must respect aggregate boundaries.
9. Creation, participant assignment and activation are separate lifecycle stages.
10. Do not modify the domain merely to accommodate transport problems.
````

---
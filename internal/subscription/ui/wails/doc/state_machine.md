**production-grade SaaS subscription state machine**.
This is the **exact pattern used by** companies like Stripe, Netflix, and GitHub for billing lifecycles.

This will make your **automatic billing + provisioning system extremely reliable**.

---

# 1️⃣ The Core Subscription Lifecycle

A robust SaaS subscription typically moves through these states:

```
pending
   ↓
trialing
   ↓
active
   ↓
past_due
   ↓
grace_period
   ↓
suspended
   ↓
canceled
```

You won’t necessarily use every state immediately, but designing for them **prevents painful migrations later**.

---

# 2️⃣ State Definitions

### 🟡 `pending`

Subscription created but **not yet activated**.

Typical reasons:

* payment setup not finished
* crypto transaction not confirmed
* provisioning not done yet

Example:

```
User clicks "subscribe"
↓
payment intent created
↓
subscription status = pending
```

---

### 🔵 `trialing`

User is in **free trial**.

```
TrialEndsAt > now
```

Example:

```
Trial started
↓
trialing
↓
trial ends
↓
payment attempt
```

---

### 🟢 `active`

User has **paid and access is enabled**.

All features unlocked.

```
active = true
```

Typical conditions:

```
payment_successful
AND
provisioning_complete
```

---

### 🟠 `past_due`

Payment **failed**, but service is still allowed temporarily.

Example:

```
billing attempt failed
↓
past_due
↓
retry payment
```

This prevents **instant service shutdown**.

---

### 🟣 `grace_period`

Temporary window to recover payment.

Example:

```
past_due
↓
7 day grace
↓
retry billing
```

Many SaaS products do:

```
grace_period = 3–14 days
```

---

### 🔴 `suspended`

User access **temporarily disabled**.

Example:

```
grace period expired
↓
suspended
```

Features disabled but subscription **can still be recovered**.

---

### ⚫ `canceled`

Subscription permanently terminated.

Reasons:

* user cancellation
* fraud
* manual admin action

---

# 3️⃣ Recommended Go Enum

Instead of free-form strings, use a typed enum.

```go
type SubscriptionStatus string

const (
	StatusPending     SubscriptionStatus = "pending"
	StatusTrialing    SubscriptionStatus = "trialing"
	StatusActive      SubscriptionStatus = "active"
	StatusPastDue     SubscriptionStatus = "past_due"
	StatusGrace       SubscriptionStatus = "grace_period"
	StatusSuspended   SubscriptionStatus = "suspended"
	StatusCanceled    SubscriptionStatus = "canceled"
)
```

Your model becomes safer.

---

# 4️⃣ Example Lifecycle

Example for a **new user**.

```
Signup
↓
Subscription created
status = pending

↓
Payment success
↓
Provision vault / infra
↓
status = active
```

---

Example for **failed billing**:

```
active
↓
billing attempt fails
↓
past_due
↓
retry payment
↓
success → active
OR
fail → grace_period
↓
fail → suspended
```

---

# 5️⃣ Example Automatic Billing Worker

Your billing worker might run daily.

Pseudo-logic:

```go
func HandleBilling(sub *Subscription) {

	if sub.Status == StatusActive && time.Now().After(sub.NextBillingDate) {

		err := ChargeCustomer(sub)

		if err != nil {
			sub.Status = StatusPastDue
			sub.RetryCount++
			return
		}

		sub.NextBillingDate = nextCycle()
	}
}
```

---

# 6️⃣ Grace Period Worker

Another worker manages grace logic.

```go
if sub.Status == StatusPastDue && retryAttemptsExceeded {
    sub.Status = StatusGrace
    sub.GraceEndsAt = now + 7days
}
```

---

# 7️⃣ Suspension Worker

```go
if sub.Status == StatusGrace && now > sub.GraceEndsAt {
    sub.Status = StatusSuspended
}
```

---

# 8️⃣ Cancellation Flow

Manual cancellation:

```
active
↓
cancel requested
↓
ends_at = billing_cycle_end
↓
status remains active
↓
cycle ends
↓
canceled
```

This avoids **refund complexity**.

---

# 9️⃣ Final State Diagram

```
        pending
           │
           ▼
        trialing
           │
           ▼
         active
        /      \
       ▼        ▼
   past_due   canceled
       │
       ▼
  grace_period
       │
       ▼
    suspended
```

---

# 🔟 Why This Is Powerful

Your platform will now support:

✔ automatic billing retries
✔ payment recovery
✔ graceful shutdown
✔ trial subscriptions
✔ infrastructure provisioning delay
✔ crypto payment verification

Exactly like **enterprise SaaS billing systems**.

---

# 11️⃣ One Last Improvement

Add these fields to your model:

```go
RetryCount   int
GraceEndsAt  int64
```

These power **automatic recovery**.

---

✅ With your current architecture (event-driven activator + provisioning workers), this state machine fits **perfectly**.

---

If you'd like, I can also show you the **next big upgrade**:

**How to implement a Stripe-level billing engine in Go**
(with invoices, ledger, retries, and schedules).

That would make your system **very close to a real SaaS billing platform.** 🚀

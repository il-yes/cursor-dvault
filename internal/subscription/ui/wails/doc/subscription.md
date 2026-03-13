**Stripe / AWS / modern SaaS systems**

The goal is to **separate concerns** so the model becomes easier to maintain, scale, and evolve.

Right now your `Subscription` struct mixes **4 different concerns**:

1️⃣ Subscription lifecycle
2️⃣ Billing / payment metadata
3️⃣ Entitlements / features
4️⃣ Infrastructure provisioning

When systems grow, this becomes **very hard to maintain**.

The solution used by companies like **Stripe, GitHub, AWS, Notion, Linear** is to split the model.

---

# 1️⃣ Recommended Subscription Architecture

Instead of one giant struct, split into:

```
Subscription
SubscriptionBilling
SubscriptionFeatures
SubscriptionProvisioning
```

So the database becomes:

```
subscriptions
subscription_billing
subscription_features
subscription_provisioning
```

---

# 2️⃣ Clean Core Subscription Model

This becomes the **source of truth**.

```go
type Subscription struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`

	Tier      string `json:"tier"`
	Status    string `json:"status"` // pending, active, past_due, canceled

	Rail      string `json:"rail"` // traditional | crypto

	StartedAt time.Time `json:"started_at"`
	EndsAt    time.Time `json:"ends_at"`

	TrialEndsAt time.Time `json:"trial_ends_at"`

	BillingCycle string `json:"billing_cycle"` // monthly, yearly

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Version int64 `json:"version"`
}
```

This struct now handles **only lifecycle**.

---

# 3️⃣ Billing Model (Payment Layer)

Payment metadata belongs here.

```go
type SubscriptionBilling struct {
	ID             string `json:"id"`
	SubscriptionID string `json:"subscription_id"`

	PaymentMethod string `json:"payment_method"`

	Price float64 `json:"price"`

	NextBillingDate time.Time `json:"next_billing_date"`

	StripeSubscriptionID string `json:"stripe_subscription_id"`

	PaymentIntent string `json:"payment_intent"`

	StellarWallet string `json:"stellar_wallet"`
	StellarScheduleID string `json:"stellar_schedule_id"`

	TxHash string `json:"tx_hash"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

Now **Stripe / Stellar / crypto** logic stays here.

---

# 4️⃣ Features / Entitlements Model

This is what your **app actually reads**.

```go
type SubscriptionFeatures struct {
	SubscriptionID string `json:"subscription_id"`

	StorageGB int
	StorageType string

	CloudBackup bool
	MobileApps bool

	SharingLimit int
	UnlimitedSharing bool

	VersionHistory bool
	VersionHistoryDays int

	Tracecore bool
	GitCLI bool

	BrowserExtension bool
	ThreatDetection bool

	TeamFeatures bool
	TeamSize int

	Compliance datatypes.JSONSlice[string]

	SLA string
}
```

Notice something important:

❌ No billing data
❌ No payment intent
❌ No Stripe IDs

Only **capabilities**.

---

# 5️⃣ Provisioning Model (Infrastructure Layer)

Your system already has:

* Vault provisioning
* Pinata
* IPFS
* Tracecore

So track provisioning separately.

```go
type SubscriptionProvisioning struct {
	SubscriptionID string

	VaultID string

	IPFSInitialized bool
	PinataConfigured bool
	TracecoreEnabled bool

	ProvisionStatus string // pending, provisioning, completed, failed

	ProvisionedAt time.Time
}
```

This allows:

```
Payment success
    ↓
Subscription activated
    ↓
Provisioning worker runs
    ↓
Vault created
```

Exactly like **Netflix infrastructure provisioning**.

---

# 6️⃣ Resulting Architecture

Now your cloud system becomes:

```
Subscription
     │
     ├── SubscriptionBilling
     │
     ├── SubscriptionFeatures
     │
     └── SubscriptionProvisioning
```

Much cleaner.

---

# 7️⃣ Benefits of This Design

### Easier migrations

You can change billing providers without touching features.

### Better security

Payment metadata isolated.

### Easier feature rollout

You modify **features table only**.

### Easier provisioning retries

Provisioning status is tracked independently.

### Better analytics

Billing queries become simpler.

---

# 8️⃣ Example Full Subscription (Combined View)

Your API can still return a **combined view**:

```json
{
  "subscription": { ... },
  "billing": { ... },
  "features": { ... },
  "provisioning": { ... }
}
```

But internally the system stays **cleanly separated**.

---

# 9️⃣ This Matches How Big Systems Work

Examples:

| Company | Equivalent models                          |
| ------- | ------------------------------------------ |
| Stripe  | Subscription + Invoice + PaymentIntent     |
| AWS     | Account + Billing + Entitlements           |
| GitHub  | Account + Plan + Billing                   |
| Notion  | Workspace + Plan + Billing                 |
| Netflix | Subscription + Entitlements + Provisioning |

You're building something **very similar**.

---

# 10️⃣ Final Advice

Your current model **works**, but as the platform grows you’ll thank yourself for splitting it.

A good rule:

```
Subscription → lifecycle
Billing → money
Features → capabilities
Provisioning → infrastructure
```

---



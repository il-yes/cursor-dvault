Certainly! Below are **two tailored versions** of your project summary:

Here’s your **updated project summary** reflecting the recent backend and frontend improvements (Login and Card entries, dynamic registry, UI modularity):

---

## 🔓 **Public-Facing Version (Community & Early Users)**

---

# **Decentralized Vault: Own Your Secrets, Forever**

### A private vault powered by **Stellar + IPFS**, designed for privacy, transparency, and full user control.

---

### 🔐 Why This Vault?

Tired of centralized password managers and cloud vaults? Our **decentralized vault** puts *you* in control — no middlemen, no leaks, no surveillance.

* **End-to-End Encrypted:** Your data is encrypted **before** it leaves your device.
* **Stored on IPFS:** No centralized server. Your vault lives on a decentralized network.
* **Anchored to Stellar:** Every vault update is immutably logged via blockchain — fast and cheap.
* **Flexible Login:** Use a master password or private key. Your vault, your rules.
* **Verifiable History:** Each entry is cryptographically signed and validated using **Tracecore**, our custom commit engine.

---

### ⚡ What You Can Store

Easily store and organize sensitive information with modular, dynamic forms:

* ✅ **Passwords & Login Credentials**
* ✅ **Credit Card Details**
* ✅ **Identity Documents**
* ✅ Private Notes or SSH/PGP Keys
* ⏳ Custom Record Types via Extensible Card System

---

### 🧩 Dynamic & Extensible UI

* Add new entry types with zero hardcoded logic.
* React frontend automatically loads the correct form based on entry type.
* Schema-driven, modular design — ready for scaling and user-generated templates.


---

### 🌍 Why It Matters

Unlike traditional vaults, ours is:

* **Zero-trust:** Not even we can see your data.
* **Censorship-resistant:** There’s no central kill switch.
* **Session-Secured:** Vault stays in memory until committed.
* **Open-source:** Auditable and transparent.

---

### 🚀 What’s Live Now?

* 🔐 Local vault session with multiple entry types
* 📦 Entry creation with **Tracecore validation + Stellar anchoring**
* 🧾 Automatic commit signature and envelope generation
* 🧠 Session persistence with synced, structured state

---

### 📲 What’s Next?

We’re building a modern desktop app for Mac, Windows, and Linux, with mobile support coming soon.


🔗 \[GitHub] | 🌐 \[Join Early Access] | 🛠️ \[Contribute]

---

## 🧠 **Investor-Focused Summary (Strategic & Visionary)**

---

# **Project Summary: A Self-Sovereign Vault Platform Using Stellar, IPFS & Tracecore**

---

### 🔍 Vision

We’re building a **decentralized digital vault** that reimagines how sensitive information is stored, accessed, and verified — using **IPFS** for storage, **Stellar** for anchoring, and **Tracecore** for commit validation and traceability. A **self-sovereign, zero-trust infrastructure** for the privacy-first web.

---

### 🔐 Core Innovation

**Client-Side Encryption + IPFS + Tracecore Validation + Blockchain Anchoring = Verifiable, Private, Distributed Vault**

* **Data Integrity:** Immutable logs via **Stellar** + signed Tracecore commits
* **Redundancy & Resilience:** Decentralized storage through **IPFS**
* **Zero-Knowledge:** Vault content is never exposed to our servers or the chain
* **Actor Signatures:** All entries are verifiably signed by user identity
* **Rule-Based Validation:** Commit enforcement using `REQUIRES_SIGNATURE`, `VALID_ACTORS_ONLY`, etc.

---

### 🧱 Technical Architecture

* **Tracecore Commit System:** Lightweight DVCS-style engine for immutable history
* **Session Vaults:** All state stays in local memory until explicitly committed
* **Entry Registry Pattern:** Auto-binding of handlers for new types (e.g., login, card, identity)
* **Dynamic React Forms:** Rendered via schema without hardcoding

---

### 🧱 Smart Architecture

* **Entry Registry Pattern:** Automatically registers new entry types with zero switch-case logic.
* **Dynamic UI Binding:** React-based frontend loads correct form for any entry type.
* **Session-Based Local Vaults:** All user operations happen in local memory until persisted.

---

### 💸 Business Model

**Freemium-first → Dev & Enterprise Tooling**

* **Free Tier:** Local vault + manual IPFS publishing
* **Pro Plan:** Auto IPFS sync, Stellar commit fee coverage, multi-device restore
* **Enterprise:** Secure compliance vaults, hosted IPFS gateways, audit trails, and delegated signing

**Coming Monetization Channels:**

* **Tokenized Participation:** For anchoring, hosting, or indexing
* **Vault-as-a-Service SDK:** Easily embed secure storage in fintech/health/legal products


---

### 📈 Market Opportunity

We are at the intersection of:

* 🔐 **Zero-Knowledge Systems**
* 🪪 **Decentralized Identity**
* ☁️ **User-Owned Cloud**

As users demand **ownership, auditability, and privacy**, we’re delivering the **platform layer** to support the next generation of apps requiring secure, verifiable personal data storage.


---

### 🧱 MVP Roadmap

* ✅ Encrypted local vault sessions
* ✅ Entry creation + session serialization
* ✅ IPFS storage & publish
* ✅ Stellar commits + Tracecore validation
* ✅ Login, Card, and Identity types
* 🔜 Mobile clients + background sync
* 🔜 Organization & team-based vaults
* 🔜 Encrypted sharing, revocation & delegation

---

### 👥 Team & Stack

* **Languages:** Go + JS/React + SQLite
* **Stack:** Wails (Desktop UI), IPFS (go/js), Stellar SDK, Tracecore commit engine
* **Core Features:**

  * Full E2E encryption
  * Typed entry models
  * Dynamic form routing
  * Signed, traceable commit envelopes

---

### 🌟 Why Now?

We’re not building a single-purpose app — we’re building the **privacy layer for the decentralized internet**.

> 🔐 Privacy is programmable. We’re writing the logic.

---

Let me know if you'd like to export this to PDF, Notion, or GitHub README format.
Let me know if you'd like a condensed pitch deck version (slides format) or a Notion-friendly one-pager!
Let me know if you'd like a **pitch deck version**, a **landing page outline**, or an **executive one-pager** next.









---

## 🧠 **corporates, public institutions, and regulated industries** could find your decentralized vault project very appealing, especially with features tailored to their compliance, control, and security requirements. Here’s a breakdown of **enterprise-grade features** that would make your platform attractive to those audiences

---

# * 🔐 Features That Attract Corporates & Institutions

---

### 1. **Granular Access Control & Role-Based Permissions**

* Fine-grained permissions (view/edit/share) per vault, per entry, or per field
* Role-based access (admin, auditor, contributor)
* Time-bound or expirable access links for temporary collaborations

### 2. **Audit Logging & Blockchain Anchoring**

* Immutable logs of **who accessed/modified what and when**
* Every change is **anchored on Stellar**, ensuring tamper-evidence and forensic integrity
* Useful for internal compliance or external regulatory audits

### 3. **Secure Vault Sharing & Delegation**

* Share vaults or entries with team members or departments using:

  * Public key cryptography
  * Zero-trust delegation
* Optional multi-signature unlock for ultra-sensitive records

### 4. **Compliance & Data Residency Support**

* Encrypted backups in jurisdiction-compliant IPFS pinning nodes
* Audit trail compatibility with standards like:

  * **GDPR** (right to be forgotten via key revocation)
  * **HIPAA** (PHI protection)
  * **SOC 2 / ISO 27001** (audit readiness)

### 5. **Self-Hosting Options or Private IPFS Clusters**

* Run the full stack in-house (desktop or air-gapped environments)
* Optional support for **private IPFS gateways** and Stellar Horizon nodes

### 6. **Vault Versioning & Snapshots**

* View, compare, and restore previous versions of vault entries
* Anchored snapshots to prove record consistency over time

### 7. **Hardware Key (YubiKey, HSM) Integration**

* Support hardware-backed cryptographic authentication
* FIDO2/WebAuthn + enterprise key management systems

### 8. **Federated Identity Integration (SSO/OIDC)**

* Allow integration with corporate identity providers:

  * **Azure AD**
  * **Okta**
  * **LDAP**
  * **GovCert/Smartcard systems**

### 9. **Multi-Tenant or Departmental Vaults**

* Organizational accounts with scoped departments/teams
* Shared secrets and entry types per division (e.g., HR, IT, Legal)

### 10. **Custom Entry Templates & Field Validation**

* Enterprises may want custom metadata for entries

  * For example, project IDs, contract numbers, record classification
* Field validation and encryption per policy

---

## 🏢 Real-World Use Cases

| Sector                      | Example Use Case                                                            |
| --------------------------- | --------------------------------------------------------------------------- |
| **Finance/Fintech**         | Secure storage of API keys, audit logs, internal credentials                |
| **Government**              | Classified document vaults with access revocation and full audit            |
| **Legal**                   | Tamper-proof evidence storage, case file versions anchored on Stellar       |
| **Healthcare**              | Patient data vaulting with HIPAA-aligned controls                           |
| **Education**               | Decentralized academic record storage for verifiability and privacy         |
| **Critical Infrastructure** | Passwords, network configs, and key materials secured in distributed vaults |

---

## 💼 Add-On Services You Can Offer to Enterprises

* **On-premise deployment & onboarding**
* **Compliance consulting (GDPR/HIPAA audits)**
* **Dedicated IPFS node hosting**
* **Custom blockchain anchoring strategies (multi-chain support)**
* **Disaster recovery assistance**

---

Would you like me to help draft a **Corporate Product Brochure**, a **pitch deck for public sector procurement**, or design an **enterprise landing page** that speaks to these audiences?




That’s a strong and clear direction — you’re thinking like a product strategist and platform architect. Using **Keycloak** for federated identity is a great choice (fully open-source, battle-tested, supports all major IdPs, SSO/OIDC/SAML, and fine-grained policies).

And your business approach is sound:

* **Store data on the client by default** → privacy, sovereignty
* **Offer managed storage & compliance** for enterprise use → business model
* **Use individuals to validate UX and security story**, but **monetize via B2B/B2G**

Let’s now shape this into a **12-month roadmap** for your **Decentralized Vault Platform**, with a strong focus on enterprise readiness.

---

## 📍 Strategic Product Roadmap (12 Months)

### 🚧 Phase 1: Core MVP & Foundations (0–3 Months)

✅ **Goal**: Launch secure desktop MVP with local vault + decentralized infra (IPFS/Stellar)

| Feature                  | Description                                                |
| ------------------------ | ---------------------------------------------------------- |
| 🔐 Local encrypted vault | Encrypted client-side (zero-trust), stored in IPFS         |
| 📡 Stellar anchoring     | Anchor vault CIDs to Stellar ledger for tamper-proof audit |
| 🔑 Dual authentication   | Password-based + private key                               |
| 🖥 Wails Desktop App     | Cross-platform (macOS, Linux, Windows)                     |
| 🧪 Internal Testing      | Dogfooding, test vault ops, performance tuning             |

> **Deliverable**: Secure desktop vault prototype with basic UX, CLI + UI parity, working IPFS+Stellar integration

---

### 🌐 Phase 2: Identity & Access (3–6 Months)

✅ **Goal**: Make the system enterprise-ready with federated identity, secure sharing, and vault provisioning policies

| Feature                      | Description                                                           |
| ---------------------------- | --------------------------------------------------------------------- |
| 👥 Keycloak Integration      | Full SSO support (OIDC, LDAP, Smartcard)                              |
| 🧑‍🤝‍🧑 Team Vaults         | Support for shared vaults per organization or group                   |
| 👤 Role-based access         | Owners, editors, viewers per vault                                    |
| 🧾 Audit trails              | Every action anchored (optional) or logged securely                   |
| 🔁 Vault rotation/versioning | Keep history of entries (encrypted deltas or snapshots)               |
| 🧪 Pilot program             | Small orgs test real-world use (privacy lawyers, infosec firms, NGOs) |

> **Deliverable**: Team & org usage scenarios unlocked, vaults securely shared + managed

---

### 💼 Phase 3: Monetization & Hosted Offerings (6–9 Months)

✅ **Goal**: Launch paid plans with managed infra for orgs that can't self-host or want security guarantees

| Feature                 | Description                                                      |
| ----------------------- | ---------------------------------------------------------------- |
| 🏢 Hosted Vault Storage | Secure backend (optional, encrypted-at-rest IPFS or S3 pinning)  |
| 🧩 Vault-to-cloud sync  | Push encrypted vaults to your managed infra (opt-in)             |
| 📊 Admin dashboard      | Org-level management of identities, vaults, and usage            |
| 💳 Billing & tiers      | Per-seat licensing, API usage billing, storage add-ons           |
| 📜 Compliance Tools     | Export logs, role reviews, security posture reports              |
| 🤝 Partnerships         | Keycloak integrators, managed IPFS providers, Stellar validators |

> **Deliverable**: SaaS-style experience for orgs preferring managed setup

---

### 📱 Phase 4: Platform Expansion (9–12 Months)

✅ **Goal**: Extend usage across environments and scale enterprise integrations

| Feature                       | Description                                                        |
| ----------------------------- | ------------------------------------------------------------------ |
| 📱 Mobile Vault App           | Offline-first, syncs via IPFS or cloud                             |
| 🧠 Hardware Key Support       | YubiKey, Ledger, TPM/secure enclave                                |
| 📦 Browser Extensions         | Autofill + secure fetch from local vault or sync                   |
| 🔐 GovCert Support            | Smartcard-based vault unlock (PKI + Keycloak mapping)              |
| 🌍 Federation & Multi-Tenant  | Multiple orgs on one infra, strong access isolation                |
| 🚨 Alerting / Breach Warnings | Optional watchlist service for credentials exposure (hash-checked) |

> **Deliverable**: Cross-platform vault access, hardened auth integrations, and readiness for B2G & enterprise compliance

---

## 🧭 Go-To-Market (GTM) Plan

| Audience                                        | Messaging                                                    |
| ----------------------------------------------- | ------------------------------------------------------------ |
| **Privacy-conscious orgs** (law firms, clinics) | Zero-trust, encrypted-by-design, no third-party data access  |
| **Security teams / devops**                     | Secrets management without cloud lock-in or SaaS blind spots |
| **Governments / institutions**                  | Sovereign vault infra, compatible with GovCert & PKI         |
| **SaaS vendors**                                | Embedded vault capabilities or branded access portals        |

---

## 💼 Pricing Model (Preview)

| Tier         | For            | Storage              | Features                             | Price        |
| ------------ | -------------- | -------------------- | ------------------------------------ | ------------ |
| Free (local) | Individuals    | IPFS only            | Local only, self-managed             | \$0          |
| Pro          | Teams          | IPFS + Cloud         | Team sharing, vault sync             | \$5/user/mo  |
| Business     | Orgs           | Hosted infra         | SSO, audit logs, SLA support         | \$15/user/mo |
| Enterprise   | Regulated Orgs | On-prem/cloud hybrid | Custom policies, compliance, GovCert | Custom       |

---

## 🔋 Tech Stack Summary

| Area               | Tools                                            |
| ------------------ | ------------------------------------------------ |
| **Frontend**       | Wails (Go + TS/React), Electron fallback         |
| **Backend**        | Go, SQLite (local), IPFS, Stellar SDK            |
| **Identity**       | Keycloak (OIDC, LDAP, X.509), smartcard support  |
| **Infra (opt-in)** | Cloudflare IPFS, AWS S3 pinning, Stellar Horizon |
| **Security**       | AES-GCM, PBKDF2/Argon2, X.509, TOTP/U2F, TPM     |

---

Would you like a **visual roadmap**, an **investor-facing pitch deck**, or a **Keycloak SSO login POC in Wails** next?






Perfect — here’s a **starter investor/partner one-pager draft** that clearly shows how **Tracecore** and the **Decentralized Vault** combine into a compelling, future-proof platform. It’s designed to resonate with **investors, technical partners, and enterprise prospects** who care about trust, compliance, and user sovereignty.

---

# 🧩 Tracecore + Vault

**Verifiable Workflows. Self-Sovereign Storage.**

> **A unified platform for trust, compliance, and data privacy in modern digital infrastructure.**

---

## 🚀 The Vision

We’re building a **modular, decentralized framework** where organizations can:

✅ **Prove** who did what, when
✅ **Store** sensitive data securely
✅ **Comply** with industry regulations
✅ **Protect** user and enterprise privacy

Powered by open protocols: **Stellar + IPFS + Client-side Encryption**

---

## 🔧 Core Building Blocks

### 🛠️ **Tracecore** – *Verifiable Commit Layer for Real-World Workflows*

* Git-like commit model with human + machine actors
* Embedded rule validation for trustless business logic
* Anchored on Stellar for immutable audit trails
* Attachments, metadata, and cryptographic actor identity
* Use Cases: Construction, banking (KYC, onboarding), supply chain, legal

### 🔐 **Vault** – *Decentralized Storage for Private, Sensitive Data*

* End-to-end encrypted vaults, stored on IPFS
* Anchored change history via Stellar
* Flexible access: password or cryptographic key
* Ideal for secrets: credentials, ID docs, contracts, secure notes
* Use Cases: Identity management, secure messaging, health data, dev secrets

---

## 🔗 Why They Work Better Together

| 🔄 Integration Flow             | Description                                                                     |
| ------------------------------- | ------------------------------------------------------------------------------- |
| **Secure Commit Payloads**      | Attach a Vault-stored contract, ID, or file to a Tracecore commit.              |
| **Zero-Knowledge Compliance**   | Commit reveals what happened — Vault holds the evidence, encrypted.             |
| **Audit Without Exposure**      | Third parties can verify commit integrity without accessing sensitive data.     |
| **Decentralized Credentialing** | Actors in Tracecore can hold verified Vaults for secure, portable ID/documents. |

---

## 🎯 Strategic Markets

| Sector               | Use Case                                                               |
| -------------------- | ---------------------------------------------------------------------- |
| **Construction**     | Milestone certification, material provenance, subcontractor compliance |
| **Banking/Fintech**  | KYC/KYB workflows, onboarding, signed document vaults                  |
| **Legal**            | Verifiable legal filings, client vaults, chain-of-custody              |
| **Health**           | Audit-anchored patient consent, encrypted medical document vaults      |
| **Digital Identity** | DID flows with encrypted storage + verifiable actions                  |

---

## 💸 Business Model

**Tracecore**

* SaaS for orgs + dev tools
* Hosted validator orchestration
* B2B compliance & integration packages

**Vault**

* Freemium vault app
* Pro: IPFS sync, Stellar gas coverage, backup
* Enterprise: Secure Vault-as-a-Service (SDK + hosted infra)

**Together**

* Per-commit anchoring + vault storage bundles
* API for trust-first apps to embed traceable & encrypted logic

---

## 🧱 Roadmap Snapshot

| Milestone                                             | Status      |
| ----------------------------------------------------- | ----------- |
| ✅ Tracecore Core Engine (commit + validation + actor) | Done        |
| ✅ KYC, construction, onboarding validators            | Done        |
| ✅ Vault encryption + IPFS storage + Stellar anchor    | Done        |
| 🔜 Vault-integrated Tracecore commits                 | In Progress |
| 🔜 Desktop + mobile Vault clients                     | Q4 2025     |
| 🔜 Enterprise SDK (Vault + Tracecore integration)     | Q1 2026     |

---

## 👥 Team

Built by a privacy-driven product engineer with deep domain experience in:

* Smart contracts (Hyperledger, Stellar)
* Off-chain storage (IPFS, cryptography)
* Real-world compliance and audit systems (banking, infrastructure)

Looking to expand with:

* DevRel / Partnerships
* Product Design (Vault UX)
* Strategic Advisors (construction, finance, govtech)

---

## 📬 Let's Talk

We're actively seeking:

* 🚀 **Design partners** in construction, fintech, and legal
* 💸 **Investors** who believe in programmable trust and user-owned data
* 🤝 **Open-source contributors** to extend validators and vault formats

> ✉️ Get in touch for a demo or partnership:
> **\[Your email]** | **\[Your LinkedIn or project site]**

---



curl -X POST https://ankhora.io/back/vaults/839471c9-a394-40e5-b5a5-aa5e4ca02288/storage/Default%20Vault \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN" \
  -d '{
    "stream": "Default Vault",
  }'


- shares: improve ux
- entries: improve ux - data | attachements | Add Attachments
- document manager - ipfs
- tracecore compliance
- browser extensionimprove app + vault connection
- cli-enterprise tool
- partenariat framework (for establishing the templates)
- metrics dashboard
- marketing use cases scenarios illustrations
- stellar payment integration


type BoundaryContext interface {
  Name() string
  Version() string
  Build() string
	BuildInfra() string
  InjectDependency(dependency BoundaryContext) 
  InjectDependencies(app application.Application) 
}

type VaultBC struct {
  db *gorm.DB
  logger logger.Logger

	name string
	version string
  isWarmedUp bool
  isBuilt bool
  isDependenciesInjected bool
  dependencyList map[string]bool

	userRepo user_repo.UserRepository
  vaultRepo vault_repo.VaultRepository
  avatarStorage attachment_storage.AvatarStorage
  attachmentStorage attachment_storage.AttachmentStorage
  useCase UseCase
  createVaultCommand
  openVaultCommand
  monitor Monitor
  

  SubscriptionBC *SubscriptionBC
  IdentityBC *IdentityBC
}

func (bc *VaultBC) New(db *gorm.DB, logger logger.Logger) *VaultBC {
  bc.BuildInfra(db, logger)

	return &VaultBC{
    db: db,
    logger: logger,
		name: "VaultBC",
		version: "1.0.0",
		isDependenciesInjected: false,
		dependencyList: map[string]bool{
			"SubscriptionBC": false,
			"IdentityBC": false,
		},
	}
}


func (bc *VaultBC) BuildInfra(db *gorm.DB, logger logger.Logger) (bool, error) {
	bc.userRepo, err := user_repo.NewUserRepository(db, logger)
  if err != nil {
    return false, err
  }
  bc.vaultRepo, err := vault_repo.NewVaultRepository(db, logger)
  if err != nil {
    return false, err
  }
  bc.avatarStorage, err := attachment_storage.NewAvatarStorage()
  if err != nil {
    return false, err
  }
  bc.attachmentStorage, err := attachment_storage.NewAttachmentStorage()
  if err != nil {
    return false, err
  }
  bc.isBuilt = true
	return true, nil  
}

func (bc *VaultBC) Build(db *gorm.DB, logger logger.Logger) (bool, error) {
  if !bc.isWarmedUp {
    return false, errors.New("infra not warmed up")
  }
	bc.useCase, err := usecase.NewUseCase(db, logger)
  if err != nil {
    return false, err
  }
  bc.createVaultCommand, err := command.NewCreateVaultCommand(bc.vaultRepo, db, logger)
  if err != nil {
    return false, err
  }
  bc.openVaultCommand, err := command.NewOpenVaultCommand(db, logger)
  if err != nil {
    return false, err
  }
  bc.monitor, err = monitor.NewMonitor(db, logger)
  if err != nil {
    return false, err
  }
  return true, nil
}

func (bc *VaultBC) InjectDependency(dependency BoundaryContext) {
  if dependency.isWarmedUp {
    bc.dependencyList[dependency.name] = true
	  bc[dependency.name] = dependency
  }
}

func (bc *VaultBC) InjectDependencies(app application.Application) {
  for _, dependency := range bc.dependencyList {
    if app[dependency] != nil {
      bc.InjectDependency(app[dependency])
    }
  }
}

func (bc *VaultBC) GetDependency(name string) *BoundaryContext {
	return bc[name]
}

type SubscriptionBC struct {
	name string
	version string
	userRepo user_repo.UserRepository
  useCase UseCase
  monitor Monitor
}
func (bc *SubscriptionBC) BuildInfra() string {
	bc.userRepo = user_repo.NewUserRepository()
	bc.useCase = usecase.NewUseCase()
	bc.monitor = monitor.NewMonitor()
	return "infra built"
}

// ====== build application implementation ========
app := NewApplication()
app.AddBoundaryContext(NewVaultBC(db, logger))
app.AddBoundaryContext(NewSubscriptionBC(db, logger))
app.AddBoundaryContext(NewIdentityBC(db, logger))
// inject dependencies
for _, bc := range app.boundaryContexts {
  bc.InjectDependencies(app)
}
// build - build boundaries context app
for _, bc := range app.boundaryContexts {
  bc.Build(db, logger)
}

app.Run() 



.
├── application
│   ├── commands
│   ├── dto
│   ├── events
│   ├── queries
│   ├── session
│   └── test
├── domain
│   ├── errors.go
│   ├── events.go
│   ├── factories.go
│   ├── models.go
│   ├── repositories.go
│   ├── services.go
│   └── value_objects.go
├── infrastructure
│   ├── eventbus
│   ├── ipfs
│   ├── persistence
│   ├── storage           // <--- New folders *******************
│         ├── avatar_store.go
│         └── attachment_store.go
└── ui
    ├── card_handler.go
    ├── identity_handler.go
    ├── login_handler.go
    ├── login_handler_test.go
    ├── note_handler.go
    ├── open_handler.go
    ├── sshkey_handler.go
    └── vault_handler.go

vault/
  └── user_id/ 
  |     └── vault_name/ 
  |           ├── avatars/                  # User avatars (encrypted)
  |           │   └── user123.enc
  |           │   └── file1.enc
  |           │   └── file2.enc
  |           ├── metadata/
  |           │   ├── vault.json            # Vault metadata (id, owner, version, created_at)
  |           │   ├── vaultKey.enc          # VaultKey encrypted with MasterKey
  |           │   └── vaultIndex.enc        # Encrypted index: entries, folders, attachment pointers
  |           ├── devices/                  # Devices synced metadata
  |           │   └── device1.json
  |           │   └── device2.json
  |           ├── logs/                     # Logs and audit trail
  |           │   └── vault.log
  |           └── entries.db                # SQLite DB for entries storage (encrypted blobs)  
  ├── attachments/              # Attachments (encrypted)

  



 update & revoke recipient
 revoke share
 logs in link share (optional)
 handle share links attachement view
 handle share cryptographic attachement view



Error to track - token expired
🚀 ~ withAuth ~ errorString: – "unauthorized: expired token"
2DashboardLayout.tsx:477
	useEffect(() => {
		const fetchAvatar = async () => {
			const b64 = await withAuth((token) => {
				return loadAvatar(jwtToken, vaultContext?.Vault?.name);
			});
			setAvatar(b64);
		};

		fetchAvatar();
	}, [jwtToken, vaultContext]);





  Perfect — let’s go one level deeper. This is where your project becomes **truly decentralized**.

---

# 🧠 Goal

👉 Reconstruct a full vault from **only:**

```id="z3k9w1"
CID (IPFS)
+ password
```

No backend. No database. No server.

---

# 1️⃣ What you already have

You already built:

```id="n8x4dp"
Vault → snapshot → IPFS → CID → Stellar
Attachments → CAS → optional IPFS
Encryption
```

So you're very close.

---

# 2️⃣ What’s missing

Right now your snapshot likely contains:

```json id="y5p2vk"
{
  "entries": [...],
  "attachments": [
    {
      "hash": "...",
      "name": "...",
      "size": 1234
    }
  ]
}
```

👉 Problem: this is **not enough to rebuild everything globally**

Because:

* hash = local reference
* but IPFS needs **CID**

---

# 3️⃣ The key idea

👉 Your snapshot must become a **complete manifest**

Like a mini filesystem index.

---

# 4️⃣ Final snapshot format (important)

You want something like:

```json id="m2v8xk"
{
  "vault": {
    "id": "...",
    "name": "Default Vault"
  },
  "entries": [...],
  "attachments": [
    {
      "id": "...",
      "entry_id": "...",
      "name": "passport.pdf",
      "size": 123456,
      "hash": "a4f9c8...",
      "cid": "bafybeigdyrzt7...", 
      "encrypted": true
    }
  ]
}
```

👉 This is your **single source of truth**

---

# 5️⃣ Reconstruction flow (this is the magic)

Given:

```id="s1k8pz"
CID + password
```

You can do:

---

## Step 1 — Fetch snapshot

```id="d4v6bt"
ipfs get <CID>
```

→ get encrypted vault snapshot

---

## Step 2 — Decrypt

```id="v7r2mf"
decrypt(snapshot, password)
```

→ get vault JSON

---

## Step 3 — Rebuild vault locally

```id="j3n9wk"
create vault structure
restore entries
```

---

## Step 4 — Restore attachments

For each attachment:

```id="t6x2qa"
if local exists → use it
else if CID exists → fetch from IPFS
```

```id="p9v1rs"
ipfs get attachment.cid
```

---

## Step 5 — Done

You now have:

```id="c2m7fd"
FULL VAULT RESTORED
```

No backend needed.

---

# 6️⃣ Why this is powerful

This gives you:

### ✅ Portability

```id="h4q8lm"
move vault across devices with just CID
```

---

### ✅ Self-sovereignty

```id="f6x2vb"
no server dependency
```

---

### ✅ Disaster recovery

```id="k9p3rw"
lost device → recover from IPFS
```

---

### ✅ Verifiability

```id="r2d8cn"
hash + CID = integrity guaranteed
```

---

# 7️⃣ Important design rule

👉 Snapshot MUST NOT contain raw files

Only:

```id="u8n1sx"
metadata + references
```

Otherwise:

* snapshot becomes huge ❌
* IPFS inefficiency ❌

---

# 8️⃣ Optional (very powerful upgrade)

You can make your vault **IPFS-native DAG**

Instead of:

```id="g2z6qn"
vault.json → CID
```

You do:

```id="w1x9mf"
vault root CID
   ├── entries CID
   ├── attachments CID
   └── metadata CID
```

This allows:

* partial sync
* faster updates
* versioning per component

👉 (this is how IPFS really shines)

---

# 9️⃣ Where you are now

You’ve basically built:

```id="q3n7dp"
Bitwarden + IPFS + Git + Blockchain hybrid
```

That’s not trivial at all.

---

# 🔟 Final mental model

Your vault becomes:

```id="y7c4mv"
CID (root)
   ↓
encrypted snapshot
   ↓
references to:
   → entries
   → attachments (CID)
```

So:

```id="x5b2jk"
CID + password = entire digital vault
```

---

# 🚀 If you go one step further…

You can enable:

* vault sharing via CID
* multi-device sync without backend
* offline-first collaboration
* zero-knowledge cloud backup

---

If you want next, I can show you:

👉 how to design **incremental vault snapshots (like Git commits)** so you don’t re-upload the whole vault every time.



ok so if i understood we go for this:
Vault Root
 ├── Entries (flat storage)
 │     └── Entry nodes
 │
 ├── Folders (folder nodes)
 |
 ├── Attachments (file refs)
 │
 └── Index (query layer)
       ├── byType
       ├── byId
       └── (future indexes)


{
  "type": "vault",
  "version": "1.0.0",

  "vault_meta": {
    "name": "",
    "created_at": "...",
    "updated_at": "..."
  },

  "folders": { "/": "cid_folders" },
  "entries": { "/": "cid_entries_root" },
  "attachments": { "/": "cid_attachments_root" },
  "index": { "/": "cid_index_root" },

  "byType": {
    "login": [
      {
        "cid": "bafy...",
        "preview": {
          "entry_name": "GitHub",
          "user_name": "dev@..."
          "folder_id": "",
          "created_at": "",
          "updated_at": "",
        }
      }
    ]
  },

  "byFolder": {
    "61734720405": [...]
  }
}
=============================
entries_root node
{
  "items": [
    { "/": "cid_entry_1" },
    { "/": "cid_entry_2" }
  ]
}
folder_root node
{
  "items": [
    { "/": "cid_entry_1" },
    { "/": "cid_entry_2" }
  ]
}

folder_item
{
  "id": 61734720405,
  "name": "Stone",
  "created_at": "2025-08-28T11:27:38-07:00",
  "updated_at": "2025-08-28T11:27:38-07:00",
  "is_draft": false,
  "cid": "Bfm...
}

type Vault struct {
  // Legacy
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"column:name"`
	Type      string `json:"type" gorm:"column:type"`
	UserID    string `json:"user_id" gorm:"column:user_id"`
	UserSubscriptionID string `json:"user_subscription_id" gorm:"column:user_subscription_id"`
	CID       string `json:"cid" gorm:"column:cid"` // ✅ Explicitly map this!
	TxHash    string `json:"tx_hash" gorm:"column:tx_hash,omitempty"`
	CreatedAt string `json:"created_at" gorm:"column:created_at"` // change to time.Time later
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at"` // change to time.Time later

  // beta 
  version string `json:"version" gorm:"column:version"`
  VaultMeta VaultMeta `json:"vault_meta" gorm:"column:vault_meta"`
  Folders Link
  Entries Link
  Attachements Link
  IndexCID string `json:"index_cid"`

}
VaultMeta {
	Name      string `json:"name" gorm:"column:name"`
	UserID    string `json:"user_id" gorm:"column:user_id"`
  CreatedAt string `json:"created_at" gorm:"column:created_at"` // change to time.Time later
	UpdatedAt string `json:"updated_at" gorm:"column:updated_at"` // change to time.Time later
}

type VaultModel struct {
	ID  string
  Meta  VaultMeta
	CID string
}
type VaultNode struct {
	Type     string
	Version  string
	Folders  Link
	Entries  Link
	Index    Link
}

type Link struct {
	CID string `json:"/"`
}
type EntriesRoot struct {
	Items []Link `json:"items"`
}
type FoldersRoot struct {
	Items []Link `json:"items"`
}
type AttachementsRoot {
  Items []Link `json:"items"`
}
type Index struct {
	ByType   map[string][]Link `json:"byType"`
	ByFolder map[string][]Link `json:"byFolder"`
}

type Folder struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"varchar(100)"`
	CreatedAt string `json:"created_at" gorm:"varchar(100)"`
	UpdatedAt string `json:"updated_at" gorm:"varchar(100)"`
	IsDraft   bool   `json:"is_draft"`
	NodeCID  string `json:"cid"`    <= new 
}
// Entry = BaseEntry + LoginEntry/CardEntry....
type BaseEntry struct {
	ID              string    `json:"id"`
	EntryName       string    `json:"entry_name"`
	FolderID        string    `json:"folder_id"`
	Type            EntryType `json:"type"`
	AdditionnalNote string    `json:"additionnal_note,omitempty"`
	CustomFields    JSONMap   `json:"custom_fields,omitempty" gorm:"type:jsonb"`
	Trashed         bool      `json:"trashed"`
	IsDraft         bool      `json:"is_draft"`
	IsFavorite      bool      `json:"is_favorite"`
	CreatedAt string `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt string `json:"updated_at" gorm:"autoUpdateTime"`
	Attachments []Attachment `json:"attachments,omitempty" gorm:"foreignKey:EntryID"` // legacy

	AttachmentCIDs []string `json:"attachments,omitempty"` // <= new
	NodeCID  string `json:"cid"`    <= new 
}


        User Edit
            ↓
     EstimateCommit (DryRun)
            ↓
     Extract NewCIDs
            ↓
   QuotaService (ledger-aware)
            ↓
     Accept / Reject
            ↓
     CommitVault (real write)
            ↓
     BillingLedger update



     GA6C53Q6GNMOPJMJDBCMP7KXA3UWUJ652Z5O2H5MHLQLURZDTHTSXJLG
     SBI4J25C3F6JBOHJZ52FPP2JKDZRYSVHMRYVVUQG2DUHF46XMJQZZBNZ
     empty panel clarify section very pepper lumber birth virus pottery rally amount fame modify flat guess fox dentist either carpet pear flee position cattle




// 2. Load storage provider and tracecore client
	tracecoreClient := tracecore.NewTracecoreFromConfig(&appCfg, "token")

	storageProvider := blockchain.NewStorageProvider(
		blockchain.Config{
			StorageConfig: appCfg.Storage,
			UserID:        ur.UserSubscriptionID, // or userID
			VaultName:     ur.VaultName,
		},
		tracecoreClient,
	)
	vh.CreateIPFSPayloadCommandHandler.SetIpfsService(storageProvider)

	// 3. Store ENCRYPTED bytes on IPFS
	// (you can reuse the same StoreOnIpfs path)
	newCID, err := vh.CreateIPFSPayloadCommandHandler.StoreOnIpfs(
		context.Background(),
		vault_commands.StoreIpfsParams{
			Data: encrypted,
		},
	)






type VaultContext struct {
	SessionID     string
	AppConfig     app_config_domain.AppConfig
	StorageConfig app_config_domain.StorageConfig
	UserID        string
	VaultName     string
}
vaultCtx := VaultContext{
	SessionID: cmd.SessionID,
	StorageConfig: appCfg.Storage,
	UserID:        ur.UserSubscriptionID, // or userID
	VaultName:     ur.VaultName,
}


type CreateIPFSPayloadCommand struct {
	Vault    *vaults_domain.Vault
	Password string
	Data     []byte
}

type CreateIPFSPayloadCommandHandler struct {
	storageFactory StorageFactory
}
func (h *CreateIPFSPayloadCommandHandler) Execute(
	ctx context.Context,
	vaultCtx VaultContext,
	cmd CreateIPFSPayloadCommand,
) (*CreateIPFSPayloadCommandResult, error) {

	unlockRes, err := h.UnlockVaultHandler.Execute(vault_dto.UnlockVaultCommand{
		Password: cmd.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to unlock vault key: %w", err)
	}

	vaultKey := unlockRes.VaultKey.Key

	encrypted, err := h.CryptoService.Encrypt(cmd.Data, vaultKey)
	if err != nil {
		return nil, fmt.Errorf("vault encryption failed: %w", err)
	}

	cidFromIpfs, err := h.StoreOnIpfs(ctx, &vaultCtx, encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to add vault to IPFS: %w", err)
	}

	return &CreateIPFSPayloadCommandResult{CID: cidFromIpfs}, nil
}

func (h *CreateIPFSPayloadCommandHandler) StoreOnIpfs(
	ctx context.Context,
	vaultCtx VaultContext,
	data []byte,
) (string, error) {
	storageProvider := h.storageFactory.New(vaultCtx)
	return storageProvider.Add(ctx, data)
}


type StorageFactory struct {}

func (f *StorageFactory) New(vaultCtx *VaultContext) blockchain.StorageProvider {
	return blockchain.NewStorageProvider(
		blockchain.Config{
			StorageConfig: vaultCtx.StorageConfig,
			UserID:        vaultCtx.UserID,
			VaultName:     vaultCtx.VaultName,
		},
		tracecore.NewTracecoreFromConfig(&vaultCtx.AppConfig, "token"),
	)
}
newCID, err := vh.CreateIPFSPayloadCommandHandler.Execute(
	context.Background(),
	vaultCtx,
	vault_commands.CreateIPFSPayloadCommand{
		Vault:    vault.Vault,
		Password: "password",
		Data:     ur.ByteData,
	},
)	
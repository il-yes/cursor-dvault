# D-Vault — Self-Sovereign Digital Vault

A decentralized digital vault built on **IPFS** and **Stellar** for secure, verifiable, and self-sovereign data storage.

## 🏗️ Architecture Overview

**Frontend**: React + TypeScript + Tailwind CSS  
**Backend**: Golang (handles all IPFS, Stellar, Tracecore logic)  
**Storage**: IPFS (content-addressed, distributed)  
**Blockchain**: Stellar (immutable proof anchoring)  
**Verification**: Tracecore (commit validation and traceability)

## 🚀 Quick Start

```bash
npm install
npm run dev
```

## 📋 Next Steps for Backend Integration

The frontend is currently using **mock data** for preview/development. To connect to your Golang backend:

### 1. Configure API Base URL

Update the environment variable in `.env`:

```env
VITE_API_URL=http://localhost:8080
```

### 2. Backend Endpoints to Implement

The frontend expects the following REST API endpoints:

#### Vault Entries
- `POST /api/vault/entries` — Create new vault entry
- `GET /api/vault/entries` — List all vault entries
- `GET /api/vault/entries/:id` — Get specific entry
- `PUT /api/vault/entries/:id` — Update entry
- `DELETE /api/vault/entries/:id` — Delete entry
- `POST /api/vault/entries/:id/share` — Share entry (generate access token)

#### IPFS Operations
- `POST /api/ipfs/upload` — Upload content to IPFS
- `GET /api/ipfs/:cid` — Retrieve content from IPFS by CID

#### Stellar Blockchain
- `POST /api/stellar/anchor` — Anchor hash to Stellar blockchain
- `GET /api/stellar/verify/:tx` — Verify Stellar transaction

#### Tracecore Verification
- `POST /api/tracecore/commit` — Create Tracecore commit
- `GET /api/tracecore/verify/:id` — Verify commit integrity

### 3. Update API Service

Replace mock implementations in `src/services/api.ts` with actual `fetch()` calls. All placeholder comments are marked with `// TODO: Replace with actual API call`.

### 4. Expected Data Models

#### VaultEntry
```typescript
{
  id: string;
  title: string;
  content: string;
  category: string;
  ipfsHash: string;
  stellarTxHash: string;
  tracecoreCommitId: string;
  createdAt: string; // ISO 8601
  updatedAt: string; // ISO 8601
  encrypted: boolean;
}
```

## 🎨 Design System

- **Background**: White `#ffffff` / Light Gray `#f9fafb`
- **Text**: Dark Slate `#111827`
- **Primary**: Teal `#00cfcf`
- **Accent**: Orange `#FD871F`
- **Shadows**: Ultra-soft 8–16px blur
- **Typography**: Inter font, weights 300–600

## 🔒 Security Principles

- **Zero-trust architecture** — encryption on device
- **Self-sovereign identity** — users control their keys
- **Decentralized storage** — no central servers
- **Blockchain verification** — immutable proof layer

## 📦 Tech Stack

- **React 18** with TypeScript
- **Vite** for build tooling
- **Tailwind CSS** for styling
- **shadcn/ui** for UI components
- **React Router** for navigation
- **Tanstack Query** for data fetching

## 📖 User Flow

1. **Homepage** — Hero section with CTA to start vault
2. **Dashboard** — Manage vault entries (create, edit, delete, share)
3. **Entry Details** — View IPFS hash, Stellar transaction, Tracecore commit ID

## 🤝 Contributing

This is a frontend scaffold ready for integration with the Golang backend. All IPFS, Stellar, and Tracecore operations should be handled server-side.

---

**D-Vault © 2025 — Built for the Self-Sovereign Web.**

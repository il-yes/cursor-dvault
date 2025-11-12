# Wails Integration Notes

## Overview
This application supports both browser (REST API) and desktop (Wails) environments. The Wails bridge automatically detects the runtime and uses the appropriate communication method.

## Backend Response Contract

When a vault is created or updated, the backend (Go/Wails) returns a `VaultContext` payload:

```json
{
  "user_id": "uuid",
  "role": "owner",
  "Vault": {
    "version": "1.0",
    "name": "My Vault",
    "folders": [...],
    "entries": {
      "login": [...],
      "card": [...],
      "note": [...],
      "sshkey": [...],
      "identity": [...]
    }
  },
  "LastCID": "QmXxx...",
  "Dirty": false,
  "LastSynced": "2025-01-01T12:00:00Z",
  "LastUpdated": "2025-01-01T12:00:00Z",
  "vault_runtime_context": {
    "CurrentUser": {
      "id": "uuid",
      "role": "owner",
      "stellar_account": {
        "public_key": "G...",
        "private_key": "S..." // Only in dev mode
      }
    },
    "AppSettings": {
      "encryption_policy": "AES-256-GCM",
      "blockchain": {
        "stellar": { ... },
        "ipfs": { ... }
      },
      "auto_sync_enabled": true,
      "commit_rules": { ... }
    },
    "WorkingBranch": "main",
    "LoadedEntries": [...]
  }
}
```

## Wails Integration Methods

### Method 1: Direct Response (Recommended)
Return the `VaultContext` directly in the HTTP response body from `/api/vault/create`:

```go
func CreateVault(ctx context.Context, payload VaultPayload) (*VaultResponse, error) {
    // Create vault...
    vaultContext := BuildVaultContext(vault)
    
    return &VaultResponse{
        Success: true,
        VaultContext: vaultContext,
    }, nil
}
```

Frontend will automatically hydrate from response:
```typescript
const response = await createVault(payload);
if (response.vaultContext) {
    hydrateVault(response.vaultContext);
}
```

### Method 2: Wails Event Bridge
Send vault context via Wails events:

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

func CreateVault(a *App, payload VaultPayload) error {
    // Create vault...
    vaultContext := BuildVaultContext(vault)
    
    // Emit to frontend
    runtime.EventsEmit(a.ctx, "vault:context:updated", vaultContext)
    return nil
}
```

Frontend automatically listens:
```typescript
// Already configured in src/services/wailsBridge.ts
wailsBridge.onContextReceived((context) => {
    hydrateVault(context);
});
```

### Method 3: REST Fallback
Frontend can request context explicitly:

```go
func GetVaultContext(userId string) (*VaultContext, error) {
    // Load and return vault context
}
```

Frontend usage:
```typescript
const context = await wailsBridge.requestContext();
if (context) {
    hydrateVault(context);
}
```

## Security Requirements

### Sensitive Field Handling
1. **Storage**: All sensitive fields MUST be encrypted at rest
2. **Decryption**: Only decrypt on explicit user request via `/api/entry/decrypt`
3. **Transport**: Use TLS for REST, secure IPC for Wails
4. **Timeout**: Frontend auto-masks after 15 seconds (configurable)
5. **Audit**: Log all decrypt events to `/api/audit/log`

### Decryption Endpoint
```go
func DecryptField(entryId, fieldName, challenge string) (*DecryptResponse, error) {
    // Verify session/challenge
    // Decrypt field server-side only
    // Return plaintext with expiry
    
    return &DecryptResponse{
        FieldName: fieldName,
        Plaintext: decrypted,
        ExpiresIn: 15, // seconds
    }, nil
}
```

## Offline Support

### Local Persistence
Frontend uses IndexedDB via `src/services/localVault.ts`:
- Saves vault drafts when offline
- Retries sync when connection restored
- Exports encrypted backups

### Backend Integration
Implement sync endpoint:
```go
func SyncOfflineDraft(draft VaultDraft) (*VaultContext, error) {
    // Validate draft
    // Create/update vault
    // Return full VaultContext
}
```

## Testing Wails Integration

Create a test page to simulate Wails responses:

```typescript
// Test page component
function WailsTest() {
  const { hydrateVault } = useVault();
  
  const testPayload: VaultContext = {
    // ... full VaultContext
  };
  
  return (
    <Button onClick={() => hydrateVault(testPayload)}>
      Simulate Wails Response
    </Button>
  );
}
```

## Environment Detection
```typescript
// Automatically detected in src/services/wailsBridge.ts
const isWails = typeof window !== 'undefined' && !!window.backend;
```

## API Endpoints Required

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/vault/create` | POST | Create vault, return VaultContext |
| `/api/vault/context` | GET | Fetch VaultContext (fallback) |
| `/api/stellar/setup` | POST | Generate Stellar keypair |
| `/api/entry/decrypt` | POST | Decrypt sensitive field |
| `/api/audit/log` | POST | Log audit event |
| `/api/payment/upgrade` | POST | Upgrade plan |

## Next Steps

1. Implement backend endpoints listed above
2. Return `VaultContext` in `/api/vault/create` response
3. Implement decrypt endpoint with timeout
4. Add audit logging to Tracecore/Stellar
5. Test with both browser and Wails runtime

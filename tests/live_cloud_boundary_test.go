package tests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"vault-app/internal/blockchain"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	workspace_domain "vault-app/internal/workspace/domain"
)

const (
	CloudBaseURL  = "http://localhost:4001/api"
	CloudHostURL  = "http://localhost:4001"
	StellarPubKey = "GDN5XIDYDONFVV2K5VVR7W34CV43Z7YFDS34CU227MYM46MWBRNW4H4G"
	StellarSecKey = "SB477IOJ3VDUYWLYCSMKP5PMVFVVLW22DFLK5B4KGWCZXSSHMEDMTGI7"
	VaultID       = "fbffb9ad-d4de-4581-b7af-7b8e160f63bb"
	MySQLDSN      = "user:userpassword@tcp(127.0.0.1:3306)/widgets"
)

func TestLiveCloudBoundary_FullVerticalSlice(t *testing.T) {
	ctx := context.Background()

	// Connect to live Cloud MySQL database to verify persistence effects
	db, err := sql.Open("mysql", MySQLDSN)
	if err != nil {
		t.Skipf("Skipping live Cloud boundary test: cannot connect to MySQL (%v)", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping live Cloud boundary test: MySQL ping failed (%v)", err)
		return
	}

	client := tracecore.NewTracecoreClient(CloudBaseURL, "", "", CloudHostURL)

	// Step 1: Stellar Challenge Request
	t.Log("▶ Step 1: Requesting Stellar Challenge...")
	challengeStr, err := client.RequestStellarChallenge(ctx, StellarPubKey)
	require.NoError(t, err, "Step 1: RequestStellarChallenge failed")
	require.NotEmpty(t, challengeStr, "Step 1: Challenge string must not be empty")
	t.Logf("  ✓ Stellar Challenge received: %s", challengeStr)

	// Step 2: Sign Challenge with real Stellar key & Authenticate
	t.Log("▶ Step 2: Signing challenge with Stellar private key & Authenticating...")
	sig, err := blockchain.SignActorWithStellarPrivateKey(StellarSecKey, challengeStr)
	require.NoError(t, err, "Step 2: SignActorWithStellarPrivateKey failed")

	authResp, err := client.StellarAuthenticate(ctx, tracecore_types.StellarAuthenticateRequest{
		PublicKey: StellarPubKey,
		Signature: sig,
	})
	require.NoError(t, err, "Step 2: StellarAuthenticate failed")
	require.NotEmpty(t, authResp.Token, "Step 2: Cloud bearer token must not be empty")
	t.Logf("  ✓ Cloud Bearer Token issued: %s", authResp.Token)

	// Step 3: Token Hydration
	client.SetToken(authResp.Token)
	t.Log("  ✓ Client token hydrated")

	// Step 4: Vault Delegation & Registration
	t.Log("▶ Step 4: Requesting Vault Challenge & Registering Delegation...")
	vaultChallenge, err := client.RequestVaultChallenge(ctx, VaultID)
	require.NoError(t, err, "Step 4: RequestVaultChallenge failed")
	require.NotEmpty(t, vaultChallenge.Data.ChallengeID, "Step 4: Vault challenge_id must not be empty")

	vaultSig, err := blockchain.SignActorWithStellarPrivateKey(StellarSecKey, vaultChallenge.Data.SigningPayload)
	require.NoError(t, err, "Step 4: Sign vault challenge payload failed")

	regReq := tracecore.VaultRegisterRequest{
		ChallengeID:    vaultChallenge.Data.ChallengeID,
		Signature:      vaultSig,
		VaultID:        VaultID,
		OrganizationID: "personal",
		SigningKey:     StellarPubKey,
		EncryptionKey:  StellarPubKey,
		Endpoint:       "local",
		Capabilities:   []string{"storage"},
		VaultAddress:   "local",
	}

	_, regErr := client.RegisterVaultIdentity(ctx, regReq)
	if regErr != nil && !errors.Is(regErr, tracecore.ErrDelegationAlreadyExists) {
		require.NoError(t, regErr, "Step 4: RegisterVaultIdentity failed")
	}
	t.Log("  ✓ Vault Delegation verified active in Cloud")

	// Verify persistence effect in MySQL user_vault_identities table
	var delegationCount int
	err = db.QueryRowContext(ctx, "SELECT count(*) FROM user_vault_identities WHERE vault_id = ? AND revoked_at IS NULL", VaultID).Scan(&delegationCount)
	require.NoError(t, err)
	require.Greater(t, delegationCount, 0, "MySQL must contain an active delegation row for VaultID")
	t.Logf("  ✓ MySQL Persistence Verified: %d active delegation row(s) found in user_vault_identities", delegationCount)

	// Step 5: Create Workspace
	t.Log("▶ Step 5: Creating Workspace against running Cloud...")
	wsName := fmt.Sprintf("Live Boundary Workspace %d", time.Now().UnixNano())
	wsDesc := "Created via live boundary acceptance test"

	createReq := workspace_domain.CreateRequest{
		VaultID: VaultID,
		Workspace: workspace_domain.Workspace{
			VaultID:     VaultID,
			Name:        wsName,
			Description: wsDesc,
			OwnerID:     VaultID,
		},
	}

	createdWs, err := client.CreateWorkspace(ctx, createReq)
	require.NoError(t, err, "Step 5: CreateWorkspace failed HTTP boundary call")
	require.Equal(t, 201, createdWs.Status, "Step 5: HTTP status must be 201 Created")
	require.NotEmpty(t, createdWs.Data.ID, "Step 5: Workspace ID must not be empty")
	t.Logf("  ✓ Workspace Created (ID: %s, Name: %s)", createdWs.Data.ID, createdWs.Data.Name)

	// Verify persistence effect in MySQL workspaces table
	var dbName string
	err = db.QueryRowContext(ctx, "SELECT name FROM workspaces WHERE id = ?", createdWs.Data.ID).Scan(&dbName)
	require.NoError(t, err, "Step 5: Query workspace row from MySQL failed")
	require.Equal(t, wsName, dbName, "Step 5: MySQL row name must match created workspace name")
	t.Logf("  ✓ MySQL Persistence Verified: Workspace ID %s stored in MySQL with name '%s'", createdWs.Data.ID, dbName)

	// Step 6: Read / List Workspaces
	t.Log("▶ Step 6: Querying Workspaces list for Vault...")
	workspaces, err := client.ListWorkspaces(ctx, VaultID)
	require.NoError(t, err, "Step 6: ListWorkspaces failed")

	found := false
	for _, ws := range workspaces {
		if ws.ID == createdWs.Data.ID {
			found = true
			break
		}
	}
	require.True(t, found, "Step 6: Created workspace must appear in ListWorkspaces response")
	t.Logf("  ✓ Read Workspace Verified: Found workspace ID %s in ListWorkspaces result (%d total workspaces)", createdWs.Data.ID, len(workspaces))

	t.Log("\n🎉 PASS — real desktop client → Docker Cloud → MySQL boundary verified!")
}

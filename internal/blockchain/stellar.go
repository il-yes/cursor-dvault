package blockchain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	utils "vault-app/internal"
	"vault-app/internal/tracecore"

	// "time"

	"crypto/sha256"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

// In-memory map for demo (you'd want a DB or Redis in production)
var ChallengeStore = map[string]string{}

type ChallengeRequest struct {
	PublicKey string `json:"public_key"`
}

type ChallengeResponse struct {
	Challenge string `json:"challenge"`
	ExpiresAt string `json:"expires_at"`
}

type SignatureVerification struct {
	PublicKey string `json:"public_key"`
	Challenge string `json:"challenge"`
	Signature string `json:"signature"`
}

func GenerateChallenge(publicKey string) string {
	challenge := fmt.Sprintf("Login to dvault at %s", time.Now().UTC().Format(time.RFC3339))
	ChallengeStore[publicKey] = challenge
	return challenge
}

func VerifySignature(publicKey, challenge, signatureB64 string) bool {
	kp, err := keypair.Parse(publicKey)
	if err != nil {
		fmt.Println("Invalid public key:", err)
		return false
	}

	fullKP, ok := kp.(*keypair.FromAddress)
	if !ok {
		fmt.Println("Invalid keypair type")
		return false
	}

	sig, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		fmt.Println("Invalid signature encoding:", err)
		return false
	}

	return fullKP.Verify([]byte(challenge), sig) == nil
}

// SubmitCID anchors an IPFS CID on the Stellar Testnet using ManageData.
func SubmitCID(secretKey, ipfsCID string) (string, error) {
	kp, err := keypair.ParseFull(secretKey)
	if err != nil {
		return "", fmt.Errorf("invalid secret key: %w", err)
	}

	client := horizonclient.DefaultTestNetClient

	acctReq := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(acctReq)
	if err != nil {
		return "", fmt.Errorf("failed to load account: %w", err)
	}

	op := &txnbuild.ManageData{
		Name:  "vault_cid_+",
		Value: []byte(ipfsCID),
	}

	// 🔁 Corrected: set minTime = 0, maxTime = now + 5 min
	maxTime := time.Now().Add(5 * time.Minute).Unix()

	txParams := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		BaseFee:              txnbuild.MinBaseFee,
		Operations:           []txnbuild.Operation{op},
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewTimebounds(0, maxTime),
		},
	}

	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return "", fmt.Errorf("failed to build transaction: %w", err)
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	resp, err := client.SubmitTransaction(tx)
	if err != nil {
		return "", fmt.Errorf("transaction submission failed: %w", err)
	}

	return resp.Hash, nil
}

// CreateStellarAccount generates a new keypair and funds it using the testnet friendbot.
func CreateStellarAccount() (publicKey string, secretKey string, transactionID string, err error) {
	kp, err := keypair.Random()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate keypair: %w", err)
	}

	// Use Stellar Testnet friendbot to fund the account
	friendbotURL := fmt.Sprintf("https://friendbot.stellar.org/?addr=%s", kp.Address())

	resp, err := http.Get(friendbotURL)
	if err != nil {
		return "", "", "", fmt.Errorf("friendbot funding failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", "", "", fmt.Errorf("friendbot returned error: %s", string(body))
	}

	// Optional: Parse friendbot response for debug
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", "", fmt.Errorf("failed to parse friendbot response: %w", err)
	}

	// Extract transaction hash (aka transaction ID)
	txHash, ok := result["hash"].(string)
	if !ok {
		txHash = "unknown"
	}
	fmt.Println("Public key: ", kp)
	fmt.Println("Secret key: ", kp)
	return kp.Address(), kp.Seed(), txHash, nil

}

func VerifyStellarSignature(message, pubKey, signatureBase64 string) bool {
	// Decode the Stellar keypair
	kp, err := keypair.Parse(pubKey)
	if err != nil {
		fmt.Println("Invalid public key:", err)
		return false
	}

	full, ok := kp.(*keypair.FromAddress)
	if !ok {
		fmt.Println("Key is not a public address")
		return false
	}

	// Decode the signature from base64
	sigBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		fmt.Println("Failed to decode signature:", err)
		return false
	}

	// Hash the message
	hash := sha256.Sum256([]byte(message))

	// Verify the signature
	err = full.Verify(hash[:], sigBytes)
	if err != nil {
		fmt.Println("Signature invalid:", err)
		return false
	}

	return true
}

func SignActorWithStellarPrivateKey(privateKey string, message string) (string, error) {
	utils.LogPretty("private key", privateKey)
	kp, err := keypair.ParseFull(privateKey)
	if err != nil {
		return "", fmt.Errorf("❌ failed to retrieve keypair from this private key %s: %w", privateKey, err)
	}
	sig, err := kp.Sign([]byte(message))
	if err != nil {
		return "", fmt.Errorf("❌ failed to sign the message: %s - %w", message, err)
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func SignWithDvaultPrivateKey(cp tracecore.CommitPayload) (string, error) {
	data, err := json.Marshal(cp)
	if err != nil {
		return "", fmt.Errorf("marshal commit payload failed: %w", err)
	}

	kp, err := keypair.ParseFull(os.Getenv("TRACECORE_SECRETKEY"))
	if err != nil {
		return "", fmt.Errorf("parse key failed: %w", err)
	}
	sig, err := kp.Sign(data)
	if err != nil {
		return "", fmt.Errorf("signing failed: %w", err)
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}


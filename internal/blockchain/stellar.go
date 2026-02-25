package blockchain

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"vault-app/internal/logger/logger"
	tracecore_models "vault-app/internal/tracecore/models"
	utils "vault-app/internal/utils"

	// "time"

	"crypto/sha256"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

type StellarBlockchainService interface {
	CreateKeypair() (string, string, string, error)
	VerifyStellarSignature(message, pubKey, signatureBase64 string) bool
	SignActorWithStellarPrivateKey(privateKey string, message string) (string, error)
	SignWithDvaultPrivateKey(cp tracecore_models.CommitPayload) (string, error)
	CreateAccount(plainPassword string) (*CreateAccountRes, error)
	CreateAccountWithFriendbotFunding(plainPassword string) (*CreateAccountRes, error)
}

// --------------------- Stellar Service ---------------------
type StellarService struct {
	Logger *logger.Logger
	HTTPClient *http.Client
}

func NewStellarService(logger *logger.Logger) *StellarService {
	return &StellarService{
		Logger: logger,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type CreateAccountRes struct {
	PublicKey   string `json:"public_key"`
	PrivateKey  string `json:"private_key"`
	Salt        []byte `json:"salt"`
	EncNonce    []byte `json:"enc_nonce"`
	EncPassword []byte `json:"enc_password"`
	TxID        string `json:"tx_id"`
}

func (s *StellarService) CreateKeypair() (publicKey string, secretKey string, transactionID string, err error) {
	kp, err := keypair.Random()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate keypair: %w", err)
	}
	return kp.Address(), kp.Seed(), "", nil
}

// CreateAccount creates a new Stellar account with funding from friendbot
func (s *StellarService) CreatAccountWithFriendbotFunding(plainPassword string, ANCHORA_SECRET string) (*CreateAccountRes, error) {
	if os.Getenv("ANCHORA_SECRET") == "" {
		s.Logger.Error("❌ StellarService: CreatAccountWithFriendbotFunding - Anchora secret not found")
		return nil, errors.New("anchora secret not found")
	}

	pub, secret, txID, err := s.CreateKeypair()
	if err != nil {
		s.Logger.Warn("⚠️ Stellar account creation failed: %v", err)
		return nil, err
	}
	s.Logger.Info("✅ CreatAccountWithFriendbotFunding: Stellar account created: %s -  tx:", pub, txID)

	// Fund the account with 10 000 XLM from friendbot: 
	// tls: failed to verify certificate: x509: certificate signed by unknown authority
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://friendbot.stellar.org/?addr="+pub, nil)
	if err != nil {
		s.Logger.Error("❌ StellarService: CreatAccountWithFriendbotFunding - Failed to create request: %v", err)
		return nil, err
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		s.Logger.Error("❌ StellarService: CreatAccountWithFriendbotFunding - Failed to send request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Cryptography - Encrypt the user password with Stellar secret
	salt, nonce, encPassword, err := EncryptPasswordWithStellarSecure(plainPassword, secret)
	if err != nil {
		s.Logger.Error("❌ StellarService: CreatAccountWithFriendbotFunding - Failed to encrypt password with Stellar secret: %v", err)
		return nil, err
	}
	s.Logger.Info("✅ Stellar account created: %s -  tx:", pub, txID)

	// Cryptography - encrypt the Stellar private key before storing (server-side master key or KMS)
	encryptedPrivateKey, err := Encrypt([]byte(secret), ANCHORA_SECRET)
	if err != nil {
		s.Logger.Error("❌ StellarService: CreatAccountWithFriendbotFunding - Failed to encrypt Stellar private key: %v", err)
		return nil, err
	}

	return &CreateAccountRes{
		PublicKey:   pub,
		PrivateKey:  base64.StdEncoding.EncodeToString(encryptedPrivateKey),
		Salt:        salt,
		EncNonce:    nonce,
		EncPassword: encPassword,
		TxID:        txID,
	}, nil
}

// CreateAccount creates a new Stellar account with no funding
func (s *StellarService) CreateAccount(plainPassword string) (*CreateAccountRes, error) {
	pub, secret, txID, err := s.CreateKeypair()
	if err != nil {
		s.Logger.Warn("⚠️ Stellar account creation failed: %v", err)
		return nil, err
	}

	salt, nonce, encPassword, err := EncryptPasswordWithStellarSecure(plainPassword, secret)
	if err != nil {
		s.Logger.Error("❌ StellarService: CreateAccount - Failed to encrypt password with Stellar secret: %v", err)
		return nil, err
	}
	s.Logger.Info("✅ Stellar account created: %s -  tx:", pub, txID)

	encrypted, err := Encrypt([]byte(secret), plainPassword)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}
	s.Logger.Info("✅ Stellar account created: %s -  tx:", pub, txID, encrypted)

	return &CreateAccountRes{
		PublicKey:   pub,
		PrivateKey:  secret, // base64.StdEncoding.EncodeToString(encryptedPrivateKey),
		Salt:        salt,
		EncNonce:    nonce,
		EncPassword: encPassword,
	}, nil
}
func (s *StellarService) OnGenerateApiKey(password string) (*CreateAccountRes, error) {
	res, err := s.CreateAccount(password)
	if err != nil {
		s.Logger.Error("❌ GenerateApiKey: OnGenerateApiKey -Stellar account creation failed: %v", err)
		return nil, err
	}
	return res, nil
}

// --------------------- Challenge Service ---------------------
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
	utils.LogPretty("VerifySignature - publicKey", publicKey)
	kp, err := keypair.Parse(publicKey)
	if err != nil {
		fmt.Println("Invalid public key:", err)
		return false
	}
	utils.LogPretty("VerifySignature - kp", kp)
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
	utils.LogPretty("VerifySignature - sig", sig)
	return fullKP.Verify([]byte(challenge), sig) == nil
}

// SubmitCID anchors an IPFS CID on the Stellar Testnet using ManageData.
func SubmitCID(secretKey, ipfsCID string) (string, error) {
	utils.LogPretty("SubmitCID - secretKey", secretKey)
	kp, err := keypair.ParseFull(secretKey)
	if err != nil {
		return "", fmt.Errorf("invalid secret key: %w", err)
	}

	client := horizonclient.DefaultTestNetClient
	acctReq := horizonclient.AccountRequest{AccountID: kp.Address()}
	utils.LogPretty("SubmitCID - acctReq", acctReq)
	sourceAccount, err := client.AccountDetail(acctReq)
	if err != nil {
		return "", fmt.Errorf("failed to load account: %w", err)
	}

	op := &txnbuild.ManageData{
		Name:  "vault_cid_+",
		Value: []byte(ipfsCID),
	}
	utils.LogPretty("SubmitCID - op", op)

	// 🔁 Corrected: set minTime = 0, maxTime = now + 5 min
	maxTime := time.Now().Add(5 * time.Minute).Unix()
	utils.LogPretty("SubmitCID - maxTime", maxTime)

	txParams := txnbuild.TransactionParams{
		SourceAccount:        &sourceAccount,
		IncrementSequenceNum: true,
		BaseFee:              txnbuild.MinBaseFee,
		Operations:           []txnbuild.Operation{op},
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewTimebounds(0, maxTime),
		},
	}
	utils.LogPretty("SubmitCID - txParams", txParams)

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

func SignWithDvaultPrivateKey(cp tracecore_models.CommitPayload) (string, error) {
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

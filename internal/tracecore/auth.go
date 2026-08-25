package tracecore

import (
	"context"
	"fmt"
	"log"
	tracecore_types "vault-app/internal/tracecore/types"
)

func (c *TracecoreClient) Login(ctx context.Context, req tracecore_types.LoginRequest) (*tracecore_types.LoginResponse, error) {
	var resp tracecore_types.LoginResponse

	log.Printf("[CLOUD-AUTH] Authenticating email=%s", req.Email)

	err := c.doRequest(ctx, "POST", "/api/authenticate", req, &resp)
	if err != nil {
		log.Printf("[CLOUD-AUTH] Authentication failed for email=%s: %v", req.Email, err)
		return nil, err
	}

	// save auth token
	if resp.AuthenticationToken != nil {
		resp.Token = resp.AuthenticationToken.Token
	}

	return &resp, nil
}

func (c *TracecoreClient) RequestStellarChallenge(ctx context.Context, publicKey string) (string, error) {
	req := tracecore_types.StellarChallengeRequest{
		PublicKey: publicKey,
	}
	var resp tracecore_types.StellarChallengeResponse

	log.Printf("[CLOUD-AUTH] Requesting Stellar challenge for public_key=%s", publicKey)

	err := c.doRequest(ctx, "POST", "/api/stellar/public-challenge", req, &resp)
	if err != nil {
		log.Printf("[CLOUD-AUTH] Stellar challenge request failed for public_key=%s: %v", publicKey, err)
		return "", err
	}

	challenge := resp.Challenge
	if challenge == "" {
		challenge = resp.Data.Challenge
	}
	if challenge == "" {
		return "", fmt.Errorf("empty challenge returned from cloud for public_key=%s", publicKey)
	}

	return challenge, nil
}

func (c *TracecoreClient) StellarAuthenticate(ctx context.Context, req tracecore_types.StellarAuthenticateRequest) (*tracecore_types.LoginResponse, error) {
	var resp tracecore_types.LoginResponse

	log.Printf("[CLOUD-AUTH] Authenticating Stellar public_key=%s", req.PublicKey)

	err := c.doRequest(ctx, "POST", "/api/stellar/authenticate", req, &resp)
	if err != nil {
		log.Printf("[CLOUD-AUTH] Stellar authentication failed for public_key=%s: %v", req.PublicKey, err)
		return nil, err
	}

	if resp.AuthenticationToken != nil {
		resp.Token = resp.AuthenticationToken.Token
	}

	return &resp, nil
}

func (c *TracecoreClient) Logout(ctx context.Context) error {
	c.Token = ""
	return nil
}


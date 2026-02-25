package tracecore

import (
	"context"
	tracecore_types "vault-app/internal/tracecore/types"
)

func (c *TracecoreClient) Login(ctx context.Context, req tracecore_types.LoginRequest) (*tracecore_types.LoginResponse, error) {
	var resp tracecore_types.LoginResponse

	err := c.doRequest(ctx, "POST", "/authenticate", req, &resp)
	if err != nil {	
		return nil, err
	}

	// save auth token
	if resp.AuthenticationToken != nil {
		resp.Token = resp.AuthenticationToken.Token
	}

	return &resp, nil
}

func (c *TracecoreClient) Logout(ctx context.Context) error {
	c.Token = ""
	return nil
}

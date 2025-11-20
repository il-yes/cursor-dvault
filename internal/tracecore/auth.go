package tracecore

import "context"

func (c *TracecoreClient) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
    var resp LoginResponse
    err := c.doRequest(ctx, "POST", "/authenticate", req, &resp)
    if err != nil {
        return nil, err
    }
    c.Token = resp.Token
    return &resp, nil
}

func (c *TracecoreClient) Logout(ctx context.Context) error {
    c.Token = ""
    return nil
}

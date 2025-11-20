package tracecore

import "context"

func (c *TracecoreClient) ListReceivedShares(ctx context.Context) ([]ShareEntry, error) {
    var resp []ShareEntry
    err := c.doRequest(ctx, "GET", "/api/shares/received", nil, &resp)
    return resp, err
}

func (c *TracecoreClient) AcceptShare(ctx context.Context, shareID string) error {
    return c.doRequest(ctx, "POST", "/api/shares/"+shareID+"/accept", nil, nil)
}

func (c *TracecoreClient) RejectShare(ctx context.Context, shareID string) error {
    return c.doRequest(ctx, "POST", "/api/shares/"+shareID+"/reject", nil, nil)
}

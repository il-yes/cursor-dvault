package tracecore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

// ---------------------------------------------------------------------------
// TrustGroupRepository implementation on TracecoreClient (Cloud contract C1)
// ---------------------------------------------------------------------------

// compile-time proof that TracecoreClient satisfies the desktop
// TrustGroupRepository port.
var _ trustgroup_domain.TrustGroupRepository = (*TracecoreClient)(nil)

// GetTrustGroup fetches a trust group from Ankhora Cloud
// (GET /api/trustgroups/{id}). The Cloud response mirrors the desktop
// trust_group.TrustGroup wire contract, so it is decoded directly.
func (c *TracecoreClient) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	if req == nil || req.TrustGroupID == "" {
		return nil, fmt.Errorf("trust group id is required")
	}

	url := c.AnkhoraCloudUrl + "/trustgroups/" + req.TrustGroupID
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("failed to decode trust group response: %w", err)
	}
	if cloudResp.Data.ID == "" {
		return nil, fmt.Errorf("trust group not found: %s", req.TrustGroupID)
	}

	return &cloudResp, nil
}

// CreateTrustGroup has no Cloud endpoint yet. It fails explicitly rather
// than fabricating local state.
func (c *TracecoreClient) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, fmt.Errorf("CreateTrustGroup is not supported by Cloud yet")
}

// GetTrustGroupMember has no Cloud endpoint yet.
func (c *TracecoreClient) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, fmt.Errorf("GetTrustGroupMember is not supported by Cloud yet")
}

// ListTrustGroups has no Cloud endpoint yet.
func (c *TracecoreClient) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, fmt.Errorf("ListTrustGroups is not supported by Cloud yet")
}

// UpdateTrustGroup has no Cloud endpoint yet.
func (c *TracecoreClient) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, fmt.Errorf("UpdateTrustGroup is not supported by Cloud yet")
}

// DeleteTrustGroup has no Cloud endpoint yet.
func (c *TracecoreClient) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, fmt.Errorf("DeleteTrustGroup is not supported by Cloud yet")
}

// AddMemberToTrustGroup has no Cloud endpoint yet.
func (c *TracecoreClient) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, fmt.Errorf("AddMemberToTrustGroup is not supported by Cloud yet")
}

// RemoveMemberFromTrustGroup has no Cloud endpoint yet.
func (c *TracecoreClient) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, fmt.Errorf("RemoveMemberFromTrustGroup is not supported by Cloud yet")
}

// RotateTrustGroupKEK has no Cloud endpoint yet.
func (c *TracecoreClient) RotateTrustGroupKEK(ctx context.Context, req *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, fmt.Errorf("RotateTrustGroupKEK is not supported by Cloud yet")
}

// ---------------------------------------------------------------------------
// ShareEntryRepository implementation on TracecoreClient (Cloud contracts C2/C3)
// ---------------------------------------------------------------------------

var _ trustgroup_domain.TrustGroupRepository = (*TracecoreClient)(nil)

// ---------------------------------------------------------------------------
// C3 ShareEntry Cloud transport (contracts C2/C3)
// ---------------------------------------------------------------------------

// CreateShareEntryDirect persists a C3 share entry on Ankhora Cloud
// (POST /api/c3/share-entries). The Cloud validates the referenced trust
// group and its current KEK version, assigns the authoritative ID when the
// client did not provide one, and returns the persisted entry.
func (c *TracecoreClient) CreateShareEntryDirect(ctx context.Context, entry c3_asset_domain.ShareEntry) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	body, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	url := c.AnkhoraCloudUrl + "/c3/share-entries"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("failed to decode share entry response: %w", err)
	}
	if cloudResp.Data.ID == "" {
		return nil, fmt.Errorf("Cloud returned a share entry without an ID")
	}

	return &cloudResp, nil
}

// GetShareEntryDirect fetches a persisted C3 share entry from Ankhora Cloud
// (GET /api/c3/share-entries/{id}).
func (c *TracecoreClient) GetShareEntryDirect(ctx context.Context, shareEntryID string) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	url := c.AnkhoraCloudUrl + "/c3/share-entries/" + shareEntryID
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("failed to decode share entry response: %w", err)
	}
	if cloudResp.Data.ID == "" {
		return nil, fmt.Errorf("share entry not found: %s", shareEntryID)
	}

	return &cloudResp, nil
}

// CloudShareEntryRepository adapts TracecoreClient to the
// c3_asset_domain.ShareEntryRepository port. The adapter exists because
// TracecoreClient already carries a legacy GetShareEntry method for the
// vault cryptoshare context; the two domains must not collide.
type CloudShareEntryRepository struct {
	client *TracecoreClient
}

func NewCloudShareEntryRepository(client *TracecoreClient) *CloudShareEntryRepository {
	return &CloudShareEntryRepository{client: client}
}

var _ c3_asset_domain.ShareEntryRepository = (*CloudShareEntryRepository)(nil)

func (r *CloudShareEntryRepository) CreateShareEntry(ctx context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	if req == nil {
		return nil, fmt.Errorf("share entry request is required")
	}
	return r.client.CreateShareEntryDirect(ctx, req.ShareEntry)
}

func (r *CloudShareEntryRepository) GetShareEntry(ctx context.Context, req *c3_asset_domain.GetShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	if req == nil || req.ShareEntryID == "" {
		return nil, fmt.Errorf("share entry id is required")
	}
	return r.client.GetShareEntryDirect(ctx, req.ShareEntryID)
}

// UpdateShareEntry has no Cloud endpoint yet.
func (r *CloudShareEntryRepository) UpdateShareEntry(ctx context.Context, req *c3_asset_domain.UpdateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, fmt.Errorf("UpdateShareEntry is not supported by Cloud yet")
}

// DeleteShareEntry has no Cloud endpoint yet.
func (r *CloudShareEntryRepository) DeleteShareEntry(ctx context.Context, req *c3_asset_domain.DeleteShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, fmt.Errorf("DeleteShareEntry is not supported by Cloud yet")
}

package tracecore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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

func (c *TracecoreClient) getCloudBaseURL() string {
	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	endpoint := strings.TrimRight(baseUrl, "/")
	if !strings.HasSuffix(endpoint, "/api") {
		endpoint += "/api"
	}
	return endpoint
}

// GetTrustGroup fetches a trust group from Ankhora Cloud
// (GET /api/trustgroups/{id}). The Cloud response mirrors the desktop
// trust_group.TrustGroup wire contract, so it is decoded directly.
func (c *TracecoreClient) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	if req == nil || req.TrustGroupID == "" {
		return nil, fmt.Errorf("trust group id is required")
	}

	url := c.getCloudBaseURL() + "/trustgroups/" + req.TrustGroupID
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

	memberCount := len(cloudResp.Data.MemberCIDs)
	memberCIDsStr := strings.Join(cloudResp.Data.MemberCIDs, ", ")
	fmt.Printf("[C3][MEMBERSHIP][LOAD] trustGroupID=%s memberCount=%d MemberCIDs=[%s]\n", req.TrustGroupID, memberCount, memberCIDsStr)

	return &cloudResp, nil
}

// CreateTrustGroup creates a new trust group via POST /api/trustgroups.
func (c *TracecoreClient) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	if req == nil {
		return nil, fmt.Errorf("create trust group request is required")
	}

	workspaceID := req.TrustGroup.ChannelID
	if workspaceID == "" {
		workspaceID = "default_workspace"
	}

	payload := map[string]interface{}{
		"workspace_id": workspaceID,
		"name":         req.TrustGroup.Name,
		"members":      []interface{}{},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := c.getCloudBaseURL() + "/trustgroups"

	log.Printf("[TRUSTGROUP][CREATE] HTTP_REQUEST method=POST url=%s payload=%s", url, string(body))

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
		log.Printf("[TRUSTGROUP][CREATE] HTTP_ERROR error=%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[TRUSTGROUP][CREATE] HTTP_RESPONSE status=%d body=%s", resp.StatusCode, string(respBytes))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("failed to decode trust group response: %w", err)
	}

	return &cloudResp, nil
}

// GetTrustGroupMember fetches a specific member.
func (c *TracecoreClient) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	tgResp, err := c.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: req.TrustGroupID})
	if err != nil {
		return nil, err
	}
	for _, env := range tgResp.Data.KeyEnvelopes {
		if env.MemberID == req.MemberID || env.ID == req.MemberID {
			return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember]{
				Data: trustgroup_domain.TrustGroupMember{
					ID:       env.ID,
					VaultID:  env.MemberID,
					Role:     "member",
					JoinedAt: env.CreatedAt,
				},
			}, nil
		}
	}
	return nil, fmt.Errorf("member not found: %s", req.MemberID)
}

// ListTrustGroups fetches all trust groups from GET /api/trustgroups.
func (c *TracecoreClient) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=tgcrud-001][18] layer=CLOUD_CLIENT file=internal/tracecore/c3_cloud_repository.go function=TracecoreClient.ListTrustGroups input.ChannelID=%s\n", req.ChannelID)
	url := c.getCloudBaseURL() + "/trustgroups"
	if req != nil && req.ChannelID != "" {
		url += "?workspace_id=" + req.ChannelID
	}
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=tgcrud-001][19] layer=HTTP_REQUEST method=GET path=%s status=SENDING\n", url)
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

	var cloudResp tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("failed to decode trust groups list response: %w", err)
	}

	totalMembers := 0
	for _, tg := range cloudResp.Data {
		totalMembers += len(tg.MemberCIDs)
	}
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=tgcrud-001][23] layer=HTTP_RESPONSE status=%d trustGroupCount=%d totalMembers=%d payload=%s\n", resp.StatusCode, len(cloudResp.Data), totalMembers, string(respBytes))

	return &cloudResp, nil
}

// UpdateTrustGroup updates a trust group via PUT /api/trustgroups/{id}.
func (c *TracecoreClient) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	if req == nil || req.TrustGroup.ID == "" {
		return nil, fmt.Errorf("trust group id is required")
	}

	fmt.Printf("[C3][REAL-E2E][08] UpdateTrustGroup outgoing membersCount=%d envelopesCount=%d trustGroupID=%s\n", len(req.TrustGroup.MemberCIDs), len(req.TrustGroup.KeyEnvelopes), req.TrustGroup.ID)
	fmt.Printf("[C3][ADD_MEMBER][STEP_18] Calling TracecoreClient.UpdateTrustGroup trustGroupID=%s envelopesCount=%d\n", req.TrustGroup.ID, len(req.TrustGroup.KeyEnvelopes))
	fmt.Printf("[C3][ENVELOPE][HTTP][OUT] trustGroupID=%s envelopesCount=%d\n", req.TrustGroup.ID, len(req.TrustGroup.KeyEnvelopes))
	for i, env := range req.TrustGroup.KeyEnvelopes {
		fmt.Printf("  -> envelope[%d]: memberID=%s deviceID=%s kekVersion=%d wrappedKEKLen=%d\n", i, env.MemberID, env.DeviceID, env.KEKVersion, len(env.WrappedKEK))
	}

	body, err := json.Marshal(req.TrustGroup)
	if err != nil {
		return nil, err
	}

	url := c.getCloudBaseURL() + "/trustgroups/" + req.TrustGroup.ID
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
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

	fmt.Printf("[C3][REAL-E2E][09] Cloud Update received membersCount=%d envelopesCount=%d trustGroupID=%s\n", len(cloudResp.Data.MemberCIDs), len(cloudResp.Data.KeyEnvelopes), cloudResp.Data.ID)
	fmt.Printf("[C3][ADD_MEMBER][STEP_19] TracecoreClient.UpdateTrustGroup return success trustGroupID=%s envelopesCount=%d\n", cloudResp.Data.ID, len(cloudResp.Data.KeyEnvelopes))
	fmt.Printf("[C3][ENVELOPE][CLIENT][RECEIVED] trustGroupID=%s envelopesCount=%d\n", cloudResp.Data.ID, len(cloudResp.Data.KeyEnvelopes))
	for i, env := range cloudResp.Data.KeyEnvelopes {
		fmt.Printf("  -> envelope[%d]: memberID=%s deviceID=%s kekVersion=%d wrappedKEKLen=%d\n", i, env.MemberID, env.DeviceID, env.KEKVersion, len(env.WrappedKEK))
	}

	return &cloudResp, nil
}

// DeleteTrustGroup deletes a trust group via DELETE /api/trustgroups/{id}.
func (c *TracecoreClient) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	if req == nil || req.TrustGroupID == "" {
		return nil, fmt.Errorf("trust group id is required")
	}

	url := c.getCloudBaseURL() + "/trustgroups/" + req.TrustGroupID
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
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

	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Success: true}, nil
}

// AddMemberToTrustGroup adds a member via POST /api/trustgroups/{id}/members.
func (c *TracecoreClient) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	if req == nil || req.TrustGroupID == "" {
		return nil, fmt.Errorf("trust group id is required")
	}
	if req.VaultID == "" {
		return nil, fmt.Errorf("vault id is required")
	}
	if req.Role == "" {
		req.Role = "member"
	}

	fmt.Printf("[C3][ADD_MEMBER][STEP_03] TracecoreClient.AddMemberToTrustGroup enter trustGroupID=%s vaultID=%s role=%s\n", req.TrustGroupID, req.VaultID, req.Role)
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=tgcrud-001][05] layer=CLOUD_CLIENT file=internal/tracecore/c3_cloud_repository.go function=TracecoreClient.AddMemberToTrustGroup input.TrustGroupID=%s input.VaultID=%s input.Role=%s\n", req.TrustGroupID, req.VaultID, req.Role)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := c.getCloudBaseURL() + "/trustgroups/" + req.TrustGroupID + "/members"
	fmt.Printf("[C3][TRACE][ADD_MEMBER][trace=tgcrud-001][06] layer=HTTP_REQUEST method=POST path=%s payload=%s\n", url, string(body))
	log.Printf("[C3][MEMBERSHIP][HTTP_REQUEST] trustGroupID=%s vaultID=%s role=%s payload=%s", req.TrustGroupID, req.VaultID, req.Role, string(body))
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

	var cloudResp tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("failed to decode trust group response: %w", err)
	}

	memberCount := len(cloudResp.Data.MemberCIDs)
	memberCIDsStr := strings.Join(cloudResp.Data.MemberCIDs, ", ")
	fmt.Printf("[C3][ADD_MEMBER][STEP_04] TracecoreClient.AddMemberToTrustGroup return success trustGroupID=%s memberCount=%d MemberCIDs=[%s]\n", req.TrustGroupID, memberCount, memberCIDsStr)
	fmt.Printf("[C3][MEMBERSHIP][PERSIST] AddMemberToTrustGroup trustGroupID=%s memberCount=%d MemberCIDs=[%s]\n", req.TrustGroupID, memberCount, memberCIDsStr)

	return &cloudResp, nil
}

// RemoveMemberFromTrustGroup revokes/removes a member via DELETE /api/trustgroups/{id}/members/{vaultID}.
func (c *TracecoreClient) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	if req == nil || req.TrustGroupID == "" || req.MemberID == "" {
		return nil, fmt.Errorf("trust group id and member id are required")
	}

	url := c.getCloudBaseURL() + "/trustgroups/" + req.TrustGroupID + "/members/" + req.MemberID
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
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

	return &cloudResp, nil
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

	url := c.getCloudBaseURL() + "/c3/share-entries"
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
	url := c.getCloudBaseURL() + "/c3/share-entries/" + shareEntryID
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

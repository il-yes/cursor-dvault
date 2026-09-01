package tracecore

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	channel_domain "vault-app/internal/channel/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

func TestChannelUpdate_SerializationRoundTrip(t *testing.T) {
	// 1. Domain Channel with Slots, Assignments, and Properties
	original := channel_domain.Channel{
		ID:          "ch_test_123",
		WorkspaceID: "ws_456",
		Title:       "Contract Execution",
		Slots: []channel_domain.Slot{
			{
				ID:      "1",
				Name:    "contract_draft",
				Role:    "primary",
				VaultID: "vault_legal",
				Gated:   false,
				Order:   1,
			},
			{
				ID:      "2",
				Name:    "financial_clearance",
				Role:    "participant",
				VaultID: "vault_finance",
				Gated:   true,
				Order:   2,
			},
		},
		Assignments: []channel_domain.Assignment{
			{
				SlotID:       "1",
				OwnerID:      "vault_legal",
				PublicKey:    "pk_legal_01",
				VaultAddress: "vault_legal",
			},
			{
				SlotID:       "2",
				OwnerID:      "vault_finance",
				PublicKey:    "pk_finance_01",
				VaultAddress: "vault_finance",
			},
		},
		Properties: []channel_domain.ChannelProperty{
			{
				Key:   "contract_id",
				Value: "CTR-2026-009",
			},
			{
				Key:   "currency",
				Value: "USD",
			},
		},
	}

	// 2. Map to Cloud DTOs (Outbound Payload Construction)
	cloudSlots := make([]tracecore_types.CloudChannelSlot, 0, len(original.Slots))
	for _, s := range original.Slots {
		cloudSlots = append(cloudSlots, tracecore_types.CloudChannelSlot{
			ID:      s.ID,
			Name:    s.Name,
			Role:    s.Role,
			VaultID: s.VaultID,
			Gated:   s.Gated,
			Order:   s.Order,
		})
	}

	cloudAssignments := make([]tracecore_types.CloudChannelAssignment, 0, len(original.Assignments))
	for _, a := range original.Assignments {
		cloudAssignments = append(cloudAssignments, tracecore_types.CloudChannelAssignment{
			SlotID:       a.SlotID,
			OwnerID:      a.OwnerID,
			PublicKey:    a.PublicKey,
			VaultAddress: a.VaultAddress,
		})
	}

	cloudProperties := make([]tracecore_types.CloudChannelProperty, 0, len(original.Properties))
	for _, p := range original.Properties {
		cloudProperties = append(cloudProperties, tracecore_types.CloudChannelProperty{
			Key:   p.Key,
			Value: p.Value,
		})
	}

	cloudDTO := tracecore_types.CloudChannelDTO{
		ID:          original.ID,
		WorkspaceID: original.WorkspaceID,
		Title:       original.Title,
		Slots:       cloudSlots,
		Assignments: cloudAssignments,
		Properties:  cloudProperties,
	}

	// 3. Serialize to JSON (Cloud Wire HTTP Payload)
	jsonBytes, err := json.Marshal(cloudDTO)
	if err != nil {
		t.Fatalf("Failed to marshal CloudChannelDTO to JSON: %v", err)
	}

	// 4. Unmarshal from JSON (Cloud Response Body)
	var returnedDTO tracecore_types.CloudChannelDTO
	if err := json.Unmarshal(jsonBytes, &returnedDTO); err != nil {
		t.Fatalf("Failed to unmarshal JSON to CloudChannelDTO: %v", err)
	}

	// 5. Map back to Domain Channel aggregate (mapCloudChannelDTO)
	reconstructed := mapCloudChannelDTO(returnedDTO)

	// 6. Verification Assertions
	if len(reconstructed.Slots) != 2 {
		t.Fatalf("Expected 2 slots after round-trip, got %d", len(reconstructed.Slots))
	}
	if reconstructed.Slots[0].Name != "contract_draft" || reconstructed.Slots[0].VaultID != "vault_legal" {
		t.Errorf("Slot[0] corrupted after round-trip: %+v", reconstructed.Slots[0])
	}
	if reconstructed.Slots[1].Name != "financial_clearance" || reconstructed.Slots[1].VaultID != "vault_finance" {
		t.Errorf("Slot[1] corrupted after round-trip: %+v", reconstructed.Slots[1])
	}

	if len(reconstructed.Assignments) != 2 {
		t.Fatalf("Expected 2 assignments after round-trip, got %d", len(reconstructed.Assignments))
	}
	if reconstructed.Assignments[0].SlotID != "1" || reconstructed.Assignments[0].OwnerID != "vault_legal" {
		t.Errorf("Assignment[0] corrupted after round-trip: %+v", reconstructed.Assignments[0])
	}

	if len(reconstructed.Properties) != 2 {
		t.Fatalf("Expected 2 properties after round-trip, got %d", len(reconstructed.Properties))
	}
	if reconstructed.Properties[0].Key != "contract_id" || reconstructed.Properties[0].Value != "CTR-2026-009" {
		t.Errorf("Property[0] corrupted after round-trip: %+v", reconstructed.Properties[0])
	}

	t.Logf("✓ Round-trip test passed: Slots(%d), Assignments(%d), Properties(%d) fully preserved.",
		len(reconstructed.Slots), len(reconstructed.Assignments), len(reconstructed.Properties))
}

func TestOutboundChannelUpdate_SnakeCaseJSONContract(t *testing.T) {
	var capturedBody map[string]interface{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "expected PUT", http.StatusMethodNotAllowed)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&capturedBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": 200,
			"data": tracecore_types.CloudChannelDTO{
				ID:    "c9fdacfb-d3da-45a4-9c94-c18923dfd781",
				Title: "Marketing Channel",
				Slots: []tracecore_types.CloudChannelSlot{
					{
						ID:      "slot_1788190352422_7ekj",
						Name:    "Marketing_operation",
						Role:    "participant",
						VaultID: "vault_finance",
						Gated:   true,
						Order:   4,
					},
				},
			},
			"message": "success",
			"success": true,
		})
	}))
	defer ts.Close()

	tc := NewTracecoreClient(ts.URL, "token_123", ts.URL, ts.URL)
	req := &channel_domain.UpdateChannelRequest{
		Channel: channel_domain.Channel{
			ID:    "c9fdacfb-d3da-45a4-9c94-c18923dfd781",
			Title: "Marketing Channel",
			Slots: []channel_domain.Slot{
				{
					ID:      "slot_1788190352422_7ekj",
					Name:    "Marketing_operation",
					Role:    "participant",
					VaultID: "vault_finance",
					Gated:   true,
					Order:   4,
				},
			},
		},
	}

	resp, err := tc.UpdateChannel(context.Background(), req)
	if err != nil {
		t.Fatalf("UpdateChannel failed: %v", err)
	}
	if resp == nil || !resp.Success {
		t.Fatalf("UpdateChannel returned unsuccessful response: %+v", resp)
	}

	slotsRaw, ok := capturedBody["slots"].([]interface{})
	if !ok || len(slotsRaw) != 1 {
		t.Fatalf("Expected 1 slot in captured JSON body, got: %+v", capturedBody["slots"])
	}

	slotMap, ok := slotsRaw[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected slot object map in JSON body, got: %+v", slotsRaw[0])
	}

	// 1. Verify exact snake_case JSON keys exist
	expectedFields := map[string]interface{}{
		"id":       "slot_1788190352422_7ekj",
		"name":     "Marketing_operation",
		"role":     "participant",
		"vault_id": "vault_finance",
		"gated":    true,
		"order":    float64(4),
	}

	for k, expectedVal := range expectedFields {
		val, exists := slotMap[k]
		if !exists {
			t.Errorf("MISSING snake_case key '%s' in outbound HTTP PUT JSON payload", k)
		} else if val != expectedVal {
			t.Errorf("Mismatch for key '%s': expected %v (%T), got %v (%T)", k, expectedVal, expectedVal, val, val)
		}
	}

	// 2. Verify PascalCase JSON keys are ABSENT
	forbiddenKeys := []string{"ID", "Name", "Role", "VaultID", "Gated", "Order"}
	for _, k := range forbiddenKeys {
		if _, exists := slotMap[k]; exists {
			t.Errorf("FORBIDDEN PascalCase key '%s' STILL PRESENT in outbound HTTP PUT JSON payload!", k)
		}
	}
}

func TestUpdateChannel_FourSlotsRoundTrip(t *testing.T) {
	fourSlots := []channel_domain.Slot{
		{ID: "slot_1", Name: "slot_one", Role: "drafter", VaultID: "vault_1", Gated: false, Order: 1},
		{ID: "slot_2", Name: "slot_two", Role: "approver", VaultID: "vault_2", Gated: true, Order: 2},
		{ID: "slot_3", Name: "slot_three", Role: "signatory", VaultID: "vault_3", Gated: true, Order: 3},
		{ID: "slot_4", Name: "slot_four", Role: "reviewer", VaultID: "vault_4", Gated: false, Order: 4},
	}

	req := &channel_domain.UpdateChannelRequest{
		Channel: channel_domain.Channel{
			ID:    "ch_four_slots_test",
			Title: "Four Slots Channel",
			Slots: fourSlots,
		},
	}

	var capturedOutboundSlots []map[string]interface{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "expected PUT", http.StatusMethodNotAllowed)
			return
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if rawSlots, ok := body["slots"].([]interface{}); ok {
			for _, s := range rawSlots {
				if sm, ok := s.(map[string]interface{}); ok {
					capturedOutboundSlots = append(capturedOutboundSlots, sm)
				}
			}
		}

		// Cloud returns response aggregate with whatever slots were received and decoded
		returnedSlots := make([]tracecore_types.CloudChannelSlot, 0)
		if rawSlots, ok := body["slots"].([]interface{}); ok {
			for _, s := range rawSlots {
				if sm, ok := s.(map[string]interface{}); ok {
					id, _ := sm["id"].(string)
					name, _ := sm["name"].(string)
					role, _ := sm["role"].(string)
					vaultID, _ := sm["vault_id"].(string)
					gated, _ := sm["gated"].(bool)
					var order int
					if o, ok := sm["order"].(float64); ok {
						order = int(o)
					}
					returnedSlots = append(returnedSlots, tracecore_types.CloudChannelSlot{
						ID:      id,
						Name:    name,
						Role:    role,
						VaultID: vaultID,
						Gated:   gated,
						Order:   order,
					})
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": 200,
			"data": tracecore_types.CloudChannelDTO{
				ID:    "ch_four_slots_test",
				Title: "Four Slots Channel",
				Slots: returnedSlots,
			},
			"message": "success",
			"success": true,
		})
	}))
	defer ts.Close()

	tc := NewTracecoreClient(ts.URL, "token_123", ts.URL, ts.URL)
	resp, err := tc.UpdateChannel(context.Background(), req)

	// Boundary 1 Assertion: Before cloud request (outbound serialization)
	t.Logf("BOUNDARY 1 (Outbound Request Payload): captured %d slots in HTTP PUT body", len(capturedOutboundSlots))
	if len(capturedOutboundSlots) != 4 {
		t.Fatalf("FAILED AT BOUNDARY 1 (Outbound Request): Expected 4 slots sent in HTTP request body, got %d", len(capturedOutboundSlots))
	}

	// Response check
	if err != nil {
		t.Fatalf("UpdateChannel HTTP request failed: %v", err)
	}

	// Boundary 2 Assertion: Read-back / Decoded Channel (after response mapping)
	readBackSlots := resp.Data.Slots
	t.Logf("BOUNDARY 2 (Read-Back Decoded Aggregate): received %d slots in response domain channel", len(readBackSlots))
	if len(readBackSlots) != 4 {
		t.Fatalf("FAILED AT BOUNDARY 2 (Cloud decoding / Read-back mapping): Expected 4 slots in read-back channel, got %d", len(readBackSlots))
	}

	// Verify all 4 slots present and uncorrupted
	for i, expected := range fourSlots {
		actual := readBackSlots[i]
		if actual.ID != expected.ID || actual.Name != expected.Name || actual.Role != expected.Role || actual.VaultID != expected.VaultID {
			t.Errorf("Slot[%d] mismatch: expected %+v, got %+v", i, expected, actual)
		}
	}
}


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

func TestGetChannel_ZeroSlotsInDB_ReturnsZeroSlots(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "expected GET", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": 200,
			"data": tracecore_types.CloudChannelDTO{
				ID:          "ch_zero_slots",
				Title:       "Zero Slots Channel",
				Slots:       []tracecore_types.CloudChannelSlot{},
				Assignments: []tracecore_types.CloudChannelAssignment{},
				Properties:  []tracecore_types.CloudChannelProperty{},
			},
			"message": "success",
			"success": true,
		})
	}))
	defer ts.Close()

	tc := NewTracecoreClient(ts.URL, "token_123", ts.URL, ts.URL)
	resp, err := tc.GetChannel(context.Background(), &channel_domain.GetChannelRequest{
		ChannelID: "ch_zero_slots",
	})

	if err != nil {
		t.Fatalf("GetChannel HTTP request failed: %v", err)
	}

	readBackSlots := resp.Data.Slots
	t.Logf("READ BOUNDARY VERIFICATION: received %d slots in decoded domain channel", len(readBackSlots))

	if len(readBackSlots) != 0 {
		t.Fatalf("FAILED AT READ BOUNDARY: Expected 0 slots from DB, got %d", len(readBackSlots))
	}
}

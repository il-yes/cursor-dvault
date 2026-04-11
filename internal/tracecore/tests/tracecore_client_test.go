package tracecore_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	tracecore_types "vault-app/internal/tracecore/types"
)

func TestDirectAddToIPFSCall(t *testing.T) {
    // Mimic your exact request
    req := tracecore_types.SyncVaultStreamRequest{
        UserID:    "3f0b3c1e-ec7e-4a69-8d06-c4a37e5107ef",
        VaultName: "Guam",
        // Metadata:  map[string]string{},
        // Stream: someTestBytes,   // set to some non‑empty slice if you want
    }

    // Build JSON body
    body := &bytes.Buffer{}
    if err := json.NewEncoder(body).Encode(req); err != nil {
        t.Fatalf("encode req: %v", err)
    }

    // Construct URL
    baseURL := "http://localhost:4001/api"
    url := fmt.Sprintf("%s/vaults/%s/storage/%s", baseURL, req.UserID, req.VaultName)

    t.Logf("TEST: URL = %s", url)
    t.Logf("TEST: request body = %s", body.String())

    // Build request
    httpRequest, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, body)
    if err != nil {
        t.Fatalf("new request: %v", err)
    }
    httpRequest.Header.Set("Content-Type", "application/json")

    // Use the same HTTP client your app uses, or just http.DefaultClient
    client := &http.Client{}
    resp, err := client.Do(httpRequest)
    if err != nil {
        t.Fatalf("do: %v", err)
    }
    defer resp.Body.Close()

    t.Logf("TEST: status = %s", resp.Status)

    respBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Fatalf("read body: %v", err)
    }

    t.Logf("TEST: raw body (hex) = %x", respBytes)
    if len(respBytes) > 0 {
        t.Logf("TEST: raw body (string) = %s", string(respBytes))
    } else {
        t.Logf("TEST: raw body = (empty)")
    }

    if len(respBytes) == 0 {
        t.Fatal("empty body")
    }

    var cloudResp tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse]
    if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
        t.Fatalf("unmarshal: %v\nraw body: %s", err, respBytes)
    }

    t.Logf("TEST: cloudResp = %+v", cloudResp)
}
package shared_realtime

import "encoding/json"

type Message struct {
    Version int             `json:"version"`
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"`
}
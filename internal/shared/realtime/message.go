package shared_realtime

import "encoding/json"

type Message struct {
    Version int             `json:"version"`
    Type    string          `json:"type"`
    Seq     uint64          `json:"seq"`
    Payload json.RawMessage `json:"payload"`
}
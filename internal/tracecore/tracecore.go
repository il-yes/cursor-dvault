package tracecore



type CommitPayload struct {
	RepoID          string         `json:"repo_id"`
	Branch          string         `json:"branch"`
	Metadata        CommitMetadata `json:"metadata"`
	ValidationRules []string       `json:"validation_rules"`
}

type CommitMetadata struct {
	Message      string            `json:"message"`
	Content      map[string]any    `json:"content"`
	Context      map[string]string `json:"context"`
	StatusChange StatusChange      `json:"status_change"`
	Actor        Actor             `json:"actor"`
	Attachments  []string          `json:"attachments,omitempty"`
	Signature    string            `json:"signature,omitempty"` // actor-level
}


type StatusChange struct {
	Old string `json:"old"`
	New string `json:"new"`
}

type Actor struct {
	ID            string `json:"id"`
	Role          string `json:"role"`
	DelegatedFrom string `json:"delegated_from"` // Optional org/unit
	Signature     string `json:"signature"`
}
type CommitResponse struct {
	CommitID    string `json:"commit_id"`
	Status      int `json:"status"`
	Anchored    bool   `json:"anchored"`
	CID         string `json:"cid"`
	TxID        string `json:"tx_id"`
	RuleResults map[string]struct {
		Passed bool `json:"passed"`
		Steps  []struct {
			Name   string `json:"name"`
			Passed bool   `json:"passed"`
		} `json:"steps"`
	} `json:"rule_results"`
}

type CommitEnvelope struct {
	Commit    CommitPayload `json:"commit"`
	Signature string        `json:"signature"` // Overall envelope signature for the app/org signature.
    Meta      struct {
        AppID     string `json:"app_id"`
        Timestamp string `json:"timestamp"`
    } `json:"meta,omitempty"`
}




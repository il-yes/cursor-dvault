package tracecore_types

type ProjectOverviewDTO struct {
	ID                    string   `json:"id,omitempty"`
	ProjectID             string   `json:"project_id,omitempty"`
	Code                  string   `json:"code,omitempty"`
	ProjectReference      string   `json:"project_reference,omitempty"`
	Name                  string   `json:"name,omitempty"`
	ProjectName           string   `json:"project_name,omitempty"`
	ContractID            string   `json:"contract_id,omitempty"`
	Type                  string   `json:"type,omitempty"`
	ProjectType           string   `json:"project_type,omitempty"`
	Sector                string   `json:"sector,omitempty"`
	Status                string   `json:"status,omitempty"`
	Location              string   `json:"location,omitempty"`
	Description           string   `json:"description,omitempty"`
	CurrentPhase          string   `json:"current_phase,omitempty"`
	ProgressPercent       int      `json:"progress_percent,omitempty"`
	OpenIssuesCount       int      `json:"open_issues_count,omitempty"`
	PendingDecisionsCount int      `json:"pending_decisions_count,omitempty"`
	ActiveDelay           string   `json:"active_delay,omitempty"`
	TargetCompletion      string   `json:"target_completion,omitempty"`
	RecentActivity        string   `json:"recent_activity,omitempty"`
	Image                 string   `json:"image,omitempty"`
	ConnectedOrgs         []string `json:"connected_orgs,omitempty"`
	IsAuthoritative       bool     `json:"is_authoritative,omitempty"`
	BudgetSpentPercent    int      `json:"budget_spent_percent,omitempty"`
	ScheduleDay           int      `json:"schedule_day,omitempty"`
}

package tracecore_types

import (
	"encoding/json"
	"strings"
)

type LocationCoordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type FlexLocation struct {
	Address     string               `json:"address,omitempty"`
	City        string               `json:"city,omitempty"`
	Country     string               `json:"country,omitempty"`
	Coordinates *LocationCoordinates `json:"coordinates,omitempty"`
	Value       string               `json:"-"`
}

func (fl *FlexLocation) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		fl.Value = s
		return nil
	}
	type rawLoc FlexLocation
	var rl rawLoc
	if err := json.Unmarshal(b, &rl); err != nil {
		return err
	}
	*fl = FlexLocation(rl)
	return nil
}

func (fl FlexLocation) MarshalJSON() ([]byte, error) {
	city := fl.City
	country := fl.Country
	address := fl.Address

	if city == "" && country == "" && address == "" && fl.Value != "" {
		parts := strings.Split(fl.Value, ",")
		if len(parts) >= 2 {
			city = strings.TrimSpace(parts[0])
			country = strings.TrimSpace(strings.Join(parts[1:], ", "))
		} else {
			city = strings.TrimSpace(fl.Value)
		}
	}

	type locStruct struct {
		Address     string               `json:"address"`
		City        string               `json:"city"`
		Country     string               `json:"country"`
		Coordinates *LocationCoordinates `json:"coordinates,omitempty"`
	}

	return json.Marshal(locStruct{
		Address:     address,
		City:        city,
		Country:     country,
		Coordinates: fl.Coordinates,
	})
}

func (fl FlexLocation) String() string {
	if fl.Value != "" {
		return fl.Value
	}
	parts := []string{}
	if fl.Address != "" {
		parts = append(parts, fl.Address)
	}
	if fl.City != "" {
		parts = append(parts, fl.City)
	}
	if fl.Country != "" {
		parts = append(parts, fl.Country)
	}
	return strings.Join(parts, ", ")
}

type ProjectOverviewDTO struct {
	ID                    string       `json:"id,omitempty"`
	ProjectID             string       `json:"project_id,omitempty"`
	WorkspaceID           string       `json:"workspace_id,omitempty"`
	Code                  string       `json:"code,omitempty"`
	ProjectReference      string       `json:"project_reference,omitempty"`
	Name                  string       `json:"name,omitempty"`
	ProjectName           string       `json:"project_name,omitempty"`
	ContractID            string       `json:"contract_id,omitempty"`
	Type                  string       `json:"type,omitempty"`
	ProjectType           string       `json:"project_type,omitempty"`
	Sector                string       `json:"sector,omitempty"`
	Status                string       `json:"status,omitempty"`
	Location              FlexLocation `json:"location,omitempty"`
	Description           string       `json:"description,omitempty"`
	CurrentPhase          string       `json:"current_phase,omitempty"`
	ProgressPercent       int          `json:"progress_percent,omitempty"`
	ProgressPercentage    int          `json:"progress_percentage,omitempty"`
	OpenIssuesCount       int          `json:"open_issues_count,omitempty"`
	PendingDecisionsCount int          `json:"pending_decisions_count,omitempty"`
	ActiveDelay           string       `json:"active_delay,omitempty"`
	TargetCompletion      string       `json:"target_completion,omitempty"`
	RecentActivity        string       `json:"recent_activity,omitempty"`
	StatusSummary         string       `json:"status_summary,omitempty"`
	Image                 string       `json:"image,omitempty"`
	ConnectedOrgs         []string     `json:"connected_orgs,omitempty"`
	IsAuthoritative       bool         `json:"is_authoritative,omitempty"`
	BudgetSpentPercent    int          `json:"budget_spent_percent,omitempty"`
	ScheduleDay           int          `json:"schedule_day,omitempty"`
}

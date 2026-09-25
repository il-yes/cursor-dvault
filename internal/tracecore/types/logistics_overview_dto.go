package tracecore_types

// LogisticsOverviewDTO is the composite application read model for delivery operations
// and transport delay explanations.
type LogisticsOverviewDTO struct {
	Delivery    DeliveryOverviewDTO     `json:"delivery"`
	Project     ProjectReferenceDTO     `json:"project"`
	Requirement RequirementReferenceDTO  `json:"requirement"`
	Material    MaterialReferenceDTO    `json:"material"`
	Supplier    SupplierReferenceDTO    `json:"supplier"`
	Offer       OfferReferenceDTO       `json:"offer"`
	Site        SiteReferenceDTO        `json:"site"`
	Transport   TransportOverviewDTO    `json:"transport"`

	Issue      *IssueOverviewDTO      `json:"issue,omitempty"`
	Decision   *DecisionOverviewDTO   `json:"decision,omitempty"`
	Inspection *InspectionOverviewDTO `json:"inspection,omitempty"`
}

type DeliveryOverviewDTO struct {
	ID                  string `json:"id,omitempty"`
	Reference           string `json:"reference,omitempty"`
	Status              string `json:"status,omitempty"`
	PlannedDeliveryDate string `json:"planned_delivery_date,omitempty"`
	ETA                 string `json:"eta,omitempty"`
	ActualDeliveryDate  string `json:"actual_delivery_date,omitempty"`
	Quantity            int    `json:"quantity,omitempty"`
	Unit                string `json:"unit,omitempty"`
	DeliveryNotes       string `json:"delivery_notes,omitempty"`
}

type ProjectReferenceDTO struct {
	ID   string `json:"id,omitempty"`
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type RequirementReferenceDTO struct {
	ID          string `json:"id,omitempty"`
	Code        string `json:"code,omitempty"`
	TargetPhase string `json:"target_phase,omitempty"`
}

type MaterialReferenceDTO struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Specification string `json:"specification,omitempty"`
}

type SupplierReferenceDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type OfferReferenceDTO struct {
	ID                   string `json:"id,omitempty"`
	Reference            string `json:"reference,omitempty"`
	TotalPrice           string `json:"total_price,omitempty"`
	UnitPrice            string `json:"unit_price,omitempty"`
	PromisedDeliveryDate string `json:"promised_delivery_date,omitempty"`
}

type SiteReferenceDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type TransportOverviewDTO struct {
	ID               string   `json:"id,omitempty"`
	Reference        string   `json:"reference,omitempty"`
	Vehicle          string   `json:"vehicle,omitempty"`
	Driver           string   `json:"driver,omitempty"`
	Status           string   `json:"status,omitempty"`
	Origin           string   `json:"origin,omitempty"`
	Destination      string   `json:"destination,omitempty"`
	Route            string   `json:"route,omitempty"`
	PlannedDeparture string   `json:"planned_departure,omitempty"`
	ActualDeparture  string   `json:"actual_departure,omitempty"`
	PlannedArrival   string   `json:"planned_arrival,omitempty"`
	ETA              string   `json:"eta,omitempty"`
	ActualArrival    string   `json:"actual_arrival,omitempty"`
	Constraints      []string `json:"constraints,omitempty"`
	DelayReason      string   `json:"delay_reason,omitempty"`
}

type IssueOverviewDTO struct {
	ID                 string   `json:"id,omitempty"`
	Reference          string   `json:"reference,omitempty"`
	Title              string   `json:"title,omitempty"`
	Description        string   `json:"description,omitempty"`
	Severity           string   `json:"severity,omitempty"`
	Status             string   `json:"status,omitempty"`
	ReportedBy         string   `json:"reported_by,omitempty"`
	ReportedAt         string   `json:"reported_at,omitempty"`
	Impact             string   `json:"impact,omitempty"`
	EvidenceReferences []string `json:"evidence_references,omitempty"`
}

type DecisionOverviewDTO struct {
	ID                    string   `json:"id,omitempty"`
	Reference             string   `json:"reference,omitempty"`
	Subject               string   `json:"subject,omitempty"`
	Context               string   `json:"context,omitempty"`
	TechnicalAssessment   string   `json:"technical_assessment,omitempty"`
	RisksIdentified       []string `json:"risks_identified,omitempty"`
	ParticipantsConsulted []string `json:"participants_consulted,omitempty"`
	OptionsConsidered     []string `json:"options_considered,omitempty"`
	Decision              string   `json:"decision,omitempty"`
	DecidedBy             string   `json:"decided_by,omitempty"`
	DecisionDate          string   `json:"decision_date,omitempty"`
	Consequence           string   `json:"consequence,omitempty"`
	Status                string   `json:"status,omitempty"`
}

type InspectionOverviewDTO struct {
	ID                 string   `json:"id,omitempty"`
	Reference          string   `json:"reference,omitempty"`
	Type               string   `json:"type,omitempty"`
	Inspector          string   `json:"inspector,omitempty"`
	InspectionDate     string   `json:"inspection_date,omitempty"`
	Status             string   `json:"status,omitempty"`
	Result             string   `json:"result,omitempty"`
	Criteria           []string `json:"criteria,omitempty"`
	Findings           []string `json:"findings,omitempty"`
	Notes              string   `json:"notes,omitempty"`
	EvidenceReferences []string `json:"evidence_references,omitempty"`
}

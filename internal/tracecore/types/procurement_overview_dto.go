package tracecore_types

type ProcurementOverviewDTO struct {
	ID                    string `json:"id,omitempty"`
	RequirementID         string `json:"requirement_id,omitempty"`
	ProjectID             string `json:"project_id,omitempty"`
	ProjectCode           string `json:"project_code,omitempty"`
	ProjectName           string `json:"project_name,omitempty"`
	Code                  string `json:"code,omitempty"`
	MaterialID            string `json:"material_id,omitempty"`
	MaterialName          string `json:"material_name,omitempty"`
	Specification         string `json:"specification,omitempty"`
	Quantity              int    `json:"quantity,omitempty"`
	Unit                  string `json:"unit,omitempty"`
	TargetPhase           string `json:"target_phase,omitempty"`
	Status                string `json:"status,omitempty"`
	RequiredDate          string `json:"required_date,omitempty"`
	Priority              string `json:"priority,omitempty"`
	SiteID                string `json:"site_id,omitempty"`
	SiteName              string `json:"site_name,omitempty"`
	InvitedSuppliersCount int    `json:"invited_suppliers_count,omitempty"`
	OffersReceivedCount   int    `json:"offers_received_count,omitempty"`
	OfferID               string `json:"offer_id,omitempty"`
	OfferReference        string `json:"offer_reference,omitempty"`
	SupplierID            string `json:"supplier_id,omitempty"`
	SupplierName          string `json:"supplier_name,omitempty"`
	TotalPrice            string `json:"total_price,omitempty"`
	UnitPrice             string `json:"unit_price,omitempty"`
	PromisedDeliveryDate  string `json:"promised_delivery_date,omitempty"`
	OfferStatus           string `json:"offer_status,omitempty"`
	IsVerifiedSupplier    bool   `json:"is_verified_supplier,omitempty"`
}

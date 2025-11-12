package registry

import (
	app_config "vault-app/internal/config"
	"vault-app/internal/tracecore"
)

type VaultTemplate struct {
	Name            string
	Description     string
	Tracecore       ToolConfig
	BranchingModel  string
	CommitRules     []app_config.CommitRule
	Actors          []string
	Layout          VaultLayout
	SharingDefaults []SharingRule
	StellarAccount ToolConfig
}

type ToolConfig struct {
	Enabled bool
	Level   string // e.g. "minimal", "compliance", "enterprise"
}

type VaultLayout struct {
	Folders []string
}

// type CommitRule struct {
// 	Rule   string
// 	Actors []string
// }

type SharingRule struct {
	Role     string
	Provider string
	Group    string
}

var VaultTemplates = map[string]VaultTemplate{
	"personal": {
		Name:        "Personal Vault",
		Description: "Standalone secure vault for individual use",
		Tracecore: ToolConfig{
			Enabled: false,
			Level:   "minimal",
		},
		BranchingModel: "single",
		Actors:         []string{"user"},
		Layout: VaultLayout{
			Folders: []string{"Documents", "Passwords", "Notes"},
		},
		CommitRules:     []app_config.CommitRule{},
		SharingDefaults: []SharingRule{},
		StellarAccount: ToolConfig{
			Enabled: true,
			Level:   "minimal",		
		},
	},
	"regulated": {
		Name:        "Regulated Vault",
		Description: "Vault for compliance-heavy use cases",
		Tracecore: ToolConfig{
			Enabled: true,
			Level:   "compliance",
		},
		BranchingModel: "multi",
		Actors:         []string{"user", "compliance_bot"},
		Layout: VaultLayout{
			Folders: []string{"KYC", "Audits", "Transactions"},
		},
		CommitRules: []app_config.CommitRule{
			{Rule: "signature_required", Actors: []string{"compliance_bot"}},
		},
		SharingDefaults: []SharingRule{
			{Role: "auditor", Provider: "keycloak", Group: "auditors"},
		},
		StellarAccount: ToolConfig{
			Enabled: true,
			Level:   "pro",		
		},
	},
	"team": {
		Name:        "Team Vault",
		Description: "Multi-user team vault with approvals",
		Tracecore: ToolConfig{
			Enabled: true,
			Level:   "enterprise",
		},
		BranchingModel: "multi",
		Actors:         []string{"owner", "reviewer", "signer"},
		Layout: VaultLayout{
			Folders: []string{"Projects", "Contracts", "Reviews"},
		},
		CommitRules: []app_config.CommitRule{
			{Rule: "requires_review", Actors: []string{"reviewer"}},
			{Rule: "requires_signature", Actors: []string{"signer"}},
		},
		SharingDefaults: []SharingRule{
			{Role: "member", Provider: "keycloak", Group: "legal-team"},
		},
		StellarAccount: ToolConfig{
			Enabled: false,
			Level:   "org",		
		},
	},
	"individual:public": {
		Name:        "Personal Vault",
		Description: "Standalone secure vault for individual use",
		Tracecore: ToolConfig{
			Enabled: false,
			Level:   "minimal",
		},
		BranchingModel: "single",
		Actors:         []string{"user"},
		Layout: VaultLayout{
			Folders: []string{"Documents", "Passwords", "Notes"},
		},
		CommitRules:     []app_config.CommitRule{},
		SharingDefaults: []SharingRule{},
		StellarAccount: ToolConfig{
			Enabled: true,
			Level:   "minimal",		
		},
	},
	"individual:pro": {
		Name:        "Regulated Vault",
		Description: "Vault for compliance-heavy use cases",
		Tracecore: ToolConfig{
			Enabled: true,
			Level:   "compliance",
		},
		BranchingModel: "multi",
		Actors:         []string{"user", "compliance_bot"},
		Layout: VaultLayout{
			Folders: []string{"KYC", "Audits", "Transactions"},
		},
		CommitRules: []app_config.CommitRule{
			{Rule: "signature_required", Actors: []string{"compliance_bot"}},
		},
		SharingDefaults: []SharingRule{
			{Role: "auditor", Provider: "keycloak", Group: "auditors"},
		},
		StellarAccount: ToolConfig{
			Enabled: true,
			Level:   "pro",		
		},
	},
	"organization:institution": {
		Name:        "Team Vault",
		Description: "Multi-user team vault with approvals",
		Tracecore: ToolConfig{
			Enabled: true,
			Level:   "enterprise",
		},
		BranchingModel: "multi",
		Actors:         []string{"owner", "reviewer", "signer"},
		Layout: VaultLayout{
			Folders: []string{"Projects", "Contracts", "Reviews"},
		},
		CommitRules: []app_config.CommitRule{
			{Rule: "requires_review", Actors: []string{"reviewer"}},
			{Rule: "requires_signature", Actors: []string{"signer"}},
		},
		SharingDefaults: []SharingRule{
			{Role: "member", Provider: "keycloak", Group: "legal-team"},
		},
		StellarAccount: ToolConfig{
			Enabled: false,
			Level:   "org",		
		},
	},
	"organization:enterprise": {
		Name:        "Team Vault",
		Description: "Multi-user team vault with approvals",
		Tracecore: ToolConfig{
			Enabled: true,
			Level:   "enterprise",
		},
		BranchingModel: "multi",
		Actors:         []string{"owner", "reviewer", "signer"},
		Layout: VaultLayout{
			Folders: []string{"Projects", "Contracts", "Reviews"},
		},
		CommitRules: []app_config.CommitRule{
			{Rule: "requires_review", Actors: []string{"reviewer"}},
			{Rule: "requires_signature", Actors: []string{"signer"}},
		},
		SharingDefaults: []SharingRule{
			{Role: "member", Provider: "keycloak", Group: "legal-team"},
		},
		StellarAccount: ToolConfig{
			Enabled: false,
			Level:   "org",		
		},
	},
}

func (v *VaultTemplate) ShouldBootstrapTracecoreRepo(tc *tracecore.TracecoreClient) (*string, error) {
	if v.Tracecore.Enabled {
		return tc.CreateRepo()
	}
	return nil, nil
}

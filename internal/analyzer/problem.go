// Package analyzer defines configuration audit results and analysis logic.
package analyzer

// Severity describes the risk associated with a finding.
type Severity string

const (
	Low    Severity = "LOW"
	Medium Severity = "MEDIUM"
	High   Severity = "HIGH"
)

// Problem describes a finding without exposing the configuration's secret values.
type Problem struct {
	RuleID         string   `json:"rule_id"`
	Path           string   `json:"path"`
	Severity       Severity `json:"severity"`
	Explanation    string   `json:"explanation"`
	Recommendation string   `json:"recommendation"`
}

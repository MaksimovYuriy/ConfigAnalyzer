package analyzer

type Severity string

const (
	Low    Severity = "LOW"
	Medium Severity = "MEDIUM"
	High   Severity = "HIGH"
)

type Problem struct {
	RuleID         string   `json:"rule_id"`
	Path           string   `json:"path"`
	Severity       Severity `json:"severity"`
	Explanation    string   `json:"explanation"`
	Recommendation string   `json:"recommendation"`
}

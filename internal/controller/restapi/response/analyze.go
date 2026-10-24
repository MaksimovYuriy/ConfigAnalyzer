package response

import "github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"

type Problem struct {
	RuleID         string `json:"rule_id"`
	Path           string `json:"path"`
	Severity       string `json:"severity"`
	Explanation    string `json:"explanation"`
	Recommendation string `json:"recommendation"`
}

type Analyze []Problem

func AnalyzeFromEntities(problems []entity.Problem) Analyze {
	response := make(Analyze, 0, len(problems))
	for _, problem := range problems {
		response = append(response, Problem{
			RuleID:         problem.RuleID,
			Path:           problem.Path,
			Severity:       string(problem.Severity),
			Explanation:    problem.Explanation,
			Recommendation: problem.Recommendation,
		})
	}
	return response
}

type Error struct {
	Error string `json:"error"`
}

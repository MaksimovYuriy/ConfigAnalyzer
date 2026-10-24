package permissions

import "github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"

func Check(file entity.FileMetadata) []entity.Problem {
	if !file.Mode.IsRegular() || file.Mode.Perm()&0002 == 0 {
		return nil
	}
	return []entity.Problem{{
		RuleID:         "world-writable-config",
		Path:           file.Path,
		Severity:       entity.High,
		Explanation:    "Файл конфигурации доступен для записи остальным пользователям.",
		Recommendation: "Уберите право записи для остальных пользователей: chmod o-w <файл>.",
	}}
}

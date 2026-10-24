package request

import "github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"

type Analyze struct {
	Content string
}

func (request Analyze) ToEntity() entity.Configuration {
	return entity.Configuration{Content: request.Content}
}

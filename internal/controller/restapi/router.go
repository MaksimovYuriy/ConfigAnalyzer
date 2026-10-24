package restapi

import (
	"net/http"

	v1 "github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/restapi/v1"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
)

func NewRouter(configuration usecase.Configuration) http.Handler {
	controller := v1.NewController(configuration)
	router := http.NewServeMux()
	router.HandleFunc("POST /analyze", controller.Analyze)
	return router
}

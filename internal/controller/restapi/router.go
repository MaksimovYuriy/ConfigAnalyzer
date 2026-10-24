package restapi

import "net/http"

func NewRouter(controller *Controller) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /analyze", controller.analyze)

	return router
}

package v1

import (
	"mime"
	"net/http"
	"strings"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/restapi/v1/request"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/restapi/v1/response"
)

const maxRequestBytes = 1 << 20

func (controller *Controller) Analyze(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "multipart/form-data" {
		writeJSON(w, http.StatusUnsupportedMediaType, response.Error{Error: "Content-Type must be multipart/form-data"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
	defer r.Body.Close()
	if err := r.ParseMultipartForm(maxRequestBytes); err != nil {
		writeRequestError(w, err)
		return
	}
	defer r.MultipartForm.RemoveAll()
	input := request.Analyze{Content: r.PostForm.Get("content")}
	if strings.TrimSpace(input.Content) == "" {
		writeJSON(w, http.StatusBadRequest, response.Error{Error: "content must be a nonempty text field"})
		return
	}

	problems, err := controller.configuration.Analyze(r.Context(), input.ToEntity())
	if err != nil {
		writeAnalysisError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, response.AnalyzeFromEntities(problems))
}

package restapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/restapi"
	responsedto "github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/restapi/response"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase/configuration"
)

func TestRouter(t *testing.T) {
	for _, tt := range []struct {
		method, path string
		status       int
	}{
		{http.MethodPost, "/analyze", http.StatusOK},
		{http.MethodGet, "/analyze", http.StatusMethodNotAllowed},
		{http.MethodPut, "/analyze", http.StatusMethodNotAllowed},
		{http.MethodPost, "/unknown", http.StatusNotFound},
	} {
		t.Run(tt.method+tt.path, func(t *testing.T) {
			request := newFormRequest(t, tt.method, tt.path, map[string]string{"content": "password: do-not-echo"})
			response := httptest.NewRecorder()
			restapi.NewRouter(restapi.NewController(configuration.New(rules.Default()...))).ServeHTTP(response, request)
			if response.Code != tt.status {
				t.Fatalf("status=%d, want=%d", response.Code, tt.status)
			}
			if strings.Contains(response.Body.String(), "do-not-echo") {
				t.Fatal("request body leaked into response")
			}
			if tt.status == http.StatusOK {
				var body []responsedto.Problem
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil || len(body) != 1 || body[0].RuleID != "plaintext-password" {
					t.Fatalf("invalid analysis response: %s", response.Body.String())
				}
				if response.Header().Get("Content-Type") != "application/json" {
					t.Fatal("missing JSON content type")
				}
			}
		})
	}
}

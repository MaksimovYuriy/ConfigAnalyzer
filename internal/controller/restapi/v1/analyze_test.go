package v1_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	controller "github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/restapi/v1"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase/configuration"
)

type configurationStub struct {
	calls    int
	input    entity.Configuration
	problems []entity.Problem
	err      error
}

func (stub *configurationStub) Analyze(ctx context.Context, input entity.Configuration) ([]entity.Problem, error) {
	stub.calls++
	stub.input = input
	return stub.problems, stub.err
}

func newFormRequest(t *testing.T, method, path string, fields map[string]string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for name, value := range fields {
		if err := writer.WriteField(name, value); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(method, path, &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request
}

func TestAnalyzeRequestValidation(t *testing.T) {
	for _, tt := range []struct {
		name   string
		fields map[string]string
		status int
	}{
		{"missing content", map[string]string{}, 400},
		{"empty content", map[string]string{"content": " "}, 400},
		{"query content ignored", map[string]string{}, 400},
		{"too large", map[string]string{"content": strings.Repeat("a", 1<<20)}, 413},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stub := &configurationStub{}
			request := newFormRequest(t, http.MethodPost, "/analyze?content=5", tt.fields)
			response := httptest.NewRecorder()
			controller.NewController(stub).Analyze(response, request)
			if response.Code != tt.status || stub.calls != 0 {
				t.Fatalf("status=%d, calls=%d", response.Code, stub.calls)
			}
		})
	}
}

func TestAnalyzeInvalidMultipart(t *testing.T) {
	for _, tt := range []struct {
		contentType string
		status      int
	}{
		{"", 415},
		{"application/json", 415},
		{"application/yaml", 415},
		{"multipart/form-data", 400},
		{"multipart/form-data; boundary=test", 400},
	} {
		t.Run(tt.contentType, func(t *testing.T) {
			stub := &configurationStub{}
			request := httptest.NewRequest(http.MethodPost, "/analyze", strings.NewReader("invalid multipart"))
			request.Header.Set("Content-Type", tt.contentType)
			response := httptest.NewRecorder()
			controller.NewController(stub).Analyze(response, request)
			if response.Code != tt.status || stub.calls != 0 {
				t.Fatalf("status=%d, calls=%d", response.Code, stub.calls)
			}
		})
	}
}

func TestAnalyzeDTOConversion(t *testing.T) {
	stub := &configurationStub{problems: []entity.Problem{{RuleID: "rule", Path: "key", Severity: entity.High, Explanation: "reason", Recommendation: "fix"}}}
	request := newFormRequest(t, http.MethodPost, "/analyze", map[string]string{"content": "key: value\n  nested: text", "extra": "ignored"})
	response := httptest.NewRecorder()
	controller.NewController(stub).Analyze(response, request)
	if stub.input.Content != "key: value\n  nested: text" || stub.calls != 1 {
		t.Fatalf("input=%+v, calls=%d", stub.input, stub.calls)
	}
	want := `[{"rule_id":"rule","path":"key","severity":"HIGH","explanation":"reason","recommendation":"fix"}]` + "\n"
	if response.Code != 200 || response.Body.String() != want {
		t.Fatalf("status=%d, body=%s", response.Code, response.Body.String())
	}
}

func TestAnalyzeErrorResponses(t *testing.T) {
	for _, tt := range []struct {
		err    error
		status int
	}{
		{usecase.ErrInvalidConfiguration, 400},
		{context.Canceled, 408},
		{errors.New("secret internal details"), 500},
	} {
		stub := &configurationStub{err: tt.err}
		request := newFormRequest(t, http.MethodPost, "/analyze", map[string]string{"content": "5"})
		response := httptest.NewRecorder()
		controller.NewController(stub).Analyze(response, request)
		if response.Code != tt.status || strings.Contains(response.Body.String(), "secret internal details") {
			t.Fatalf("status=%d, body=%s", response.Code, response.Body.String())
		}
	}
}

func TestAnalyzeNoFindingsReturnsEmptyArray(t *testing.T) {
	request := newFormRequest(t, http.MethodPost, "/analyze", map[string]string{"content": "5"})
	response := httptest.NewRecorder()
	controller.NewController(configuration.New(entity.RuleSet{})).Analyze(response, request)
	if response.Code != 200 || response.Body.String() != "[]\n" {
		t.Fatalf("status=%d, body=%s", response.Code, response.Body.String())
	}
}

func TestAnalyzeFormConfigurationFormats(t *testing.T) {
	for _, content := range []string{
		"logging:\n  level: debug\n",
		`{"logging":{"level":"debug"}}`,
	} {
		request := newFormRequest(t, http.MethodPost, "/analyze", map[string]string{"content": content})
		response := httptest.NewRecorder()
		controller.NewController(configuration.New(rules.Default())).Analyze(response, request)
		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "debug-logging") {
			t.Fatalf("status=%d, body=%s", response.Code, response.Body.String())
		}
	}
}

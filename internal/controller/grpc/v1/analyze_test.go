package v1_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	grpcv1 "github.com/MaksimovYuriy/ConfigAnalyzer/docs/proto/v1"
	controller "github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/grpc/v1"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase/configuration"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type configurationStub struct {
	calls    int
	ctx      context.Context
	input    entity.Configuration
	problems []entity.Problem
	err      error
}

func (stub *configurationStub) Analyze(ctx context.Context, input entity.Configuration) ([]entity.Problem, error) {
	stub.calls++
	stub.ctx = ctx
	stub.input = input
	return stub.problems, stub.err
}

func TestAnalyzeRejectsEmptyContent(t *testing.T) {
	for _, input := range []*grpcv1.AnalyzeRequest{nil, {}, {Content: " \n\t"}} {
		stub := &configurationStub{}
		result, err := controller.NewController(stub).Analyze(context.Background(), input)
		if status.Code(err) != codes.InvalidArgument || result != nil || stub.calls != 0 {
			t.Fatalf("result=%v, error=%v, calls=%d", result, err, stub.calls)
		}
	}
}

func TestAnalyzeConvertsRequestAndResponse(t *testing.T) {
	for _, tt := range []struct {
		severity entity.Severity
		want     grpcv1.Severity
	}{
		{entity.Low, grpcv1.Severity_SEVERITY_LOW},
		{entity.Medium, grpcv1.Severity_SEVERITY_MEDIUM},
		{entity.High, grpcv1.Severity_SEVERITY_HIGH},
		{entity.Severity("unknown"), grpcv1.Severity_SEVERITY_UNSPECIFIED},
	} {
		t.Run(string(tt.severity), func(t *testing.T) {
			stub := &configurationStub{problems: []entity.Problem{{
				RuleID: "rule", Path: "key", Severity: tt.severity,
				Explanation: "reason", Recommendation: "fix",
			}}}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			content := "logging:\n  level: debug\n"
			result, err := controller.NewController(stub).Analyze(ctx, &grpcv1.AnalyzeRequest{Content: content})
			if err != nil {
				t.Fatal(err)
			}
			if stub.calls != 1 || stub.input.Content != content || stub.ctx != ctx {
				t.Fatalf("calls=%d, input=%+v, context preserved=%v", stub.calls, stub.input, stub.ctx == ctx)
			}
			want := &grpcv1.AnalyzeResponse{Problems: []*grpcv1.Problem{{
				RuleId: "rule", Path: "key", Severity: tt.want,
				Explanation: "reason", Recommendation: "fix",
			}}}
			if !proto.Equal(result, want) {
				t.Fatalf("result=%v, want=%v", result, want)
			}
		})
	}
}

func TestAnalyzeMapsErrors(t *testing.T) {
	for _, tt := range []struct {
		err  error
		code codes.Code
	}{
		{fmt.Errorf("secret details: %w", usecase.ErrInvalidConfiguration), codes.InvalidArgument},
		{fmt.Errorf("wrapped: %w", context.Canceled), codes.Canceled},
		{fmt.Errorf("wrapped: %w", context.DeadlineExceeded), codes.DeadlineExceeded},
		{errors.New("secret internal details"), codes.Internal},
	} {
		stub := &configurationStub{err: tt.err}
		result, err := controller.NewController(stub).Analyze(context.Background(), &grpcv1.AnalyzeRequest{Content: "5"})
		if status.Code(err) != tt.code || result != nil || strings.Contains(err.Error(), "secret") {
			t.Fatalf("result=%v, error=%v", result, err)
		}
	}
}

func TestAnalyzeWithRealUseCase(t *testing.T) {
	service := controller.NewController(configuration.New(rules.Default()...))
	for _, tt := range []struct {
		content  string
		findings int
		code     codes.Code
	}{
		{"logging:\n  level: debug", 1, codes.OK},
		{`{"password":"do-not-echo"}`, 1, codes.OK},
		{"5", 0, codes.OK},
		{"log: [", 0, codes.InvalidArgument},
	} {
		result, err := service.Analyze(context.Background(), &grpcv1.AnalyzeRequest{Content: tt.content})
		if status.Code(err) != tt.code {
			t.Fatalf("error=%v, want=%v", err, tt.code)
		}
		if tt.code != codes.OK {
			continue
		}
		if result == nil || len(result.Problems) != tt.findings || strings.Contains(result.String(), "do-not-echo") {
			t.Fatalf("unexpected response: %v", result)
		}
	}
}

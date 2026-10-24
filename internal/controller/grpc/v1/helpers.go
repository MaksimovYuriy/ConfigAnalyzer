package v1

import (
	"context"
	"errors"

	grpcv1 "github.com/MaksimovYuriy/ConfigAnalyzer/docs/proto/v1"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toResponse(problems []entity.Problem) *grpcv1.AnalyzeResponse {
	response := &grpcv1.AnalyzeResponse{
		Problems: make([]*grpcv1.Problem, 0, len(problems)),
	}
	for _, problem := range problems {
		response.Problems = append(response.Problems, &grpcv1.Problem{
			RuleId:         problem.RuleID,
			Path:           problem.Path,
			Severity:       toSeverity(problem.Severity),
			Explanation:    problem.Explanation,
			Recommendation: problem.Recommendation,
		})
	}
	return response
}

func toSeverity(severity entity.Severity) grpcv1.Severity {
	switch severity {
	case entity.Low:
		return grpcv1.Severity_SEVERITY_LOW
	case entity.Medium:
		return grpcv1.Severity_SEVERITY_MEDIUM
	case entity.High:
		return grpcv1.Severity_SEVERITY_HIGH
	default:
		return grpcv1.Severity_SEVERITY_UNSPECIFIED
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, usecase.ErrInvalidConfiguration):
		return status.Error(codes.InvalidArgument, "invalid configuration")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request timed out")
	default:
		return status.Error(codes.Internal, "configuration analysis failed")
	}
}

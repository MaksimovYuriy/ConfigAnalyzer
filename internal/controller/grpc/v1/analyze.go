package v1

import (
	"context"
	"strings"

	grpcv1 "github.com/MaksimovYuriy/ConfigAnalyzer/docs/proto/v1"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (controller *Controller) Analyze(ctx context.Context, input *grpcv1.AnalyzeRequest) (*grpcv1.AnalyzeResponse, error) {
	if strings.TrimSpace(input.GetContent()) == "" {
		return nil, status.Error(codes.InvalidArgument, "content must be nonempty")
	}
	configuration := entity.Configuration{Content: input.GetContent()}
	problems, err := controller.configuration.Analyze(ctx, configuration)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return toResponse(problems), nil
}

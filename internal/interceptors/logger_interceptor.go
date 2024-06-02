package interceptors

import (
	"context"

	"github.com/google/uuid"
	"github.com/ochiengotieno304/oneotp/internal/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ServerLogInterceptor(logger utils.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		traceID := uuid.New().String()

		ctx = context.WithValue(ctx, "trace_id", traceID)

		md, _ := metadata.FromIncomingContext(ctx)
		delete(md, "authorization")
		delete(md, "secret_key")
		delete(md, "api_key")
		
		logger.Debug(ctx, "", utils.Fields{"method": info.FullMethod, "metadata": md, "trace_id": traceID})
		// logger.Info(ctx, "", utils.Fields{"method": info.FullMethod, "metadata": md, "trace_id": traceID})
		m, err := handler(ctx, req)
		if err != nil {
			logger.Error(ctx, "Server Error", utils.Fields{"method": info.FullMethod, "error": err.Error(), "trace_id": traceID})
			return nil, err
		}

		return m, err
	}
}
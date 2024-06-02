package interceptors

import (
	"context"
	"strings"

	"github.com/ochiengotieno304/oneotp/internal/utils/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip if method if registration
		x := strings.Split(info.FullMethod, "/") // split full method into array
		if x[len(x)-1] == "CreateAccount" {
			m, err := handler(ctx, req)

			return m, err
		}

		md, _ := metadata.FromIncomingContext(ctx)
		secret := md.Get("secret")
		apiKey := md.Get("api_key")

		if secret == nil {
			return nil, errors.ErrMissingSecret
		}

		if apiKey == nil {
			return nil, errors.ErrMissingAPIKey
		}

		m, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}
		return m, err
	}
}
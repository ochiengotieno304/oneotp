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
		var (
			requireSecret bool
			requireAPiKey bool
		)

		// Skip if method if registration
		x := strings.Split(info.FullMethod, "/") // split full method into array
		switch x[len(x)-1] {
		case "CreateAccount":
			requireSecret = false
			requireAPiKey = false
		case "GenerateCredentials":
			requireSecret = false
			requireAPiKey = true
		}

		md, _ := metadata.FromIncomingContext(ctx)
		secret := md.Get("secret")
		apiKey := md.Get("api_key")

		if requireAPiKey {
			if apiKey == nil {
				return nil, errors.ErrMissingAPIKey
			}
		}

		if requireSecret {
			if secret == nil {
				return nil, errors.ErrMissingSecret
			}
		}

		m, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return m, err
	}
}

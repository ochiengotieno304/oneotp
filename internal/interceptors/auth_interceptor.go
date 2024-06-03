package interceptors

import (
	"context"
	"strings"

	"github.com/ochiengotieno304/oneotp/internal/helpers/auth"
	"github.com/ochiengotieno304/oneotp/internal/utils/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var (
			requireSecret   bool
			requireClientID bool
		)

		// Skip if method if registration
		x := strings.Split(info.FullMethod, "/") // split full method into array
		switch x[len(x)-1] {
		case "CreateAccount":
			requireSecret = false
			requireClientID = false
		case "GenerateCredentials":
			requireSecret = false
			requireClientID = true
		case "RequestOTP", "VerifyOTP":
			requireSecret = true
			requireClientID = true
		}

		md, _ := metadata.FromIncomingContext(ctx)
		secret := md.Get("secret_key")
		clientID := md.Get("client_id")

		if requireClientID {
			if clientID == nil {
				return nil, errors.ErrMissingClientID
			}

			if clientID[0] == "" {
				return nil, errors.ErrBlankClientID
			}
		}

		if requireSecret {
			if secret == nil {
				return nil, errors.ErrMissingSecret
			}

			if secret[0] == "" {
				return nil, errors.ErrBlankSecretKey
			}
		}

		if requireClientID && requireSecret {
			err := auth.ValidateRequest(clientID[0], secret[0])
			if err != nil {
				return nil, err
			}
		}

		ctx = context.WithValue(ctx, `secretKey`, secret)
		ctx = context.WithValue(ctx, "clientID", clientID)

		m, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return m, err
	}
}

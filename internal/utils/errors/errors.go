package errors

import "errors"

var (
	ErrMissingSecret   error = errors.New("missing 'secret_key' header")
	ErrMissingClientID error = errors.New("missing 'client_id' header")
	ErrBlankSecretKey  error = errors.New("secret_key cannot be blank")
	ErrBlankClientID   error = errors.New("client_d cannot be blank")
)

package errors

import "errors"

var ErrMissingSecret error = errors.New("missing 'secret_key' header")
var ErrMissingAPIKey error = errors.New("missing 'api_key' header")
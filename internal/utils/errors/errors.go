package errors

import "errors"

var ErrMissingSecret error = errors.New("missing 'secret_key' header")
var ErrMissingClientID error = errors.New("missing 'client_id' header")

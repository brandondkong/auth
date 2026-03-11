package jwt

import "errors"

var ErrUnknownClaimsType = errors.New("unknown claims type")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

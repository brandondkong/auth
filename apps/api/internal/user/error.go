package user

import "errors"

var USER_NOT_FOUND_ERROR_CODE string = "user_not_found"

var ErrUserNotFound = errors.New("user not found")


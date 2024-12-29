package common

import "errors"

var (
	ValidationError     = errors.New("Request format is invalid")
	InternalServerError = errors.New("Something wrong happened! We're working on it")
	BadRequestError     = errors.New("Bad malformat request")
	UnauthorizedError   = errors.New("Unauthorized access")
	UserNotFoundError	= errors.New("User not found")
)

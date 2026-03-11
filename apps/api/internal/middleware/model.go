package middleware

type ErrMalformedRequest struct {
	status	int
	message	string
}

func (err *ErrMalformedRequest) Error() string {
	return err.message
}

func (err *ErrMalformedRequest) Status() int {
	return err.status
}

type ResponseOptions[T any] struct {
	Code	int
	Error	*string
	Message	string
	Data	T
}

type JsonResponse[T any] struct {
	Success bool	`json:"success"`
	Error	*string	`json:"error"`
	Message	string	`json:"message"`
	Data	T		`json:"data"`
}


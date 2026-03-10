package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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

func DecodeJsonRequestBody[K any](w http.ResponseWriter, r *http.Request, payload *K) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			return &ErrMalformedRequest{
				status: http.StatusUnsupportedMediaType,
				message: msg,
			}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(payload)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &ErrMalformedRequest{
				status:	http.StatusBadRequest,
				message: msg,
			}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &ErrMalformedRequest{
				status: http.StatusBadRequest,
				message: msg,
			}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &ErrMalformedRequest{
				status: http.StatusBadRequest,
				message: msg,
			}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)

			return &ErrMalformedRequest{
				status: http.StatusBadRequest,
				message: msg,
			}
		

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &ErrMalformedRequest{
				status: http.StatusBadRequest,
				message: msg,
			}

		case errors.As(err, &maxBytesError):
			msg := fmt.Sprintf("Request body must not be larger than %d bytes", maxBytesError.Limit)
			return &ErrMalformedRequest{
				status: http.StatusBadRequest,
				message: msg,
			}

		default:
			return err
		}
		
	}

	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return &ErrMalformedRequest{
			status: http.StatusBadRequest,
			message: msg,
		}
	}

	return nil
} 

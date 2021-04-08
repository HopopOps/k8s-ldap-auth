package server

import (
	"errors"
	"net/http"
)

type ServerError struct {
	e error
	s int
}

var (
	ErrServerError = &ServerError{
		e: errors.New(http.StatusText(http.StatusInternalServerError)),
		s: http.StatusInternalServerError,
	}
	// ErrNotAcceptable means that the request is not acceptable because of it's content-type, language or encoding
	ErrNotAcceptable = &ServerError{
		e: errors.New(http.StatusText(http.StatusNotAcceptable)),
		s: http.StatusNotAcceptable,
	}
	// ErrDecodeFailed means the request body could not be decoded
	ErrDecodeFailed = &ServerError{
		e: errors.New("Failed Decoding Request Body"),
		s: http.StatusBadRequest,
	}
	// ErrMalformedCredentials
	ErrMalformedCredentials = &ServerError{
		e: errors.New("Malformed Credential Object"),
		s: http.StatusBadRequest,
	}
	// ErrMalformedToken
	ErrMalformedToken = &ServerError{
		e: errors.New("Malformed Token Object"),
		s: http.StatusBadRequest,
	}
	// ErrUnauthorized
	ErrUnauthorized = &ServerError{
		e: errors.New(http.StatusText(http.StatusUnauthorized)),
		s: http.StatusUnauthorized,
	}
	// ErrForbidden
	ErrForbidden = &ServerError{
		e: errors.New(http.StatusText(http.StatusForbidden)),
		s: http.StatusForbidden,
	}
)

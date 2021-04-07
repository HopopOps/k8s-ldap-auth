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
	// ErrMethodNotAllowed means that the request method is not allowed on this route
	ErrMethodNotAllowed = &ServerError{
		e: errors.New("method not allowed"),
		s: http.StatusMethodNotAllowed,
	}
	// ErrNotAcceptable means that the request is not acceptable because of it's content-type, language or encoding
	ErrNotAcceptable = &ServerError{
		e: errors.New("not acceptable"),
		s: http.StatusNotAcceptable,
	}
	// ErrDecodeFailed means the request body could not be decoded
	ErrDecodeFailed = &ServerError{
		e: errors.New("failed decoding request body"),
		s: http.StatusBadRequest,
	}
	// ErrMalformedCredentials
	ErrMalformedCredentials = &ServerError{
		e: errors.New("malformed credential object"),
		s: http.StatusBadRequest,
	}
	// ErrMalformedToken
	ErrMalformedToken = &ServerError{
		e: errors.New("malformed token object"),
		s: http.StatusBadRequest,
	}
	// ErrUnauthorized
	ErrUnauthorized = &ServerError{
		e: errors.New("unauthorized"),
		s: http.StatusUnauthorized,
	}
	// ErrForbidden
	ErrForbidden = &ServerError{
		e: errors.New("forbidden"),
		s: http.StatusForbidden,
	}
)

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"bouchaud.org/legion/kubernetes/k8s-ldap-auth/ldap"
	"bouchaud.org/legion/kubernetes/k8s-ldap-auth/types"
)

// ContentTypeHeader is the HTTP header name that contains the message content type
const ContentTypeHeader = "Content-Type"

// ContentTypeJSON is the content type header value for JSON content
const ContentTypeJSON = "application/json"

var (
	// ErrMethodNotAllowed means that the request method is not allowed on this route
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrNotAcceptable means that the request is not acceptable because of it's content-type, language or encoding
	ErrNotAcceptable = errors.New("not acceptable")
	// ErrDecodeFailed means the request body could not be decoded
	ErrDecodeFailed = errors.New("failed decoding request body")
	// ErrMalformedCredentials
	ErrMalformedCredentials = errors.New("malformed credential object")
	// ErrMalformedToken
	ErrMalformedToken = errors.New("malformed token object")
	// ErrUnauthorized
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden
	ErrForbidden = errors.New("forbidden")
)

// Instance for managing pipelines with HTTP
type Instance struct {
	l *ldap.Ldap
	m []mux.MiddlewareFunc
}

func NewInstance(opts ...Option) *Instance {
	s := &Instance{
		m: []mux.MiddlewareFunc{},
	}

	r := mux.NewRouter()

	r.HandleFunc("/auth", s.authenticate()).Methods("POST")
	r.HandleFunc("/token", s.validate()).Methods("POST")
	r.Use(s.m...)

	http.Handle("/", r)

	return s
}

func (s *Instance) Start(addr string) error {
	if err := http.ListenAndServe(addr, nil); err != http.ErrServerClosed {
		return fmt.Errorf("Server stopped unexpectedly, %w", err)
	}

	return nil
}

func (s *Instance) authenticate() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get(ContentTypeHeader) != ContentTypeJSON {
			res.WriteHeader(http.StatusNotAcceptable)
			res.Write([]byte(ErrNotAcceptable.Error()))
			return
		}

		decoder := json.NewDecoder(req.Body)
		var credentials types.Credentials
		if err := decoder.Decode(&credentials); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(ErrDecodeFailed.Error()))
			return
		}
		defer req.Body.Close()

		if !credentials.IsValid() {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(ErrMalformedCredentials.Error()))
			return
		}

		_, err := s.l.Search(credentials.Username, credentials.Password)
		if err != nil {
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte(ErrUnauthorized.Error()))
			return
		}

		// TODO: implement

		res.Header().Set(ContentTypeHeader, ContentTypeJSON)
		json.NewEncoder(res).Encode(nil)
	}
}

func (s *Instance) validate() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get(ContentTypeHeader) != ContentTypeJSON {
			res.WriteHeader(http.StatusNotAcceptable)
			res.Write([]byte(ErrNotAcceptable.Error()))
			return
		}

		decoder := json.NewDecoder(req.Body)
		var token types.Token
		if err := decoder.Decode(&token); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(ErrDecodeFailed.Error()))
			return
		}
		defer req.Body.Close()

		if !token.IsValid() {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(ErrMalformedToken.Error()))
			return
		}

		// TODO: implement

		res.Header().Set(ContentTypeHeader, ContentTypeJSON)
		json.NewEncoder(res).Encode(nil)
	}
}

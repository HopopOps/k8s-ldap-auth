package server

import (
	"encoding/json"
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

func writeError(res http.ResponseWriter, s *ServerError) {
	res.WriteHeader(s.s)
	res.Write([]byte(s.e.Error()))

}

func (s *Instance) authenticate() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get(ContentTypeHeader) != ContentTypeJSON {
			writeError(res, ErrNotAcceptable)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var credentials types.Credentials
		if err := decoder.Decode(&credentials); err != nil {
			writeError(res, ErrDecodeFailed)
			return
		}
		defer req.Body.Close()

		if !credentials.IsValid() {
			writeError(res, ErrMalformedCredentials)
			return
		}

		_, err := s.l.Search(credentials.Username, credentials.Password)
		if err != nil {
			writeError(res, ErrUnauthorized)
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
			writeError(res, ErrNotAcceptable)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var token types.Token
		if err := decoder.Decode(&token); err != nil {
			writeError(res, ErrDecodeFailed)
			return
		}
		defer req.Body.Close()

		if !token.IsValid() {
			writeError(res, ErrMalformedToken)
			return
		}

		// TODO: implement

		res.Header().Set(ContentTypeHeader, ContentTypeJSON)
		json.NewEncoder(res).Encode(nil)
	}
}

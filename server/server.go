package server

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	auth "k8s.io/api/authentication/v1"
	machinery "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"

	"vbouchaud/k8s-ldap-auth/ldap"
	"vbouchaud/k8s-ldap-auth/types"
)

const ContentTypeHeader = "Content-Type"
const ContentTypeJSON = "application/json"

type Instance struct {
	l   *ldap.Ldap
	m   []mux.MiddlewareFunc
	k   *rsa.PrivateKey
	ttl int64
}

func NewInstance(opts ...Option) (*Instance, error) {
	s := &Instance{
		m: []mux.MiddlewareFunc{},
	}

	log.Info().Msg("Applying extra options.")
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	r := mux.NewRouter()

	log.Info().Msg("Registering route handlers.")
	r.HandleFunc("/auth", s.authenticate()).Methods("POST")
	r.HandleFunc("/token", s.validate()).Methods("POST")

	log.Info().Msg("Applying middlewares.")
	r.Use(s.m...)

	http.Handle("/", r)

	return s, nil
}

func (s *Instance) Start(addr string) error {
	if err := http.ListenAndServe(addr, nil); err != http.ErrServerClosed {
		return fmt.Errorf("Server stopped unexpectedly, %w", err)
	}

	return nil
}

func writeExecCredentialError(res http.ResponseWriter, s *ServerError) {
	res.WriteHeader(s.s)

	ec := client.ExecCredential{
		Spec: client.ExecCredentialSpec{
			// Response: &client.Response{
			//		Code: s.s,
			// },
		},
	}

	res.Header().Set(ContentTypeHeader, ContentTypeJSON)
	json.NewEncoder(res).Encode(ec)
}

func (s *Instance) authenticate() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get(ContentTypeHeader) != ContentTypeJSON {
			writeExecCredentialError(res, ErrNotAcceptable)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var credentials types.Credentials
		if err := decoder.Decode(&credentials); err != nil {
			writeExecCredentialError(res, ErrDecodeFailed)
			return
		}
		defer req.Body.Close()

		if !credentials.IsValid() {
			writeExecCredentialError(res, ErrMalformedCredentials)
			return
		}

		log.Debug().Str("username", credentials.Username).Msg("Received valid authentication request.")
		user, err := s.l.Search(credentials.Username, credentials.Password)
		if err != nil {
			writeExecCredentialError(res, ErrUnauthorized)
			return
		}

		log.Debug().Str("username", credentials.Username).Msg("Successfully authenticated.")

		token, err := types.NewToken(user, s.ttl)
		if err != nil {
			writeExecCredentialError(res, ErrServerError)
			return
		}

		tokenData, err := token.Payload(s.k)
		if err != nil {
			writeExecCredentialError(res, ErrServerError)
			return
		}

		tokenExp, err := token.Expiration()
		if err != nil {
			writeExecCredentialError(res, ErrServerError)
			return
		}

		log.Debug().Str("username", credentials.Username).Str("token", string(tokenData)).Msg("Sending back token.")

		res.Header().Set(ContentTypeHeader, ContentTypeJSON)
		json.NewEncoder(res).Encode(client.ExecCredential{
			Status: &client.ExecCredentialStatus{
				Token: string(tokenData),
				ExpirationTimestamp: &machinery.Time{
					Time: tokenExp,
				},
			},
		})
	}
}

func writeError(res http.ResponseWriter, s *ServerError) {
	res.WriteHeader(s.s)
	res.Write([]byte(s.e.Error()))
}

func writeTokenReviewError(res http.ResponseWriter, s *ServerError, tr auth.TokenReview) {
	res.WriteHeader(s.s)

	tr.Status.Authenticated = false
	tr.Status.Error = s.e.Error()

	res.Header().Set(ContentTypeHeader, ContentTypeJSON)
	json.NewEncoder(res).Encode(tr)
}

func (s *Instance) validate() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Debug().Msg("Got a request.")

		if req.Header.Get(ContentTypeHeader) != ContentTypeJSON {
			writeError(res, ErrNotAcceptable)
			return
		}

		log.Debug().Msg("Request is in JSON.")

		decoder := json.NewDecoder(req.Body)
		var tr auth.TokenReview
		if err := decoder.Decode(&tr); err != nil {
			writeError(res, ErrDecodeFailed)
			return
		}
		defer req.Body.Close()

		log.Debug().Str("token", tr.Spec.Token).Msg("Request is a TokenReview.")

		token, err := types.Parse([]byte(tr.Spec.Token), s.k)
		if err != nil {
			log.Debug().Str("err", err.Error()).Msg("Failed to parse token")

			writeTokenReviewError(res, ErrMalformedToken, tr)
			return
		}

		log.Debug().Msg("TokenReview was parsed.")

		if token.IsValid() == false {
			log.Debug().Msg("TokenReview is not valid.")
			tr.Status.Authenticated = false
		} else {
			user, err := token.GetUser()
			if err != nil {
				log.Debug().Str("error", err.Error()).Msg("Could not extract user.")

				writeTokenReviewError(res, ErrServerError, tr)
				return
			}

			log.Debug().Msg("Got user from token.")

			tr.Status.Authenticated = true
			tr.Status.User = *user
		}

		res.Header().Set(ContentTypeHeader, ContentTypeJSON)
		json.NewEncoder(res).Encode(tr)
	}
}

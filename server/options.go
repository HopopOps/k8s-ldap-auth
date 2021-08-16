package server

import (
	"crypto/rsa"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"vbouchaud/k8s-ldap-auth/ldap"
	"vbouchaud/k8s-ldap-auth/server/middlewares"
	"vbouchaud/k8s-ldap-auth/types"
)

// Option function for configuring a server instance
type Option func(*Instance) error

// WithLdap bind a ldap object to a server instance
func WithLdap(
	ldapURL,
	bindDN,
	bindPassword,
	searchBase,
	searchScope,
	searchFilter,
	memberofProperty,
	usernameProperty string,
	extraAttributes []string) Option {
	return func(i *Instance) error {
		i.l = ldap.NewInstance(
			ldapURL,
			bindDN,
			bindPassword,
			searchBase,
			searchScope,
			searchFilter,
			memberofProperty,
			usernameProperty,
			extraAttributes,
			append(extraAttributes, memberofProperty, usernameProperty),
		)

		return nil
	}
}

// WithMiddleware will bind the given middleware function to the root of the router
func WithMiddleware(m mux.MiddlewareFunc) Option {
	return func(i *Instance) error {
		i.m = append(i.m, m)

		return nil
	}
}

// WithAccessLogs add an access log middleware to the server
func WithAccessLogs() Option {
	return WithMiddleware(middlewares.AccessLog)
}

// WithLdap bind a ldap object to a server instance
func WithKey(privateKeyFile, publicKeyFile string) Option {
	return func(i *Instance) error {
		var (
			key *rsa.PrivateKey
			err error
		)

		if privateKeyFile != "" && publicKeyFile != "" {
			log.Info().Msg("privateKeyFile and publicKeyFile were provided, loading key.")
			key, err = types.LoadKey(privateKeyFile, publicKeyFile)
		} else {
			log.Info().Msg("No key provided, generating a new one.")
			key, err = types.GenerateKey()
		}

		i.k = key

		return err
	}
}

// WithLdap bind a ldap object to a server instance
func WithTTL(ttl int64) Option {
	return func(i *Instance) error {
		i.ttl = ttl

		return nil
	}
}

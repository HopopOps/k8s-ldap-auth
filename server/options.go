package server

import (
	"github.com/gorilla/mux"

	"bouchaud.org/legion/kubernetes/k8s-ldap-auth/ldap"
	"bouchaud.org/legion/kubernetes/k8s-ldap-auth/server/middlewares"
)

// Option function for configuring a server instance
type Option func(*Instance)

// WithLdap bind a ldap object to a server instance
func WithLdap(ldapURL, bindDN, bindPassword, searchBase, searchScope, searchFilter, memberOfProperty string, searchAttributes []string) Option {
	return func(i *Instance) {
		i.l = ldap.NewInstance(ldapURL, bindDN, bindPassword, searchBase, searchScope, searchFilter, memberOfProperty, searchAttributes)
	}
}

// WithMiddleware will bind the given middleware function to the root of the router
func WithMiddleware(m mux.MiddlewareFunc) Option {
	return func(i *Instance) {
		i.m = append(i.m, m)
	}
}

// WithAccessLogs add an access log middleware to the server
func WithAccessLogs() Option {
	return WithMiddleware(middlewares.AccessLog)
}

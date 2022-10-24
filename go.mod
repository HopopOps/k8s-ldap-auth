module vbouchaud/k8s-ldap-auth

go 1.16

require (
	github.com/adrg/xdg v0.4.0
	github.com/etherlabsio/healthcheck/v2 v2.0.0
	github.com/go-ldap/ldap/v3 v3.4.4
	github.com/gorilla/mux v1.8.0
	github.com/lestrrat-go/jwx v1.2.25
	github.com/mattn/go-isatty v0.0.16
	github.com/rs/zerolog v1.28.0
	github.com/urfave/cli/v2 v2.20.3
	github.com/zalando/go-keyring v0.2.1
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	k8s.io/api v0.25.1
	k8s.io/apimachinery v0.25.1
	k8s.io/client-go v0.25.1
)

module vbouchaud/k8s-ldap-auth

go 1.15

require (
	github.com/adrg/xdg v0.4.0
	github.com/etherlabsio/healthcheck/v2 v2.0.0
	github.com/go-ldap/ldap/v3 v3.4.2
	github.com/gorilla/mux v1.8.0
	github.com/lestrrat-go/jwx v1.2.20
	github.com/mattn/go-isatty v0.0.14
	github.com/rs/zerolog v1.26.1
	github.com/urfave/cli/v2 v2.4.0
	github.com/zalando/go-keyring v0.2.0
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
)

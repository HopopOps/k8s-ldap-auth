module vbouchaud/k8s-ldap-auth

go 1.15

require (
	github.com/adrg/xdg v0.4.0
	github.com/etherlabsio/healthcheck/v2 v2.0.0
	github.com/go-ldap/ldap/v3 v3.4.1
	github.com/gorilla/mux v1.8.0
	github.com/lestrrat-go/jwx v1.2.10
	github.com/mattn/go-isatty v0.0.14
	github.com/rs/zerolog v1.26.0
	github.com/urfave/cli/v2 v2.3.0
	github.com/zalando/go-keyring v0.1.1
	golang.org/x/term v0.0.0-20210406210042-72f3dc4e9b72
	k8s.io/api v0.22.4
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.3
)

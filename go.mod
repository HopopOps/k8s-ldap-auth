module vbouchaud/k8s-ldap-auth

go 1.15

require (
	github.com/adrg/xdg v0.3.3
	github.com/etherlabsio/healthcheck/v2 v2.0.0
	github.com/go-ldap/ldap v3.0.3+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/lestrrat-go/jwx v1.2.6
	github.com/mattn/go-isatty v0.0.13
	github.com/rs/zerolog v1.25.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/term v0.0.0-20210406210042-72f3dc4e9b72
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.1
)

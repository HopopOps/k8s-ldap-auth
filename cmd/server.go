package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"vbouchaud/k8s-ldap-auth/server"
)

func getServerCmd() *cli.Command {
	return &cli.Command{
		Name:     "server",
		Aliases:  []string{"s", "serve"},
		Usage:    "start the authentication server",
		HideHelp: false,
		Flags: []cli.Flag{
			// server configuration
			&cli.StringFlag{
				Name:    "host",
				Value:   "",
				EnvVars: []string{"HOST"},
				Usage:   "The `HOST` the server will listen on.",
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   3000,
				EnvVars: []string{"PORT"},
				Usage:   "The `PORT` the server will listen to.",
			},

			// ldap server configuration
			&cli.StringFlag{
				Name:    "ldap-host",
				Value:   "ldap://localhost",
				EnvVars: []string{"LDAP_ADDR"},
				Usage:   "The ldap `HOST` (and scheme) the server will authenticate against.",
			},

			// bind dn configuration
			&cli.StringFlag{
				Name:     "bind-dn",
				EnvVars:  []string{"LDAP_BINDDN"},
				Required: true,
				Usage:    "The service account `DN` to do the ldap search.",
			},
			&cli.StringFlag{
				Name:     "bind-credentials",
				EnvVars:  []string{"LDAP_BINDCREDENTIALS"},
				FilePath: "/etc/k8s-ldap-auth/ldap/password",
				Usage:    "The service account `PASSWORD` to do the ldap search, can be located in '/etc/k8s-ldap-auth/ldap/password'.",
			},

			// user search configuration
			&cli.StringFlag{
				Name:    "search-base",
				EnvVars: []string{"LDAP_USER_SEARCHBASE"},
				Usage:   "The `DN` where the ldap search will take place.",
			},
			&cli.StringFlag{
				Name:    "search-filter",
				Value:   "(&(objectClass=inetOrgPerson)(uid=%s))",
				EnvVars: []string{"LDAP_USER_SEARCHFILTER"},
				Usage:   "The `FILTER` to select users.",
			},
			&cli.StringFlag{
				Name:    "memberof-property",
				Value:   "ismemberof",
				EnvVars: []string{"LDAP_USER_MEMBEROFPROPERTY"},
				Usage:   "The `PROPERTY` that will be used to fetch groups. Usually memberof or ismemberof.",
			},
			&cli.StringFlag{
				Name:    "username-property",
				Value:   "uid",
				EnvVars: []string{"LDAP_USER_USERNAMEPROPERTY"},
				Usage:   "The `PROPERTY` that will be used as username in the TokenReview.",
			},
			&cli.StringSliceFlag{
				Name:    "extra-attributes",
				EnvVars: []string{"LDAP_USER_EXTRAATTR"},
				Usage:   "Repeatable. User `PROPERTY` to fetch. Those will be stored in extra values in the UserInfo object.",
			},
			&cli.StringFlag{
				Name:    "search-scope",
				Value:   "sub",
				EnvVars: []string{"LDAP_USER_SEARCHSCOPE"},
				Usage:   "The `SCOPE` of the search. Can take to values base object: 'base', single level: 'single' or whole subtree: 'sub'.",
			},

			// jtw signing configuration
			&cli.StringFlag{
				Name:    "private-key-file",
				Usage:   "The `PATH` to the private key file",
				EnvVars: []string{"PRIVATE_KEY_FILE"},
			},
			&cli.StringFlag{
				Name:    "public-key-file",
				Usage:   "The `PATH` to the public key file",
				EnvVars: []string{"PUBLIC_KEY_FILE"},
			},
			&cli.Int64Flag{
				Name:    "token-ttl",
				Value:   43200,
				EnvVars: []string{"TTL"},
				Usage:   "The `TTL` for newly generated tokens, in seconds",
			},
		},
		Action: func(c *cli.Context) error {
			var (
				port = c.Int("port")
				host = c.String("host")

				ldapURL          = c.String("ldap-host")
				bindDN           = c.String("bind-dn")
				bindPassword     = c.String("bind-credentials")
				searchBase       = c.String("search-base")
				searchScope      = c.String("search-scope")
				searchFilter     = c.String("search-filter")
				extraAttributes  = c.StringSlice("extra-attributes")
				memberofProperty = c.String("memberof-property")
				usernameProperty = c.String("username-property")

				privateKeyFile = c.String("private-key-file")
				publicKeyFile  = c.String("public-key-file")

				ttl = c.Int64("token-ttl")
			)

			addr := fmt.Sprintf("%s:%d", host, port)

			s, err := server.NewInstance(
				server.WithLdap(
					ldapURL,
					bindDN,
					bindPassword,
					searchBase,
					searchScope,
					searchFilter,
					memberofProperty,
					usernameProperty,
					extraAttributes,
				),
				server.WithAccessLogs(),
				server.WithKey(
					privateKeyFile,
					publicKeyFile,
				),
				server.WithTTL(ttl),
			)
			if err != nil {
				return fmt.Errorf("There was an error instanciation the server, %w", err)
			}

			if err := s.Start(addr); err != nil {
				return fmt.Errorf("There was an error starting the server, %w", err)
			}

			return nil
		},
	}
}

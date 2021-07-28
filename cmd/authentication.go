package cmd

import (
	"fmt"
	"path"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v2"

	"vbouchaud/k8s-ldap-auth/client"
)

func getAuthenticationCmd() *cli.Command {
	passwordFile := path.Join(xdg.ConfigHome, "k8s-ldap-auth", "password")

	return &cli.Command{
		Name:     "authenticate",
		Aliases:  []string{"a", "auth"},
		Usage:    "perform an authentication through a /auth endpoint",
		HideHelp: false,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "endpoint",
				Required: true,
				Usage:    "The full `URI` the client will authenticate against.",
			},
			&cli.StringFlag{
				Name:    "user",
				EnvVars: []string{"USER"},
				Usage:   "The `USER` the client will connect as.",
			},
			&cli.StringFlag{
				Name:     "password",
				EnvVars:  []string{"PASSWORD"},
				FilePath: passwordFile,
				Usage:    "The `PASSWORD` the client will connect with, can be located in '" + passwordFile + "'.",
			},
		},
		Action: func(c *cli.Context) error {
			var (
				addr     = c.String("endpoint")
				username = c.String("user")
				password = c.String("password")
			)

			err := client.Auth(addr, username, password)
			if err != nil {
				return fmt.Errorf("There was an error authenticating, %w", err)
			}

			return nil
		},
	}
}

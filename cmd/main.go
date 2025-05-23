package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"k8s-ldap-auth/version"
)

type action func(*cli.Context) error

// Start the cmd application
func Start() error {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(version.Version())
	}

	app := cli.NewApp()
	app.Name = "k8s-ldap-auth"
	app.Usage = "A client/server for kubernetes webhook authentication."
	app.Version = version.VERSION
	app.Compiled = version.Compiled()
	app.Authors = []*cli.Author{
		{
			Name:  "Vianney Bouchaud",
			Email: "vianney@hopopops.org",
		},
	}

	app.UseShortOptionHandling = true
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    "verbose",
			Value:   int(zerolog.ErrorLevel),
			EnvVars: []string{"VERBOSE"},
			Usage:   "The verbosity `LEVEL` - (rs/zerolog#Level values).",
		},
	}

	app.Before = func(c *cli.Context) error {
		var (
			verbose = zerolog.Level(c.Int("verbose"))
		)

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.SetGlobalLevel(verbose)

		if verbose < zerolog.InfoLevel {
			log.Logger = log.With().Logger()
		}

		return nil
	}

	app.Commands = []*cli.Command{
		getServerCmd(),
		getAuthenticationCmd(),
		getResetCmd(),
	}

	return app.Run(os.Args)
}

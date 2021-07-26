package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"bouchaud.org/legion/kubernetes/k8s-ldap-auth/client"
)

func getResetCmd() *cli.Command {
	return &cli.Command{
		Name:     "reset",
		Aliases:  []string{"r"},
		Usage:    "delete the cached ExecCredential to force authentication at next invocation",
		HideHelp: false,
		Action: func(c *cli.Context) error {
			// ignore error
			os.Remove(client.GetCacheFilePath())

			return nil
		},
	}
}

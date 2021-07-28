package main

import (
	"vbouchaud/k8s-ldap-auth/cmd"
)

func main() {
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}

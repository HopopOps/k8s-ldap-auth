package client

import (
	"encoding/json"
	"fmt"
	"time"

	client "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

func Auth(addr, user, pass string) error {
	var (
		err   error
		token []byte
		ec    client.ExecCredential
	)

	token = getCachedToken()

	// TODO: warn log
	err = json.Unmarshal(token, &ec)
	if err != nil || ec.Status.ExpirationTimestamp.Time.Unix() < time.Now().Unix() {
		token, err = performAuth(addr, user, pass)
		if err != nil {
			return err
		}

		err = json.Unmarshal(token, &ec)
		if err != nil {
			return err
		}
	}
	refreshCache(ec)

	fmt.Printf("%s\n", token)

	return nil
}

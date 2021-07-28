package client

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"vbouchaud/k8s-ldap-auth/types"
)

func Auth(addr, user, pass string) error {
	var (
		err   error
		token []byte
		ec    types.ExecCredential
	)

	token = getCachedToken()

	log.Info().Str("KUBERNETES_EXEC_INFO", os.Getenv("KUBERNETES_EXEC_INFO")).Msg("Testing variable of interest.")

	err = json.Unmarshal(token, &ec)
	if err != nil || ec.Status.ExpirationTimestamp.Time.Unix() < time.Now().Unix() {
		if err != nil {
			log.Warn().Err(err).Send()
		} else {
			log.Warn().Int64("expirationTimestamp", ec.Status.ExpirationTimestamp.Time.Unix()).Msg("ExecCredential expired.")
		}

		token, err = performAuth(addr, user, pass)
		if err != nil {
			log.Error().Err(err).Msg("Could not perform authentication.")
			return err
		}
		log.Info().Msg("Got credentials")
		log.Debug().RawJSON("token", token).Send()

		err = json.Unmarshal(token, &ec)
		if err != nil {
			log.Error().Err(err).Msg("Could not parse token.")
			return err
		}

		log.Info().Msg("Token parsed successfully")
	}

	token, err = ec.Marshal("")
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal ExecCredential.")
	}

	if shouldCache := os.Getenv("KUBERNETES_AUTHENTICATION_CACHE"); shouldCache != "false" && shouldCache != "0" {
		err = refreshCache(token)
		if err != nil {
			log.Error().Err(err).Msg("Could not refresh cache.")
		}
	}

	log.Info().Msg("End of authentication.")
	fmt.Printf("%s\n", token)

	return nil
}

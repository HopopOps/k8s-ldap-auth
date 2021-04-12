package client

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/adrg/xdg"

	client "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

func getCacheDirPath() string {
	dir := path.Join(xdg.CacheHome, "k8s-ldap-auth")
	return dir
}

func getCacheFilePath() string {
	file := path.Join(getCacheDirPath(), "token")
	return file
}

func getCachedToken() []byte {
	token, err := ioutil.ReadFile(getCacheFilePath())
	if err != nil {
		return nil
	}

	return token
}

func refreshCache(ec client.ExecCredential) error {
	if err := os.MkdirAll(getCacheDirPath(), 0750); err != nil {
		return err
	}

	data, err := json.Marshal(ec)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(getCacheFilePath(), data, 0640)
	if err != nil {
		return err
	}

	return nil
}

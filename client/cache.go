package client

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/adrg/xdg"
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

func refreshCache(data []byte) error {
	if err := os.MkdirAll(getCacheDirPath(), 0750); err != nil {
		return err
	}

	if err := ioutil.WriteFile(getCacheFilePath(), data, 0640); err != nil {
		return err
	}

	return nil
}

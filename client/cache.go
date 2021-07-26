package client

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/adrg/xdg"
)

func GetCacheDirPath() string {
	dir := path.Join(xdg.CacheHome, "k8s-ldap-auth")
	return dir
}

func GetCacheFilePath() string {
	file := path.Join(GetCacheDirPath(), "token")
	return file
}

func getCachedToken() []byte {
	token, err := ioutil.ReadFile(GetCacheFilePath())
	if err != nil {
		return nil
	}

	return token
}

func refreshCache(data []byte) error {
	if err := os.MkdirAll(GetCacheDirPath(), 0750); err != nil {
		return err
	}

	if err := ioutil.WriteFile(GetCacheFilePath(), data, 0640); err != nil {
		return err
	}

	return nil
}

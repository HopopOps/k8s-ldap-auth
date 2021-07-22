package types

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/rs/zerolog/log"
)

func GenerateKey() (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return key, nil
}

var (
	ErrPrivKeyNotFound    = errors.New("No RSA private key found")
	ErrPrivKeyNotReadable = errors.New("Unable to parse private key")
	ErrPubKeyNotFound     = errors.New("No RSA private key found")
	ErrPubKeyNotReadable  = errors.New("Unable to parse public key")
)

// The following is heavily inspired from https://gist.github.com/jshap70/259a87a7146393aab5819873a193b88c
func LoadKey(rsaPrivateKeyLocation, rsaPublicKeyLocation string) (*rsa.PrivateKey, error) {
	priv, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		log.Error().Msg("Private key file was not found.")
		return nil, ErrPrivKeyNotFound
	}

	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte
	if privPem.Type != "RSA PRIVATE KEY" {
		log.Warn().Str("pem_type", privPem.Type).Msg("RSA private key has the wrong type")
	}
	privPemBytes = privPem.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		log.Error().Err(err).Msg("Could not parse to PKCS1 key.")
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			log.Error().Err(err).Msg("Could not parse to PKCS8 key.")
			return nil, ErrPrivKeyNotReadable
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Error().Msg("Could not parse to PKCS8 key.")
		return nil, ErrPrivKeyNotReadable
	}

	pub, err := ioutil.ReadFile(rsaPublicKeyLocation)
	if err != nil {
		log.Error().Msg("Public key file was not found.")
		return nil, ErrPubKeyNotFound
	}

	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {
		log.Error().Msg("Could not decode pem public key.")
		return nil, ErrPubKeyNotReadable
	}

	if pubPem.Type != "PUBLIC KEY" {
		log.Error().Str("pem_type", pubPem.Type).Msg("Public key has the wrong type.")
		return nil, ErrPubKeyNotReadable
	}

	if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
		log.Error().Err(err).Msg("Could not parse to PKIX public key.")
		return nil, ErrPubKeyNotReadable
	}

	var pubKey *rsa.PublicKey
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		log.Error().Err(err).Msg("Could not parse public key to rsa.")
		return nil, ErrPubKeyNotReadable
	}

	privateKey.PublicKey = *pubKey

	return privateKey, nil
}

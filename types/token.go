package types

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/rs/zerolog/log"
)

func Key() (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return key, nil
}

type Token struct {
	token jwt.Token
}

func NewToken(data []byte) *Token {
	now := time.Now()

	t := jwt.New()
	t.Set(jwt.IssuedAtKey, now.Unix())
	t.Set(jwt.ExpirationKey, now.Add(12*time.Hour).Unix())
	t.Set("user", data)

	token := &Token{
		token: t,
	}

	return token
}

func Parse(payload []byte, key *rsa.PrivateKey) (*Token, error) {
	t, err := jwt.Parse(
		payload,
		jwt.WithVerify(jwa.RS256, &key.PublicKey),
		jwt.WithValidate(true),
	)

	if err != nil {
		return nil, err
	}

	token := &Token{
		token: t,
	}

	return token, nil
}

func (t *Token) GetUser() (*User, error) {
	if v, ok := t.token.Get("user"); ok {
		var user User

		log.Debug().Str("data", fmt.Sprintf("%v", v)).Msg("Got user data.")

		data, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(fmt.Sprintf("%v", v))
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &user)
		if err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, fmt.Errorf("Could not get user attribute of jwt token")
}

func (t *Token) IsValid() bool {
	exp, err := t.Expiration()

	if err != nil {
		log.Debug().Str("err", err.Error()).Msg("token validation")
	} else {
		log.Debug().Str("exp", exp.String()).Bool("stillvalid", time.Now().Unix() < exp.Unix()).Msg("token validation")
	}

	return err == nil && time.Now().Unix() < exp.Unix()
}

func (t *Token) Expiration() (time.Time, error) {
	if v, ok := t.token.Get(jwt.ExpirationKey); ok {
		return v.(time.Time), nil
	}

	return time.Time{}, fmt.Errorf("Could not get jwt expiration time")
}

func (t *Token) Payload(key *rsa.PrivateKey) ([]byte, error) {
	signed, err := jwt.Sign(t.token, jwa.RS256, key)
	if err != nil {
		return nil, err
	}

	return signed, nil
}

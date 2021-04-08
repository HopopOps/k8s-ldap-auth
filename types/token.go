package types

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
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

		err := json.Unmarshal([]byte(fmt.Sprintf("%v", v)), &user)
		if err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, fmt.Errorf("Could not get user attribute of jwt token")
}

func (t *Token) IsValid() bool {
	exp, err := t.Expiration()
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

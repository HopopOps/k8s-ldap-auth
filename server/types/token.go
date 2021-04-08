package types

import (
	"crypto/rand"
	"crypto/rsa"
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

func (t *Token) IsValid() bool {
	return true
}

func (t *Token) Groups() []string {
	return nil
}

func (t *Token) Uid() string {
	return ""
}

func (t *Token) User() string {
	return ""
}

func (t *Token) Expiration() time.Time {
	return time.Now()
}

func (t *Token) Payload(key *rsa.PrivateKey) ([]byte, error) {
	signed, err := jwt.Sign(t.token, jwa.RS256, key)
	if err != nil {
		return nil, err
	}

	return signed, nil
}

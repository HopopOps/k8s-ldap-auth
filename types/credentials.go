package types

type Credentials struct {
	Username string
	Password string
}

func (c *Credentials) IsValid() bool {
	return len(c.Username) != 0 && len(c.Password) != 0
}

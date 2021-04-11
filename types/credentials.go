package types

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *Credentials) IsValid() bool {
	return len(c.Username) != 0 && len(c.Password) != 0
}

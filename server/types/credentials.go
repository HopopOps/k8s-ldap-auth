package types

type Credentials struct {
	Target   string
	Username string
	Password string
}

func (c *Credentials) IsValid( /* configuration */ ) bool {
	// TODO: check if group exist in configuration
	return len(c.Target) != 0 && len(c.Username) != 0 && len(c.Password) != 0
}

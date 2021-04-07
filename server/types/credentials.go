package types

type Credentials struct {
	Target   string
	Username string
	Password string
}

func (*Credentials) IsValid() bool {
	return true // TODO: implement
}

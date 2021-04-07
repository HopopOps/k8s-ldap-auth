package types

type Token struct {
	token string
}

func (*Token) IsValid() bool {
	return true // TODO: implement
}

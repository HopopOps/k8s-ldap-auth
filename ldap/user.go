package ldap

type User struct {
	Uid    string
	DN     string
	Groups []string
}

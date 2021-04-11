package types

type User struct {
	Uid    string   `json:"uid"`
	DN     string   `json:"dn"`
	Groups []string `json:"groups"`
}

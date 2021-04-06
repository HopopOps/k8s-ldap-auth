package ldap

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap"
)

type Ldap struct {
	ldapURL          string
	bindDN           string
	bindPassword     string
	searchBase       string
	searchScope      string
	searchFilter     string
	memberOfProperty string
	searchAttributes []string
}

func toLower(a []string) []string {
	var res []string

	for _, item := range a {
		res = append(res, strings.ToLower(item))
	}

	return res
}

func NewInstance(ldapURL, bindDN, bindPassword, searchBase, searchScope, searchFilter, memberOfProperty string, searchAttributes []string) *Ldap {
	s := &Ldap{
		ldapURL:          ldapURL,
		bindDN:           bindDN,
		bindPassword:     bindPassword,
		searchBase:       searchBase,
		searchScope:      searchScope,
		searchFilter:     searchFilter,
		memberOfProperty: memberOfProperty,
	}

	return s
}

func (s *Ldap) Search(username, password string) (*User, error) {
	l, err := ldap.DialURL(s.ldapURL)
	if err != nil {
		return nil, err
	}

	defer l.Close()

	err = l.Bind(s.bindDN, s.bindPassword)
	if err != nil {
		return nil, err
	}

	// Execute LDAP Search request
	searchRequest := ldap.NewSearchRequest(
		s.searchBase,
		scopeMap[s.searchScope],
		ldap.NeverDerefAliases, // Dereference aliases
		0,                      // Size limit (0 = no limit)
		0,                      // Time limit (0 = no limit)
		false,                  // Types only
		fmt.Sprintf(s.searchFilter, username),
		s.searchAttributes,
		nil, // Additional 'Controls'
	)
	result, err := l.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	// If LDAP Search produced a result, return UserInfo, otherwise, return nil
	if len(result.Entries) == 0 {
		return nil, nil
	} else if len(result.Entries) > 1 {
		return nil, fmt.Errorf("Too many entries returned")
	}

	// Bind as the user to verify their password
	err = l.Bind(result.Entries[0].DN, password)
	if err != nil {
		return nil, err
	}

	// Rebinding as the read only user for any further queries is not necessary since the ldap connection will be closed.

	return &User{
		Uid:    strings.ToLower(result.Entries[0].GetAttributeValue("uid")),
		DN:     strings.ToLower(result.Entries[0].DN),
		Groups: toLower(result.Entries[0].GetAttributeValues(s.memberOfProperty)),
	}, nil
}

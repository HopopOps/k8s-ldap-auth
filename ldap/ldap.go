package ldap

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap"
	"github.com/rs/zerolog/log"

	auth "k8s.io/api/authentication/v1"
)

type Ldap struct {
	ldapURL          string
	bindDN           string
	bindPassword     string
	searchBase       string
	searchScope      string
	searchFilter     string
	memberofProperty string
	usernameProperty string
	extraAttributes  []string
	searchAttributes []string
}

func sanitize(a []string) []string {
	var res []string

	for _, item := range a {
		res = append(res, strings.ToLower(item))
	}

	return res
}

func NewInstance(
	ldapURL,
	bindDN,
	bindPassword,
	searchBase,
	searchScope,
	searchFilter,
	memberofProperty,
	usernameProperty string,
	extraAttributes,
	searchAttributes []string,
) *Ldap {
	s := &Ldap{
		ldapURL:          ldapURL,
		bindDN:           bindDN,
		bindPassword:     bindPassword,
		searchBase:       searchBase,
		searchScope:      searchScope,
		searchFilter:     searchFilter,
		memberofProperty: memberofProperty,
		usernameProperty: usernameProperty,
		extraAttributes:  extraAttributes,
		searchAttributes: searchAttributes,
	}

	return s
}

func (s *Ldap) Search(username, password string) (*auth.UserInfo, error) {
	l, err := ldap.DialURL(s.ldapURL)
	if err != nil {
		return nil, err
	}
	log.Debug().Msg("Successfully dialed ldap.")

	defer l.Close()

	err = l.Bind(s.bindDN, s.bindPassword)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("Successfully authenticated to ldap.")

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
		return nil, fmt.Errorf("User not found")
	} else if len(result.Entries) > 1 {
		return nil, fmt.Errorf("Too many entries returned")
	}

	// Bind as the user to verify their password
	err = l.Bind(result.Entries[0].DN, password)
	if err != nil {
		return nil, err
	}

	var extra map[string]auth.ExtraValue

	for _, item := range s.extraAttributes {
		extra[item] = result.Entries[0].GetAttributeValues(item)
	}

	user := &auth.UserInfo{
		UID:      strings.ToLower(result.Entries[0].DN),
		Username: strings.ToLower(result.Entries[0].GetAttributeValue(s.usernameProperty)),
		Groups:   sanitize(result.Entries[0].GetAttributeValues(s.memberofProperty)),
		Extra:    extra,
	}

	log.Debug().Str("uid", user.UID).Strs("groups", user.Groups).Str("username", user.Username).Msg("Research returned a result.")

	return user, nil
}

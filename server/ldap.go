package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"gopkg.in/ldap.v3"
)

type LDAPConnInfo struct {
	HOST            string
	PORT            int
	BASE_DN         string
	USER_DN         string
	USER_RDN_ATTR   string
	USER_LOGIN_ATTR string
	ADMIN_GROUP     string
	GROUP_FILTER    string
}

var LDAPConfig = LDAPConnInfo{
	`keymaker.lab.cucyber.net`,
	636,
	`dc=lab,dc=cucyber,dc=net`,
	`cn=users,cn=accounts`,
	`uid`,
	`uid`,
	`admins`,
	`(&(cn=*)(memberUid=%s))`,
}

func LDAPError(err error) string {
	var ldapError *ldap.Error
	if errors.As(err, &ldapError) {
		switch ldapError.ResultCode {
		case ldap.LDAPResultInvalidCredentials:
			return "Invalid credentials."
		default:
			return "LDAP Error."
		}
	}
	return ""
}

func LDAPIsAdmin(username string) bool {
	l, err := ldap.DialTLS("tcp",
		fmt.Sprintf("%s:%d", LDAPConfig.HOST, LDAPConfig.PORT),
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	searchRequest := ldap.NewSearchRequest(
		LDAPConfig.BASE_DN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(LDAPConfig.GROUP_FILTER, username),
		[]string{"dn", "cn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		panic(err)
	}

	for _, entry := range sr.Entries {
		if entry.GetAttributeValue("cn") == LDAPConfig.ADMIN_GROUP {
			return true
		}
	}

	return false
}

func LDAPAuthenticateSearchBind(username, password string) error {
	return errors.New("LDAP Authenticate via Search Bind: Not Implemented.")
}

func LDAPAuthenticateDirectBind(username, password string) error {
	l, err := ldap.DialTLS("tcp",
		fmt.Sprintf("%s:%d", LDAPConfig.HOST, LDAPConfig.PORT),
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	err = l.Bind(
		fmt.Sprintf("%s=%s,%s,%s",
			LDAPConfig.USER_RDN_ATTR,
			username,
			LDAPConfig.USER_DN,
			LDAPConfig.BASE_DN,
		), password,
	)
	if err != nil {
		return err
	}

	return nil
}

func LDAPAuthenticate(username, password string) error {
	if LDAPConfig.USER_RDN_ATTR == LDAPConfig.USER_LOGIN_ATTR {
		return LDAPAuthenticateDirectBind(username, password)
	} else {
		return LDAPAuthenticateSearchBind(username, password)
	}
}

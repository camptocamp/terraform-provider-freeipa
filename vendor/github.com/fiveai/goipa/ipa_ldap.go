package ipa

import (
	"gopkg.in/ldap.v2"
	"fmt"
	"crypto/tls"
	"errors"
)

type LdapClient struct {
	BaseDN     string
	Connection *ldap.Conn
}

func LdapConnect(host string, baseDn string, username string, password string) (*LdapClient, error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, 389))
	if err != nil {
		return nil, err
	}

	// TODO: Fix needing InsecureSkipVerify
	config := &tls.Config{ServerName: host, RootCAs: ipaCertPool, InsecureSkipVerify: true}
	config.BuildNameToCertificate()

	// Reconnect with TLS
	err = l.StartTLS(config)
	if err != nil {
		l.Close()
		return nil, err
	}

	err = l.Bind(fmt.Sprintf("uid=%s,cn=users,cn=accounts,%s", username, baseDn), password)
	if err != nil {
		l.Close()
		return nil, err
	}

	return &LdapClient{BaseDN: baseDn, Connection: l}, nil
}

func (c *LdapClient) Close() {
	c.Connection.Close()
}

func (c *LdapClient) Search(childDn string, filter string, attributes []string) (*ldap.SearchResult, error) {
	dn := c.BaseDN
	if len(childDn) > 0 {
		dn = fmt.Sprintf("%s,%s", childDn, c.BaseDN)
	}

	searchRequest := ldap.NewSearchRequest(
		dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil,
	)

	return c.Connection.Search(searchRequest)
}

func (c *LdapClient) GetUserForUUID(uuid string) (*string, error) {
	sr, err := c.Search("cn=users,cn=accounts",
		fmt.Sprintf("(ipaUniqueID=%s)", uuid),
		[]string{"uid"})

	if err != nil {
		return nil, err
	}

	if len(sr.Entries) != 1 {
		return nil, errors.New("too many entries returned")
	}

	uid := sr.Entries[0].GetAttributeValue("uid")

	return &uid, nil
}


func (c *LdapClient) UserExistsForUUID(uuid string) (bool, error) {
	sr, err := c.Search("cn=users,cn=accounts",
		fmt.Sprintf("(ipaUniqueID=%s)", uuid),
		[]string{})

	if err != nil {
		return false, err
	}

	if len(sr.Entries) > 1 {
		return false, errors.New("too many entries returned")
	}

	return len(sr.Entries) == 1, nil
}

func (c *LdapClient) GetGroupForUUID(uuid string) (*string, error) {
	sr, err := c.Search("cn=groups,cn=accounts",
		fmt.Sprintf("(ipaUniqueID=%s)", uuid),
		[]string{"cn"})

	if err != nil {
		return nil, err
	}

	if len(sr.Entries) != 1 {
		return nil, errors.New("too many entries returned")
	}

	gid := sr.Entries[0].GetAttributeValue("cn")

	return &gid, nil
}

func (c *LdapClient) GroupExistsForUUID(uuid string) (bool, error) {
	sr, err := c.Search("cn=groups,cn=accounts",
		fmt.Sprintf("(ipaUniqueID=%s)", uuid),
		[]string{})

	if err != nil {
		return false, err
	}

	if len(sr.Entries) > 1 {
		return false, errors.New("too many entries returned")
	}

	return len(sr.Entries) == 1, nil
}


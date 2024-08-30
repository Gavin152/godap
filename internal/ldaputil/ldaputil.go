package ldaputil

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

// LDAPServerConfig holds the LDAP server configuration
type LDAPServerConfig struct {
	Protocol string
	URL      string
	Port     int
	BaseDN   string
	BindUser string
	BindPass string
}

// DefaultConfig returns a default LDAP server configuration
func DefaultConfig() LDAPServerConfig {
	return LDAPServerConfig{
		Protocol: "ldap",
		URL:      "10.21.0.10",                                          // Replace with your LDAP server URL
		Port:     389,                                                   // Use 636 for LDAPS (secure)
		BaseDN:   "DC=development,DC=prolion",                           // Replace with your Base DN
		BindUser: "CN=Administrator,CN=Users,DC=development,DC=prolion", // Replace with bind DN
		BindPass: "DLism3xU!",                                           // Replace with the bind password
	}
}

// AuthenticateUser authenticates the user with the provided username and password
func AuthenticateUser(username, password string) (bool, error) {
	config := DefaultConfig()

	// Connect to the LDAP server
	l, err := ldap.DialURL(fmt.Sprintf("%s://%s:%d", config.Protocol, config.URL, config.Port), ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		log.Printf("Failed to connect to LDAP server: %v", err)
		return false, err
	}
	defer l.Close()

	// Bind with service account credentials
	err = l.Bind(config.BindUser, config.BindPass)
	if err != nil {
		log.Printf("Failed to bind with service account: %v", err)
		return false, err
	}

	// Search for the user DN
	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=user)(userPrincipalName=%s))", username), // sAMAccountName is the AD username
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Printf("Failed to search for user: %v", err)
		return false, err
	}

	if len(sr.Entries) != 1 {
		return false, fmt.Errorf("user not found or too many entries returned: %d", len(sr.Entries))
	}

	userDN := sr.Entries[0].DN

	// Attempt to bind as the user
	err = l.Bind(userDN, password)
	if err != nil {
		log.Printf("Failed to authenticate user: %v", err)
		return false, err
	}

	// Authentication successful
	return true, nil
}

// GetUserGroups fetches the groups for the given username
func GetUserGroups(username string) ([]string, error) {
	config := DefaultConfig()

	// Connect to the LDAP server
	l, err := ldap.DialURL(fmt.Sprintf("%s://%s:%d", config.Protocol, config.URL, config.Port), ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		log.Printf("Failed to connect to LDAP server: %v", err)
		return nil, err
	}
	defer l.Close()

	// Bind with service account credentials
	err = l.Bind(config.BindUser, config.BindPass)
	if err != nil {
		log.Printf("Failed to bind with service account: %v", err)
		return nil, err
	}

	// Search for the user DN
	searchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", username),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Printf("Failed to search for user: %v", err)
		return nil, err
	}

	if len(sr.Entries) != 1 {
		return nil, fmt.Errorf("user not found or too many entries returned")
	}

	userDN := sr.Entries[0].DN

	// Search for groups the user is a member of
	groupSearchRequest := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=group)(member=%s))", userDN),
		[]string{"cn"},
		nil,
	)

	groupResult, err := l.Search(groupSearchRequest)
	if err != nil {
		log.Printf("Failed to search for user groups: %v", err)
		return nil, err
	}

	// Collect group names
	var groups []string
	for _, entry := range groupResult.Entries {
		groups = append(groups, entry.GetAttributeValue("cn"))
	}

	return groups, nil
}

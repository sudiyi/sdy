package utils

import (
	"net/url"
	"strings"
)

// Dsn Parse
// parse: redis://[:password@]host[:port][/db-number][?option=value]
// eg:    redis://localhost:6379/10
//		  redis://:password@localhost:6379/0

// parse: mysql://[username:password@]host[:port]/db-number
func DsnParse(dns string) (string, string, string, error) {
	password := ""
	db := ""
	u, err := url.Parse(dns)
	if err != nil {
		return "", "", "", err
	}
	if u.User != nil {
		p, _ := u.User.Password()
		password = p
	}
	if u.Path != "" {
		arr := strings.Split(u.Path, "/")
		db = arr[1]
	}
	return u.Host, password, db, nil
}

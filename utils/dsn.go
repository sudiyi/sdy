package utils

import (
	"net/url"
	"strings"
)

// Dsn Parse
// parse: redis://[:password@]host[:port][/db-number][?option=value]
// eg:    redis://localhost:6379/10
//		  redis://:password@localhost:6379/0

// parse: mysql://[username:password@]host[:port]/db
func DsnParse(dns string) (string, string, string, error) {
	var password, db string
	u, err := url.Parse(dns)
	if err != nil {
		return "", "", "", err
	}

	if u.User != nil {
		password, _ = u.User.Password()
	}

	arr := strings.Split(u.Path, "/")
	if len(arr) >= 1 {
		db = arr[1]
	}

	return u.Host, password, db, nil
}

package utils

import (
	"net/url"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// Dsn Parse
// parse: redis://[:password@]host[:port][/db-number][?option=value]
// eg:    redis://localhost:6379/10
//		  redis://:password@localhost:6379/0
func DsnParse(dsn string) (string, string, string, error) {
	var password, db string
	u, err := url.Parse(dsn)
	if err != nil {
		return "", "", "", err
	}

	if u.User != nil {
		password, _ = u.User.Password()
	}

	arr := strings.Split(u.Path, "/")
	if len(arr) >= 2 {
		db = arr[1]
	}

	return u.Host, password, db, nil
}

func MysqlDsnParse(dsn string) (string, string, string, error) {
	my, err := mysql.ParseDSN(dsn)
	if nil != err {
		return "", "", "", err
	}
	return my.Addr, my.Passwd, my.DBName, err
}

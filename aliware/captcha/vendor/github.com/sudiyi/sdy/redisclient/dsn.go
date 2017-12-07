package redisclient

import (
	"github.com/huhongda/GoToolBox/utils"
	"net/url"
	"strings"
)

// Dsn Parse
// parse: redis://[:password@]host[:port][/db-number][?option=value]
// eg:    redis://localhost:6379/10
//		  redis://:password@localhost:6379/0
func dsnParse(dns string) (string, string, int) {
	password := ""
	db := 0
	u, err := url.Parse(dns)
	if err != nil {
		panic(err)
	}
	if u.User != nil {
		p, _ := u.User.Password()
		password = p
	}
	if u.Path != "" {
		arr := strings.Split(u.Path, "/")
		db = utils.StringToInt(arr[1])
	}
	return u.Host, password, db
}

# Query String Convert

you can convert query string for struct or json string

## Usage


```
import "github.com/sudiyi/sdy/utils/querystring"

fmt.Println(querystring.New(`{"query":"bicycle", "size": "50x50", "china": "中国", "baidu": "%!中国"}`).Build())
```



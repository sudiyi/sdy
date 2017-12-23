package querystring

import (
	"fmt"
	"testing"
)

type TestStruct struct {
	Query string
	Size  string
	China string
	BaiDu string
}

func TestQueryString_Build(t *testing.T) {
	fmt.Println(New(`{"query":"bicycle", "size": "50x50", "china": "中国", "baidu": "%!中国"}`).Build())
	s := TestStruct{Query: "bicycle", Size: "50x50", China: "中国", BaiDu: "%!中国"}
	fmt.Println(New(s).Build())

	uri := "baidu=%25%21%E4%B8%AD%E5%9B%BD&china=%E4%B8%AD%E5%9B%BD&query=bicycle&size=50x50"
	fmt.Println(New(uri).Build())
}

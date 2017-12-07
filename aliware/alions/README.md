# AliOns

this is aliware http ons client packages for golang (NOTE: will be remove on 2018.04.01)

## Usage

```
import(
    "github.com/sudiyi/sdy/aliware/alions"
)

alions.AlionsConfig.HeaderTimeout
alions.CurrentTimeForMillisSecond

responseBody := alions.Post(string(bodyJson), alions.AlionsConfig.Tag, key)

```

## References

[HTTP下线通告](https://help.aliyun.com/document_detail/61438.html?spm=5176.doc29532.6.604.rJYdrK)
package querystring

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type QueryString struct {
	QueryData url.Values
	Errors    []error
	Content   interface{}
}

func New(content interface{}) *QueryString {
	return &QueryString{QueryData: url.Values{}, Errors: nil, Content: content}
}

func (s *QueryString) query() *QueryString {
	switch v := reflect.ValueOf(s.Content); v.Kind() {
	case reflect.String:
		s.queryString(v.String())
	case reflect.Struct:
		s.queryStruct(v.Interface())
	case reflect.Map:
		s.queryMap(v.Interface())
	default:
	}
	return s
}

func (s *QueryString) queryString(content string) *QueryString {
	var val map[string]string
	if err := json.Unmarshal([]byte(content), &val); err == nil {
		for k, v := range val {
			s.QueryData.Add(k, v)
		}
	} else {
		if queryData, err := url.ParseQuery(content); err == nil {
			for k, queryValues := range queryData {
				for _, queryValue := range queryValues {
					s.QueryData.Add(k, string(queryValue))
				}
			}
		} else {
			s.Errors = append(s.Errors, err)
		}
	}
	return s
}

func (s *QueryString) queryStruct(content interface{}) *QueryString {
	if marshalContent, err := json.Marshal(content); err != nil {
		s.Errors = append(s.Errors, err)
	} else {
		var val map[string]interface{}
		if err := json.Unmarshal(marshalContent, &val); err != nil {
			s.Errors = append(s.Errors, err)
		} else {
			for k, v := range val {
				k = strings.ToLower(k)
				var queryVal string
				switch t := v.(type) {
				case string:
					queryVal = t
				case float64:
					queryVal = strconv.FormatFloat(t, 'f', -1, 64)
				case time.Time:
					queryVal = t.Format(time.RFC3339)
				default:
					j, err := json.Marshal(v)
					if err != nil {
						continue
					}
					queryVal = string(j)
				}
				s.QueryData.Add(k, queryVal)
			}
		}
	}
	return s
}

func (s *QueryString) queryMap(content interface{}) *QueryString {
	return s.queryStruct(content)
}

func (s *QueryString) Build() string {
	s.query()
	var req *http.Request
	req, _ = http.NewRequest("", "", nil)
	q := req.URL.Query()
	for k, v := range s.QueryData {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}
	return q.Encode()
}

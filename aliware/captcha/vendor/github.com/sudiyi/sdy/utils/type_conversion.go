package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func IntToBool(i int8) bool {
	if i >= 1 {
		return true
	}
	return false
}

func JsonToMap(str string) (map[string]interface{}, error) {
	parsed := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &parsed)
	return parsed, err
}

func EachMapToInt(str []string) []int {
	var newInt = []int{}
	for _, i := range str {
		j, _ := strconv.Atoi(i)
		newInt = append(newInt, j)
	}
	return newInt
}

func StringToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func InterfaceToString(value interface{}) (string, error) {
	switch value := value.(type) {
	case nil:
		return "NULL", nil
	case int, uint:
		return fmt.Sprintf("%d", value), nil
	case float64, float32:
		return fmt.Sprintf("%d", int(value.(float64))), nil
	case bool:
		if value {
			return "TRUE", nil
		}
		return "FALSE", nil
	case string:
		return fmt.Sprintf("%s", value), nil
	default:
		return ``, fmt.Errorf("not in array")
	}
}

func StructToMap(in interface{}, tag string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tagv := fi.Tag.Get(tag); tagv != "" {
			// set key of map to value in struct field
			out[tagv] = v.Field(i).Interface()
		}
	}
	return out, nil
}

func FloatToString(i float64) string {
	if i*100 == math.Floor(i*100) {
		return strconv.FormatFloat(i, 'f', 2, 64)
	}
	return strconv.FormatFloat(i, 'f', -1, 64)
}

func FloatToCeil(num float64) int {
	return int(math.Ceil(num))
}

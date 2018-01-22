package utils

func GenerateFormatResponds(code int, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{"data": data, "code": code, "message": message}
}

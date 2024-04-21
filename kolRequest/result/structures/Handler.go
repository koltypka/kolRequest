package structures

import (
	"encoding/xml"
	"reflect"

	"github.com/koltypka/kolRequest/kolRequest/result/structures/types"
)

func Handler(currentType string, Body []byte) any {
	switch currentType {
	case "RSS":
		return handleRss(Body)
	default:
		return nil
	}
}

func handleRss(Body []byte) any {
	result := anyStructHandle(types.Channel{}, Body)
	return structToMap(result)
}

func structToMap(s interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			fieldValue := val.Field(i).Interface()
			result[field.Name] = fieldValue
		}
	}

	return result
}

func anyStructHandle[T any](XmlStruct T, Body []byte) T {
	xml.Unmarshal(Body, &XmlStruct)

	return XmlStruct
}

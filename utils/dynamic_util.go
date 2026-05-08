package utils

import (
	"fmt"
	"reflect"
)

func GetNestedValue(data map[string]interface{}, path string) interface{} {
	parts := splitPath(path)
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			return current[part]
		}
		if val, ok := current[part]; ok {
			if m, ok := val.(map[string]interface{}); ok {
				current = m
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return nil
}

func SetNestedValue(data map[string]interface{}, path string, value interface{}) {
	parts := splitPath(path)
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			current[part] = value
			return
		}
		if _, ok := current[part]; !ok {
			current[part] = map[string]interface{}{}
		}
		current = current[part].(map[string]interface{})
	}
}

func splitPath(path string) []string {
	var parts []string
	current := ""
	inBracket := false

	for _, c := range path {
		if c == '[' && !inBracket {
			inBracket = true
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else if c == ']' && inBracket {
			inBracket = false
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else if c == '.' && !inBracket {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		var result int
		_, err := fmt.Sscanf(v, "%d", &result)
		return result, err
	default:
		return 0, fmt.Errorf("cannot convert %v to int", reflect.TypeOf(value))
	}
}

func ToFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		var result float64
		_, err := fmt.Sscanf(v, "%f", &result)
		return result, err
	default:
		return 0, fmt.Errorf("cannot convert %v to float64", reflect.TypeOf(value))
	}
}

func ToBool(value interface{}) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return v == "true" || v == "1", nil
	case int:
		return v != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %v to bool", reflect.TypeOf(value))
	}
}

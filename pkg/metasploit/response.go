package metasploit

import (
	"fmt"
	"strconv"
)

func readString(resp map[string]interface{}, key string, required bool) (string, error) {
	v, ok := resp[key]
	if !ok {
		if required {
			return "", fmt.Errorf("missing %q in response", key)
		}
		return "", nil
	}

	switch typed := v.(type) {
	case string:
		return typed, nil
	case []byte:
		return string(typed), nil
	default:
		return "", fmt.Errorf("field %q has invalid type %T", key, v)
	}
}

func readInt(resp map[string]interface{}, key string, required bool) (int64, error) {
	v, ok := resp[key]
	if !ok {
		if required {
			return 0, fmt.Errorf("missing %q in response", key)
		}
		return 0, nil
	}

	switch typed := v.(type) {
	case int:
		return int64(typed), nil
	case int8:
		return int64(typed), nil
	case int16:
		return int64(typed), nil
	case int32:
		return int64(typed), nil
	case int64:
		return typed, nil
	case uint:
		return int64(typed), nil
	case uint8:
		return int64(typed), nil
	case uint16:
		return int64(typed), nil
	case uint32:
		return int64(typed), nil
	case uint64:
		if typed > uint64(^uint64(0)>>1) {
			return 0, fmt.Errorf("field %q overflows int64", key)
		}
		return int64(typed), nil
	case float32:
		return int64(typed), nil
	case float64:
		return int64(typed), nil
	case string:
		parsed, err := strconv.ParseInt(typed, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("field %q has invalid int string %q", key, typed)
		}
		return parsed, nil
	case []byte:
		parsed, err := strconv.ParseInt(string(typed), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("field %q has invalid int bytes %q", key, string(typed))
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("field %q has invalid type %T", key, v)
	}
}

func checkRPCFailure(resp map[string]interface{}) error {
	if result, ok := resp["result"]; ok && fmt.Sprintf("%v", result) == "failure" {
		if msg, exists := resp["error_message"]; exists {
			return fmt.Errorf("%v", msg)
		}
		return fmt.Errorf("rpc call returned failure")
	}
	if msg, ok := resp["error_message"]; ok {
		return fmt.Errorf("%v", msg)
	}
	return nil
}

package utils

import (
	"fmt"
)

func NumericInterfToInt(a interface{}) (int64, error) {
	switch a.(type) {
	case float32:
		val, _ := a.(float32)
		return int64(val), nil
	case float64:
		val, _ := a.(float64)
		return int64(val), nil
	case int:
		val, _ := a.(int)
		return int64(val), nil
	case int8:
		val, _ := a.(int8)
		return int64(val), nil
	case int16:
		val, _ := a.(int16)
		return int64(val), nil
	case int32:
		val, _ := a.(int32)
		return int64(val), nil
	case int64:
		val, _ := a.(int64)
		return val, nil
	case uint:
		val, _ := a.(uint)
		return int64(val), nil
	case uint8:
		val, _ := a.(uint8)
		return int64(val), nil
	case uint16:
		val, _ := a.(uint16)
		return int64(val), nil
	case uint32:
		val, _ := a.(uint32)
		return int64(val), nil
	case uint64:
		val, _ := a.(uint64)
		return int64(val), nil
	default:
		err := fmt.Errorf("cannot convert interface to any numeric type")
		return 0, err
	}
}

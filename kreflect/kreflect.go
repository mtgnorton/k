// kreflect 提供了一些对reflect.Value类型的一些常用操作
//
// 主要功能:
//   - IsNil: 判断任意类型是否为nil
//   - ToString: 将任意类型转换为string类型
package kreflect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// IsNil 判断一个值是否为nil
//
// 参数说明:
//   - a: 任意类型的值(any)
//
// 返回值说明:
//   - bool: 如果值为nil返回true,否则返回false
//
// 注意事项:
//   - 支持判断基础类型及引用类型(chan/map/slice/func/interface/pointer等)
//   - 对于基础类型(int/string等)始终返回false
//   - 可以处理reflect.Value类型的输入
func IsNil(a any) bool {
	if a == nil {
		return true
	}
	var rv reflect.Value
	if v, ok := a.(reflect.Value); ok {
		rv = v
	} else {
		rv = reflect.ValueOf(a)
	}
	switch rv.Kind() {
	case reflect.Chan,
		reflect.Map,
		reflect.Slice,
		reflect.Func,
		reflect.Interface,
		reflect.UnsafePointer,
		reflect.Ptr:
		return !rv.IsValid() || rv.IsNil()
	default:
		return false
	}
}

// ToString 将任意类型转换为string类型
//
// 参数说明:
//   - a: 任意类型的值(any)
//
// 返回值说明:
//   - string: 转换后的字符串
func ToString(a any) string {
	if a == nil {
		return ""
	}
	switch value := a.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	default:
		// Empty checks.
		if value == nil {
			return ""
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return ToString(rv.Elem().Interface())
		}
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

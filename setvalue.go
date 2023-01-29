package json2struct

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func setValue(rv reflect.Value, v interface{}) error {
	switch rv.Kind() {
	case reflect.Int64:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int:
		i, e := toInt64(v)
		if e != nil {
			return e
		}
		rv.SetInt(i)

	case reflect.Uint64:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint:
		u, e := toUint64(v)
		if e != nil {
			return e
		}
		rv.SetUint(u)

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		f, e := toFloat(v)
		if e != nil {
			return e
		}
		rv.SetFloat(f)

	case reflect.String:
		s, e := toString(v)
		if e != nil {
			return e
		}
		rv.SetString(s)

	default:
		var kindStr string
		if rv.Kind() == reflect.Ptr {
			kindStr = fmt.Sprintf("ptr of %s", rv.Elem().Kind())
		} else {
			kindStr = fmt.Sprintf("%s", rv.Kind())
		}
		return errors.New(fmt.Sprintf("unsupport set value(%v) to type(%s)", v, kindStr))
	}

	return nil
}

func toInt64(v interface{}) (int64, error) {
	switch v.(type) {
	case json.Number:
		return v.(json.Number).Int64()
	case string:
		i, e := strconv.ParseInt(v.(string), 10, 0)
		return i, e
	case int64:
		return v.(int64), nil
	case int32:
		return int64(v.(int32)), nil
	case int16:
		return int64(v.(int16)), nil
	case int8:
		return int64(v.(int8)), nil
	case int:
		return int64(v.(int)), nil
	default:
		return 0, errors.New(fmt.Sprintf("toInt64 failed! 类型错误，value:%v", v))
	}
}

func toUint64(v interface{}) (uint64, error) {
	switch v.(type) {
	case json.Number:
		i, e := v.(json.Number).Int64()
		return uint64(i), e
	case string:
		u, e := strconv.ParseUint(v.(string), 10, 0)
		return u, e
	case uint64:
		return v.(uint64), nil
	case uint32:
		return uint64(v.(uint32)), nil
	case uint16:
		return uint64(v.(uint16)), nil
	case uint8:
		return uint64(v.(uint8)), nil
	case uint:
		return uint64(v.(uint)), nil
	default:
		return 0, errors.New(fmt.Sprintf("toUint64 failed! type error，value:%v", v))
	}
}

func toFloat(v interface{}) (float64, error) {
	switch v.(type) {
	case json.Number:
		return v.(json.Number).Float64()
	case string:
		f, e := strconv.ParseFloat(v.(string), 64)
		return f, e
	case float64:
		return v.(float64), nil
	case float32:
		return float64(v.(float32)), nil
	default:
		return 0, errors.New(fmt.Sprintf("toFloat failed! type error，value:%v", v))
	}
}

func toString(value interface{}) (string, error) {
	var key string
	if value == nil {
		return key, nil
	}
	switch value.(type) {
	case json.Number:
		key = value.(json.Number).String()
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		return "", errors.New(fmt.Sprintf("toString failed! type error，value:%v", value))
	}
	return key, nil
}

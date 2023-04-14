package util

import (
	. "github.com/infrago/base"
)

// MapString 从 map 读字串
func MapString(data Map, keys ...string) string {
	for _, key := range keys {
		if vv, ok := data[key].(string); ok {
			return vv
		}
	}
	return ""
}

func MapFloat(data Map, keys ...string) float64 {
	for _, key := range keys {
		if vv, ok := data[key].(float64); ok {
			return vv
		} else if vv, ok := data[key].(int64); ok {
			return float64(vv)
		} else if vv, ok := data[key].(int); ok {
			return float64(vv)
		}
	}
	return float64(0)
}
func MapInt(data Map, keys ...string) int64 {
	for _, key := range keys {
		if vv, ok := data[key].(int64); ok {
			return vv
		} else if vv, ok := data[key].(int); ok {
			return int64(vv)
		} else if vv, ok := data[key].(float64); ok {
			return int64(vv)
		}
	}
	return int64(0)
}

func MapBool(data Map, keys ...string) bool {
	for _, key := range keys {
		if vv, ok := data[key].(bool); ok {
			return vv
		} else if vv, ok := data[key].(int); ok {
			return vv > 0
		} else if vv, ok := data[key].(int64); ok {
			return vv > 0
		} else if vv, ok := data[key].(float64); ok {
			return vv > 0
		} else if vv, ok := data[key].(string); ok {
			return vv == "T" || vv == "t" || vv == "true" || vv == "TRUE"
		}
	}
	return false
}

// DeepMapping
// 深度复制Map对象
func DeepMapping(src Map, dests ...Map) Map {
	var dest Map
	if len(dests) > 0 {
		dest = dests[0]
	} else {
		dest = Map{}
	}
	for key, val := range src {
		switch value := val.(type) {
		case Map:
			dest[key] = DeepMapping(value)
		case []Map:
			temps := []Map{}
			for _, srcMap := range value {
				temps = append(temps, DeepMapping(srcMap))
			}
			dest[key] = temps
		default:
			dest[key] = value
		}
	}

	return dest
}

//

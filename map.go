package util

import (
	. "github.com/infrago/base"
)

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

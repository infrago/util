package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	. "github.com/infrago/base"
)

func MergeString(args ...string) string {
	buffer := bytes.Buffer{}
	for i := 0; i < len(args); i++ {
		buffer.WriteString(args[i])
	}

	return buffer.String()
}

func Ip2Num(ip string) int64 {
	ipSegs := strings.Split(ip, ".")
	var ipInt int64 = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.ParseInt(ipSeg, 10, 64)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

func Num2Ip(ipInt int64) string {
	ipSegs := make([]string, 4)
	var len int = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < len; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[len-i-1] = strconv.FormatInt(tempInt, 10)
		ipInt = ipInt >> 8
	}
	for i := 0; i < len; i++ {
		buffer.WriteString(ipSegs[i])
		if i < len-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}
func Split(s string) []string {

	s = strings.TrimSpace(s)

	arr := []string{}
	if s != "" {
		if strings.Contains(s, "|") {
			arr = strings.Split(s, "|")
		} else if strings.Contains(s, ";") {
			arr = strings.Split(s, ";")
		} else if strings.Contains(s, ",") {
			arr = strings.Split(s, ",")
		} else if strings.Contains(s, "/") {
			arr = strings.Split(s, "/")
		} else if strings.Contains(s, "-") {
			arr = strings.Split(s, "-")
		} else if strings.Contains(s, ":") {
			arr = strings.Split(s, ":")
		} else {
			arr = append(arr, s)
		}
	}

	return arr
}

func AnyToString(val Any) string {
	sv := ""
	switch v := val.(type) {
	case string:
		sv = v
	case int:
		sv = strconv.Itoa(v)
	case int64:
		sv = strconv.FormatInt(v, 10)
	case bool:
		sv = strconv.FormatBool(v)
	case Map:
		d, e := json.Marshal(v)
		if e == nil {
			sv = string(d)
		} else {
			sv = "{}"
		}
	case []Map:
		d, e := json.Marshal(v)
		if e == nil {
			sv = string(d)
		} else {
			sv = "[]"
		}
	case []int, []int8, []int16, []int32, []int64, []float32, []float64, []string, []bool, []Any:
		d, e := json.Marshal(v)
		if e == nil {
			sv = string(d)
		} else {
			sv = "[]"
		}
	default:
		sv = fmt.Sprintf("%v", v)
	}

	return sv
}

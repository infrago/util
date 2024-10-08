package util

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	. "github.com/infrago/base"
)

const (
	NoDuration = -99999
)

func ParseDurationConfig(config Map, field string) time.Duration {
	if expire, ok := config[field].(string); ok {
		dur, err := ParseDuration(expire)
		if err == nil {
			return dur
		}
	}
	if expire, ok := config[field].(int); ok {
		return time.Second * time.Duration(expire)
	}
	if expire, ok := config[field].(int64); ok {
		return time.Second * time.Duration(expire)
	}
	if expire, ok := config[field].(float64); ok {
		return time.Second * time.Duration(expire)
	}

	return NoDuration
}

// GenerateUUID is used to generate a random UUID
func UUID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err == nil {
		return fmt.Sprintf(
			"%08x-%04x-%04x-%04x-%12x",
			buf[0:4],
			buf[4:6],
			buf[6:8],
			buf[8:10],
			buf[10:16],
		)
	}
	return "77777777-7777-7777-7777-777777777777"
}

// sha1加密
func Sha1(str string) string {
	sha1Ctx := sha1.New()
	sha1Ctx.Write([]byte(str))
	cipherStr := sha1Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// sha1加密文件
func Sha1File(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := sha1.New()
		if _, e := io.Copy(h, f); e == nil {
			return fmt.Sprintf("%x", h.Sum(nil))
		}
	}
	return ""
}

func Sha1BaseFile(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := sha1.New()
		if _, e := io.Copy(h, f); e == nil {
			return base64.URLEncoding.EncodeToString(h.Sum(nil))
		}
	}
	return ""
}

// sha256加密
func Sha256(str string) string {
	sha256Ctx := sha256.New()
	sha256Ctx.Write([]byte(str))
	cipherStr := sha256Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// sha256加密文件
func Sha256File(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := sha256.New()
		if _, e := io.Copy(h, f); e == nil {
			return fmt.Sprintf("%x", h.Sum(nil))
		}
	}
	return ""
}

func Sha256BaseFile(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := sha256.New()
		if _, e := io.Copy(h, f); e == nil {
			return base64.URLEncoding.EncodeToString(h.Sum(nil))
		}
	}
	return ""
}

// md5加密
func Md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// md5加密文件
func Md5File(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := md5.New()
		if _, e := io.Copy(h, f); e == nil {
			return fmt.Sprintf("%x", h.Sum(nil))
		}
	}
	return ""
}

//密码加密格式
// func passcode(str string) string {
// 	return md5str(str)
// }

var sizeMap = map[string]int64{
	"B": int64(1),
	"K": int64(1024),
	"M": int64(1024 * 1024),
	"G": int64(1024 * 1024 * 1024),
	"T": int64(1024 * 1024 * 1024 * 1024),
	"P": int64(1024 * 1024 * 1024 * 1024 * 1024),
	"E": int64(1024 * 1024 * 1024 * 1024 * 1024 * 1024),
}

func ParseSize(s string) int64 {
	// orig := s
	var d int64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0
	}
	if s == "" {
		return 0
	}

	s = strings.ToUpper(s)

	for s != "" {
		var (
			v, f  int64       // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if !(s[0] == '.' || '0' <= s[0] && s[0] <= '9') {
			return 0
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, scale, s = leadingFraction(s)
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return 0
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0
		}
		u := s[:i]
		s = s[i:]
		unit, ok := sizeMap[u]
		if !ok {
			return 0
		}
		if v > (1<<63-1)/unit {
			// overflow
			return 0
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += int64(float64(f) * (float64(unit) / scale))
			if v < 0 {
				// overflow
				return 0
			}
		}
		d += v
		if d < 0 {
			// overflow
			return 0
		}
	}

	if neg {
		d = -d
	}
	return d
}

//------------------------- 自定义的duration解析，加入d, w 两个 -----

var unitMap = map[string]int64{
	"ns": int64(time.Nanosecond),
	"us": int64(time.Microsecond),
	"µs": int64(time.Microsecond), // U+00B5 = micro symbol
	"μs": int64(time.Microsecond), // U+03BC = Greek letter mu
	"ms": int64(time.Millisecond),
	"s":  int64(time.Second),
	"m":  int64(time.Minute),
	"h":  int64(time.Hour),
	"d":  int64(time.Hour * 24),
	"w":  int64(time.Hour * 24 * 7),
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParseDuration(s string) (time.Duration, error) {
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
	var d int64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, errors.New("time: invalid duration " + orig)
	}
	for s != "" {
		var (
			v, f  int64       // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if !(s[0] == '.' || '0' <= s[0] && s[0] <= '9') {
			return 0, errors.New("time: invalid duration " + orig)
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, errors.New("time: invalid duration " + orig)
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, scale, s = leadingFraction(s)
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return 0, errors.New("time: invalid duration " + orig)
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, errors.New("time: missing unit in duration " + orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, errors.New("time: unknown unit " + u + " in duration " + orig)
		}
		if v > (1<<63-1)/unit {
			// overflow
			return 0, errors.New("time: invalid duration " + orig)
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += int64(float64(f) * (float64(unit) / scale))
			if v < 0 {
				// overflow
				return 0, errors.New("time: invalid duration " + orig)
			}
		}
		d += v
		if d < 0 {
			// overflow
			return 0, errors.New("time: invalid duration " + orig)
		}
	}

	if neg {
		d = -d
	}
	return time.Duration(d), nil
}

var errLeadingInt = errors.New("time: bad [0-9]*") // never printed

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x int64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > (1<<63-1)/10 {
			// overflow
			return 0, "", errLeadingInt
		}
		x = x*10 + int64(c) - '0'
		if x < 0 {
			// overflow
			return 0, "", errLeadingInt
		}
	}
	return x, s[i:], nil
}

// leadingFraction consumes the leading [0-9]* from s.
// It is used only for fractions, so does not return an error on overflow,
// it just stops accumulating precision.
func leadingFraction(s string) (x int64, scale float64, rem string) {
	i := 0
	scale = 1
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if overflow {
			continue
		}
		if x > (1<<63-1)/10 {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		y := x*10 + int64(c) - '0'
		if y < 0 {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return x, scale, s[i:]
}

package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
)

// sha1加密
func SHA1(str string) string {
	sha1Ctx := sha1.New()
	sha1Ctx.Write([]byte(str))
	cipherStr := sha1Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// sha1加密文件
func SHA1File(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := sha1.New()
		if _, e := io.Copy(h, f); e == nil {
			return hex.EncodeToString(h.Sum(nil))
		}
	}
	return ""
}

func SHA1FileBase64(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := sha1.New()
		if _, e := io.Copy(h, f); e == nil {
			return base64.URLEncoding.EncodeToString(h.Sum(nil))
		}
	}
	return ""
}

// SHA256
func SHA256(str string) string {
	sha256Ctx := sha256.New()
	sha256Ctx.Write([]byte(str))
	cipherStr := sha256Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// sha256加密文件
func SHA256File(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := sha256.New()
		if _, e := io.Copy(h, f); e == nil {
			return hex.EncodeToString(h.Sum(nil))
		}
	}
	return ""
}

func SHA256FileBase64(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := sha256.New()
		if _, e := io.Copy(h, f); e == nil {
			return base64.URLEncoding.EncodeToString(h.Sum(nil))
		}
	}
	return ""
}

// MD5
func MD5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// MD5File
func MD5File(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := md5.New()
		if _, e := io.Copy(h, f); e == nil {
			return hex.EncodeToString(h.Sum(nil))
		}
	}
	return ""
}

func MD5FileBase64(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := md5.New()
		if _, e := io.Copy(h, f); e == nil {
			return base64.URLEncoding.EncodeToString(h.Sum(nil))
		}
	}
	return ""
}

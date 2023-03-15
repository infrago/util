package util

import (
	. "github.com/infrago/base"

	"github.com/BurntSushi/toml"
)

// ParseTOML
// 解析toml文本得到配置
func ParseTOML(s string) (Map, error) {
	var config Map
	_, err := toml.Decode(s, &config)
	return config, err
}

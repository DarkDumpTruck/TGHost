package tghost

import (
	"os"
	"strings"
)

type Script struct {
	Name string
	Code string
}

func LoadScriptFromFile(path string) (*Script, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	code := string(buf)
	name := "未命名游戏"
	for _, c := range strings.Split(code, "\n") {
		if strings.HasPrefix(c, "//!name=") {
			name = strings.TrimPrefix(c, "//!name=")
		} else {
			break
		}
	}
	return &Script{
		Name: name,
		Code: code,
	}, nil
}

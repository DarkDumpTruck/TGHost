package tghost

import (
	"fmt"
	"os"
	"strings"
)

type Script struct {
	Name      string
	Code      string
	PlayerNum int
}

func LoadScriptFromFile(path string) (*Script, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	code := string(buf)
	name := "未命名游戏"
	playerNum := 1
	for _, c := range strings.Split(code, "\n") {
		if strings.HasPrefix(c, "//!name=") {
			name = strings.TrimPrefix(c, "//!name=")
		} else if strings.HasPrefix(c, "//!player=") {
			fmt.Sscanf(c, "//!player=%d", &playerNum)
		} else {
			break
		}
	}
	return &Script{
		Name:      name,
		Code:      code,
		PlayerNum: playerNum,
	}, nil
}

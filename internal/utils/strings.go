package utils

import "strings"

func EitherTrimEqual(s string, prefix ...string) (string, bool) {
	// 任一前缀匹配则返回剩余部分
	for _, p := range prefix {
		if strings.TrimSpace(s) == p {
			return "", true
		}
	}
	return s, false
}

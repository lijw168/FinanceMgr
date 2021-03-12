package utils

import (
	"fmt"
	"strings"
)

func GenLog(kv map[string]interface{}) string {
	sli := make([]string, 0, len(kv))
	for k, v := range kv {
		switch v.(type) {
		case string:
			sli = append(sli, fmt.Sprintf("%s:\"%s\"", k, strings.Replace(v.(string), "\"", "\\\"", -1)))
		default:
			sli = append(sli, fmt.Sprintf("%s:%v", k, v))
		}
	}
	return "[" + strings.Join(sli, ", ") + "]"
}

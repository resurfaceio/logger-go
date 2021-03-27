package logger

import (
	"strings"
)

// method for converting message details to string format
func msgStringify(msg [][]string) string {
	stringified := ""
	n := len(msg)
	for i, val := range msg {
		stringified += "[\"" + strings.Join(val, "\", \"") + "\"]"
		if i != n-1 {
			stringified += ","
		}
	}
	stringified = "[" + stringified + "]"
	return stringified
}

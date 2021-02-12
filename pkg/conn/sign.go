package conn

import (
	"strconv"
)

var tmp = map[string]int{}

// generateSign 0-98
func generateSign(connSign string) string {
	d := tmp[connSign]
	if d > 98 || d < 10 {
		tmp[connSign] = 10
		return "10"
	}
	tmp[connSign] = d + 1
	return strconv.Itoa(d + 1)
}

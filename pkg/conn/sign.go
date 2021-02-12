package conn

import (
	"strconv"
)

var tmp = map[string]int{}

// generateSign 0-98
func generateSign(destination string) string {
	d := tmp[destination]
	tmp[destination] = d + 1
	return strconv.Itoa(d + 1)
}

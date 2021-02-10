package broker

import "strconv"

var tmp = map[string][]int{}

func generateSign(destination string) string {
	d := tmp[destination]
	if d == nil {
		tmp[destination][0] = 0
		return "0"
	}

	last := len(d) - 1
	//lVal := d[last]

	if last == 99 {
		tmp[destination] = []int{}
		tmp[destination][0] = 0
		return "0"
	}

	d[last] = 0

	return strconv.Itoa(last + 1)
}

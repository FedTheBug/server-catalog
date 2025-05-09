package utils

import (
	"strconv"
	"strings"
)

func ParseRAMValues(ramStr string) []int {
	ramValues := strings.Split(ramStr, ",")
	result := make([]int, 0, len(ramValues))

	for _, v := range ramValues {
		v = strings.TrimSpace(v)
		v = strings.TrimSuffix(v, "GB")
		if num, err := strconv.Atoi(v); err == nil {
			result = append(result, num)
		}
	}

	return result
}

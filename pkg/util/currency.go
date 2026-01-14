package util

import (
	"fmt"
	"strings"
)

func FormatRupiah(amount float64) string {
	s := fmt.Sprintf("%.2f", amount)

	parts := strings.Split(s, ".")
	intPart := parts[0]
	decimalPart := parts[1]

	var result []string
	n := len(intPart)
	for i, v := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ".")
		}
		result = append(result, string(v))
	}

	return "Rp " + strings.Join(result, "") + "," + decimalPart
}

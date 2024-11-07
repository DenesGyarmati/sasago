package src

import (
	"fmt"
	"strconv"
	"strings"
)

func parseFloat(s string) float64 {
	val, _ := strconv.ParseFloat(trimSymbol(s), 64)
	return val
}

func getInt(s string) (int, error) {
	return strconv.Atoi(trimSymbol(s))
}

func trimSymbol(s string) string {
	return strings.TrimSpace(s)
}
func trimTrailingZeros(value float64) float64 {
	formattedValue := fmt.Sprintf("%.3f", value)
	if formattedValue[len(formattedValue)-1] == '0' {
		formattedValue = formattedValue[:len(formattedValue)-1]
	}
	if formattedValue[len(formattedValue)-1] == '.' {
		formattedValue = formattedValue[:len(formattedValue)-1]
	}
	result, _ := strconv.ParseFloat(formattedValue, 64)
	return result
}

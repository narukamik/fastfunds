package util

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// DecimalStringToPennies parses a decimal string with up to 2 fractional digits into pennies (int64).
// Examples: "10" -> 1000, "10.2" -> 1020, "10.23" -> 1023. Rejects more than 2 decimals.
func DecimalStringToPennies(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty amount")
	}
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}
	parts := strings.SplitN(s, ".", 3)
	if len(parts) > 2 {
		return 0, errors.New("invalid amount format")
	}
	whole := parts[0]
	frac := ""
	if len(parts) == 2 {
		frac = parts[1]
		if len(frac) > 2 {
			return 0, errors.New("too many decimal places; max 2")
		}
	}
	if whole == "" {
		whole = "0"
	}
	wholeNum, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return 0, errors.New("invalid whole part")
	}
	var fracNum int64
	if frac != "" {
		if len(frac) == 1 {
			frac += "0"
		}
		f, err := strconv.ParseInt(frac, 10, 64)
		if err != nil {
			return 0, errors.New("invalid fractional part")
		}
		fracNum = f
	}
	pennies := wholeNum*100 + fracNum
	if neg {
		pennies = -pennies
	}
	return pennies, nil
}

// PenniesToDecimalString formats pennies (int64) into a string with 2 decimal places.
func PenniesToDecimalString(pennies int64) string {
	neg := pennies < 0
	if neg {
		pennies = -pennies
	}
	whole := pennies / 100
	frac := pennies % 100
	s := fmt.Sprintf("%d.%02d", whole, frac)
	if neg {
		return "-" + s
	}
	return s
}

// SafeMulPercent computes amount * percent (basis points) safely using integers.
func SafeMulPercent(pennies int64, basisPoints int64) int64 {
	// round half up
	return int64(math.Round(float64(pennies*basisPoints) / 10000.0))
}

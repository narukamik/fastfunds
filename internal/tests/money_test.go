package tests

import (
	"fastfunds/internal/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecimalStringToPennies(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    int64
		wantErr string
	}{
		{"int_zero", "0", 0, ""},
		{"int_positive", "10", 1000, ""},
		{"one_decimal", "10.2", 1020, ""},
		{"two_decimals", "10.23", 1023, ""},
		{"pad_one_decimal", "1.5", 150, ""},
		{"negative", "-1.23", -123, ""},
		{"spaces", "  2.50 \t", 250, ""},
		{"too_many_decimals", "1.234", 0, "too many decimal places; max 2"},
		{"invalid_whole", "a.23", 0, "invalid whole part"},
		{"invalid_frac", "1.ab", 0, "invalid fractional part"},
		{"empty", "", 0, "empty amount"},
		{"just_dot", ".", 0, "invalid format"},
		{"leading_dot", ".5", 50, ""},
	}

	converter := util.DefaultMoneyConverter{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := converter.DecimalStringToPennies(tc.in)
			if tc.wantErr != "" {
				assert.Error(t, err)
				if tc.wantErr != "any error" {
					assert.Contains(t, err.Error(), tc.wantErr)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPenniesToDecimalString(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0.00"},
		{1, "0.01"},
		{10, "0.10"},
		{123, "1.23"},
		{-123, "-1.23"},
		{100000, "1000.00"},
	}

	converter := util.DefaultMoneyConverter{}
	for _, tc := range cases {
		got := converter.PenniesToDecimalString(tc.in)
		assert.Equal(t, tc.want, got)
	}
}

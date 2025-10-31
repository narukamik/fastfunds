package util

// MoneyConverter abstracts conversions between decimal strings and pennies.
type MoneyConverter interface {
    DecimalStringToPennies(s string) (int64, error)
    PenniesToDecimalString(pennies int64) string
}

// DefaultMoneyConverter is a concrete adapter that delegates to package-level functions.
type DefaultMoneyConverter struct{}

func (DefaultMoneyConverter) DecimalStringToPennies(s string) (int64, error) {
    return DecimalStringToPennies(s)
}

func (DefaultMoneyConverter) PenniesToDecimalString(pennies int64) string {
    return PenniesToDecimalString(pennies)
}

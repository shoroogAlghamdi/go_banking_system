package util

const (
	SAR = "SAR"
	USD = "USD"
	EUR = "EUR"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case SAR, USD, EUR:
		return true
	}
	return false
}
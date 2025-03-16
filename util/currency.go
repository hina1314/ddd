package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	CNY = "CNY"
)

//IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, CNY:
		return true
	}
	return false
}

package util

const (
	USD = "USD"
	EUR = "EUR"
	UZS = "UZS"
	HUF = "HUF"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case UZS, USD, EUR, HUF:
		return true
	default:
		return false
	}
}

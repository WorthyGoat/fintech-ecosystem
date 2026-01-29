package currency

import (
	"fmt"
	"strings"
)

var supportedCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"GBP": true,
	"JPY": true,
	"EGP": true,
	"CAD": true,
	"AUD": true,
}

// IsSupported checks if the currency code is supported.
func IsSupported(code string) bool {
	return supportedCurrencies[strings.ToUpper(code)]
}

// Validate returns an error if the currency is not supported.
func Validate(code string) error {
	if !IsSupported(code) {
		return fmt.Errorf("unsupported currency: %s", code)
	}
	return nil
}

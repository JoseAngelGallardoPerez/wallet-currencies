package controllers

func Providers() []interface{} {
	return []interface{}{
		NewCurrencies,
		NewRatesController,
		NewSettingsController,
	}
}

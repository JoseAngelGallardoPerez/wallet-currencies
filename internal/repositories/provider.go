package repositories

func Providers() []interface{} {
	return []interface{}{
		NewCurrencies,
		NewRates,
		NewSettings,
	}
}

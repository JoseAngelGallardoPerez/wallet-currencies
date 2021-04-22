package serializers

func Providers() []interface{} {
	return []interface{}{
		Settings,
		Rates,
		Currency,
	}
}

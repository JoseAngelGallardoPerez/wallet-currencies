package config

func Providers() []interface{} {
	return []interface{}{
		NewGeneral,
		NewDb,
		NewCors,
		NewProtoBuf,
	}
}

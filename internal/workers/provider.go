package workers

func Providers() []interface{} {
	return []interface{}{
		NewJobs,
		NewUpdateRates,
	}
}

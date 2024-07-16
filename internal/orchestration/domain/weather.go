package domain

type LocaleResult struct {
	Locale string
}

type WeatherResult struct {
	TempC float64
}

type Result struct {
	TempC float64
	TempF float64
	TempK float64
}

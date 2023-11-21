package converter

// var multiplier = float64(config.GetConfig().Multiplier)
const multiplier = 100000

func ConvertToCent(amount float64) int {
	return int(amount * multiplier)
}

func ConvertFromCent(amount int) float64 {
	return float64(amount) / multiplier
}

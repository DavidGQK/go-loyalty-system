package converter

import "github.com/DavidGQK/go-loyalty-system/internal/config"

var multiplier = float64(config.GetConfig().Multiplier)

func ConvertToCent(amount float64) int {
	return int(amount * multiplier)
}

func ConvertFromCent(amount int) float64 {
	return float64(amount) / multiplier
}

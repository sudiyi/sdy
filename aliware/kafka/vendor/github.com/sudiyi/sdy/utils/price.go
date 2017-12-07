package utils

func Yuan2Fen(yuan float64) int {
	return int(yuan * 100.0)
}

func Fen2Yuan(fen int) float64 {
	return float64(fen) / 100.0
}

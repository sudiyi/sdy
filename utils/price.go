package utils

func YuanToFen(yuan float64) int {
	return int(yuan * 100.0)
}

func FenToYuan(fen int) float64 {
	return float64(fen) / 100.0
}
